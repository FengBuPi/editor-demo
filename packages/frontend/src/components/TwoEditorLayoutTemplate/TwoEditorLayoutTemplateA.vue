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
  border-right: 1px solid #ddd;
}
</style>
