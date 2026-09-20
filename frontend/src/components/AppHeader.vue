<template>
  <header class="app-header">
    <div>
      <p class="eyebrow">Skill Swap Campus</p>
      <h1>{{ APP_NAME }}</h1>
      <p>{{ APP_TAGLINE }}</p>
    </div>
    <div class="header-actions">
      <el-radio-group :model-value="active" size="default" @update:model-value="onTab">
        <el-radio-button value="dashboard">总览</el-radio-button>
        <el-radio-button value="match">智能匹配</el-radio-button>
      </el-radio-group>
      <el-badge :value="unread" class="message-badge">
        <el-button type="primary">站内消息</el-button>
      </el-badge>
    </div>
  </header>
</template>

<script setup lang="ts">
import { APP_NAME, APP_TAGLINE } from '../constants/app.constants';

defineProps<{ unread: number; active: 'dashboard' | 'match' }>();

const emit = defineEmits<{
  (event: 'navigate', view: 'dashboard' | 'match'): void;
}>();

function onTab(value: string | number | boolean): void {
  emit('navigate', value === 'match' ? 'match' : 'dashboard');
}
</script>

<style scoped>
.header-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}
</style>
