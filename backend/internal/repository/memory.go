package repository

import (
	"cyskillswap/internal/constants"
	"cyskillswap/internal/model"
)

func ListAccounts() []model.Account {
	return []model.Account{
		{Name: "林澈", Major: "新闻传播 2023", Campus: "东校区", CreditScore: 91, CreditLevel: constants.CreditGold},
		{Name: "孟野", Major: "音乐表演 2022", Campus: "东校区", CreditScore: 88, CreditLevel: constants.CreditSilver},
		{Name: "周芮", Major: "统计学 2021", Campus: "中心校区", CreditScore: 93, CreditLevel: constants.CreditGold},
		{Name: "许安", Major: "社会学 2023", Campus: "中心校区", CreditScore: 85, CreditLevel: constants.CreditSilver},
		{Name: "韩沙", Major: "视觉传达 2022", Campus: "西校区", CreditScore: 89, CreditLevel: constants.CreditSilver},
		{Name: "唐鹿", Major: "法语 2021", Campus: "东校区", CreditScore: 76, CreditLevel: constants.CreditBronze},
		{Name: "沈霜", Major: "法学 2022", Campus: "东校区", CreditScore: 83, CreditLevel: constants.CreditSilver},
		{Name: "姜澜", Major: "摄影 2020", Campus: "东校区", CreditScore: 87, CreditLevel: constants.CreditSilver},
		{Name: "岑戈", Major: "体育新闻 2021", Campus: "东校区", CreditScore: 85, CreditLevel: constants.CreditSilver},
		{Name: "韦祺", Major: "市场营销 2022", Campus: "东校区", CreditScore: 84, CreditLevel: constants.CreditSilver},
		{Name: "白珩", Major: "工艺美术 2020", Campus: "东校区", CreditScore: 86, CreditLevel: constants.CreditSilver},
	}
}

func ListSkills() []model.Skill {
	return []model.Skill{
		{ID: 1, Owner: "林澈", Title: "毕业照人像摄影", Category: "摄影", Level: 92, Campus: "东校区", Description: "提供构图、修图和毕业季跟拍，可交换吉他入门课。", TimeSlots: []string{"周三晚上", "周六上午", "周六下午"}, Rewards: []string{"技能交换", "请吃饭"}, Portfolio: "12组校园人像作品"},
		{ID: 2, Owner: "周芮", Title: "Python 数据分析", Category: "编程", Level: 88, Campus: "中心校区", Description: "pandas、可视化、论文数据清洗辅导，接受小额报酬。", TimeSlots: []string{"周二晚上", "周日上午", "周日下午"}, Rewards: []string{"技能交换", "小额报酬"}, Portfolio: "3份课程项目证书"},
		{ID: 3, Owner: "孟野", Title: "民谣吉他陪练", Category: "乐器", Level: 81, Campus: "东校区", Description: "节奏型、弹唱和舞台经验分享，想找人拍宣传照。", TimeSlots: []string{"周三晚上", "周六上午", "周六下午"}, Rewards: []string{"技能交换", "无偿"}, Portfolio: "校园音乐节演出视频"},
		{ID: 4, Owner: "韩沙", Title: "乐队宣传照精修", Category: "摄影", Level: 84, Campus: "西校区", Description: "擅长室外调色与海报精修，可跨校区接单。", TimeSlots: []string{"周六上午", "周六下午"}, Rewards: []string{"技能交换", "小额报酬"}, Portfolio: "2支乐队海报作品"},
		{ID: 5, Owner: "唐鹿", Title: "法语口语陪练", Category: "外语", Level: 79, Campus: "东校区", Description: "DELF B2 口语训练和发音纠正，想用摄影课练构图。", TimeSlots: []string{"周一晚上", "周六下午"}, Rewards: []string{"技能交换"}, Portfolio: "DELF B2 证书"},
		{ID: 6, Owner: "沈霜", Title: "法语论文润色", Category: "外语", Level: 90, Campus: "东校区", Description: "法学法语双语，擅长学术写作与答辩口语，接受技能交换。", TimeSlots: []string{"周六下午"}, Rewards: []string{"技能交换", "请吃饭"}, Portfolio: "2篇法语期刊润色记录"},
		{ID: 7, Owner: "姜澜", Title: "活动纪实摄影", Category: "摄影", Level: 86, Campus: "东校区", Description: "社团活动与宣传照纪实拍摄，周六下午可约。", TimeSlots: []string{"周六下午", "周日晚上"}, Rewards: []string{"技能交换", "小额报酬"}, Portfolio: "8场活动纪实相册"},
		{ID: 8, Owner: "岑戈", Title: "运动场跟拍", Category: "摄影", Level: 82, Campus: "东校区", Description: "擅长运动瞬间抓拍，但未来一周仅工作日上午有空。", TimeSlots: []string{"周一上午", "周二上午", "周四上午", "周五上午"}, Rewards: []string{"技能交换", "请吃饭"}, Portfolio: "校运会跟拍合集"},
		{ID: 9, Owner: "韦祺", Title: "社团形象照拍摄", Category: "摄影", Level: 80, Campus: "东校区", Description: "擅长个人与社团形象照，周三晚和周六上午可交换。", TimeSlots: []string{"周三晚上", "周六上午"}, Rewards: []string{"技能交换"}, Portfolio: "3套社团形象照"},
		{ID: 10, Owner: "白珩", Title: "静物与人像后期", Category: "摄影", Level: 83, Campus: "东校区", Description: "以精修为主，周六上午、下午都可线上交付。", TimeSlots: []string{"周六上午", "周六下午"}, Rewards: []string{"技能交换", "小额报酬"}, Portfolio: "15组后期作品"},
	}
}

func ListNeeds() []model.Need {
	return []model.Need{
		{ID: 1, Requester: "孟野", Title: "找人帮忙拍乐队宣传照", Category: "摄影", Campus: "东校区", ExpectTime: "本周六下午", AvailableSlots: []string{"周三晚上", "周六上午", "周六下午"}, BudgetType: "技能交换", Description: "可交换 3 次吉他课，希望会调色和室外构图。", Responses: 5},
		{ID: 2, Requester: "许安", Title: "求教 Python 数据分析", Category: "编程", Campus: "中心校区", ExpectTime: "周二晚上", AvailableSlots: []string{"周二晚上", "周日上午"}, BudgetType: "小额报酬", Description: "论文问卷数据需要清洗和画图，最好有 pandas 经验。", Responses: 8},
		{ID: 3, Requester: "林澈", Title: "想学吉他扫弦入门", Category: "乐器", Campus: "东校区", ExpectTime: "周三晚上", AvailableSlots: []string{"周三晚上", "周六上午"}, BudgetType: "技能交换", Description: "用摄影课交换吉他基础，希望同校区或线上。", Responses: 3},
		{ID: 4, Requester: "林澈", Title: "求法语发音陪练", Category: "外语", Campus: "东校区", ExpectTime: "周六下午", AvailableSlots: []string{"周三晚上", "周六下午"}, BudgetType: "技能交换", Description: "准备 DELF 入门，希望同校区同学帮忙纠音。", Responses: 2},
	}
}

func ListAppointments() []model.Appointment {
	return []model.Appointment{
		{ID: 1, Pair: "林澈 ↔ 孟野", Time: "周六 15:00", Place: "东校区湖边", Status: "双方已确认", Agenda: "先拍宣传照，再约 2 次吉他课"},
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
