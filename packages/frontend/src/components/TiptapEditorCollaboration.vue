<template>
  <div class="tiptap-editor">
    <div class="toolbar">
      <button
        @click="editor?.chain().focus().toggleBold().run()"
        :class="{ 'is-active': editor?.isActive('bold') }"
        type="button">
        <strong>B</strong>
      </button>
      <button
        @click="editor?.chain().focus().toggleItalic().run()"
        :class="{ 'is-active': editor?.isActive('italic') }"
        type="button">
        <em>I</em>
      </button>
      <div class="divider"></div>
      <button
        @click="editor?.chain().focus().toggleHeading({ level: 1 }).run()"
        :class="{ 'is-active': editor?.isActive('heading', { level: 1 }) }"
        type="button">
        H1
      </button>
      <button
        @click="editor?.chain().focus().toggleHeading({ level: 2 }).run()"
        :class="{ 'is-active': editor?.isActive('heading', { level: 2 }) }"
        type="button">
        H2
      </button>
    </div>
    <div class="editor-wrapper">
      <drag-handle v-if="editor" :editor="editor">
        <div class="drag-handle-icon">⋮⋮</div>
      </drag-handle>
      <editor-content v-if="editor" :editor="editor" class="editor-content" />
    </div>
    <div class="collaboration-info">
      <h3>协作信息</h3>
      <p>文档 ID: {{ documentId }}</p>
      <p>
        连接状态: <span :class="statusClass">{{ connectionStatus }}</span>
      </p>
      <p>在线用户数: {{ onlineUsersCount }}</p>
      <div v-if="onlineUsers.length > 0" class="online-users">
        <h4>在线用户</h4>
        <div class="user-list">
          <div v-for="(user, index) in onlineUsers" :key="index" class="user-item">
            <span class="user-color-indicator" :style="{ backgroundColor: user.color }"></span>
            <span class="user-name">{{ user.name }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, ref, computed, onMounted, shallowRef } from "vue";
import { Editor, EditorContent } from "@tiptap/vue-3";
import StarterKit from "@tiptap/starter-kit";
import Placeholder from "@tiptap/extension-placeholder";
import Collaboration from "@tiptap/extension-collaboration";
import CollaborationCaret from "@tiptap/extension-collaboration-caret";
import DragHandle from "@tiptap/extension-drag-handle-vue-3";
import * as Y from "yjs";
import { WebsocketProvider } from "y-websocket";

const documentId = ref("demo-document-1");

const connectionStatus = ref("未连接");
const onlineUsersCount = ref(0);
const onlineUsers = ref<Array<{ name: string; color: string }>>([]);
const statusClass = computed(() =>
  connectionStatus.value === "已连接" ? "status-connected" : "status-disconnected",
);

const WS_URL = "ws://localhost:3001";

// 生成更易读的用户颜色（避免太亮或太暗）
const generateUserColor = () => {
  const colors = [
    "#FF6B6B", // 红色
    "#4ECDC4", // 青色
    "#45B7D1", // 蓝色
    "#FFA07A", // 浅橙色
    "#98D8C8", // 薄荷绿
    "#F7DC6F", // 黄色
    "#BB8FCE", // 紫色
    "#85C1E2", // 天蓝色
    "#F8B739", // 橙色
    "#52BE80", // 绿色
  ];
  return colors[Math.floor(Math.random() * colors.length)];
};

const userInfo = {
  name: `用户 ${Math.floor(Math.random() * 1000)}`,
  color: generateUserColor(),
};

