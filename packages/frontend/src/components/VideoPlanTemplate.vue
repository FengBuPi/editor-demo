<!-- 视频规划模板组件 - Tiptap NodeView -->
<template>
  <node-view-wrapper class="video-plan-template">
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
          <div class="video-plan-field-input">
            <a-textarea
              v-model="fieldValues.remark"
              placeholder="请输入备注"
              :auto-size="{ minRows: 2, maxRows: 4 }"
              @change="handleFieldChange('remark', $event)" />
          </div>
        </div>
        <div class="video-plan-field-row">
          <div class="video-plan-field-label">画面描述</div>
          <div class="video-plan-field-input">
            <a-textarea
              v-model="fieldValues.description"
              placeholder="请输入画面描述"
              :auto-size="{ minRows: 2, maxRows: 4 }"
              @change="handleFieldChange('description', $event)" />
          </div>
        </div>
        <div class="video-plan-field-row">
          <div class="video-plan-field-label">口播内容</div>
          <div class="video-plan-field-input">
            <a-textarea
              v-model="fieldValues.voiceover"
              placeholder="请输入口播内容"
              :auto-size="{ minRows: 2, maxRows: 4 }"
              @change="handleFieldChange('voiceover', $event)" />
          </div>
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
        <div class="video-plan-table-cell">
          <a-textarea
            v-model="fieldValues.scene"
            placeholder="请输入场景"
            :auto-size="{ minRows: 2, maxRows: 4 }"
            @change="handleFieldChange('scene', $event)" />
        </div>
        <div class="video-plan-table-cell">
          <a-textarea
            v-model="fieldValues.props"
            placeholder="请输入道具"
            :auto-size="{ minRows: 2, maxRows: 4 }"
            @change="handleFieldChange('props', $event)" />
        </div>
        <div class="video-plan-table-cell">
          <a-textarea
            v-model="fieldValues.actors"
            placeholder="请输入演员"
            :auto-size="{ minRows: 2, maxRows: 4 }"
            @change="handleFieldChange('actors', $event)" />
        </div>
        <div class="video-plan-table-cell">
          <a-textarea
            v-model="fieldValues.costume"
            placeholder="请输入服装"
            :auto-size="{ minRows: 2, maxRows: 4 }"
            @change="handleFieldChange('costume', $event)" />
        </div>
        <div class="video-plan-table-cell">
          <a-textarea
            v-model="fieldValues.materials"
            placeholder="请输入物料"
            :auto-size="{ minRows: 2, maxRows: 4 }"
            @change="handleFieldChange('materials', $event)" />
        </div>
      </div>
    </div>
  </node-view-wrapper>
</template>

<script setup lang="ts">
import { NodeViewWrapper, nodeViewProps } from "@tiptap/vue-3";
import { onMounted, watch, reactive } from "vue";

const props = defineProps(nodeViewProps);

const fieldNames = [
  "remark",
  "description",
  "voiceover",
  "scene",
  "props",
  "actors",
  "costume",
  "materials",
];

const fieldValues = reactive<Record<string, string>>({
  remark: "",
  description: "",
  voiceover: "",
  scene: "",
  props: "",
  actors: "",
  costume: "",
  materials: "",
});

const handleFieldChange = (fieldName: string, value: string) => {
  if (props.updateAttributes) {
    props.updateAttributes({
      [fieldName]: value || "",
    });
  }
};

onMounted(() => {
  fieldNames.forEach((field) => {
    fieldValues[field] = props.node.attrs[field] || "";
  });
});

watch(
  () => props.node.attrs,
  (newAttrs) => {
    fieldNames.forEach((field) => {
      const newValue = newAttrs[field] || "";
      if (fieldValues[field] !== newValue) {
        fieldValues[field] = newValue;
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
  display: flex;
  align-items: stretch;
}

.video-plan-field-input .arco-textarea {
  width: 100%;
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
  display: flex;
  align-items: stretch;
}

.video-plan-table-cell .arco-textarea {
  width: 100%;
}
</style>
