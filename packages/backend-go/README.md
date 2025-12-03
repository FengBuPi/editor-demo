# Go后端实时协作服务器

基于Go + WebSocket + Redis的Yjs实时协作后端服务。

## 功能特性

- ✅ 支持多文档实时同步
- ✅ Redis持久化历史更新
- ✅ 新用户自动同步历史状态
- ✅ 心跳检测和连接管理
- ✅ 自动清理空闲房间
- ✅ 优雅关闭

## 快速开始

### 1. 安装依赖

```bash
cd packages/backend-go
go mod tidy
```

### 2. 启动Redis

确保Redis运行在 `localhost:6379`，无密码。

```bash
# Windows
redis-server

# Linux/Mac
redis-server
```

### 3. 启动服务器

```bash
go run main.go
```

或编译后运行：

```bash
go build -o yjs-server
./yjs-server
```

### 4. 验证

服务器启动后会显示：
```
Redis连接成功
WebSocket服务器运行在 ws://localhost:3001
连接格式: ws://localhost:3001/文档名
```

## 配置

默认配置：
- 端口：3001
- Redis地址：localhost:6379
- 最大历史更新数：10000
- 心跳间隔：30秒
- 连接超时：60秒

修改配置请编辑 `main.go` 中的常量。

## 协议

兼容 `y-websocket` 协议：

- 连接格式：`ws://localhost:3001/文档名`
- 消息类型：二进制（BinaryMessage）
- 数据格式：Yjs编码的二进制更新

## 架构

```
前端 (Tiptap + Yjs)
    ↓ WebSocket
Go服务器
    ↓ 存储/读取
Redis (历史更新)
```

## 注意事项

1. 确保Redis已启动
2. 前端连接地址需匹配：`ws://localhost:3001`
3. 文档ID通过URL路径传递

