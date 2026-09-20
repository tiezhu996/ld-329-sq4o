import { defineStore } from 'pinia';
import { computed, ref } from 'vue';
import {
  acceptInvitation as acceptInvitationApi,
  createInvitation as createInvitationApi,
  fetchMatchBoard,
} from '../services/match.service';
import {
  INVITATION_PENDING,
  VIEWER_STORAGE_KEY,
} from '../constants/match.constants';
import type { Invitation, Match, MatchBoard } from '../types/domain';

const DEFAULT_VIEWER = '林澈';

function loadViewer(): string {
  return localStorage.getItem(VIEWER_STORAGE_KEY) ?? DEFAULT_VIEWER;
}

export const useMatchStore = defineStore('match', () => {
  const viewer = ref<string>(loadViewer());
  const board = ref<MatchBoard | null>(null);
  const loading = ref(false);
  const acting = ref(false);
  const error = ref('');
  const notice = ref('');

  const readOnly = computed(() => board.value?.readOnly ?? false);
  const pendingInvitations = computed<Invitation[]>(() =>
    (board.value?.invitations ?? []).filter((item) => item.status === INVITATION_PENDING),
  );
  const terminalInvitations = computed<Invitation[]>(() =>
    (board.value?.invitations ?? []).filter((item) => item.terminal),
  );

  async function load(): Promise<void> {
    loading.value = true;
    error.value = '';
    try {
      board.value = await fetchMatchBoard(viewer.value);
    } catch (err) {
      error.value = err instanceof Error ? err.message : '匹配数据加载失败';
    } finally {
      loading.value = false;
    }
  }

  async function switchViewer(name: string): Promise<void> {
    viewer.value = name;
    localStorage.setItem(VIEWER_STORAGE_KEY, name);
    notice.value = '';
    await load();
  }

  // 操作成功后整页重新回读，确保匹配依据、过滤原因与邀请终态一致。
  async function sendInvitation(match: Match, slot: string): Promise<void> {
    acting.value = true;
    error.value = '';
    notice.value = '';
    try {
      await createInvitationApi({ viewer: viewer.value, needId: match.needId, proposedSlot: slot });
      notice.value = `已向 ${match.learner} 发起邀请，等待对方确认`;
      await load();
    } catch (err) {
      error.value = err instanceof Error ? err.message : '发起邀请失败';
    } finally {
      acting.value = false;
    }
  }

  async function acceptInvitationTarget(invitation: Invitation): Promise<void> {
    acting.value = true;
    error.value = '';
    notice.value = '';
    try {
      const result = await acceptInvitationApi(invitation.id, viewer.value);
      notice.value = result.message;
      await load();
    } catch (err) {
      error.value = err instanceof Error ? err.message : '确认邀请失败';
    } finally {
      acting.value = false;
    }
  }

  return {
    viewer, board, loading, acting, error, notice,
    readOnly, pendingInvitations, terminalInvitations,
    load, switchViewer, sendInvitation, acceptInvitationTarget,
  };
});
