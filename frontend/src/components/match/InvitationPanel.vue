<template>
  <el-table :data="invitations" size="small" :row-key="(row: Invitation) => row.id">
    <el-table-column label="邀请" min-width="220">
      <template #default="{ row }">
        <strong>{{ row.fromUser }}</strong>
        <span class="muted"> → {{ row.toUser }}</span>
        <p class="muted invitation-skill">{{ row.offerSkill }} ↔ {{ row.wantedSkill }}</p>
      </template>
    </el-table-column>
    <el-table-column prop="campus" label="校区" width="96" />
    <el-table-column prop="proposedSlot" label="提议时段" width="120" />
    <el-table-column label="状态 / 终态" width="150">
      <template #default="{ row }">
        <el-tag :type="statusMeta(row.status).type" size="small">
          {{ statusMeta(row.status).label }}
        </el-tag>
        <p v-if="row.terminal" class="terminal-text">终态，不可再变更</p>
      </template>
    </el-table-column>
    <el-table-column label="备注" min-width="160">
      <template #default="{ row }">
        <span class="muted">{{ row.note || '—' }}</span>
      </template>
    </el-table-column>
    <el-table-column label="操作" width="120" fixed="right">
      <template #default="{ row }">
        <el-button
          v-if="canAccept(row)"
          type="primary"
          size="small"
          :loading="acting"
          @click="emit('accept', row)"
        >
          接受
        </el-button>
        <span v-else-if="row.toUser !== viewer" class="muted">仅需求方可操作</span>
        <span v-else class="muted">已结束</span>
      </template>
    </el-table-column>
  </el-table>
</template>

<script setup lang="ts">
import {
  INVITATION_PENDING,
  INVITATION_STATUS_TYPE,
  type InvitationStatusType,
} from '../../constants/match.constants';
import type { Invitation } from '../../types/domain';

const props = defineProps<{
  invitations: Invitation[];
  viewer: string;
  acting: boolean;
}>();

const emit = defineEmits<{
  (event: 'accept', invitation: Invitation): void;
}>();

function statusMeta(status: string) {
  return (
    INVITATION_STATUS_TYPE[status as InvitationStatusType] ?? {
      type: 'info' as const,
      label: status,
    }
  );
}

function canAccept(row: Invitation): boolean {
  return row.status === INVITATION_PENDING && row.toUser === props.viewer;
}
</script>

<style scoped>
.invitation-skill {
  margin: 2px 0 0;
  font-size: 12px;
}
.terminal-text {
  margin: 4px 0 0;
  font-size: 12px;
  color: #98a2b3;
}
</style>
