<script setup lang="ts">
import TiptapEditor from './components/TiptapEditor.vue'
</script>

<template>
  <div class="app">
    <h1>Tiptap 编辑器 Demo</h1>
    <TiptapEditor />
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
</style>
