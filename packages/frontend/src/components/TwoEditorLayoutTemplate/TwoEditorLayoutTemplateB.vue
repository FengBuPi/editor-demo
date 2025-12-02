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
}
</style>
