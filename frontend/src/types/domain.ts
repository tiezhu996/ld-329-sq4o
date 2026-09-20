export interface Skill {
  id: number;
  owner: string;
  title: string;
  category: string;
  level: number;
  campus: string;
  description: string;
  timeSlots: string[];
  rewards: string[];
  portfolio: string;
}

export interface Need {
  id: number;
  requester: string;
  title: string;
  category: string;
  campus: string;
  expectTime: string;
  budgetType: string;
  description: string;
  responses: number;
}

export interface Match {
  id: number;
  needId: number;
  provider: string;
  learner: string;
  campus: string;
  offerSkill: string;
  wantedSkill: string;
  category: string;
  score: number;
  commonSlots: string[];
  commonBlocks: SlotBlock[];
  reward: string;
  basis: string[];
  recommendation: string;
  viewerInvolved: boolean;
  viewerCanInvite: boolean;
}

export interface SlotBlock {
  weekday: string;
  parts: string[];
  slotText: string;
}

export interface Account {
  name: string;
  major: string;
  campus: string;
  creditScore: number;
  creditLevel: string;
}

export type FilterReasonCode = 'LOW_CREDIT' | 'CAMPUS_MISMATCH' | 'NO_COMMON_SLOT';

export interface FilteredMatch {
  provider: string;
  learner: string;
  campus: string;
  offerSkill: string;
  wantedSkill: string;
  reasonCode: FilterReasonCode;
  reason: string;
  commonSlots: string[];
}

export type InvitationStatus = '待确认' | '已接受' | '已失效';

export interface Invitation {
  id: number;
  matchId: number;
  needId: number;
  fromUser: string;
  toUser: string;
  offerSkill: string;
  wantedSkill: string;
  campus: string;
  proposedSlot: string;
  status: InvitationStatus;
  terminal: boolean;
  createdAtUnix: number;
  decidedAtUnix?: number;
  note: string;
}

export interface MatchBoard {
  viewer: string;
  viewerCampus: string;
  viewerCredit: number;
  readOnly: boolean;
  weekStart: string;
  weekEnd: string;
  creditThreshold: number;
  maxActive: number;
  accounts: Account[];
  matches: Match[];
  filtered: FilteredMatch[];
  invitations: Invitation[];
  activeCount: number;
  activeByNeed: Record<string, number>;
}

export interface InvitationActionResult {
  invitation: Invitation;
  accepted: boolean;
  message: string;
}

export interface BusinessError {
  code: string;
  message: string;
}

export interface Appointment {
  id: number;
  pair: string;
  time: string;
  place: string;
  status: string;
  agenda: string;
}

export interface Review {
  id: number;
  from: string;
  to: string;
  rating: number;
  content: string;
}

export interface Conversation {
  id: number;
  withUser: string;
  unread: number;
  messages: string[];
}

export interface Profile {
  name: string;
  major: string;
  creditScore: number;
  creditLevel: string;
  skillWall: Skill[];
  radar: Record<string, number>;
  history: string[];
  reviews: Review[];
}

export interface Overview {
  service: string;
  categories: string[];
  metrics: Record<string, number>;
  skills: Skill[];
  needs: Need[];
  matches: Match[];
  appointments: Appointment[];
  reviews: Review[];
  messages: Conversation[];
  profile: Profile;
}
