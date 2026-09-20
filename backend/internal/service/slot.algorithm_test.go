package service

import (
	"testing"
)

func TestIntersectSlots(t *testing.T) {
	a := []string{"周三晚上", "周六上午", "周六下午"}
	b := []string{"周六上午", "周六下午", "周日晚上"}
	got := IntersectSlots(a, b)
	if len(got) != 2 {
		t.Fatalf("交集长度 %d，期望 2 (%v)", len(got), got)
	}
}

func TestHasContinuousBlock(t *testing.T) {
	cases := []struct {
		name   string
		slots  []string
		want   bool
		blocks int
		maxRun int
	}{
		{
			name:   "相邻时段合并为连续块",
			slots:  []string{"周六上午", "周六下午"},
			want:   true,
			blocks: 1,
			maxRun: 2,
		},
		{
			name:   "同一天三段全连续",
			slots:  []string{"周日上午", "周日下午", "周日晚上"},
			want:   true,
			blocks: 1,
			maxRun: 3,
		},
		{
			name:   "跨天各自成块",
			slots:  []string{"周三晚上", "周六上午", "周六下午"},
			want:   true,
			blocks: 2,
			maxRun: 2,
		},
		{
			name:   "单个时段也算长度为1的可交换块",
			slots:  []string{"周六下午"},
			want:   true,
			blocks: 1,
			maxRun: 1,
		},
		{
			name:   "无法识别的标签不产生块",
			slots:  []string{"下周三傍晚"},
			want:   false,
			blocks: 0,
			maxRun: 0,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ok, blocks := HasContinuousBlock(tc.slots)
			if ok != tc.want || len(blocks) != tc.blocks {
				t.Fatalf("slots=%v ok=%v blocks=%d，期望 ok=%v blocks=%d",
					tc.slots, ok, len(blocks), tc.want, tc.blocks)
			}
			maxRun := 0
			for _, b := range blocks {
				if len(b.Parts) > maxRun {
					maxRun = len(b.Parts)
				}
			}
			if maxRun != tc.maxRun {
				t.Fatalf("最长连续段 %d，期望 %d", maxRun, tc.maxRun)
			}
		})
	}
}
