package constants

const (
	ServiceName = "cyskillswap"
	APIPrefix   = "/api"
)

var SkillCategories = []string{"摄影", "编程", "乐器", "外语", "平面设计", "健身指导"}

const (
	CreditBronze = "青铜互助者"
	CreditSilver = "白银协作者"
	CreditGold   = "黄金导师"
)

// 匹配与邀请规则
const (
	// CreditBrowseThreshold 信用分低于该值的账号仅保留浏览权限
	CreditBrowseThreshold = 80
	// MaxActiveInvitationsPerNeed 每个需求最多保留的有效（待确认）邀请数
	MaxActiveInvitationsPerNeed = 3
	// SlotWindowDays 可交换时段的匹配窗口：未来一周
	SlotWindowDays = 7
)

// 邀请终态
const (
	InvitationPending  = "待确认"
	InvitationAccepted = "已接受"
	InvitationExpired  = "已失效"
)

// 匹配推荐终态
const (
	MatchStateOpen    = "可邀请"
	MatchStateInvited = "已邀请"
	MatchStateClosed  = "已达成"
)

// 过滤原因
const (
	FilterReasonCredit = "信用分低于80，仅保留浏览权限"
	FilterReasonCampus = "校区不一致"
	FilterReasonSlot   = "未来一周无共同连续可交换时段"
)
