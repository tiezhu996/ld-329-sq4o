<template>
  <div class="viewer-bar">
    <div class="viewer-info">
      <span class="muted">当前视角</span>
      <el-select
        :model-value="viewer"
        size="default"
        style="width: 220px"
        @update:model-value="emit('switch', String($event))"
      >
        <el-option
          v-for="account in accounts"
          :key="account.name"
          :label="`${account.name}（${account.campus} · 信用 ${account.creditScore}）`"
          :value="account.name"
        >
          <span>{{ account.name }}</span>
          <span class="muted option-sub">
            {{ account.campus }} · {{ account.creditLevel }} · {{ account.creditScore }}
          </span>
        </el-option>
      </el-select>
    </div>
    <el-tag v-if="readOnly" type="danger" effect="dark">仅浏览权限（信用分 &lt; 80）</el-tag>
    <el-tag v-else type="success" effect="plain">可推荐 / 可发起邀请</el-tag>
  </div>
</template>

<script setup lang="ts">
import type { Account } from '../../types/domain';

defineProps<{
  viewer: string;
  accounts: Account[];
  readOnly: boolean;
}>();

const emit = defineEmits<{
  (event: 'switch', name: string): void;
}>();
</script>

<style scoped>
.viewer-bar {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}
.viewer-info {
  display: flex;
  align-items: center;
  gap: 10px;
}
.option-sub {
  float: right;
  font-size: 12px;
}
</style>
