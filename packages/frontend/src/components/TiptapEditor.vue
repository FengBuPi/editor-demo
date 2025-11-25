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
      <button
        @click="editor?.chain().focus().toggleStrike().run()"
        :class="{ 'is-active': editor?.isActive('strike') }"
        type="button">
        <s>S</s>
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
      <button
        @click="editor?.chain().focus().toggleHeading({ level: 3 }).run()"
        :class="{ 'is-active': editor?.isActive('heading', { level: 3 }) }"
        type="button">
        H3
      </button>
      <div class="divider"></div>
      <button
        @click="editor?.chain().focus().toggleBulletList().run()"
        :class="{ 'is-active': editor?.isActive('bulletList') }"
        type="button">
        • 列表
      </button>
      <button
        @click="editor?.chain().focus().toggleOrderedList().run()"
        :class="{ 'is-active': editor?.isActive('orderedList') }"
        type="button">
        1. 列表
      </button>
      <button
        @click="editor?.chain().focus().toggleBlockquote().run()"
        :class="{ 'is-active': editor?.isActive('blockquote') }"
        type="button">
        " 引用
      </button>
      <div class="divider"></div>
      <button
        @click="editor?.chain().focus().toggleCodeBlock().run()"
        :class="{ 'is-active': editor?.isActive('codeBlock') }"
        type="button">
        { }
      </button>
      <div class="divider"></div>
      <button
        @click="editor?.chain().focus().setTextAlign('left').run()"
        :class="{ 'is-active': editor?.isActive({ textAlign: 'left' }) }"
        type="button"
        title="左对齐">
        ⬅
      </button>
      <button
        @click="editor?.chain().focus().setTextAlign('center').run()"
        :class="{ 'is-active': editor?.isActive({ textAlign: 'center' }) }"
        type="button"
        title="居中">
        ⬌
      </button>
      <button
        @click="editor?.chain().focus().setTextAlign('right').run()"
        :class="{ 'is-active': editor?.isActive({ textAlign: 'right' }) }"
        type="button"
        title="右对齐">
        ➡
      </button>
      <div class="divider"></div>
      <button @click="editor?.chain().focus().setHorizontalRule().run()" type="button">─</button>
      <button
        @click="editor?.chain().focus().undo().run()"
        :disabled="!editor?.can().undo()"
        type="button">
        ↶
      </button>
      <button
        @click="editor?.chain().focus().redo().run()"
        :disabled="!editor?.can().redo()"
        type="button">
        ↷
      </button>
    </div>
    <div class="editor-wrapper">
      <drag-handle v-if="editor" :editor="editor">
        <div class="drag-handle-icon">⋮⋮</div>
      </drag-handle>
      <editor-content :editor="editor" class="editor-content" />
    </div>
    <div class="output">
      <h3>输出内容（JSON）：</h3>
      <pre>{{ JSON.stringify(editor?.getJSON(), null, 2) }}</pre>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount } from "vue";
import { useEditor, EditorContent } from "@tiptap/vue-3";
import StarterKit from "@tiptap/starter-kit";
import Placeholder from "@tiptap/extension-placeholder";
import TextAlign from "@tiptap/extension-text-align";
import DragHandle from "@tiptap/extension-drag-handle-vue-3";

const editor = useEditor({
  extensions: [
    StarterKit,
    Placeholder.configure({
      placeholder: "开始输入内容...",
    }),
    TextAlign.configure({
      types: ["heading", "paragraph"],
    }),
  ],
  content: "<p>欢迎使用 Tiptap 编辑器！</p><p>试试选中文字并点击工具栏按钮来格式化文本。</p>",
  editorProps: {
    attributes: {
      class: "prose prose-sm sm:prose lg:prose-lg xl:prose-2xl mx-auto focus:outline-none",
    },
  },
});

onBeforeUnmount(() => {
  editor.value?.destroy();
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

.toolbar button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
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

.drag-handle-icon:active {
  cursor: grabbing;
}

.editor-content :deep(.ProseMirror) {
  outline: none;
  min-height: 300px;
}

.editor-content :deep(.ProseMirror p.is-editor-empty:first-child::before) {
  content: attr(data-placeholder);
  float: left;
  color: #adb5bd;
  pointer-events: none;
  height: 0;
}

.editor-content :deep(.ProseMirror h1) {
  font-size: 2em;
  font-weight: bold;
  margin: 0.67em 0;
}

.editor-content :deep(.ProseMirror h2) {
  font-size: 1.5em;
  font-weight: bold;
  margin: 0.75em 0;
}

.editor-content :deep(.ProseMirror h3) {
  font-size: 1.17em;
  font-weight: bold;
  margin: 0.83em 0;
}

.editor-content :deep(.ProseMirror ul),
.editor-content :deep(.ProseMirror ol) {
  padding-left: 1.5em;
  margin: 1em 0;
}

.editor-content :deep(.ProseMirror blockquote) {
  border-left: 4px solid #ddd;
  padding-left: 1em;
  margin: 1em 0;
  color: #666;
}

.editor-content :deep(.ProseMirror code) {
  background: #f4f4f4;
  padding: 2px 6px;
  border-radius: 3px;
  font-family: "Courier New", monospace;
  font-size: 0.9em;
}

.editor-content :deep(.ProseMirror pre) {
  background: #f4f4f4;
  padding: 12px;
  border-radius: 4px;
  overflow-x: auto;
  margin: 1em 0;
}

.editor-content :deep(.ProseMirror pre code) {
  background: none;
  padding: 0;
}

.editor-content :deep(.ProseMirror hr) {
  border: none;
  border-top: 2px solid #ddd;
  margin: 2em 0;
}

.editor-content :deep(.ProseMirror p[style*="text-align: left"]),
.editor-content :deep(.ProseMirror h1[style*="text-align: left"]),
.editor-content :deep(.ProseMirror h2[style*="text-align: left"]),
.editor-content :deep(.ProseMirror h3[style*="text-align: left"]) {
  text-align: left;
}

.editor-content :deep(.ProseMirror p[style*="text-align: center"]),
.editor-content :deep(.ProseMirror h1[style*="text-align: center"]),
.editor-content :deep(.ProseMirror h2[style*="text-align: center"]),
.editor-content :deep(.ProseMirror h3[style*="text-align: center"]) {
  text-align: center;
}

.editor-content :deep(.ProseMirror p[style*="text-align: right"]),
.editor-content :deep(.ProseMirror h1[style*="text-align: right"]),
.editor-content :deep(.ProseMirror h2[style*="text-align: right"]),
.editor-content :deep(.ProseMirror h3[style*="text-align: right"]) {
  text-align: right;
}

.output {
  margin-top: 24px;
  padding: 16px;
  background: #f9f9f9;
  border-radius: 8px;
  border: 1px solid #e0e0e0;
  text-align: left;
}

.output h3 {
  margin: 0 0 12px 0;
  font-size: 16px;
  color: #333;
  text-align: left;
}

.output pre {
  margin: 0;
  padding: 12px;
  background: white;
  border-radius: 4px;
  overflow-x: auto;
  font-size: 12px;
  line-height: 1.5;
  max-height: 300px;
  overflow-y: auto;
  text-align: left;
}
</style>
