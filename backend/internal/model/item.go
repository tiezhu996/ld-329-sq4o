package model

type Skill struct {
	ID          int      `json:"id"`
	Owner       string   `json:"owner"`
	Title       string   `json:"title"`
	Category    string   `json:"category"`
	Level       int      `json:"level"`
	Campus      string   `json:"campus"`
	Description string   `json:"description"`
	TimeSlots   []string `json:"timeSlots"`
	Rewards     []string `json:"rewards"`
	Portfolio   string   `json:"portfolio"`
}

type Need struct {
	ID          int    `json:"id"`
	Requester   string `json:"requester"`
	Title       string `json:"title"`
	Category    string `json:"category"`
	Campus      string `json:"campus"`
	ExpectTime  string `json:"expectTime"`
	BudgetType  string `json:"budgetType"`
	Description string `json:"description"`
	Responses   int    `json:"responses"`
}

type Match struct {
	ID             int      `json:"id"`
	NeedID         int      `json:"needId"`
	Provider       string   `json:"provider"`
	Learner        string   `json:"learner"`
	OfferSkill     string   `json:"offerSkill"`
	WantedSkill    string   `json:"wantedSkill"`
	Score          int      `json:"score"`
	CommonSlots    []string `json:"commonSlots"`
	Recommendation string   `json:"recommendation"`
	Basis          []string `json:"basis"`
	Status         string   `json:"status"`
}

// User 平台账号，信用分决定能否进入推荐或发起邀请
type User struct {
	Name        string `json:"name"`
	Campus      string `json:"campus"`
	CreditScore int    `json:"creditScore"`
	CreditLevel string `json:"creditLevel"`
}

// FilteredCandidate 被规则过滤的候选及原因
type FilteredCandidate struct {
	Pair   string `json:"pair"`
	Skill  string `json:"skill"`
	Need   string `json:"need"`
	Reason string `json:"reason"`
}

// Invitation 需求下的交换邀请，状态为终态机：待确认→已接受/已失效
type Invitation struct {
	ID        int    `json:"id"`
	NeedID    int    `json:"needId"`
	NeedTitle string `json:"needTitle"`
	FromUser  string `json:"fromUser"`
	ToUser    string `json:"toUser"`
	Status    string `json:"status"`
	Note      string `json:"note"`
}

// MatchBoard 匹配页数据：推荐、过滤原因、邀请终态与规则说明
type MatchBoard struct {
	Viewer          string             `json:"viewer"`
	ViewerCredit    int                `json:"viewerCredit"`
	CanInvite       bool               `json:"canInvite"`
	Rules           []string           `json:"rules"`
	Recommendations []Match            `json:"recommendations"`
	Filtered        []FilteredCandidate `json:"filtered"`
	Invitations     []Invitation       `json:"invitations"`
}

type Appointment struct {
	ID     int    `json:"id"`
	Pair   string `json:"pair"`
	Time   string `json:"time"`
	Place  string `json:"place"`
	Status string `json:"status"`
	Agenda string `json:"agenda"`
}

type Review struct {
	ID      int    `json:"id"`
	From    string `json:"from"`
	To      string `json:"to"`
	Rating  int    `json:"rating"`
	Content string `json:"content"`
}

type Conversation struct {
	ID       int      `json:"id"`
	WithUser string   `json:"withUser"`
	Unread   int      `json:"unread"`
	Messages []string `json:"messages"`
}

type Profile struct {
	Name        string         `json:"name"`
	Major       string         `json:"major"`
	CreditScore int            `json:"creditScore"`
	CreditLevel string         `json:"creditLevel"`
	SkillWall   []Skill        `json:"skillWall"`
	Radar       map[string]int `json:"radar"`
	History     []string       `json:"history"`
	Reviews     []Review       `json:"reviews"`
}

type Overview struct {
	Service      string         `json:"service"`
	Categories   []string       `json:"categories"`
	Metrics      map[string]int `json:"metrics"`
	Skills       []Skill        `json:"skills"`
	Needs        []Need         `json:"needs"`
	Matches      []Match        `json:"matches"`
	Appointments []Appointment  `json:"appointments"`
	Reviews      []Review       `json:"reviews"`
	Messages     []Conversation `json:"messages"`
	Profile      Profile        `json:"profile"`
}
