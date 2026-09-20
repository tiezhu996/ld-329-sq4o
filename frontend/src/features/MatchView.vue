<template>
  <main v-loading="loading" class="page-shell">
    <AppHeader :unread="0" active="match" @navigate="onNavigate" />

    <section class="panel match-header">
      <div>
        <h2>智能匹配与交换邀请</h2>
        <p class="muted">
          未来一周 {{ weekText }} · 信用门槛 {{ board?.creditThreshold ?? 80 }} 分 ·
          每个需求最多 {{ board?.maxActive ?? 3 }} 项有效邀请
        </p>
      </div>
      <ViewerSwitcher
        v-if="board"
        :viewer="viewer"
        :accounts="board.accounts"
        :read-only="board.readOnly"
        @switch="onSwitchViewer"
      />
    </section>

    <el-alert
      v-if="board?.readOnly"
      title="当前账号信用分低于 80，仅保留浏览权限：不会进入推荐结果，也不能发起或确认邀请。"
      type="error"
      show-icon
      :closable="false"
      class="rule-alert"
    />
    <el-alert v-if="notice" :title="notice" type="success" show-icon :closable="true" class="rule-alert" />
    <el-alert v-if="error" :title="error" type="error" show-icon :closable="true" class="rule-alert" />

    <template v-if="board">
      <section class="panel">
        <div class="section-title">
          <h3>推荐匹配（匹配依据 + 连续时隙）</h3>
          <el-tag type="warning">{{ board.matches.length }} 对通过全部约束</el-tag>
        </div>
        <el-empty v-if="board.matches.length === 0" description="当前视角没有可进入推荐的互补匹配" />
        <div v-else class="match-grid">
          <MatchBasisCard
            v-for="match in board.matches"
            :key="match.id"
            :match="match"
            :viewer="viewer"
            :read-only="board.readOnly"
            :acting="acting"
            :active-by-need="board.activeByNeed"
            :max-active="board.maxActive"
            @invite="onInvite"
          />
        </div>
      </section>

      <section class="panel">
        <div class="section-title">
          <h3>被过滤候选与过滤原因</h3>
          <el-tag type="info">{{ board.filtered.length }} 对未通过约束</el-tag>
        </div>
        <FilterReasons :items="board.filtered" />
      </section>

      <section class="panel">
        <div class="section-title">
          <h3>邀请与终态（刷新后可回读）</h3>
          <div class="legend">
            <el-tag type="warning" size="small">待确认 {{ pendingCount }}</el-tag>
            <el-tag type="success" size="small">已接受</el-tag>
            <el-tag type="info" size="small">已失效</el-tag>
          </div>
        </div>
        <el-empty v-if="board.invitations.length === 0" description="暂无邀请" />
        <InvitationPanel
          v-else
          :invitations="orderedInvitations"
          :viewer="viewer"
          :acting="acting"
          @accept="onAccept"
        />
      </section>
    </template>
  </main>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue';
import { storeToRefs } from 'pinia';
import AppHeader from '../components/AppHeader.vue';
import ViewerSwitcher from '../components/match/ViewerSwitcher.vue';
import MatchBasisCard from '../components/match/MatchBasisCard.vue';
import FilterReasons from '../components/match/FilterReasons.vue';
import InvitationPanel from '../components/match/InvitationPanel.vue';
import { useMatchStore } from '../stores/match.store';
import { useAppStore } from '../stores/app.store';
import type { Invitation, Match } from '../types/domain';

const store = useMatchStore();
const app = useAppStore();
const { viewer, board, loading, acting, error, notice, pendingInvitations } = storeToRefs(store);

const weekText = computed(() =>
  board.value ? `${board.value.weekStart} 至 ${board.value.weekEnd}` : '',
);
const pendingCount = computed(() => pendingInvitations.value.length);

function onNavigate(view: 'dashboard' | 'match'): void {
  app.setView(view);
}

// 待确认在前、终态在后，同组内按创建时间倒序，便于回读最新状态。
const orderedInvitations = computed<Invitation[]>(() => {
  const list = [...(board.value?.invitations ?? [])];
  return list.sort((a, b) => {
    if (a.terminal !== b.terminal) return a.terminal ? 1 : -1;
    return b.createdAtUnix - a.createdAtUnix;
  });
});

function onSwitchViewer(name: string): void {
  void store.switchViewer(name);
}
function onInvite(match: Match, slot: string): void {
  void store.sendInvitation(match, slot);
}
function onAccept(invitation: Invitation): void {
  void store.acceptInvitationTarget(invitation);
}

onMounted(() => {
  void store.load();
});
</script>

<style scoped>
.match-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
  margin-top: 18px;
}
.rule-alert {
  margin: 14px 0;
}
.section-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
}
.section-title h3 {
  margin: 0;
}
.legend {
  display: flex;
  gap: 8px;
}
.match-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px;
}
@media (max-width: 900px) {
  .match-grid {
    grid-template-columns: 1fr;
  }
}
</style>
