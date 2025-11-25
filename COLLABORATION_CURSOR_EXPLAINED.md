# 协作光标功能实现原理详解

## 整体架构

协作光标功能基于以下技术栈：
1. **Yjs** - CRDT（无冲突复制数据类型）库，用于同步文档状态
2. **y-websocket** - WebSocket 提供器，用于实时通信
3. **Tiptap CollaborationCaret** - 协作光标扩展
4. **Awareness API** - Yjs 的感知系统，用于同步用户状态（光标位置、用户信息等）

## 工作流程

### 1. 初始化阶段

```typescript
// 1. 创建 Yjs 文档（CRDT 数据结构）
ydoc = new Y.Doc();

// 2. 创建 WebSocket 提供器，连接到服务器
provider = new WebsocketProvider(WS_URL, documentId.value, ydoc, {
  connect: true,
});

// 3. 配置编辑器，添加协作扩展
editor.value = new Editor({
  extensions: [
    Collaboration.configure({
      document: ydoc,  // 将编辑器内容与 Y.Doc 绑定
    }),
    CollaborationCaret.configure({
      provider: provider,  // 提供器用于同步光标位置
      user: userInfo,      // 当前用户信息
      render: (user) => {  // 自定义光标渲染函数
        // 创建光标 DOM 元素
      }
    }),
  ],
});
```

### 2. 光标位置同步机制

#### 前端（客户端）

**CollaborationCaret 扩展的工作流程：**

1. **监听本地光标移动**
   - 当用户在编辑器中移动光标时，Tiptap 会触发 `selectionUpdate` 事件
   - CollaborationCaret 扩展捕获这个事件，获取当前光标位置

2. **通过 Awareness API 广播光标位置**
   ```typescript
   // CollaborationCaret 内部实现（简化版）
   editor.on('selectionUpdate', () => {
     const selection = editor.state.selection;
     const position = selection.anchor; // 光标位置
     
     // 通过 awareness 更新本地状态
     awareness.setLocalStateField('cursor', {
       position: position,
       user: userInfo
     });
   });
   ```

3. **接收其他用户的光标位置**
   - Awareness API 会自动同步所有客户端的状态
   - 当其他用户的光标位置更新时，会触发 `awareness.on('change')` 事件
   - CollaborationCaret 扩展监听这些变化，并在编辑器中渲染其他用户的光标

#### 后端（服务器）

**WebSocket 服务器的职责：**

```javascript
// 服务器不直接处理光标位置，而是转发 awareness 数据
// y-websocket 协议会自动处理 awareness 同步

wss.on('connection', (ws, req) => {
  // 1. 创建或获取文档房间
  let room = rooms.get(docName);
  
  // 2. 监听文档更新（包括 awareness 数据）
  room.doc.on('update', (update, origin) => {
    // 3. 将更新转发给房间内的其他客户端
    room.clients.forEach((client) => {
      if (client !== ws && client.readyState === 1) {
        client.send(update); // 二进制数据
      }
    });
  });
});
```

**关键点：**
- 服务器只负责转发二进制数据，不解析内容
- Yjs 使用高效的二进制协议（CRDT 更新）
- Awareness 数据与文档数据一起通过 WebSocket 传输

### 3. 光标渲染机制

#### 自定义渲染函数

```typescript
CollaborationCaret.configure({
  render: (user) => {
    // 这个函数为每个用户创建一个光标 DOM 元素
    const cursor = document.createElement("span");
    cursor.classList.add("collaboration-caret");
    cursor.setAttribute("data-user", user.name);
    cursor.setAttribute("data-color", user.color);
    
    // 设置光标样式（彩色竖线）
    cursor.style.borderLeftColor = user.color;
    cursor.style.borderLeftWidth = "2px";
    cursor.style.borderLeftStyle = "solid";
    cursor.style.height = "1.2em";
    
    return cursor;
  }
})
```

**工作流程：**

1. **获取所有在线用户的状态**
   ```typescript
   const awareness = provider.awareness;
   const states = awareness.getStates(); // Map<clientId, state>
   ```

2. **为每个用户创建光标元素**
   - CollaborationCaret 遍历所有用户状态
   - 对每个用户调用 `render(user)` 函数
   - 将返回的 DOM 元素插入到编辑器的对应位置

