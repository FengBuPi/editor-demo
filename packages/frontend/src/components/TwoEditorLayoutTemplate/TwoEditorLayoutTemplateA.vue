<template>
  <node-view-wrapper class="editor-a">
    <div class="editor-a">
      <node-view-content class="content is-editable" />
    </div>
  </node-view-wrapper>
</template>

<script setup lang="ts">
import { NodeViewWrapper, nodeViewProps, NodeViewContent } from "@tiptap/vue-3";

const props = defineProps(nodeViewProps);
</script>

<style scoped>
.editor-a {
  flex: 1;
  min-height: 240px;
  border: 1px solid #ddd;
  border-right: 1px solid #ddd;
  padding: 8px;
  box-sizing: border-box;
  border-radius: 4px;
}

.editor-a .content {
  min-height: 100%;
}
</style>
