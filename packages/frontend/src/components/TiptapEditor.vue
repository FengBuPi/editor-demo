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
      <button
        @click="
          editor?.chain().focus().insertTable({ rows: 3, cols: 3, withHeaderRow: true }).run()
        "
        :disabled="!editor?.can().insertTable()"
        type="button"
        title="插入表格">
        ⧈ 表格
      </button>
      <div class="divider"></div>
      <button
        @click="editor?.chain().focus().addColumnBefore().run()"
        :disabled="!editor?.can().addColumnBefore()"
        type="button"
        title="在左侧添加列">
        ⬅ 列
      </button>
      <button
        @click="editor?.chain().focus().addColumnAfter().run()"
        :disabled="!editor?.can().addColumnAfter()"
        type="button"
        title="在右侧添加列">
        列 ➡
      </button>
      <button
        @click="editor?.chain().focus().deleteColumn().run()"
        :disabled="!editor?.can().deleteColumn()"
        type="button"
        title="删除列">
        ✕ 列
      </button>
      <div class="divider"></div>
      <button
        @click="editor?.chain().focus().addRowBefore().run()"
        :disabled="!editor?.can().addRowBefore()"
        type="button"
        title="在上方添加行">
        ⬆ 行
      </button>
      <button
        @click="editor?.chain().focus().addRowAfter().run()"
        :disabled="!editor?.can().addRowAfter()"
        type="button"
        title="在下方添加行">
        行 ⬇
      </button>
      <button
        @click="editor?.chain().focus().deleteRow().run()"
        :disabled="!editor?.can().deleteRow()"
        type="button"
        title="删除行">
        ✕ 行
      </button>
      <div class="divider"></div>
      <button
        @click="editor?.chain().focus().deleteTable().run()"
        :disabled="!editor?.can().deleteTable()"
        type="button"
        title="删除表格">
        ✕ 表格
      </button>
      <button
        @click="editor?.chain().focus().mergeCells().run()"
        :disabled="!editor?.can().mergeCells()"
        type="button"
        title="合并单元格">
        ⧉ 合并
      </button>
      <button
        @click="editor?.chain().focus().splitCell().run()"
        :disabled="!editor?.can().splitCell()"
        type="button"
        title="拆分单元格">
        ⧈ 拆分
      </button>
      <button
        @click="editor?.chain().focus().toggleHeaderColumn().run()"
        :disabled="!editor?.can().toggleHeaderColumn()"
        :class="{ 'is-active': editor?.isActive('tableHeader') }"
        type="button"
        title="切换列标题">
        ⧉ 列标题
      </button>
      <button
        @click="editor?.chain().focus().toggleHeaderRow().run()"
        :disabled="!editor?.can().toggleHeaderRow()"
        :class="{ 'is-active': editor?.isActive('tableHeader') }"
        type="button"
        title="切换行标题">
        ⧉ 行标题
      </button>
      <div class="divider"></div>
      <button @click="insertImage" type="button" title="通过 URL 插入图片">🖼️ 图片</button>
      <button @click="insertImageFromFile" type="button" title="上传本地图片">📁 上传</button>
      <input
        ref="fileInput"
        type="file"
        accept="image/*"
        style="display: none"
        @change="handleImageUpload" />
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
import { onBeforeUnmount, ref } from "vue";
import { useEditor, EditorContent } from "@tiptap/vue-3";
import StarterKit from "@tiptap/starter-kit";
import Placeholder from "@tiptap/extension-placeholder";
import TextAlign from "@tiptap/extension-text-align";
import DragHandle from "@tiptap/extension-drag-handle-vue-3";
import { Table } from "@tiptap/extension-table";
import { TableRow } from "@tiptap/extension-table-row";
import { TableHeader } from "@tiptap/extension-table-header";
import { TableCell } from "@tiptap/extension-table-cell";
import Image from "@tiptap/extension-image";

const editor = useEditor({
  extensions: [
    StarterKit,
    Placeholder.configure({
      placeholder: "开始输入内容...",
    }),
    TextAlign.configure({
      types: ["heading", "paragraph"],
    }),
    Table.configure({
      resizable: true,
      HTMLAttributes: {
        class: "tiptap-table",
      },
    }),
    TableRow.configure({
      HTMLAttributes: {
        class: "tiptap-table-row",
      },
    }),
    TableHeader.configure({
      HTMLAttributes: {
        class: "tiptap-table-header",
      },
    }),
    TableCell.configure({
      HTMLAttributes: {
        class: "tiptap-table-cell",
      },
    }),
    Image.configure({
      inline: true,
      allowBase64: true,
      HTMLAttributes: {
        class: "tiptap-image",
      },
    }),
  ],
  content: "<p>欢迎使用 Tiptap 编辑器！</p><p>试试选中文字并点击工具栏按钮来格式化文本。</p>",
  editorProps: {
    attributes: {
      class: "prose prose-sm sm:prose lg:prose-lg xl:prose-2xl mx-auto focus:outline-none",
    },
  },
});

