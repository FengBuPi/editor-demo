# Go后端实时协作完整方案

## 一、y-websocket协议详解

### 1.1 协议特点
- **二进制协议**：所有消息都是Uint8Array格式
- **URL路由**：文档ID通过URL路径传递 `ws://host:port/文档名`
- **双向同步**：客户端发送更新，服务器广播给其他客户端
- **状态同步**：新客户端连接时，服务器发送当前文档完整状态

### 1.2 消息流程

#### 客户端连接
```
客户端 → ws://localhost:3001/demo-document-1
服务器 → 发送当前文档状态（二进制）
```

#### 客户端发送更新
```
客户端 → 发送Yjs更新（二进制）
服务器 → 应用更新到文档 → 广播给其他客户端
```

#### 协议细节
- **消息类型**：WebSocket BinaryMessage
- **数据格式**：Yjs编码的二进制更新
- **状态管理**：服务器维护每个文档的Y.Doc实例

## 二、Go实现方案对比

### 方案1：纯转发服务器（不推荐）
**原理**：只转发二进制消息，不解析Yjs内容

**优点**：
- 实现简单
- 性能高

**缺点**：
- 无法发送初始状态给新客户端
- 无法持久化
- 客户端断线重连会丢失数据

**适用场景**：仅用于测试，不适合生产

---

### 方案2：Hocuspocus中间层（推荐用于快速上线）
**原理**：Go后端调用Node.js的Hocuspocus服务

**架构**：
```
前端 → Go后端 → Hocuspocus (Node.js) → 数据库
```

**优点**：
- 快速实现，Hocuspocus已实现完整功能
- 支持持久化、权限控制
- 生产级稳定性

**缺点**：
- 需要运行Node.js服务
- 增加一层调用开销

**实现步骤**：
1. 部署Hocuspocus服务（独立进程或Docker）
2. Go后端通过HTTP/WebSocket代理转发请求
3. 或Go后端直接调用Hocuspocus的WebSocket端口

---

### 方案3：Go原生实现（推荐用于长期项目）
**原理**：在Go中实现y-websocket协议

**技术栈**：
- `github.com/gorilla/websocket` - WebSocket支持
- `github.com/hashicorp/golang-lru` - LRU缓存（可选）
- 需要实现Yjs二进制协议解析

**核心挑战**：
1. **Yjs协议解析**：Go没有官方Yjs实现，需要：
   - 解析Yjs二进制格式
   - 实现CRDT合并逻辑
   - 或使用第三方Go CRDT库

2. **状态管理**：
   - 维护每个文档的完整状态
   - 新客户端连接时发送完整状态
   - 处理更新合并

**实现方案A：使用Go CRDT库**
- 查找Go生态的CRDT实现（如`github.com/ipfs/go-ds-crdt`）
- 实现Yjs协议兼容层
- 工作量：中等

**实现方案B：最小化实现**
- 只实现y-websocket的消息转发
- 使用Redis/数据库存储文档状态
- 客户端负责状态合并
- 工作量：较小

**实现方案C：完整实现**
- 移植Yjs核心逻辑到Go
- 实现完整的CRDT合并
- 工作量：大（需要深入理解Yjs）

---

## 三、推荐方案详细设计

### 方案3B：最小化实现（推荐）

#### 3.1 架构设计
```
┌─────────┐      WebSocket      ┌──────────┐
│ 前端    │ ←─────────────────→ │ Go服务器 │
│ Tiptap  │                      │          │
└─────────┘                      └────┬─────┘
                                      │
                                      │ 存储/同步
                                      ↓
                              ┌──────────────┐
                              │ Redis/数据库  │
                              └──────────────┘
```

#### 3.2 数据结构
```go
type Room struct {
    docName    string
    clients    map[*websocket.Conn]*Client
    docState   []byte  // 当前文档状态（Yjs编码）
    updates    [][]byte // 更新历史（可选，用于重放）
    mu         sync.RWMutex
}

type Client struct {
    conn      *websocket.Conn
    docName   string
    lastSync  time.Time
}
```

#### 3.3 核心流程

**连接处理**：
1. 从URL提取文档名：`/demo-document-1` → `demo-document-1`
2. 创建/获取Room
3. 发送当前文档状态（从Redis/内存获取）
4. 添加到客户端列表

**消息处理**：
1. 接收客户端二进制更新
2. 存储到Redis（可选，用于持久化）
3. 更新内存中的docState
4. 广播给房间内其他客户端

**状态同步**：
- 新客户端：发送完整docState
- 现有客户端：只发送增量更新

#### 3.4 持久化方案

