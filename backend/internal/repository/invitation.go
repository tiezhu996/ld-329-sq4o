package repository

import (
	"errors"
	"sync"

	"cyskillswap/internal/constants"
	"cyskillswap/internal/model"
)

// 邀请写入阶段的并发安全错误，由 service 层映射为业务错误码。
var (
	ErrActiveLimit = errors.New("active invitation limit reached")
	ErrDuplicate   = errors.New("duplicate active invitation")
)

// invitationStore 进程内邀请存储。使用互斥锁保证并发确认只成功一次。
type invitationStore struct {
	mu          sync.Mutex
	nextID      int
	invitations []model.Invitation
}

func newInvitationStore() *invitationStore {
	s := &invitationStore{nextID: 100}
	s.seed()
	return s
}

var invites = newInvitationStore()

// seed 预置演示与刷新回读所需的邀请数据。
func (s *invitationStore) seed() {
	// 两条针对需求 1（孟野「乐队宣传照」）的待确认邀请，
	// 其中任一被接受后，另一条应立即失效，用于演示上限与终态规则。
	s.addSeeded(model.Invitation{
		MatchID: 1, NeedID: 1, FromUser: "林澈", ToUser: "孟野",
		OfferSkill: "毕业照人像摄影", WantedSkill: "民谣吉他陪练",
		Campus: "东校区", ProposedSlot: "周六下午", Status: constants.InvitationPending,
	})
	s.addSeeded(model.Invitation{
		MatchID: 2, NeedID: 1, FromUser: "姜澜", ToUser: "孟野",
		OfferSkill: "活动纪实摄影", WantedSkill: "民谣吉他陪练",
		Campus: "东校区", ProposedSlot: "周六下午", Status: constants.InvitationPending,
	})
	// 一条已失效终态：历史邀请，因需求方接受了别人而失效，刷新后仍可回读。
	s.addSeeded(model.Invitation{
		MatchID: 3, NeedID: 2, FromUser: "林澈", ToUser: "许安",
		OfferSkill: "毕业照人像摄影", WantedSkill: "Python 数据分析",
		Campus: "中心校区", ProposedSlot: "周日上午", Status: constants.InvitationExpired, Terminal: true,
		Note: "需求方已接受其他邀请",
	})
}

func (s *invitationStore) addSeeded(inv model.Invitation) {
	inv.ID = s.nextID
	s.nextID++
	s.invitations = append(s.invitations, inv)
}

// ListInvitations 返回全部邀请（含终态），供匹配页回读。
func ListInvitations() []model.Invitation { return invites.list() }

func (s *invitationStore) list() []model.Invitation {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]model.Invitation, len(s.invitations))
	copy(out, s.invitations)
	return out
}

// 有效邀请数与重复检查由 tryCreate 在同一把锁内完成，
// 避免“先查后写”在并发下突破三项上限。

// TryCreateInvitation 包级入口，委托给进程内单例。
func TryCreateInvitation(inv model.Invitation, maxActive int) (model.Invitation, error) {
	return invites.tryCreate(inv, maxActive)
}

// tryCreate 在同一把锁内原子地完成查重、上限校验与写入，
// 保证并发发起时不会突破三项有效邀请上限，也不会写入重复邀请。
func (s *invitationStore) tryCreate(inv model.Invitation, maxActive int) (model.Invitation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	active := 0
	for _, existing := range s.invitations {
		if existing.NeedID != inv.NeedID || existing.Status != constants.InvitationPending {
			continue
		}
		if existing.FromUser == inv.FromUser {
			return model.Invitation{}, ErrDuplicate
		}
		active++
	}
	if active >= maxActive {
		return model.Invitation{}, ErrActiveLimit
	}

	inv.ID = s.nextID
	s.nextID++
	inv.Status = constants.InvitationPending
	s.invitations = append(s.invitations, inv)
	return inv, nil
}

// AcceptInvitation 包级入口，委托给进程内单例。
func AcceptInvitation(id int, operator string) (model.Invitation, bool, bool) {
	return invites.accept(id, operator)
}

// accept 幂等接受邀请。
// 仅需求方本人可操作；同一邀请重复或并发确认只成功一次；
// 接受成功后该需求其余待确认邀请立即失效。
// 返回值：当前邀请副本、是否“本次”完成接受、邀请是否存在。
func (s *invitationStore) accept(id int, operator string) (model.Invitation, bool, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := -1
	for i := range s.invitations {
		if s.invitations[i].ID == id {
			idx = i
			break
		}
	}
	if idx < 0 {
		return model.Invitation{}, false, false
	}
	target := &s.invitations[idx]
	if target.ToUser != operator {
		// 邀请存在但操作者不是需求方：交由上层判定为越权。
		return *target, false, true
	}
	if target.Status != constants.InvitationPending {
		// 终态：重复或并发确认不重复生效。
		return *target, false, true
	}

	nowUnix := nowFunc().Unix()
	target.Status = constants.InvitationAccepted
	target.Terminal = true
	target.DecidedAtUnix = nowUnix

	// 同一需求的其余待确认邀请立即失效。
	for i := range s.invitations {
		other := &s.invitations[i]
		if other.ID != target.ID && other.NeedID == target.NeedID &&
			other.Status == constants.InvitationPending {
			other.Status = constants.InvitationExpired
			other.Terminal = true
			other.DecidedAtUnix = nowUnix
			other.Note = "需求方已接受其他邀请"
		}
	}
	return *target, true, true
}
