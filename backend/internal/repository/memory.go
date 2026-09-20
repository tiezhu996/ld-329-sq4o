package repository

import (
	"sync"

	"cyskillswap/internal/constants"
	"cyskillswap/internal/model"
)

var users = []model.User{
	{Name: "林澈", Campus: "东校区", CreditScore: 91, CreditLevel: constants.CreditGold},
	{Name: "孟野", Campus: "西校区", CreditScore: 88, CreditLevel: constants.CreditSilver},
	{Name: "周芮", Campus: "中心校区", CreditScore: 93, CreditLevel: constants.CreditGold},
	{Name: "许安", Campus: "中心校区", CreditScore: 72, CreditLevel: constants.CreditBronze},
	{Name: "苏漫", Campus: "东校区", CreditScore: 85, CreditLevel: constants.CreditSilver},
	{Name: "何屿", Campus: "西校区", CreditScore: 84, CreditLevel: constants.CreditSilver},
	{Name: "秦朗", Campus: "西校区", CreditScore: 86, CreditLevel: constants.CreditSilver},
	{Name: "郑北", Campus: "东校区", CreditScore: 90, CreditLevel: constants.CreditGold},
	{Name: "陆离", Campus: "东校区", CreditScore: 83, CreditLevel: constants.CreditSilver},
}

func ListUsers() []model.User {
	out := make([]model.User, len(users))
	copy(out, users)
	return out
}

func GetUser(name string) (model.User, bool) {
	for _, u := range users {
		if u.Name == name {
			return u, true
		}
	}
	return model.User{}, false
}

func ListSkills() []model.Skill {
	return []model.Skill{
		{ID: 1, Owner: "林澈", Title: "毕业照人像摄影", Category: "摄影", Level: 92, Campus: "东校区", Description: "提供构图、修图和毕业季跟拍，可交换吉他入门课。", TimeSlots: []string{"周三晚", "周六上午"}, Rewards: []string{"技能交换", "请吃饭"}, Portfolio: "12组校园人像作品"},
		{ID: 2, Owner: "周芮", Title: "Python 数据分析", Category: "编程", Level: 88, Campus: "中心校区", Description: "pandas、可视化、论文数据清洗辅导，接受小额报酬。", TimeSlots: []string{"周二晚", "周日全天"}, Rewards: []string{"技能交换", "小额报酬"}, Portfolio: "3份课程项目证书"},
		{ID: 3, Owner: "孟野", Title: "民谣吉他陪练", Category: "乐器", Level: 81, Campus: "西校区", Description: "节奏型、弹唱和舞台经验分享，想找人拍宣传照。", TimeSlots: []string{"周三晚", "周六上午"}, Rewards: []string{"技能交换", "无偿"}, Portfolio: "校园音乐节演出视频"},
		{ID: 4, Owner: "何屿", Title: "校园活动跟拍", Category: "摄影", Level: 80, Campus: "西校区", Description: "运动会、演出和乐队宣传照跟拍，自带稳定器。", TimeSlots: []string{"周六上午", "周日下午"}, Rewards: []string{"技能交换", "请吃饭"}, Portfolio: "8场校园活动成片"},
		{ID: 5, Owner: "秦朗", Title: "静物与产品摄影", Category: "摄影", Level: 84, Campus: "西校区", Description: "静物布光和电商图拍摄，可交换手冲咖啡课。", TimeSlots: []string{"周日下午"}, Rewards: []string{"技能交换"}, Portfolio: "2套电商图案例"},
		{ID: 6, Owner: "陆离", Title: "吉他指弹入门", Category: "乐器", Level: 78, Campus: "东校区", Description: "指弹基本功和扫弦节奏训练，想换人像摄影课。", TimeSlots: []string{"周三晚", "周日上午"}, Rewards: []string{"技能交换", "无偿"}, Portfolio: "指弹翻弹合集"},
	}
}

func ListNeeds() []model.Need {
	return []model.Need{
		{ID: 1, Requester: "孟野", Title: "找人帮忙拍乐队宣传照", Category: "摄影", Campus: "西校区", ExpectTime: "本周六上午", BudgetType: "技能交换", Description: "可交换 3 次吉他课，希望会调色和室外构图。", Responses: 5},
		{ID: 2, Requester: "许安", Title: "求教 Python 数据分析", Category: "编程", Campus: "中心校区", ExpectTime: "周二晚", BudgetType: "小额报酬", Description: "论文问卷数据需要清洗和画图，最好有 pandas 经验。", Responses: 8},
		{ID: 3, Requester: "林澈", Title: "想学吉他扫弦入门", Category: "乐器", Campus: "东校区", ExpectTime: "周三晚", BudgetType: "技能交换", Description: "用摄影课交换吉他基础，希望同校区或线上。", Responses: 3},
		{ID: 4, Requester: "苏漫", Title: "求毕业照跟拍", Category: "摄影", Campus: "东校区", ExpectTime: "周六上午", BudgetType: "技能交换", Description: "宿舍四人毕业照，可换两节西班牙语口语课。", Responses: 2},
		{ID: 5, Requester: "郑北", Title: "求人像修图指导", Category: "摄影", Campus: "东校区", ExpectTime: "周五晚", BudgetType: "请吃饭", Description: "想系统学 Lightroom 人像调色流程。", Responses: 1},
	}
}

