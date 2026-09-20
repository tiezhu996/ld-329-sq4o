package service

import (
	"sort"
	"strconv"
	"time"

	"cyskillswap/internal/constants"
	"cyskillswap/internal/model"
	"cyskillswap/internal/repository"
)

// WeekRange 返回“未来一周”的起止日期（周一至周日），供匹配页展示时隙范围。
func WeekRange() (string, string) {
	now := time.Now()
	offset := (int(now.Weekday()) + 6) % 7 // 周一为一周起点
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).
		AddDate(0, 0, -offset)
	end := start.AddDate(0, 0, 6)
	const layout = "2006-01-02"
	return start.Format(layout), end.Format(layout)
}

// candidate 互补技能对的中间结构。
type candidate struct {
	provider   model.Skill
	learner    model.Account
	need       model.Need
	rewardOK   bool
	common     []string
	blocks     []model.SlotBlock
	reasonCode string
}

// findAccount 按姓名查账号。
func findAccount(accounts []model.Account, name string) (model.Account, bool) {
	for _, acc := range accounts {
		if acc.Name == name {
			return acc, true
		}
	}
	return model.Account{}, false
}

// buildCandidates 枚举“提供技能 ↔ 互补需求”候选对，
// 并依次套用信用分、同校区、未来一周连续时隙三道约束，
// 未通过的候选附带首个命中的过滤原因。
func buildCandidates() []candidate {
	skills := repository.ListSkills()
	needs := repository.ListNeeds()
	accounts := repository.ListAccounts()

	candidates := make([]candidate, 0)
	for _, skill := range skills {
		provider, ok := findAccount(accounts, skill.Owner)
		if !ok {
			continue
		}
		for _, need := range needs {
			if need.Requester == skill.Owner || need.Category != skill.Category {
				continue
			}
			learner, ok := findAccount(accounts, need.Requester)
			if !ok {
				continue
			}
			c := candidate{
				provider: skill, learner: learner, need: need,
				rewardOK: RewardCompatible(skill, need),
			}
			common := IntersectSlots(skill.TimeSlots, need.AvailableSlots)
			sort.Strings(common)
			c.common = common
			_, c.blocks = HasContinuousBlock(common)

			switch {
			case provider.CreditScore < constants.CreditReadOnlyThreshold:
				c.reasonCode = constants.ReasonLowCredit
			case learner.CreditScore < constants.CreditReadOnlyThreshold:
				c.reasonCode = constants.ReasonLowCredit
			case skill.Campus != need.Campus:
				c.reasonCode = constants.ReasonCampusMismatch
			case len(c.blocks) == 0:
				c.reasonCode = constants.ReasonNoCommonSlot
			}
			candidates = append(candidates, c)
		}
	}
	return candidates
}

func blockText(blocks []model.SlotBlock) []string {
	texts := make([]string, 0, len(blocks))
	for _, b := range blocks {
		texts = append(texts, b.SlotText)
	}
	sort.Strings(texts)
	return texts
}

// MatchBoard 组装某查看者视角的匹配页数据：
// 通过约束的互补对进入推荐，其余以过滤原因呈现；
// 低信用查看者仅可浏览，推荐结果不再标记其可发起邀请。
func MatchBoard(viewer string) model.MatchBoard {
	accounts := repository.ListAccounts()
	viewerAccount, ok := findAccount(accounts, viewer)
	if !ok {
		viewerAccount = accounts[0]
	}
	readOnly := !viewerAccount.CanInteract(constants.CreditReadOnlyThreshold)
	weekStart, weekEnd := WeekRange()

	// 预先统计每个需求的有效邀请数与提供方占用情况，用于判断是否还能发起。
	all := repository.ListInvitations()
	active := 0
	activeByNeed := map[int]int{}
	providerPending := map[string]bool{} // key: needID|fromUser
	for _, inv := range all {
		if inv.Status != constants.InvitationPending {
			continue
		}
		active++
		activeByNeed[inv.NeedID]++
		providerPending[invitationKey(inv.NeedID, inv.FromUser)] = true
	}

	matched := make([]model.Match, 0)
	filtered := make([]model.FilteredMatch, 0)
	matchID := 0

	for _, c := range buildCandidates() {
		if c.reasonCode != "" {
			filtered = append(filtered, model.FilteredMatch{
				Provider:    c.provider.Owner,
				Learner:     c.learner.Name,
				Campus:      c.provider.Campus + " / " + c.need.Campus,
				OfferSkill:  c.provider.Title,
				WantedSkill: c.need.Title,
				ReasonCode:  c.reasonCode,
				Reason:      constants.FilterReasonText[c.reasonCode],
				CommonSlots: c.common,
			})
			continue
		}
		matchID++
		score := ScoreMatch(c.provider, c.learner, c.rewardOK, len(c.blocks))
		basis := []string{constants.BasisCampus, constants.BasisSkill, constants.BasisSlot, constants.BasisCredit}
		if c.rewardOK {
			basis = append(basis, constants.BasisReward)
		}
		needFull := activeByNeed[c.need.ID] >= constants.MaxActiveInvitations
		alreadyInvited := providerPending[invitationKey(c.need.ID, c.provider.Owner)]
		canInvite := !readOnly && c.provider.Owner == viewerAccount.Name && !needFull && !alreadyInvited
		matched = append(matched, model.Match{
			ID: matchID, NeedID: c.need.ID, Provider: c.provider.Owner, Learner: c.learner.Name,
			Campus: c.provider.Campus, OfferSkill: c.provider.Title,
			WantedSkill: c.need.Title, Category: c.provider.Category, Score: score,
			CommonSlots: c.common, CommonBlocks: c.blocks,
			Reward: c.need.BudgetType, Basis: basis,
			Recommendation: buildRecommendation(c),
			ViewerInvolved: !readOnly &&
				(c.provider.Owner == viewerAccount.Name || c.learner.Name == viewerAccount.Name),
			ViewerCanInvite: canInvite,
		})
	}

	return model.MatchBoard{
		Viewer: viewerAccount.Name, ViewerCampus: viewerAccount.Campus,
		ViewerCredit: viewerAccount.CreditScore, ReadOnly: readOnly,
		WeekStart: weekStart, WeekEnd: weekEnd,
		Threshold: constants.CreditReadOnlyThreshold, MaxActive: constants.MaxActiveInvitations,
		Accounts: accounts, Matches: matched, Filtered: filtered,
		Invitations: all, ActiveCount: active, ActiveByNeed: activeByNeed,
	}
}

// invitationKey 生成「需求|提供方」维度的占用键。
func invitationKey(needID int, fromUser string) string {
	return strconv.Itoa(needID) + "|" + fromUser
}

func buildRecommendation(c candidate) string {
	reward := "回报类型不一致"
	if c.rewardOK {
		reward = "双方均接受「" + c.need.BudgetType + "」"
	}
	return "同属" + c.provider.Campus + "，" + c.provider.Title + " 互补 " +
		c.need.Title + "，" + reward + "，未来一周可约 " + blockText(c.blocks)[0] + "。"
}
