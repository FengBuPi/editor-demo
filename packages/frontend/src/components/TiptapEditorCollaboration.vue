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
const statusClass = computed(() =>
  connectionStatus.value === "已连接" ? "status-connected" : "status-disconnected",
);

const WS_URL = "ws://localhost:3001";

const userInfo = {
  name: `用户 ${Math.floor(Math.random() * 1000)}`,
  color: `#${Math.floor(Math.random() * 16777215).toString(16)}`,
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
      }), // 临时类型断言，因为某些版本的 StarterKit 类型定义可能不完整
      Placeholder.configure({
        placeholder: "开始协作编辑...",
      }),
      Collaboration.configure({
        document: ydoc,
      }),
      CollaborationCaret.configure({
        provider: provider,
        user: userInfo,
      }),
    ],
    // 使用协作编辑时，不要设置初始 content，内容从 Y.Doc 中读取
  });

  if (!provider || !ydoc) return;
  const awareness = provider?.awareness;
  if (!awareness) return;
  const updateUserCount = () => {
    onlineUsersCount.value = awareness.getStates().size;
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
</style>
