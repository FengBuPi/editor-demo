# 后端服务器

基于 Node.js 原生 WebSocket 的 Yjs 协作服务器。

## 功能

- ✅ 支持多个文档的实时同步
- ✅ 自动处理文档状态同步
- ✅ 支持客户端断线重连
- ✅ 轻量级，无额外依赖（仅需 yjs）

## 启动

```bash
# 开发模式（自动重启）
pnpm dev

# 生产模式
pnpm start
```

## 配置

默认端口：`3001`

可以通过环境变量修改：

```bash
PORT=3001 pnpm start
```

## WebSocket 连接

客户端连接地址：

```
ws://localhost:3001/ws?doc=文档ID
```

如果不提供 `doc` 参数，将使用默认文档 `default`。

## 协议

服务器使用 JSON 格式的消息：

### 同步消息（服务器 -> 客户端）
```json
{
  "type": "sync",
  "stateVector": [1, 2, 3, ...],
  "update": [1, 2, 3, ...]
}
```

### 更新消息（客户端 <-> 服务器）
```json
{
  "type": "update",
  "update": [1, 2, 3, ...]
}
```

