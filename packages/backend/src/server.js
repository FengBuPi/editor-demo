import http from 'http';
import { WebSocketServer } from 'ws';
import * as Y from 'yjs';

const PORT = process.env.PORT || 3001;

// 创建 HTTP 服务器
const server = http.createServer();

// 创建 WebSocket 服务器（兼容 y-websocket 协议）
const wss = new WebSocketServer({
  server,
  // 不设置 path，让所有路径都可以连接
});

// 存储每个文档的状态和客户端连接
// Map<docName, { doc: Y.Doc, clients: Set<WebSocket> }>
const rooms = new Map();

// 处理 WebSocket 连接
wss.on('connection', (ws, req) => {
  console.log('========================================');
  console.log('新的 WebSocket 连接');
  console.log('请求 URL:', req.url);
  console.log('请求头:', req.headers);

  // 从 URL 路径中提取文档 ID
  // y-websocket 格式: ws://host:port/文档名
  const url = new URL(req.url, `http://${req.headers.host}`);
  const docName = url.pathname.slice(1) || 'default'; // 移除前导 '/'

  console.log(`文档名: ${docName}`);
  console.log('WebSocket 状态: OPEN');

  // 获取或创建房间
  let room = rooms.get(docName);
  if (!room) {
    room = {
      doc: new Y.Doc(),
      clients: new Set(),
    };
    rooms.set(docName, room);
    console.log(`创建新文档房间: ${docName}`);
  }

  // 将客户端添加到房间
  room.clients.add(ws);

  // 创建更新函数，用于同步文档状态
  const updateHandler = (update, origin) => {
    if (origin !== ws) {
      // 将二进制更新直接转发给房间内的其他客户端
      room.clients.forEach((client) => {
        if (client !== ws && client.readyState === 1) { // 1 = OPEN
          client.send(update);
        }
      });
    }
  };

  // 监听文档更新
  room.doc.on('update', updateHandler);

  // 发送初始状态给新连接的客户端
  // y-websocket 协议：直接发送 Yjs 的二进制更新
  const stateVector = Y.encodeStateVector(room.doc);
  const update = Y.encodeStateAsUpdate(room.doc, stateVector);

  if (update.length > 0) {
    ws.send(update);
  }

  // 处理接收到的消息（y-websocket 使用二进制协议）
  ws.on('message', (message) => {
    try {
      // y-websocket 发送的是二进制 Uint8Array
      if (Buffer.isBuffer(message)) {
        const update = new Uint8Array(message);
        // 应用更新到文档（origin 设为 ws 以避免循环）
        Y.applyUpdate(room.doc, update, ws);
      } else if (message instanceof Uint8Array) {
        Y.applyUpdate(room.doc, message, ws);
      } else {
        console.warn('收到非二进制消息，忽略');
      }
    } catch (error) {
      console.error('处理消息错误:', error);
    }
  });

  // 处理连接关闭
  ws.on('close', (code, reason) => {
    console.log(`WebSocket 连接关闭 (文档: ${docName})`);
    console.log(`关闭代码: ${code}, 原因: ${reason}`);
    room.doc.off('update', updateHandler);
    room.clients.delete(ws);

    // 如果房间没有客户端了，可以选择保留或删除房间
    if (room.clients.size === 0) {
      console.log(`文档房间 ${docName} 已清空（保留文档状态）`);
      // 这里选择保留文档，以便后续连接可以继续编辑
    }
    console.log('========================================');
  });

  // 处理错误
  ws.on('error', (error) => {
    console.error('WebSocket 错误:', error);
    console.error('错误详情:', error.message);
    room.clients.delete(ws);
  });
});

// 启动服务器
server.listen(PORT, () => {
  console.log(`WebSocket 服务器运行在 ws://localhost:${PORT}`);
  console.log(`连接格式: ws://localhost:${PORT}/文档名`);
  console.log(`示例: ws://localhost:${PORT}/demo-document-1`);
});

// 优雅关闭
process.on('SIGTERM', () => {
  console.log('收到 SIGTERM，正在关闭服务器...');
  wss.close(() => {
    server.close(() => {
      console.log('服务器已关闭');
      process.exit(0);
    });
  });
});

