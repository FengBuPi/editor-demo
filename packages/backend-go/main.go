package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

const (
	port           = "3001"
	redisAddr      = "localhost:6379"
	maxHistorySize = 10000
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	rdb *redis.Client
)

type Room struct {
	docName string
	clients map[*websocket.Conn]*Client
	updates [][]byte
	mu      sync.RWMutex
	created time.Time
	updated time.Time
}

type Client struct {
	conn     *websocket.Conn
	docName  string
	joinedAt time.Time
	lastPing time.Time
	writeMu  sync.Mutex
}

var rooms = sync.Map{}

func main() {
	ctx := context.Background()

	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis连接失败: %v", err)
	}
	log.Println("Redis连接成功")

	http.HandleFunc("/", handleWebSocket)

	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		log.Printf("WebSocket服务器运行在 ws://localhost:%s", port)
		log.Printf("连接格式: ws://localhost:%s/文档名", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	go cleanupRooms()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务器...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("服务器关闭失败: %v", err)
	}

	log.Println("服务器已关闭")
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	docName := r.URL.Path[1:]
	if docName == "" {
		docName = "default"
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket升级失败: %v", err)
		return
	}
	defer conn.Close()

	room := getOrCreateRoom(docName)

	client := &Client{
		conn:     conn,
		docName:  docName,
		joinedAt: time.Now(),
		lastPing: time.Now(),
	}

	room.mu.Lock()
	room.clients[conn] = client
	room.mu.Unlock()

	log.Printf("新客户端连接: 文档=%s, 总连接数=%d", docName, len(room.clients))

	defer func() {
		room.mu.Lock()
		delete(room.clients, conn)
		room.mu.Unlock()
		log.Printf("客户端断开: 文档=%s, 剩余连接数=%d", docName, len(room.clients))
	}()

	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		client.lastPing = time.Now()
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	go pingClient(conn, client)

	if err := syncHistoryToClient(room, conn); err != nil {
		log.Printf("同步历史失败: %v", err)
		return
	}

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket错误: %v", err)
			}
			break
		}

		if messageType == websocket.BinaryMessage {
			if err := handleUpdate(room, message, conn); err != nil {
				log.Printf("处理更新失败: %v", err)
			}
		} else if messageType == websocket.PongMessage {
			client.lastPing = time.Now()
			conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		}
	}
}

func getOrCreateRoom(docName string) *Room {
	roomInterface, _ := rooms.LoadOrStore(docName, &Room{
		docName: docName,
		clients: make(map[*websocket.Conn]*Client),
		updates: make([][]byte, 0),
		created: time.Now(),
		updated: time.Now(),
	})

	room := roomInterface.(*Room)

	ctx := context.Background()
	history, err := loadHistoryFromRedis(ctx, docName)
	if err == nil && len(history) > 0 {
		room.mu.Lock()
		if len(room.updates) == 0 {
			room.updates = history
		}
		room.mu.Unlock()
	}

	return room
}

func syncHistoryToClient(room *Room, conn *websocket.Conn) error {
	room.mu.RLock()
	history := make([][]byte, len(room.updates))
	copy(history, room.updates)
	client, exists := room.clients[conn]
	room.mu.RUnlock()

	if !exists {
		return nil
	}

	client.writeMu.Lock()
	defer client.writeMu.Unlock()

	for _, update := range history {
		if err := conn.WriteMessage(websocket.BinaryMessage, update); err != nil {
			return err
		}
	}

	return nil
}

func handleUpdate(room *Room, update []byte, sender *websocket.Conn) error {
	updateCopy := make([]byte, len(update))
	copy(updateCopy, update)

	room.mu.Lock()
	room.updates = append(room.updates, updateCopy)
	if len(room.updates) > maxHistorySize {
		room.updates = room.updates[len(room.updates)-maxHistorySize:]
	}
	room.updated = time.Now()
	room.mu.Unlock()

	go saveUpdateToRedis(context.Background(), room.docName, updateCopy)

	room.mu.RLock()
	clients := make([]*websocket.Conn, 0, len(room.clients))
	for conn := range room.clients {
		if conn != sender {
			clients = append(clients, conn)
		}
	}
	room.mu.RUnlock()

	for _, conn := range clients {
		go func(c *websocket.Conn) {
			room.mu.RLock()
			client, exists := room.clients[c]
			room.mu.RUnlock()

			if !exists {
				return
			}

			client.writeMu.Lock()
			err := c.WriteMessage(websocket.BinaryMessage, update)
			client.writeMu.Unlock()

			if err != nil {
				log.Printf("广播失败: %v", err)
				room.mu.Lock()
				delete(room.clients, c)
				room.mu.Unlock()
				c.Close()
			}
		}(conn)
	}

	return nil
}

func saveUpdateToRedis(ctx context.Context, docName string, update []byte) {
	key := "ydoc:updates:" + docName
	if err := rdb.RPush(ctx, key, update).Err(); err != nil {
		log.Printf("保存更新到Redis失败: %v", err)
		return
	}

	if err := rdb.Expire(ctx, key, 7*24*time.Hour).Err(); err != nil {
		log.Printf("设置过期时间失败: %v", err)
	}
}

func loadHistoryFromRedis(ctx context.Context, docName string) ([][]byte, error) {
	key := "ydoc:updates:" + docName
	updates, err := rdb.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	result := make([][]byte, 0, len(updates))
	for _, updateStr := range updates {
		result = append(result, []byte(updateStr))
	}

	return result, nil
}

func pingClient(conn *websocket.Conn, client *Client) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if time.Since(client.lastPing) > 60*time.Second {
			conn.Close()
			return
		}
		client.writeMu.Lock()
		err := conn.WriteMessage(websocket.PingMessage, nil)
		client.writeMu.Unlock()
		if err != nil {
			return
		}
	}
}

func cleanupRooms() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rooms.Range(func(key, value interface{}) bool {
			room := value.(*Room)
			room.mu.RLock()
			clientCount := len(room.clients)
			room.mu.RUnlock()

			if clientCount == 0 && time.Since(room.updated) > 30*time.Minute {
				rooms.Delete(key)
				log.Printf("清理空闲房间: %s", room.docName)
			}
			return true
		})
	}
}
