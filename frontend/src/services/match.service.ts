import type {
  BusinessError,
  Invitation,
  InvitationActionResult,
  MatchBoard,
} from '../types/domain';
import { ERROR_MESSAGE_FALLBACK } from '../constants/match.constants';

const API_BASE = '/api';

async function parseError(response: Response): Promise<BusinessError> {
  try {
    const body = (await response.json()) as BusinessError;
    return {
      code: body.code ?? 'INTERNAL',
      message: body.message ?? ERROR_MESSAGE_FALLBACK.INVALID_PAYLOAD,
    };
  } catch {
    return { code: 'INTERNAL', message: '网络异常，请稍后重试' };
  }
}

function friendly(error: BusinessError): Error {
  return new Error(ERROR_MESSAGE_FALLBACK[error.code] ?? error.message);
}

// 获取匹配页整页数据：匹配依据、过滤原因、邀请终态均在此回读，刷新页面后重新拉取。
export async function fetchMatchBoard(viewer: string): Promise<MatchBoard> {
  const response = await fetch(`${API_BASE}/match-board?viewer=${encodeURIComponent(viewer)}`);
  if (!response.ok) {
    throw friendly(await parseError(response));
  }
  return response.json() as Promise<MatchBoard>;
}

export interface CreateInvitationPayload {
  viewer: string;
  needId: number;
  proposedSlot?: string;
}

// 发起邀请：后端强制信用分门槛、同校区/连续时隙、三项上限与重复校验。
export async function createInvitation(
  payload: CreateInvitationPayload,
): Promise<Invitation> {
  const response = await fetch(`${API_BASE}/invitations`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  });
  if (!response.ok) {
    throw friendly(await parseError(response));
  }
  const body = (await response.json()) as { invitation: Invitation };
  return body.invitation;
}

// 接受邀请：后端保证重复或并发确认只成功一次。
export async function acceptInvitation(id: number, viewer: string): Promise<InvitationActionResult> {
  const response = await fetch(
    `${API_BASE}/invitations/${id}/accept?viewer=${encodeURIComponent(viewer)}`,
    { method: 'POST' },
  );
  if (!response.ok) {
    throw friendly(await parseError(response));
  }
  return response.json() as Promise<InvitationActionResult>;
}
