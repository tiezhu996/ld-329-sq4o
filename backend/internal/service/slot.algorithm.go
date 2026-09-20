package service

import (
	"strings"

	"cyskillswap/internal/constants"
	"cyskillswap/internal/model"
)

// slotPoint 把「周三晚」这样的标签解析为星期序与时段序，序值用于连续性判定。
type slotPoint struct {
	weekday int
	part    int
	text    string
}

// ParseSlot 解析单个时隙标签，无法识别时返回 false。
func ParseSlot(label string) (slotPoint, bool) {
	label = strings.TrimSpace(label)
	for wi, weekday := range constants.WeekdayOrder {
		if !strings.HasPrefix(label, weekday) {
			continue
		}
		rest := strings.TrimPrefix(label, weekday)
		for pi, part := range constants.SlotPartOrder {
			if strings.Contains(rest, part) || (part == "晚上" && strings.Contains(rest, "晚")) {
				return slotPoint{weekday: wi, part: pi, text: label}, true
			}
		}
	}
	return slotPoint{}, false
}

// IntersectSlots 求两组时隙标签的交集。
func IntersectSlots(a, b []string) []string {
	set := make(map[string]struct{}, len(a))
	for _, label := range a {
		set[label] = struct{}{}
	}
	common := make([]string, 0)
	for _, label := range b {
		if _, ok := set[label]; ok {
			common = append(common, label)
		}
	}
	return common
}

// HasContinuousBlock 判定共同时隙中是否至少存在一段连续时隙。
// 规则：同一星期内相邻时段（上午→下午→晚上）视为连续；
// 若共同时隙只有一个时段，则其本身即构成长度为 1 的可交换块。
func HasContinuousBlock(common []string) (bool, []model.SlotBlock) {
	points := make([]slotPoint, 0, len(common))
	for _, label := range common {
		if p, ok := ParseSlot(label); ok {
			points = append(points, p)
		}
	}
	// 按星期、时段归组排序。
	grid := map[int]map[int]string{}
	for _, p := range points {
		if _, ok := grid[p.weekday]; !ok {
			grid[p.weekday] = map[int]string{}
		}
		grid[p.weekday][p.part] = p.text
	}

	blocks := make([]model.SlotBlock, 0)
	for wi, weekday := range constants.WeekdayOrder {
		parts, ok := grid[wi]
		if !ok {
			continue
		}
		run := []int{}
		flush := func() {
			if len(run) == 0 {
				return
			}
			texts := make([]string, 0, len(run))
			names := make([]string, 0, len(run))
			for _, pi := range run {
				texts = append(texts, parts[pi])
				names = append(names, constants.SlotPartOrder[pi])
			}
			blocks = append(blocks, model.SlotBlock{
				Weekday:  weekday,
				Parts:    names,
				SlotText: weekday + strings.Join(names, "、"),
			})
			run = []int{}
		}
		for pi := range constants.SlotPartOrder {
			if _, ok := parts[pi]; ok {
				run = append(run, pi)
			} else {
				flush()
			}
		}
		flush()
	}
	return len(blocks) > 0, blocks
}
