package repository

import (
	"sync"
	"testing"

	"cyskillswap/internal/constants"
	"cyskillswap/internal/model"
)

func pendingInvitation(needID int, from string) model.Invitation {
	return model.Invitation{
		NeedID:       needID,
		FromUser:     from,
		ToUser:       "孟野",
		OfferSkill:   "摄影",
		WantedSkill:  "吉他",
		Campus:       "东校区",
		ProposedSlot: "周六下午",
	}
}

// 并发发起时不允许突破三项有效邀请上限，也不允许写入重复邀请。
func TestTryCreateConcurrentLimit(t *testing.T) {
	s := &invitationStore{nextID: 1}
	const max = constants.MaxActiveInvitations

	var wg sync.WaitGroup
	var mu sync.Mutex
	success, limited, duplicated := 0, 0, 0
	start := make(chan struct{})

	// 6 个不同提供方并发，仅 max 个应成功。
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			<-start
			_, err := s.tryCreate(pendingInvitation(1, "user"+string(rune('A'+idx))), max)
			mu.Lock()
			defer mu.Unlock()
			switch err {
			case nil:
				success++
			case ErrActiveLimit:
				limited++
			case ErrDuplicate:
				duplicated++
			}
		}(i)
	}
	close(start)
	wg.Wait()

	if success != max {
		t.Fatalf("成功创建 %d，期望恰好 %d", success, max)
	}
	if limited != 6-max {
		t.Fatalf("被上限拒绝 %d，期望 %d", limited, 6-max)
	}
	if duplicated != 0 {
		t.Fatalf("不应出现重复拒绝，实际 %d", duplicated)
	}
	if active := countActive(s, 1); active != max {
		t.Fatalf("有效邀请数 %d，期望 %d", active, max)
	}
}

// 同一提供方重复发起只允许一条有效邀请。
func TestTryCreateDuplicate(t *testing.T) {
	s := &invitationStore{nextID: 1}
	if _, err := s.tryCreate(pendingInvitation(1, "林澈"), constants.MaxActiveInvitations); err != nil {
		t.Fatalf("首次创建失败: %v", err)
	}
	if _, err := s.tryCreate(pendingInvitation(1, "林澈"), constants.MaxActiveInvitations); err != ErrDuplicate {
		t.Fatalf("重复创建 err=%v，期望 ErrDuplicate", err)
	}
}

// 重复或并发确认只成功一次，且同需求其余邀请立即失效。
func TestAcceptConcurrentOnce(t *testing.T) {
	s := &invitationStore{nextID: 1}
	target, _ := s.tryCreate(pendingInvitation(1, "林澈"), constants.MaxActiveInvitations)
	other1, _ := s.tryCreate(pendingInvitation(1, "姜澜"), constants.MaxActiveInvitations)
	other2, _ := s.tryCreate(pendingInvitation(1, "白珩"), constants.MaxActiveInvitations)

	var wg sync.WaitGroup
	var mu sync.Mutex
	acceptedNow := 0
	start := make(chan struct{})
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_, now, existed := s.accept(target.ID, "孟野")
			if !existed {
				t.Error("邀请应当存在")
			}
			mu.Lock()
			if now {
				acceptedNow++
			}
			mu.Unlock()
		}()
	}
	close(start)
	wg.Wait()

	if acceptedNow != 1 {
		t.Fatalf("并发确认成功 %d 次，期望恰好 1 次", acceptedNow)
	}

	byID := mapByID(s)
	if byID[target.ID].Status != constants.InvitationAccepted || !byID[target.ID].Terminal {
		t.Fatalf("目标邀请应为已接受终态，实际 %+v", byID[target.ID])
	}
	for _, id := range []int{other1.ID, other2.ID} {
		if byID[id].Status != constants.InvitationExpired || !byID[id].Terminal {
			t.Fatalf("其余邀请应立即失效，实际 %+v", byID[id])
		}
	}

	// 终态后再次确认不会重复生效。
	if _, now, _ := s.accept(target.ID, "孟野"); now {
		t.Fatal("终态邀请再次确认不应返回 acceptedNow=true")
	}
}

// 非需求方不能接受邀请，状态保持待确认。
func TestAcceptForbidden(t *testing.T) {
	s := &invitationStore{nextID: 1}
	inv, _ := s.tryCreate(pendingInvitation(1, "林澈"), constants.MaxActiveInvitations)
	got, now, existed := s.accept(inv.ID, "韩沙")
	if !existed {
		t.Fatal("邀请应当存在")
	}
	if now || got.Status != constants.InvitationPending {
		t.Fatalf("越权接受不应生效，now=%v status=%s", now, got.Status)
	}
}

func countActive(s *invitationStore, needID int) int {
	count := 0
	for _, inv := range s.list() {
		if inv.NeedID == needID && inv.Status == constants.InvitationPending {
			count++
		}
	}
	return count
}

func mapByID(s *invitationStore) map[int]model.Invitation {
	out := map[int]model.Invitation{}
	for _, inv := range s.list() {
		out[inv.ID] = inv
	}
	return out
}