const fileInput = ref<HTMLInputElement | null>(null);

// 插入图片（通过 URL）
const insertImage = () => {
  const url = window.prompt("请输入图片 URL:");
  if (url) {
    editor.value?.chain().focus().setImage({ src: url }).run();
  }
};

// 处理图片上传
const handleImageUpload = (event: Event) => {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  if (!file) return;

  // 检查文件类型
  if (!file.type.startsWith("image/")) {
    alert("请选择图片文件");
    return;
  }

  // 使用 createObjectURL 创建临时 URL
  const imageUrl = URL.createObjectURL(file);
  if (imageUrl && editor.value) {
    // 插入图片
    editor.value.chain().focus().setImage({ src: imageUrl }).run();
  } else {
    // 如果 createObjectURL 失败，则使用 FileReader 将图片转换为 base64
    // 使用 FileReader 将图片转换为 base64
    const reader = new FileReader();
    reader.onload = (e) => {
      const base64 = e.target?.result as string;
      if (base64 && editor.value) {
        editor.value.chain().focus().setImage({ src: base64 }).run();
      }
    };
    reader.readAsDataURL(file);
  }

  // 重置 input
  if (fileInput.value) {
    fileInput.value.value = "";
  }
};

// 通过文件选择器插入图片
const insertImageFromFile = () => {
  fileInput.value?.click();
};

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

/* 表格样式 */
.editor-content :deep(.ProseMirror table) {
  border-collapse: collapse;
  margin: 1em 0;
  table-layout: fixed;
  width: 100%;
  overflow: hidden;
  border: 1px solid #ddd;
}

.editor-content :deep(.ProseMirror table td),
.editor-content :deep(.ProseMirror table th) {
  min-width: 1em;
  border: 1px solid #ddd;
  padding: 8px 12px;
  vertical-align: top;
  box-sizing: border-box;
  position: relative;
}

.editor-content :deep(.ProseMirror table th) {
  font-weight: 600;
  text-align: left;
  background-color: #f5f5f5;
  color: #333;
}

.editor-content :deep(.ProseMirror table .selectedCell:after) {
  z-index: 2;
  position: absolute;
  content: "";
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
  background: rgba(200, 200, 255, 0.4);
  pointer-events: none;
}

.editor-content :deep(.ProseMirror table .column-resize-handle) {
  position: absolute;
  right: -2px;
  top: 0;
  bottom: -2px;
  width: 4px;
  background-color: #adf;
  pointer-events: none;
  opacity: 0;
  transition: opacity 0.2s;
}

.editor-content :deep(.ProseMirror table:hover .column-resize-handle) {
  opacity: 1;
}

.editor-content :deep(.ProseMirror table p) {
  margin: 0;
  padding: 0;
}

.editor-content :deep(.ProseMirror table p.is-editor-empty:first-child::before) {
  content: attr(data-placeholder);
  float: left;
  color: #adb5bd;
  pointer-events: none;
  height: 0;
}

/* 表格行样式 */
.editor-content :deep(.ProseMirror table tr) {
  border-bottom: 1px solid #ddd;
}

.editor-content :deep(.ProseMirror table tr:last-child) {
  border-bottom: none;
}

/* 表格单元格悬停效果 */
.editor-content :deep(.ProseMirror table td:hover),
.editor-content :deep(.ProseMirror table th:hover) {
  background-color: #f9f9f9;
}

/* 表格整体样式增强 */
.editor-content :deep(.ProseMirror table) {
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border-radius: 4px;
  overflow: hidden;
}

/* 图片样式 */
.editor-content :deep(.ProseMirror img) {
  max-width: 100%;
  height: auto;
  display: block;
  margin: 0.5em 0;
  border-radius: 4px;
}

/* 表格中的图片样式 */
.editor-content :deep(.ProseMirror table img) {
  max-width: 100%;
  height: auto;
  display: block;
  margin: 0;
  border-radius: 2px;
}

/* 图片选中状态 */
.editor-content :deep(.ProseMirror img.selected) {
  outline: 2px solid #007bff;
  outline-offset: 2px;
}

/* 图片加载失败时的占位符 */
.editor-content :deep(.ProseMirror img[src=""]) {
  display: none;
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
