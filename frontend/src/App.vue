<template>
  <DashboardView v-if="view === 'dashboard'" />
  <MatchView v-else />
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia';
import DashboardView from './features/DashboardView.vue';
import MatchView from './features/MatchView.vue';
import { useAppStore } from './stores/app.store';

const app = useAppStore();
const { view } = storeToRefs(app);
</script>
