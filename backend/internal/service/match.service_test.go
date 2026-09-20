package service

import (
	"testing"

	"cyskillswap/internal/constants"
)

func reasonMap(candidates []candidate) map[string]string {
	out := map[string]string{}
	for _, c := range candidates {
		key := c.provider.Owner + "->" + c.learner.Name
		out[key] = c.reasonCode
	}
	return out
}

func TestBuildCandidatesConstraints(t *testing.T) {
	reasons := reasonMap(buildCandidates())

	// 韩沙（西校区摄影）与孟野（东校区摄影需求）时隙重合但校区不一致。
	if got := reasons["韩沙->孟野"]; got != constants.ReasonCampusMismatch {
		t.Fatalf("韩沙->孟野 原因=%q，期望 %q", got, constants.ReasonCampusMismatch)
	}

	// 唐鹿信用 76，即使时隙与林澈法语需求重合，也先被信用门槛拦截。
	if got := reasons["唐鹿->林澈"]; got != constants.ReasonLowCredit {
		t.Fatalf("唐鹿->林澈 原因=%q，期望 %q", got, constants.ReasonLowCredit)
	}

	// 岑戈（工作日上午摄影）与孟野（周三晚/周末）无共同连续时隙。
	if got := reasons["岑戈->孟野"]; got != constants.ReasonNoCommonSlot {
		t.Fatalf("岑戈->孟野 原因=%q，期望 %q", got, constants.ReasonNoCommonSlot)
	}

	// 林澈摄影与孟野东校区摄影需求同校区、时隙连续、信用达标，应通过（无过滤原因）。
	if got := reasons["林澈->孟野"]; got != "" {
		t.Fatalf("林澈->孟野 不应被过滤，实际原因=%q", got)
	}
}

func TestMatchBoardReadOnly(t *testing.T) {
	board := MatchBoard("唐鹿")
	if !board.ReadOnly {
		t.Fatal("信用 76 的唐鹿应为只读")
	}
	for _, m := range board.Matches {
		if m.Provider == "唐鹿" || m.Learner == "唐鹿" {
			t.Fatalf("只读账号不应进入推荐：%+v", m)
		}
		if m.ViewerCanInvite {
			t.Fatalf("只读视角不应允许发起邀请：%+v", m)
		}
	}
	if board.Threshold != constants.CreditReadOnlyThreshold || board.MaxActive != constants.MaxActiveInvitations {
		t.Fatalf("门槛/上限下发错误：%d/%d", board.Threshold, board.MaxActive)
	}
}

func TestMatchBoardBasisAndBlocks(t *testing.T) {
	board := MatchBoard("林澈")
	var found bool
	for _, m := range board.Matches {
		if m.Provider == "林澈" && m.Learner == "孟野" {
			found = true
			if len(m.Basis) == 0 {
				t.Fatal("匹配依据不应为空")
			}
			if len(m.CommonBlocks) == 0 {
				t.Fatal("应至少有一段连续共同时隙")
			}
			if m.Campus != "东校区" {
				t.Fatalf("通过对校区=%q，期望 东校区", m.Campus)
			}
		}
	}
	if !found {
		t.Fatal("应包含林澈->孟野 的通过匹配")
	}
}
