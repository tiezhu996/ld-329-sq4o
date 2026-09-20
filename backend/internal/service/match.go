package service

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"cyskillswap/internal/constants"
	"cyskillswap/internal/model"
	"cyskillswap/internal/repository"
)

// 中文星期标签到 time.Weekday 的映射，用于把“周六上午”落入未来一周的具体日期
var weekdayIndex = map[string]time.Weekday{
	"周一": time.Monday, "周二": time.Tuesday, "周三": time.Wednesday,
	"周四": time.Thursday, "周五": time.Friday, "周六": time.Saturday, "周日": time.Sunday,
}

// normalizeSlot 归一化时段标签，如“本周六上午”→“周六上午”
func normalizeSlot(raw string) string {
	return strings.TrimPrefix(strings.TrimSpace(raw), "本")
}

// slotDate 计算时段标签在未来一周内的下一次出现日期，返回如 “09-26”
func slotDate(slot string, now time.Time) (string, bool) {
	runes := []rune(slot)
	if len(runes) < 2 {
		return "", false
	}
	weekday, ok := weekdayIndex[string(runes[:2])]
	if !ok {
		return "", false
	}
	daysAhead := (int(weekday) - int(now.Weekday()) + constants.SlotWindowDays) % constants.SlotWindowDays
	return now.AddDate(0, 0, daysAhead).Format("01-02"), true
}

// commonSlotsWithinWeek 求技能方可交换时段与需求期望时段的交集，
// 每个共同时段都是一段连续可交换时间，并标注未来一周内的具体日期。
func commonSlotsWithinWeek(skillSlots []string, expectTime string, now time.Time) []string {
	expect := normalizeSlot(expectTime)
	var common []string
	for _, raw := range skillSlots {
		slot := normalizeSlot(raw)
		if slot != expect {
			continue
		}
		if date, ok := slotDate(slot, now); ok {
			common = append(common, fmt.Sprintf("%s · %s", slot, date))
		}
	}
	return common
}

// matchStatus 计算推荐终态：需求已达成 > 该组合已邀请 > 可邀请
func matchStatus(needID int, provider string) string {
	if repository.HasAcceptedInvitation(needID) {
		return constants.MatchStateClosed
	}
	for _, inv := range repository.ListInvitations() {
		if inv.NeedID == needID && inv.ToUser == provider && inv.Status == constants.InvitationPending {
			return constants.MatchStateInvited
		}
	}
	return constants.MatchStateOpen
}

// BuildMatchBoard 生成匹配页数据：按校区、信用分、未来一周共同连续时段
// 过滤“技能 × 需求”候选，产出推荐（含匹配依据与终态）和被过滤项（含原因）。
func BuildMatchBoard() model.MatchBoard {
	now := time.Now()
	profile := repository.GetProfile()
	viewer, _ := repository.GetUser(profile.Name)

	board := model.MatchBoard{
		Viewer:       viewer.Name,
		ViewerCredit: viewer.CreditScore,
		CanInvite:    viewer.CreditScore >= constants.CreditBrowseThreshold,
		Rules: []string{
			fmt.Sprintf("推荐双方须同校区，且未来%d天内至少共有一段连续可交换时段", constants.SlotWindowDays),
			fmt.Sprintf("信用分低于%d的账号仅保留浏览权限，不进入推荐也不能发起邀请", constants.CreditBrowseThreshold),
			fmt.Sprintf("每个需求最多保留%d项有效邀请，被接受后其余邀请立即失效", constants.MaxActiveInvitationsPerNeed),
		},
		Invitations: repository.ListInvitations(),
	}

	for _, need := range repository.ListNeeds() {
		for _, skill := range repository.ListSkills() {
			if skill.Category != need.Category {
				continue
			}
			pair := skill.Owner + " × " + need.Requester
			candidate := model.FilteredCandidate{Pair: pair, Skill: skill.Title, Need: need.Title}

			provider, providerOK := repository.GetUser(skill.Owner)
			requester, requesterOK := repository.GetUser(need.Requester)
			if !providerOK || !requesterOK {
				continue
			}
			// 规则一：信用分低于80的账号不进入推荐
			if provider.CreditScore < constants.CreditBrowseThreshold || requester.CreditScore < constants.CreditBrowseThreshold {
				candidate.Reason = constants.FilterReasonCredit
				board.Filtered = append(board.Filtered, candidate)
				continue
			}
			// 规则二：推荐双方须同校区
			if skill.Campus != need.Campus {
				candidate.Reason = constants.FilterReasonCampus
				board.Filtered = append(board.Filtered, candidate)
				continue
			}
			// 规则三：未来一周至少共有一段连续可交换时段
			common := commonSlotsWithinWeek(skill.TimeSlots, need.ExpectTime, now)
			if len(common) == 0 {
				candidate.Reason = constants.FilterReasonSlot
				board.Filtered = append(board.Filtered, candidate)
				continue
			}

			score := 70 + skill.Level/5 + len(common)*5
			if score > 99 {
				score = 99
			}
			board.Recommendations = append(board.Recommendations, model.Match{
				ID:          need.ID*100 + skill.ID,
				NeedID:      need.ID,
				Provider:    skill.Owner,
				Learner:     need.Requester,
				OfferSkill:  skill.Title,
				WantedSkill: need.Title,
				Score:       score,
				CommonSlots: common,
				Recommendation: fmt.Sprintf("%s可向%s提供「%s」，双方%s可交换。",
					skill.Owner, need.Requester, skill.Title, common[0]),
				Basis: []string{
					"技能互补：" + skill.Category + " ↔ " + need.Title,
					"同校区：" + skill.Campus,
					"共同连续时段：" + strings.Join(common, "、"),
					fmt.Sprintf("双方信用分≥%d（%d / %d）", constants.CreditBrowseThreshold, provider.CreditScore, requester.CreditScore),
				},
				Status: matchStatus(need.ID, skill.Owner),
			})
		}
	}

	sort.Slice(board.Recommendations, func(i, j int) bool {
		return board.Recommendations[i].Score > board.Recommendations[j].Score
	})
	return board
}
