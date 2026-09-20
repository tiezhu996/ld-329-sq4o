<template>
  <div class="panel match-board" v-loading="loading">
    <h2>智能匹配（校区 / 时隙 / 信用约束）</h2>

    <el-alert v-if="!board?.canInvite" type="warning" show-icon :closable="false"
      :title="`当前账号信用分 ${board?.viewerCredit ?? 0}，低于 80，仅保留浏览权限`" />
    <ul class="rule-list">
      <li v-for="rule in board?.rules ?? []" :key="rule">{{ rule }}</li>
    </ul>

    <el-tabs>
      <el-tab-pane :label="`推荐 ${board?.recommendations.length ?? 0}`">
        <el-empty v-if="!board?.recommendations.length" description="暂无符合约束的推荐" />
        <FeatureCard v-for="match in board?.recommendations ?? []" :key="match.id"
          :title="`${match.provider} × ${match.learner}`" :description="match.recommendation">
          <template #tag>
            <el-tag :type="matchStateType(match.status)">{{ match.status }}</el-tag>
            <el-tag type="warning">{{ match.score }}%</el-tag>
          </template>
          <p class="muted">{{ match.offerSkill }} ↔ {{ match.wantedSkill }}</p>
          <div class="tag-row">
            <el-tag v-for="slot in match.commonSlots" :key="slot" effect="plain">{{ slot }}</el-tag>
          </div>
          <p class="basis-title">匹配依据</p>
          <ul class="basis-list">
            <li v-for="item in match.basis" :key="item">{{ item }}</li>
          </ul>
          <el-button v-if="canInviteMatch(match)" size="small" type="primary"
            :disabled="acting" @click="onInvite(match)">发起邀请</el-button>
        </FeatureCard>
      </el-tab-pane>

      <el-tab-pane :label="`已过滤 ${board?.filtered.length ?? 0}`">
        <el-empty v-if="!board?.filtered.length" description="没有被过滤的候选" />
        <el-table v-else :data="board.filtered" size="small">
          <el-table-column prop="pair" label="候选组合" min-width="120" />
          <el-table-column prop="skill" label="技能" min-width="130" show-overflow-tooltip />
          <el-table-column prop="need" label="需求" min-width="130" show-overflow-tooltip />
          <el-table-column label="过滤原因" min-width="190">
            <template #default="{ row }"><el-tag type="info" effect="plain">{{ row.reason }}</el-tag></template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane :label="`邀请 ${board?.invitations.length ?? 0}`">
        <el-empty v-if="!board?.invitations.length" description="暂无邀请记录" />
        <el-timeline v-else>
          <el-timeline-item v-for="inv in board.invitations" :key="inv.id"
            :type="invStateType(inv.status)" :hollow="inv.status === '已失效'">
            <div class="inv-row">
              <strong>{{ inv.fromUser }} → {{ inv.toUser }}</strong>
              <el-tag size="small" :type="invStateType(inv.status)">{{ inv.status }}</el-tag>
            </div>
            <p class="muted">需求「{{ inv.needTitle }}」 · {{ inv.note }}</p>
            <el-button v-if="canAccept(inv)" size="small" type="success"
              :disabled="acting" @click="onAccept(inv)">确认接受</el-button>
          </el-timeline-item>
        </el-timeline>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { ElMessage } from 'element-plus';
import FeatureCard from './FeatureCard.vue';
import { acceptInvitation, createInvitation, fetchMatchBoard } from '../services/storage.service';
import type { Invitation, Match, MatchBoard } from '../types/domain';

type TagType = 'success' | 'warning' | 'info' | 'primary' | 'danger';

const board = ref<MatchBoard | null>(null);
const loading = ref(true);
const acting = ref(false);

const matchStateMap: Record<string, TagType> = { 可邀请: 'success', 已邀请: 'warning', 已达成: 'info' };
const invStateMap: Record<string, TagType> = { 待确认: 'warning', 已接受: 'success', 已失效: 'info' };
const matchStateType = (status: string): TagType => matchStateMap[status] ?? 'info';
const invStateType = (status: string): TagType => invStateMap[status] ?? 'info';

const canInviteMatch = (match: Match) =>
  !!board.value && board.value.canInvite && match.learner === board.value.viewer && match.status === '可邀请';
const canAccept = (inv: Invitation) => inv.toUser === board.value?.viewer && inv.status === '待确认';

async function reload() {
  board.value = await fetchMatchBoard();
}

async function onInvite(match: Match) {
  acting.value = true;
  try {
    await createInvitation(match.needId, match.provider);
    ElMessage.success('邀请已发出');
    await reload();
  } catch (err) {
    ElMessage.error(err instanceof Error ? err.message : '发起邀请失败');
  } finally {
    acting.value = false;
  }
}

async function onAccept(inv: Invitation) {
  acting.value = true;
  try {
    const result = await acceptInvitation(inv.id);
    ElMessage[result.changed ? 'success' : 'info'](result.changed ? '已确认，同需求其余邀请已失效' : '该邀请已进入终态，重复确认不再生效');
    await reload();
  } catch (err) {
    ElMessage.error(err instanceof Error ? err.message : '确认邀请失败');
  } finally {
    acting.value = false;
  }
}

onMounted(async () => {
  try {
    await reload();
  } catch (err) {
    ElMessage.error(err instanceof Error ? err.message : '加载匹配数据失败');
  } finally {
    loading.value = false;
  }
});
</script>

<style scoped>
.match-board { grid-column: 1 / -1; }
.rule-list { margin: 0 0 12px; padding-left: 18px; color: #667085; font-size: 13px; line-height: 1.8; }
.basis-title { margin: 8px 0 4px; font-size: 13px; font-weight: 600; color: #475467; }
.basis-list { margin: 0 0 8px; padding-left: 18px; color: #667085; font-size: 13px; line-height: 1.7; }
.inv-row { display: flex; align-items: center; gap: 10px; }
</style>
