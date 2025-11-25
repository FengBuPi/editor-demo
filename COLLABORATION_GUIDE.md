# Tiptap + Yjs 协作编辑指南

## 概述

本指南介绍如何使用 Yjs 和 y-websocket 在 Tiptap 编辑器中实现实时协作编辑功能。

## 安装依赖

安装所有必需的依赖包：

```bash
pnpm add @tiptap/extension-collaboration @tiptap/extension-collaboration-cursor @tiptap/y-tiptap yjs y-protocols y-websocket
```

## 实现步骤

### 1. 搭建 WebSocket 服务器

可以使用 y-websocket 提供的服务器：

```bash
# 安装 y-websocket 服务器
npm install -g y-websocket

# 启动服务器（默认端口 1234）
y-websocket
```

或者使用 Node.js 创建自定义服务器：

```javascript
// server.js
const WebSocket = require('ws');
const http = require('http');
const wss = new WebSocket.Server({ port: 1234 });

const docs = new Map();

wss.on('connection', (ws) => {
  ws.on('message', (message) => {
    // 处理 Yjs 同步消息
    wss.clients.forEach((client) => {
      if (client !== ws && client.readyState === WebSocket.OPEN) {
        client.send(message);
      }
    });
  });
});
```

### 2. 在 Vue 组件中使用

```vue
<script setup>
import { onMounted, onBeforeUnmount } from 'vue';
import { useEditor, EditorContent } from '@tiptap/vue-3';
import StarterKit from '@tiptap/starter-kit';
import Collaboration from '@tiptap/extension-collaboration';
import CollaborationCursor from '@tiptap/extension-collaboration-cursor';
import { WebsocketProvider } from 'y-websocket';
import * as Y from 'yjs';

const documentId = 'my-document-id'; // 文档标识符
const ydoc = new Y.Doc();
let provider: WebsocketProvider | null = null;

const editor = useEditor({
  extensions: [
    StarterKit.configure({
      history: false, // 禁用默认历史，使用协作版本
    }),
    Collaboration.configure({
      document: ydoc,
    }),
    CollaborationCursor.configure({
      provider: () => provider!,
      user: {
        name: '用户名',
        color: '#ff0000',
      },
    }),
  ],
});

onMounted(() => {
  // 创建 WebSocket Provider
  provider = new WebsocketProvider('ws://localhost:1234', documentId, ydoc);
  
  // 监听连接状态
  provider.on('status', (event) => {
    console.log('连接状态:', event.status);
  });
});

onBeforeUnmount(() => {
  if (provider) {
    provider.destroy();
  }
  editor.value?.destroy();
  ydoc.destroy();
});
</script>
```

## 核心概念

### 1. Y.Doc

Yjs 的核心数据结构，用于存储和同步文档状态。

```javascript
const ydoc = new Y.Doc();
```

### 2. WebsocketProvider

WebsocketProvider 负责通过 WebSocket 在不同客户端之间同步 Y.Doc 的状态。

### 3. Collaboration 扩展

Tiptap 的协作扩展，将编辑器内容与 Y.Doc 绑定。

```javascript
Collaboration.configure({
  document: ydoc,
})
```

### 4. CollaborationCursor 扩展

显示其他用户的光标位置和选择。

```javascript
CollaborationCursor.configure({
  provider: () => provider,
  user: {
    name: '用户名',
    color: '#ff0000',
  },
})
```

## 完整示例

参考 `src/components/TiptapEditorCollaboration.vue` 文件，包含：

- ✅ 基础协作编辑
- ✅ 协作光标显示
- ✅ 连接状态监控
- ✅ 在线用户数显示

## 测试协作功能

1. 启动 WebSocket 服务器：
   ```bash
   y-websocket
   ```

2. 启动你的 Vue 应用：
   ```bash
   pnpm dev
   ```

3. 打开多个浏览器标签页，访问相同的编辑器页面

4. 在不同标签页中编辑内容，观察实时同步效果

## 注意事项

1. **禁用默认 History**: 使用协作功能时，需要禁用 StarterKit 中的 history 扩展
2. **文档 ID**: 相同文档 ID 的用户会看到相同的文档内容
3. **连接状态**: 建议监听 provider 的连接状态，处理断线重连
4. **清理资源**: 组件销毁时记得清理 provider 和 ydoc

## 更多资源

- [Tiptap 协作文档](https://tiptap.dev/docs/collaboration)
- [Yjs 官方文档](https://docs.yjs.dev/)
- [y-websocket 文档](https://github.com/yjs/y-websocket)

