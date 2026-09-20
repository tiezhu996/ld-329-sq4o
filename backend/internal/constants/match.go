package constants

// 匹配与邀请相关的业务规则常量集中维护，禁止在业务代码中写魔法数字。

const (
	// CreditReadOnlyThreshold 信用分低于该值的账号只保留浏览权限，
	// 不能进入推荐结果，也不能发起邀请。
	CreditReadOnlyThreshold = 80
	// MaxActiveInvitations 每个需求最多保留的有效（待确认）邀请数量。
	MaxActiveInvitations = 3
)

// 邀请终态：仅待确认可被接受，接受/失效均为终态，重复或并发确认只成功一次。
const (
	InvitationPending  = "待确认"
	InvitationAccepted = "已接受"
	InvitationExpired  = "已失效"
)

// 可交换时间段使用「星期 + 时段」两级结构，便于判定连续时隙。
var WeekdayOrder = []string{"周一", "周二", "周三", "周四", "周五", "周六", "周日"}

var SlotPartOrder = []string{"上午", "下午", "晚上"}

// FilterReason* 为匹配过滤原因码与中文说明，匹配页需展示过滤原因。
const (
	ReasonLowCredit      = "LOW_CREDIT"
	ReasonCampusMismatch = "CAMPUS_MISMATCH"
	ReasonNoCommonSlot   = "NO_COMMON_SLOT"
)

var FilterReasonText = map[string]string{
	ReasonLowCredit:      "信用分低于 80，仅保留浏览权限",
	ReasonCampusMismatch: "双方校区不一致",
	ReasonNoCommonSlot:   "未来一周没有共同的连续可交换时段",
}

// 匹配依据说明模板，匹配页需展示匹配依据。
const (
	BasisCampus = "同校区"
	BasisSkill  = "技能互补"
	BasisSlot   = "未来一周共有连续可交换时段"
	BasisCredit = "双方信用分均达到 80"
	BasisReward = "回报类型兼容"
)

// 业务错误码。
const (
	CodeReadOnly          = "READ_ONLY_ACCOUNT"
	CodeInvitationLimit   = "INVITATION_LIMIT"
	CodeDuplicateInvite   = "DUPLICATE_INVITATION"
	CodeInvitationClosed  = "INVITATION_CLOSED"
	CodeInvitationMissing = "INVITATION_NOT_FOUND"
	CodeForbidden         = "FORBIDDEN"
)