let provider: WebsocketProvider | null = null;
const editor = shallowRef<Editor>();
let ydoc: Y.Doc | null = null;
onMounted(() => {
  ydoc = new Y.Doc();
  provider = new WebsocketProvider(WS_URL, documentId.value, ydoc, {
    connect: true,
  });
  editor.value = new Editor({
    extensions: [
      StarterKit.configure({
        // 禁用 history 扩展，因为协作编辑使用 CollaborationHistory
        undoRedo: false,
      }),
      Placeholder.configure({
        placeholder: "开始协作编辑...",
      }),
      Collaboration.configure({
        document: ydoc,
      }),
      CollaborationCaret.configure({
        provider: provider,
        user: userInfo,
        // 自定义光标渲染，使其更明显
        render: (user) => {
          const cursor = document.createElement("span");
          cursor.classList.add("collaboration-caret");
          cursor.setAttribute("data-user", user.name);
          cursor.setAttribute("data-color", user.color);
          // 设置光标样式
          cursor.style.borderLeftColor = user.color;
          // 创建用户标签（使用 CSS 变量来传递颜色）
          const style = document.createElement("style");
          style.textContent = `
            .collaboration-caret[data-color="${user.color}"]::before {
              background-color: ${user.color} !important;
            }
          `;
          if (!document.head.querySelector(`style[data-caret-color="${user.color}"]`)) {
            style.setAttribute("data-caret-color", user.color);
            document.head.appendChild(style);
          }
          return cursor;
        },
      }),
    ],
  });

  if (!provider || !ydoc) return;
  const awareness = provider?.awareness;
  if (!awareness) return;

  const updateUserCount = () => {
    const states = awareness.getStates();
    onlineUsersCount.value = states.size;

    // 更新在线用户列表
    const users: Array<{ name: string; color: string }> = [];
    states.forEach((state) => {
      if (state.user) {
        users.push({
          name: state.user.name || "未知用户",
          color: state.user.color || "#999999",
        });
      }
    });
    onlineUsers.value = users;
  };

  updateUserCount();
  awareness.on("change", updateUserCount);

  provider?.on("status", (event: { status: string }) => {
    if (event.status === "connected") {
      connectionStatus.value = "已连接";
    } else if (event.status === "disconnected") {
      connectionStatus.value = "已断开";
    } else if (event.status === "connecting") {
      connectionStatus.value = "连接中...";
    }
  });

  provider.on("sync", (isSynced: boolean) => {
    if (isSynced) {
      connectionStatus.value = "已同步";
    }
  });

  provider.on("connection-close", () => {
    connectionStatus.value = "已断开";
  });

  provider.on("connection-error", (event: Event) => {
    console.error("WebSocket 连接错误:", event);
    connectionStatus.value = "连接错误";
  });
});

onBeforeUnmount(() => {
  if (!provider || !ydoc) return;
  provider.destroy();
  provider = null;
  editor?.value?.destroy();
  ydoc.destroy();
});
</script>

<style scoped>
.tiptap-editor {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.toolbar {
  display: flex;
  gap: 4px;
  padding: 12px;
  background: #f5f5f5;
  border: 1px solid #e0e0e0;
  border-radius: 8px 8px 0 0;
  flex-wrap: wrap;
  align-items: center;
}

.toolbar button {
  padding: 8px 12px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.toolbar button:hover:not(:disabled) {
  background: #f0f0f0;
}

.toolbar button.is-active {
  background: #007bff;
  color: white;
  border-color: #007bff;
}

.divider {
  width: 1px;
  height: 24px;
  background: #ddd;
  margin: 0 4px;
}

.editor-wrapper {
  position: relative;
}

.editor-content {
  border: 1px solid #e0e0e0;
  border-top: none;
  border-radius: 0 0 8px 8px;
  min-height: 300px;
  padding: 16px;
  background: white;
}

.editor-content :deep(.ProseMirror) {
  outline: none;
  min-height: 300px;
}

/* 协作光标样式优化 */
.editor-content :deep(.collaboration-caret) {
  position: relative;
  margin-left: -1px;
  margin-right: -1px;
  border-left: 2px solid;
  word-break: normal;
  pointer-events: none;
  height: 1.2em;
  display: inline-block;
}

.editor-content :deep(.collaboration-caret::before) {
  content: attr(data-user);
  position: absolute;
  top: -1.5em;
  left: 0;
  font-size: 11px;
  font-weight: 500;
  white-space: nowrap;
  padding: 3px 8px;
  border-radius: 4px 4px 4px 0;
  color: white;
  pointer-events: none;
  opacity: 1; /* 始终显示用户名称 */
  z-index: 1000;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.15);
  transform: translateX(-50%);
  transition: opacity 0.2s;
}

.drag-handle-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  color: #999;
  font-size: 16px;
  cursor: grab;
  user-select: none;
  transition: color 0.2s;
}

.drag-handle-icon:hover {
  color: #007bff;
}

.collaboration-info {
  margin-top: 24px;
  padding: 16px;
  background: #f9f9f9;
  border-radius: 8px;
  border: 1px solid #e0e0e0;
  text-align: left;
}

.collaboration-info h3 {
  margin: 0 0 12px 0;
  font-size: 16px;
  color: #333;
}

.collaboration-info p {
  margin: 8px 0;
  font-size: 14px;
  color: #666;
}

.status-connected {
  color: #28a745;
  font-weight: bold;
}

.status-disconnected {
  color: #dc3545;
  font-weight: bold;
}

.online-users {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #e0e0e0;
}

.online-users h4 {
  margin: 0 0 12px 0;
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.user-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.user-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 8px;
  background: white;
  border-radius: 4px;
  border: 1px solid #e0e0e0;
  transition: all 0.2s;
}

.user-item:hover {
  background: #f5f5f5;
  border-color: #d0d0d0;
}

.user-color-indicator {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  flex-shrink: 0;
  border: 2px solid white;
  box-shadow: 0 0 0 1px rgba(0, 0, 0, 0.1);
}

.user-name {
  font-size: 13px;
  color: #666;
  font-weight: 500;
}
</style>
