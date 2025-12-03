# 使用说明

## 前端是否需要修改？

**答案：不需要修改！** 前端代码已经正确配置。

### 前端连接方式

前端使用 `WebsocketProvider`，它会自动构建URL：
```typescript
provider = new WebsocketProvider(WS_URL, documentId.value, ydoc)
// WS_URL = "ws://localhost:3001"
// documentId = "demo-document-1"
// 实际连接: ws://localhost:3001/demo-document-1
```

Go后端已经实现了从URL路径提取文档名，完全兼容。

## 启动步骤

### 1. 确保Redis运行

```bash
# 检查Redis是否运行
redis-cli ping
# 应该返回: PONG
```

如果没有Redis，请先安装并启动：
- Windows: 下载Redis for Windows
- Linux: `sudo apt install redis-server && redis-server`
- Mac: `brew install redis && redis-server`

### 2. 启动Go后端

```bash
cd packages/backend-go
go run main.go
```

应该看到：
```
Redis连接成功
WebSocket服务器运行在 ws://localhost:3001
连接格式: ws://localhost:3001/文档名
```

### 3. 启动前端

```bash
# 在项目根目录
pnpm dev
```

### 4. 测试

1. 打开浏览器访问前端（通常是 http://localhost:5173）
2. 打开"协作编辑器"标签页
3. 打开第二个浏览器标签页，访问相同页面
4. 在一个标签页输入内容，另一个标签页应该实时同步

## 故障排查

### Redis连接失败

**错误**：`Redis连接失败: dial tcp :6379: connectex: No connection could be made`

**解决**：
1. 确认Redis已启动：`redis-cli ping`
2. 检查Redis地址是否正确（默认 localhost:6379）
3. 如果Redis有密码，需要修改 `main.go` 中的配置

### WebSocket连接失败

**错误**：前端显示"连接错误"

**解决**：
1. 确认Go后端已启动
2. 确认端口是3001
3. 检查浏览器控制台的错误信息
4. 确认前端连接地址是 `ws://localhost:3001`

### 数据不同步

**现象**：两个用户编辑，但看不到对方的修改

**解决**：
1. 确认两个用户使用相同的文档ID
2. 检查后端日志，看是否有错误
3. 检查Redis是否正常工作
4. 清除浏览器缓存，重新连接

## 性能优化

### 历史更新过多

如果文档编辑历史很长，可能影响性能。可以：

1. **限制历史数量**：修改 `main.go` 中的 `maxHistorySize`
2. **定期清理**：已实现自动清理30分钟无活动的房间
3. **使用快照**：定期保存完整状态，而不是所有历史更新

### Redis内存占用

历史更新存储在Redis中，如果文档很多，可能占用大量内存。

**优化方案**：
1. 设置Redis过期时间（已实现，7天）
2. 定期清理旧文档
3. 使用Redis持久化配置

## 生产环境部署

### 1. 编译

```bash
cd packages/backend-go
go build -o yjs-server
```

### 2. 配置环境变量

可以修改代码中的常量，或使用配置文件：
- 端口
- Redis地址
- 最大历史数
- 超时时间

### 3. 使用进程管理器

**使用systemd（Linux）**：
```ini
[Unit]
Description=Yjs WebSocket Server
After=network.target redis.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/path/to/backend-go
ExecStart=/path/to/backend-go/yjs-server
Restart=always

[Install]
WantedBy=multi-user.target
```

**使用PM2（Node.js进程管理器，也支持Go）**：
```bash
pm2 start yjs-server --name yjs-server
```

### 4. 使用Docker

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o yjs-server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/yjs-server /yjs-server
CMD ["/yjs-server"]
```

## 监控

### 日志

后端会输出以下日志：
- 新客户端连接
- 客户端断开
- Redis操作错误
- WebSocket错误

### 指标

可以添加以下监控指标：
- 活跃连接数
- 每个文档的连接数
- 消息处理延迟
- Redis操作延迟

## 安全建议

1. **CORS配置**：生产环境需要配置 `CheckOrigin`
2. **认证**：添加JWT或Token验证
3. **限流**：防止恶意连接
4. **Redis密码**：生产环境使用Redis密码
5. **HTTPS/WSS**：使用加密连接

