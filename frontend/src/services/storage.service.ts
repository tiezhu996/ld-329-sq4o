import type { Invitation, MatchBoard, Overview } from '../types/domain';

const API_BASE = '/api';

async function parseError(response: Response, fallback: string): Promise<Error> {
  try {
    const body = (await response.json()) as { message?: string };
    if (body.message) {
      return new Error(body.message);
    }
  } catch {
    // 忽略非 JSON 响应，使用兜底文案
  }
  return new Error(fallback);
}

export async function fetchOverview(): Promise<Overview> {
  const response = await fetch(`${API_BASE}/dashboard/overview`);
  if (!response.ok) {
    throw new Error('无法加载校园技能交换数据');
  }
  return response.json() as Promise<Overview>;
}

export async function fetchMatchBoard(): Promise<MatchBoard> {
  const response = await fetch(`${API_BASE}/match/board`);
  if (!response.ok) {
    throw new Error('无法加载匹配数据');
  }
  return response.json() as Promise<MatchBoard>;
}

export async function createInvitation(needId: number, toUser: string): Promise<Invitation> {
  const response = await fetch(`${API_BASE}/invitations`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ needId, toUser }),
  });
  if (!response.ok) {
    throw await parseError(response, '发起邀请失败');
  }
  return response.json() as Promise<Invitation>;
}

export async function acceptInvitation(id: number): Promise<{ invitation: Invitation; changed: boolean }> {
  const response = await fetch(`${API_BASE}/invitations/${id}/accept`, { method: 'POST' });
  if (!response.ok) {
    throw await parseError(response, '确认邀请失败');
  }
  return response.json() as Promise<{ invitation: Invitation; changed: boolean }>;
}
