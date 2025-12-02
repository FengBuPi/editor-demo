<template>
  <node-view-wrapper class="editor-b">
    <div class="editor-b">
      <node-view-content class="content is-editable" />
    </div>
  </node-view-wrapper>
</template>

<script setup lang="ts">
import { NodeViewWrapper, nodeViewProps, NodeViewContent } from "@tiptap/vue-3";

const props = defineProps(nodeViewProps);
</script>

<style scoped>
.editor-b {
  flex: 1;
  min-height: 240px;
  border: 1px solid #ddd;
  padding: 8px;
  box-sizing: border-box;
  border-radius: 4px;
}

.editor-b .content {
  min-height: 100%;
}
</style>
