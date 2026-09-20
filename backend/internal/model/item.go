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
	ID             int      `json:"id"`
	Requester      string   `json:"requester"`
	Title          string   `json:"title"`
	Category       string   `json:"category"`
	Campus         string   `json:"campus"`
	ExpectTime     string   `json:"expectTime"`
	AvailableSlots []string `json:"availableSlots"`
	BudgetType     string   `json:"budgetType"`
	Description    string   `json:"description"`
	Responses      int      `json:"responses"`
}

// Account 平台账号，携带信用分与权限。信用分低于门槛时仅可浏览。
type Account struct {
	Name        string `json:"name"`
	Major       string `json:"major"`
	Campus      string `json:"campus"`
	CreditScore int    `json:"creditScore"`
	CreditLevel string `json:"creditLevel"`
}

func (a Account) CanInteract(threshold int) bool { return a.CreditScore >= threshold }

// SlotBlock 表示一段连续可交换时隙（同一星期内相邻时段合并）。
type SlotBlock struct {
	Weekday  string   `json:"weekday"`
	Parts    []string `json:"parts"`
	SlotText string   `json:"slotText"`
}

// Match 一条通过全部约束的互补匹配推荐。
type Match struct {
	ID             int         `json:"id"`
	NeedID         int         `json:"needId"`
	Provider       string      `json:"provider"`
	Learner        string      `json:"learner"`
	Campus         string      `json:"campus"`
	OfferSkill     string      `json:"offerSkill"`
	WantedSkill    string      `json:"wantedSkill"`
	Category       string      `json:"category"`
	Score          int         `json:"score"`
	CommonSlots    []string    `json:"commonSlots"`
	CommonBlocks   []SlotBlock `json:"commonBlocks"`
	Reward         string      `json:"reward"`
	Basis          []string    `json:"basis"`
	Recommendation string      `json:"recommendation"`
	// ViewerInvolved 标识当前查看者是否为匹配当事人，决定是否可发起邀请。
	ViewerInvolved bool `json:"viewerInvolved"`
	// ViewerCanInvite 查看者是提供方且该需求未满三项、本人未重复邀请时为真。
	ViewerCanInvite bool `json:"viewerCanInvite"`
}

// FilteredMatch 未通过约束的候选对，需在匹配页展示过滤原因。
type FilteredMatch struct {
	Provider    string   `json:"provider"`
	Learner     string   `json:"learner"`
	Campus      string   `json:"campus"`
	OfferSkill  string   `json:"offerSkill"`
	WantedSkill string   `json:"wantedSkill"`
	ReasonCode  string   `json:"reasonCode"`
	Reason      string   `json:"reason"`
	CommonSlots []string `json:"commonSlots"`
}

// Invitation 由技能提供方向需求发起的交换邀请。
type Invitation struct {
	ID            int    `json:"id"`
	MatchID       int    `json:"matchId"`
	NeedID        int    `json:"needId"`
	FromUser      string `json:"fromUser"`
	ToUser        string `json:"toUser"`
	OfferSkill    string `json:"offerSkill"`
	WantedSkill   string `json:"wantedSkill"`
	Campus        string `json:"campus"`
	ProposedSlot  string `json:"proposedSlot"`
	Status        string `json:"status"`
	Terminal      bool   `json:"terminal"`
	CreatedAtUnix int64  `json:"createdAtUnix"`
	DecidedAtUnix int64  `json:"decidedAtUnix,omitempty"`
	Note          string `json:"note"`
}

// MatchBoard 匹配页整页数据，刷新后可整体回读。
type MatchBoard struct {
	Viewer       string          `json:"viewer"`
	ViewerCampus string          `json:"viewerCampus"`
	ViewerCredit int             `json:"viewerCredit"`
	ReadOnly     bool            `json:"readOnly"`
	WeekStart    string          `json:"weekStart"`
	WeekEnd      string          `json:"weekEnd"`
	Threshold    int             `json:"creditThreshold"`
	MaxActive    int             `json:"maxActive"`
	Accounts     []Account       `json:"accounts"`
	Matches      []Match         `json:"matches"`
	Filtered     []FilteredMatch `json:"filtered"`
	Invitations  []Invitation    `json:"invitations"`
	ActiveCount  int             `json:"activeCount"`
	ActiveByNeed map[int]int     `json:"activeByNeed"`
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
