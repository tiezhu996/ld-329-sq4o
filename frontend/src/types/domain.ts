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
  offerSkill: string;
  wantedSkill: string;
  score: number;
  commonSlots: string[];
  recommendation: string;
  basis: string[];
  status: string;
}

export interface FilteredCandidate {
  pair: string;
  skill: string;
  need: string;
  reason: string;
}

export interface Invitation {
  id: number;
  needId: number;
  needTitle: string;
  fromUser: string;
  toUser: string;
  status: string;
  note: string;
}

export interface MatchBoard {
  viewer: string;
  viewerCredit: number;
  canInvite: boolean;
  rules: string[];
  recommendations: Match[];
  filtered: FilteredCandidate[];
  invitations: Invitation[];
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
