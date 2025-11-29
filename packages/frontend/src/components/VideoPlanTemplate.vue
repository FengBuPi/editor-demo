<!-- 视频规划模板组件 - Tiptap NodeView -->
<template>
  <node-view-wrapper class="video-plan-template" contenteditable="false">
    <!-- 顶部区域：参考画面、中间字段、拍摄素材 -->
    <div class="video-plan-top">
      <!-- 左侧：参考画面 -->
      <div class="video-plan-left">
        <div class="video-plan-add-icon">+</div>
        <div class="video-plan-label">参考画面</div>
      </div>
      <!-- 中间：备注、画面描述、口播内容 -->
      <div class="video-plan-middle">
        <div class="video-plan-field-row">
          <div class="video-plan-field-label">备注</div>
          <div
            ref="remarkRef"
            class="video-plan-field-input"
            contenteditable="true"
            @input="handleInput"
            data-field="remark"></div>
        </div>
        <div class="video-plan-field-row">
          <div class="video-plan-field-label">画面描述</div>
          <div
            ref="descriptionRef"
            class="video-plan-field-input"
            contenteditable="true"
            @input="handleInput"
            data-field="description"></div>
        </div>
        <div class="video-plan-field-row">
          <div class="video-plan-field-label">口播内容</div>
          <div
            ref="voiceoverRef"
            class="video-plan-field-input"
            contenteditable="true"
            @input="handleInput"
            data-field="voiceover"></div>
        </div>
      </div>
      <!-- 右侧：拍摄素材 -->
      <div class="video-plan-right">
        <div class="video-plan-add-icon">+</div>
        <div class="video-plan-label">拍摄素材</div>
      </div>
    </div>
    <!-- 底部区域：场景、道具、演员、服装、物料表格 -->
    <div class="video-plan-bottom">
      <div class="video-plan-table-header">
        <div class="video-plan-table-cell">场景</div>
        <div class="video-plan-table-cell">道具</div>
        <div class="video-plan-table-cell">演员</div>
        <div class="video-plan-table-cell">服装</div>
        <div class="video-plan-table-cell">物料</div>
      </div>
      <div class="video-plan-table-row">
        <div
          ref="sceneRef"
          class="video-plan-table-cell"
          contenteditable="true"
          @input="handleInput"
          data-field="scene"></div>
        <div
          ref="propsRef"
          class="video-plan-table-cell"
          contenteditable="true"
          @input="handleInput"
          data-field="props"></div>
        <div
          ref="actorsRef"
          class="video-plan-table-cell"
          contenteditable="true"
          @input="handleInput"
          data-field="actors"></div>
        <div
          ref="costumeRef"
          class="video-plan-table-cell"
          contenteditable="true"
          @input="handleInput"
          data-field="costume"></div>
        <div
          ref="materialsRef"
          class="video-plan-table-cell"
          contenteditable="true"
          @input="handleInput"
          data-field="materials"></div>
      </div>
    </div>
  </node-view-wrapper>
</template>

<script setup lang="ts">
import { NodeViewWrapper, nodeViewProps } from "@tiptap/vue-3";
import { onMounted, watch, ref, type Ref } from "vue";

// 使用 Tiptap 提供的 nodeViewProps，包含 editor, node, updateAttributes 等
const props = defineProps(nodeViewProps);

// 所有可编辑字段的 ref 引用
const remarkRef = ref<HTMLElement | null>(null);
const descriptionRef = ref<HTMLElement | null>(null);
const voiceoverRef = ref<HTMLElement | null>(null);
const sceneRef = ref<HTMLElement | null>(null);
const propsRef = ref<HTMLElement | null>(null);
const actorsRef = ref<HTMLElement | null>(null);
const costumeRef = ref<HTMLElement | null>(null);
const materialsRef = ref<HTMLElement | null>(null);

// 字段名到 ref 的映射，便于统一处理
const fieldRefs: Record<string, Ref<HTMLElement | null>> = {
  remark: remarkRef,
  description: descriptionRef,
  voiceover: voiceoverRef,
  scene: sceneRef,
  props: propsRef,
  actors: actorsRef,
  costume: costumeRef,
  materials: materialsRef,
};

// 处理输入事件，将内容同步到节点属性
const handleInput = (event: Event) => {
  const target = event.target as HTMLElement;
  const field = target.getAttribute("data-field");
  if (field && props.updateAttributes) {
    props.updateAttributes({
      [field]: target.textContent || "",
    });
  }
};

// 组件挂载时，从节点属性恢复内容
onMounted(() => {
  Object.keys(fieldRefs).forEach((field) => {
    const element = fieldRefs[field].value;
    if (element && props.node.attrs[field]) {
      element.textContent = props.node.attrs[field];
    }
  });
});

// 监听节点属性变化，同步更新到 DOM
watch(
  () => props.node.attrs,
  (newAttrs) => {
    Object.keys(fieldRefs).forEach((field) => {
      const element = fieldRefs[field].value;
      if (element && newAttrs[field] !== undefined) {
        if (element.textContent !== newAttrs[field]) {
          element.textContent = newAttrs[field] || "";
        }
      }
    });
  },
  { deep: true },
);
</script>

<style>
/* 模板容器 */
.video-plan-template {
  border: 1px solid #e0e0e0;
  background: white;
  border-radius: 4px;
  padding: 16px;
  margin: 1em 0;
}

/* 顶部区域：三列网格布局 */
.video-plan-top {
  display: grid;
  grid-template-columns: 1fr 2fr 1fr;
  gap: 1px;
  background: #e0e0e0;
  border: 1px solid #e0e0e0;
  margin-bottom: 16px;
}

/* 左右两侧：参考画面和拍摄素材 */
.video-plan-left,
.video-plan-right {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: white;
  padding: 24px;
  min-height: 120px;
}

/* 加号图标 */
.video-plan-add-icon {
  font-size: 32px;
  color: #999;
  margin-bottom: 8px;
  font-weight: 300;
  line-height: 1;
  user-select: none;
}

/* 标签文字 */
.video-plan-label {
  font-size: 14px;
  color: #333;
  text-align: center;
}

/* 中间区域：备注、画面描述、口播内容 */
.video-plan-middle {
  display: grid;
  grid-template-rows: repeat(3, 1fr);
  gap: 1px;
  background: #e0e0e0;
  border: 1px solid #e0e0e0;
}

/* 字段行：标签 + 输入框 */
.video-plan-field-row {
  display: grid;
  grid-template-columns: 80px 1fr;
  background: white;
}

/* 字段标签 */
.video-plan-field-label {
  font-size: 14px;
  color: #333;
  padding: 8px 12px;
  display: flex;
  justify-content: center;
  align-items: center;
  border-right: 1px solid #e0e0e0;
}

/* 可编辑输入框 */
.video-plan-field-input {
  padding: 8px 12px;
  min-height: 32px;
  outline: none;
}

/* 底部区域：表格 */
.video-plan-bottom {
  border-top: 1px solid #e0e0e0;
  padding-top: 16px;
}

/* 表格表头 */
.video-plan-table-header {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 1px;
  background: #e0e0e0;
  border: 1px solid #e0e0e0;
  margin-bottom: 1px;
}

/* 表格数据行 */
.video-plan-table-row {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 1px;
  background: #e0e0e0;
  border: 1px solid #e0e0e0;
}

/* 表格单元格（可编辑） */
.video-plan-table-cell {
  background: white;
  padding: 12px;
  min-height: 40px;
  border: none;
  outline: none;
}
</style>