3. **实时更新光标位置**
   - 当用户移动光标时，awareness 状态更新
   - CollaborationCaret 监听变化，重新计算光标位置
   - 使用 ProseMirror 的装饰系统（Decorations）将光标元素插入到正确位置

#### CSS 样式

```css
.collaboration-caret {
  position: relative;
  border-left: 2px solid;  /* 彩色竖线 */
  height: 1.2em;
  display: inline-block;
}

.collaboration-caret::before {
  content: attr(data-user);  /* 显示用户名 */
  position: absolute;
  top: -1.5em;
  background-color: var(--user-color);
  /* ... 其他样式 */
}
```

### 4. Awareness API 详解

**Awareness 是什么？**

Awareness 是 Yjs 提供的感知系统，用于同步**临时状态**（与文档内容不同）：
- 光标位置
- 用户信息（名称、颜色）
- 选择范围
- 其他临时 UI 状态

**为什么需要 Awareness？**

- 文档内容（CRDT）需要持久化和冲突解决
- 光标位置是临时状态，不需要持久化
- Awareness 数据会在用户断开连接后自动清除

**Awareness 数据结构：**

```typescript
awareness.setLocalStateField('user', {
  name: '用户 123',
  color: '#FF6B6B'
});

awareness.setLocalStateField('cursor', {
  position: 42,  // 光标在文档中的位置
  selection: { from: 40, to: 45 }  // 选择范围（如果有）
});
```

### 5. 数据同步流程

```
用户A移动光标
    ↓
CollaborationCaret 捕获事件
    ↓
更新本地 awareness 状态
    ↓
WebSocket 发送 awareness 更新（二进制）
    ↓
服务器接收并转发给其他客户端
    ↓
用户B的客户端接收更新
    ↓
Yjs 更新 awareness 状态
    ↓
触发 awareness.on('change') 事件
    ↓
CollaborationCaret 重新渲染光标
    ↓
用户B看到用户A的光标
```

### 6. 关键技术点

#### CRDT（无冲突复制数据类型）

- **特点**：多个用户同时编辑不会产生冲突
- **原理**：使用数学算法（如 Yjs 的 YATA）确保操作的可交换性
- **优势**：无需中央服务器解决冲突，每个客户端都能独立合并更新

#### 二进制协议

- Yjs 使用高效的二进制编码（比 JSON 小得多）
- 支持增量更新（只传输变化的部分）
- 自动压缩和优化

#### ProseMirror 装饰系统

- Tiptap 基于 ProseMirror
- 使用 Decorations API 在文档中插入非内容元素（如光标）
- 装饰不会影响文档内容，只是视觉展示

## 代码关键部分解析

### 1. 用户信息配置

```typescript
const userInfo = {
  name: `用户 ${Math.floor(Math.random() * 1000)}`,
  color: generateUserColor(),  // 从预定义颜色列表中选择
};
```

### 2. 在线用户列表

```typescript
const updateUserCount = () => {
  const states = awareness.getStates();
  onlineUsersCount.value = states.size;
  
  // 提取所有用户信息
  const users = [];
  states.forEach((state) => {
    if (state.user) {
      users.push({
        name: state.user.name,
        color: state.user.color,
      });
    }
  });
  onlineUsers.value = users;
};

awareness.on("change", updateUserCount);
```

### 3. 服务器转发逻辑

```javascript
// 服务器端：简单的消息转发
room.doc.on('update', (update, origin) => {
  // origin 是发送更新的 WebSocket 连接
  // 只转发给其他客户端，避免回环
  room.clients.forEach((client) => {
    if (client !== ws && client.readyState === 1) {
      client.send(update);  // 直接转发二进制数据
    }
  });
});
```

## 性能优化

1. **增量更新**：只传输变化的部分，不是整个文档
2. **二进制编码**：比 JSON 更高效
3. **本地优先**：操作立即在本地生效，然后异步同步
4. **智能合并**：CRDT 算法自动处理并发编辑

## 总结

协作光标功能的实现依赖于：
1. **Yjs CRDT** - 处理数据同步和冲突解决
2. **Awareness API** - 同步临时状态（光标位置）
3. **WebSocket** - 实时通信通道
4. **CollaborationCaret 扩展** - 监听和渲染光标
5. **ProseMirror Decorations** - 在编辑器中插入视觉元素

整个系统设计为**去中心化**的，每个客户端都是平等的，服务器只负责转发消息，不维护状态。这使得系统具有高可用性和低延迟。

