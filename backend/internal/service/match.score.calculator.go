package service

import (
	"cyskillswap/internal/model"
)

// 匹配度评分权重，合计 100。
const (
	weightSkillLevel = 35 // 提供技能熟练度
	weightSlot       = 30 // 共同连续时隙数量
	weightCredit     = 20 // 双方信用分
	weightReward     = 15 // 回报类型兼容
)

// ScoreMatch 计算互补技能对的匹配度（0-100）。
func ScoreMatch(provider model.Skill, learner model.Account, rewardMatched bool, blockCount int) int {
	levelScore := provider.Level * weightSkillLevel / 100

	slotScore := weightSlot
	switch {
	case blockCount >= 3:
		slotScore = weightSlot
	case blockCount == 2:
		slotScore = weightSlot * 8 / 10
	default:
		slotScore = weightSlot * 6 / 10
	}

	creditScore := (provider.Level + learner.CreditScore) / 2 * weightCredit / 100

	rewardScore := 0
	if rewardMatched {
		rewardScore = weightReward
	}

	total := levelScore + slotScore + creditScore + rewardScore
	if total > 100 {
		total = 100
	}
	return total
}

// RewardCompatible 判定技能方接受的回报类型是否覆盖需求方预算类型。
func RewardCompatible(provider model.Skill, need model.Need) bool {
	for _, reward := range provider.Rewards {
		if reward == need.BudgetType {
			return true
		}
	}
	return false
}