// 邀请存储：互斥锁保证重复或并发确认只成功一次
var (
	invitationMu     sync.Mutex
	nextInvitationID = 5
	invitations      = []model.Invitation{
		{ID: 1, NeedID: 1, NeedTitle: "找人帮忙拍乐队宣传照", FromUser: "孟野", ToUser: "何屿", Status: constants.InvitationAccepted, Note: "双方已确认，需求达成"},
		{ID: 2, NeedID: 1, NeedTitle: "找人帮忙拍乐队宣传照", FromUser: "孟野", ToUser: "秦朗", Status: constants.InvitationExpired, Note: "何屿的邀请已被接受，本邀请失效"},
		{ID: 3, NeedID: 1, NeedTitle: "找人帮忙拍乐队宣传照", FromUser: "孟野", ToUser: "林澈", Status: constants.InvitationExpired, Note: "何屿的邀请已被接受，本邀请失效"},
		{ID: 4, NeedID: 4, NeedTitle: "求毕业照跟拍", FromUser: "苏漫", ToUser: "林澈", Status: constants.InvitationPending, Note: "等待技能方确认"},
	}
)

func ListInvitations() []model.Invitation {
	invitationMu.Lock()
	defer invitationMu.Unlock()
	out := make([]model.Invitation, len(invitations))
	copy(out, invitations)
	return out
}

// HasAcceptedInvitation 需求是否已有被接受的邀请（即已达成）
func HasAcceptedInvitation(needID int) bool {
	invitationMu.Lock()
	defer invitationMu.Unlock()
	for _, inv := range invitations {
		if inv.NeedID == needID && inv.Status == constants.InvitationAccepted {
			return true
		}
	}
	return false
}

// 邀请写入拒绝原因
const (
	RejectClosed     = "closed"
	RejectDuplicate  = "duplicate"
	RejectLimitReach = "limit"
)

// TryAddInvitation 在去重、限额、需求未达成的检查与写入之间持同一把锁，
// 保证并发发起时有效邀请数不会超过 maxActive。
func TryAddInvitation(inv model.Invitation, maxActive int) (model.Invitation, string) {
	invitationMu.Lock()
	defer invitationMu.Unlock()
	for _, existing := range invitations {
		if existing.NeedID == inv.NeedID && existing.Status == constants.InvitationAccepted {
			return model.Invitation{}, RejectClosed
		}
	}
	for _, existing := range invitations {
		if existing.NeedID == inv.NeedID && existing.ToUser == inv.ToUser && existing.Status == constants.InvitationPending {
			return existing, RejectDuplicate
		}
	}
	if activeInvitationCountLocked(inv.NeedID) >= maxActive {
		return model.Invitation{}, RejectLimitReach
	}
	inv.ID = nextInvitationID
	nextInvitationID++
	invitations = append(invitations, inv)
	return inv, ""
}

// activeInvitationCountLocked 统计需求下仍有效（待确认）的邀请数，调用方须已持有锁
func activeInvitationCountLocked(needID int) int {
	count := 0
	for _, inv := range invitations {
		if inv.NeedID == needID && inv.Status == constants.InvitationPending {
			count++
		}
	}
	return count
}

// AcceptInvitation 原子确认邀请：仅待确认状态可被接受，
// 接受成功的同时使同需求其余待确认邀请立即失效；
// 重复或并发确认返回 changed=false，保证只成功一次。
func AcceptInvitation(id int) (model.Invitation, bool) {
	invitationMu.Lock()
	defer invitationMu.Unlock()
	for i, inv := range invitations {
		if inv.ID != id {
			continue
		}
		if inv.Status != constants.InvitationPending {
			return inv, false
		}
		invitations[i].Status = constants.InvitationAccepted
		invitations[i].Note = "双方已确认，需求达成"
		for j, other := range invitations {
			if other.NeedID == inv.NeedID && other.ID != id && other.Status == constants.InvitationPending {
				invitations[j].Status = constants.InvitationExpired
				invitations[j].Note = inv.ToUser + "的邀请已被接受，本邀请失效"
			}
		}
		return invitations[i], true
	}
	return model.Invitation{}, false
}

func ListAppointments() []model.Appointment {
	return []model.Appointment{
		{ID: 1, Pair: "林澈 ↔ 孟野", Time: "周六 10:00", Place: "东校区湖边", Status: "双方已确认", Agenda: "先拍宣传照，再约 2 次吉他课"},
		{ID: 2, Pair: "周芮 ↔ 许安", Time: "周二 19:30", Place: "线上会议室", Status: "等待对方确认", Agenda: "导入问卷 CSV 并完成基础可视化"},
	}
}

func ListReviews() []model.Review {
	return []model.Review{
		{ID: 1, From: "孟野", To: "林澈", Rating: 5, Content: "构图建议很细，成片当天就给了预览。"},
		{ID: 2, From: "林澈", To: "孟野", Rating: 5, Content: "吉他入门节奏拆得很清楚，课后还发了练习谱。"},
	}
}

func ListMessages() []model.Conversation {
	return []model.Conversation{
		{ID: 1, WithUser: "孟野", Unread: 2, Messages: []string{"周六湖边光线不错", "我带两套衣服可以吗？"}},
		{ID: 2, WithUser: "系统通知", Unread: 1, Messages: []string{"你与周芮的 Python 数据分析预约待确认。"}},
	}
}

func GetProfile() model.Profile {
	return model.Profile{
		Name: "林澈", Major: "新闻传播 2023", CreditScore: 91, CreditLevel: constants.CreditGold,
		SkillWall: ListSkills()[:1],
		Radar:     map[string]int{"摄影": 92, "修图": 86, "沟通": 90, "编程": 42, "乐器": 35},
		History:   []string{"完成毕业照拍摄交换", "响应 Python 数据分析需求", "预约吉他入门课"},
		Reviews:   ListReviews(),
	}
}
