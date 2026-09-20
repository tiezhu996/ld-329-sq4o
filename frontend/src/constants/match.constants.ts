// 匹配与邀请相关的前端常量，与后端 constants/match.go 保持一致。

export const CREDIT_READ_ONLY_THRESHOLD = 80;
export const MAX_ACTIVE_INVITATIONS = 3;

export const INVITATION_PENDING = '待确认';
export const INVITATION_ACCEPTED = '已接受';
export const INVITATION_EXPIRED = '已失效';

export type InvitationStatusType =
  | typeof INVITATION_PENDING
  | typeof INVITATION_ACCEPTED
  | typeof INVITATION_EXPIRED;

export const INVITATION_STATUS_TYPE: Record<
  InvitationStatusType,
  { type: 'warning' | 'success' | 'info'; label: string }
> = {
  [INVITATION_PENDING]: { type: 'warning', label: '待确认' },
  [INVITATION_ACCEPTED]: { type: 'success', label: '已接受（终态）' },
  [INVITATION_EXPIRED]: { type: 'info', label: '已失效（终态）' },
};

export const FILTER_REASON_TAG_TYPE: Record<string, 'danger' | 'warning' | 'info'> = {
  LOW_CREDIT: 'danger',
  CAMPUS_MISMATCH: 'warning',
  NO_COMMON_SLOT: 'info',
};

export const FILTER_REASON_LABEL: Record<string, string> = {
  LOW_CREDIT: '信用分不足',
  CAMPUS_MISMATCH: '校区不一致',
  NO_COMMON_SLOT: '无连续共同时段',
};

// 后端业务错误码 -> 面向用户的提示。
export const ERROR_MESSAGE_FALLBACK: Record<string, string> = {
  READ_ONLY_ACCOUNT: '信用分低于 80，该账号仅保留浏览权限，不能进入推荐或发起邀请。',
  INVITATION_LIMIT: '该需求最多保留 3 项有效邀请，已有邀请被接受或失效后再试。',
  DUPLICATE_INVITATION: '你已对该需求发起过有效邀请，请勿重复提交。',
  INVITATION_CLOSED: '该邀请已是终态，无法再次操作。',
  INVITATION_NOT_FOUND: '邀请不存在或已被清理，请刷新后重试。',
  FORBIDDEN: '该操作不被允许：非同校区、无连续共同时段，或你无权确认该邀请。',
  INVALID_PAYLOAD: '请求参数不完整，请检查后重试。',
};

export const VIEWER_STORAGE_KEY = 'cyskillswap.viewer';
