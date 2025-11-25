<script setup lang="ts">
import TiptapEditor from "./components/TiptapEditor.vue";
import TiptapEditorCollaboration from "./components/TiptapEditorCollaboration.vue";
import { ref } from "vue";

const activeTab = ref<"editor" | "collaboration">("collaboration");
</script>

<template>
  <div class="app">
    <h1>Tiptap 编辑器 Demo</h1>

    <div class="tabs">
      <button :class="{ active: activeTab === 'editor' }" @click="activeTab = 'editor'">
        基础编辑器
      </button>
      <button
        :class="{ active: activeTab === 'collaboration' }"
        @click="activeTab = 'collaboration'">
        协作编辑器
      </button>
    </div>

    <div class="content">
      <TiptapEditor v-if="activeTab === 'editor'" />
      <TiptapEditorCollaboration v-if="activeTab === 'collaboration'" />
    </div>
  </div>
</template>

<style scoped>
.app {
  min-height: 100vh;
  background: #fafafa;
  padding: 20px 0;
}

h1 {
  text-align: center;
  color: #333;
  margin-bottom: 30px;
  font-size: 28px;
}

.tabs {
  display: flex;
  justify-content: center;
  gap: 12px;
  margin-bottom: 30px;
}

.tabs button {
  padding: 10px 24px;
  border: 2px solid #007bff;
  background: white;
  color: #007bff;
  border-radius: 6px;
  cursor: pointer;
  font-size: 16px;
  font-weight: 500;
  transition: all 0.2s;
}

.tabs button:hover {
  background: #f0f7ff;
}

.tabs button.active {
  background: #007bff;
  color: white;
}

.content {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 20px;
}
</style>