**选项1：Redis**
```go
// 存储格式
key: "ydoc:{docName}"
value: 二进制状态（[]byte）
```

**选项2：数据库**
```sql
CREATE TABLE documents (
    id VARCHAR(255) PRIMARY KEY,
    state BLOB,
    updated_at TIMESTAMP
);
```

**选项3：内存+定期快照**
- 内存存储活跃文档
- 定期保存到文件/数据库

#### 3.5 依赖库
```go
require (
    github.com/gorilla/websocket v1.5.1
    github.com/redis/go-redis/v9 v9.3.0  // 可选
    github.com/go-sql-driver/mysql v1.7.1 // 可选
)
```

---

## 四、完整实现代码结构

```
packages/backend-go/
├── main.go              # 主入口
├── server/
│   ├── websocket.go     # WebSocket处理
│   ├── room.go          # 房间管理
│   └── storage.go       # 持久化接口
├── storage/
│   ├── memory.go        # 内存存储
│   ├── redis.go         # Redis存储
│   └── database.go      # 数据库存储
├── protocol/
│   └── yjs.go           # Yjs协议处理（如果需要）
├── go.mod
└── go.sum
```

---

## 五、关键技术点

### 5.1 WebSocket升级
```go
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true  // 生产环境需要验证
    },
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}
```

### 5.2 文档名提取
```go
// URL格式: ws://host:port/文档名
docName := strings.TrimPrefix(r.URL.Path, "/")
if docName == "" {
    docName = "default"
}
```

### 5.3 二进制消息处理
```go
messageType, message, err := conn.ReadMessage()
if messageType == websocket.BinaryMessage {
    // 处理Yjs更新
    update := make([]byte, len(message))
    copy(update, message)
    // 广播给其他客户端
}
```

### 5.4 并发安全
- 使用`sync.RWMutex`保护Room状态
- 客户端连接使用`sync.Map`或带锁的map
- 广播时使用goroutine避免阻塞

### 5.5 错误处理
- WebSocket连接断开处理
- 消息解析错误处理
- 存储失败降级（内存模式）

---

## 六、性能优化

### 6.1 连接管理
- 心跳检测：定期ping/pong
- 超时断开：长时间无活动自动断开
- 连接池：限制每个文档的最大连接数

### 6.2 消息优化
- 批量更新：合并短时间内的多个更新
- 压缩：对大型文档状态使用gzip压缩
- 增量同步：只发送差异，不发送完整状态

### 6.3 存储优化
- LRU缓存：热点文档保持在内存
- 异步持久化：更新先写内存，异步写存储
- 分片存储：大文档分块存储

---

## 七、部署方案

### 7.1 单机部署
```bash
go build -o yjs-server
./yjs-server
```

### 7.2 Docker部署
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/server /server
CMD ["/server"]
```

### 7.3 集群部署
- 使用Redis Pub/Sub实现多实例同步
- 或使用消息队列（如NATS、RabbitMQ）
- 负载均衡：Nginx/HAProxy做WebSocket代理

---

## 八、监控和调试

### 8.1 指标收集
- 连接数：每个文档的活跃连接
- 消息吞吐：每秒处理的消息数
- 延迟：消息处理延迟
- 错误率：连接错误、消息错误

### 8.2 日志
- 连接建立/断开
- 文档创建/更新
- 错误详情
- 性能指标

---

## 九、方案选择建议

### 快速上线（1-2周）
→ **方案2：Hocuspocus中间层**
- 使用现有Node.js服务
- Go做业务层，Hocuspocus做协作层

### 长期项目（1-2月）
→ **方案3B：最小化Go实现**
- 纯Go实现，无Node.js依赖
- 使用Redis做持久化
- 逐步完善功能

### 极致性能（3-6月）
→ **方案3C：完整Go实现**
- 实现完整Yjs协议
- 深度优化性能
- 支持所有高级特性

---

## 十、注意事项

1. **Yjs协议复杂性**：Yjs使用复杂的CRDT算法，简单转发可能丢失数据
2. **状态一致性**：多实例部署需要解决状态同步问题
3. **内存管理**：大量文档会占用大量内存，需要LRU淘汰
4. **安全性**：WebSocket需要验证来源，防止未授权访问
5. **扩展性**：设计时考虑水平扩展能力

---

## 十一、参考资源

- Yjs文档：https://docs.yjs.dev/
- y-websocket源码：https://github.com/yjs/y-websocket
- Hocuspocus文档：https://tiptap.dev/docs/hocuspocus
- gorilla/websocket：https://github.com/gorilla/websocket

