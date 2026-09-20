package service

import (
	"errors"
	"time"

	"cyskillswap/internal/constants"
	"cyskillswap/internal/model"
	"cyskillswap/internal/repository"
)

var (
	ErrReadOnly          = errors.New(constants.CodeReadOnly)
	ErrInvitationLimit   = errors.New(constants.CodeInvitationLimit)
	ErrDuplicateInvite   = errors.New(constants.CodeDuplicateInvite)
	ErrInvitationClosed  = errors.New(constants.CodeInvitationClosed)
	ErrInvitationMissing = errors.New(constants.CodeInvitationMissing)
	ErrForbidden         = errors.New(constants.CodeForbidden)
)

// InviteInput 发起邀请入参。
type InviteInput struct {
	Viewer       string `json:"viewer"`
	NeedID       int    `json:"needId"`
	ProposedSlot string `json:"proposedSlot"`
}

// findMatchPair 在已通过约束的推荐中定位提供方视角的匹配。
func findMatchPair(viewer string, needID int) (candidate, bool) {
	for _, c := range buildCandidates() {
		if c.reasonCode != "" {
			continue
		}
		if c.provider.Owner == viewer && c.need.ID == needID {
			return c, true
		}
	}
	return candidate{}, false
}

// CreateInvitation 发起交换邀请。
// 信用分低于门槛的账号只可浏览，不能发起；
// 每个需求最多保留三项有效邀请；同一提供方重复邀请被拒绝。
func CreateInvitation(in InviteInput) (model.Invitation, error) {
	accounts := repository.ListAccounts()
	account, ok := findAccount(accounts, in.Viewer)
	if !ok {
		return model.Invitation{}, ErrForbidden
	}
	if !account.CanInteract(constants.CreditReadOnlyThreshold) {
		return model.Invitation{}, ErrReadOnly
	}

	pair, ok := findMatchPair(in.Viewer, in.NeedID)
	if !ok {
		// 未通过同校区/连续时隙约束的组合不能发起邀请。
		return model.Invitation{}, ErrForbidden
	}

	slot := in.ProposedSlot
	if slot == "" && len(pair.blocks) > 0 {
		slot = pair.blocks[0].SlotText
	}

	inv := model.Invitation{
		MatchID:       0,
		NeedID:        in.NeedID,
		FromUser:      in.Viewer,
		ToUser:        pair.learner.Name,
		OfferSkill:    pair.provider.Title,
		WantedSkill:   pair.need.Title,
		Campus:        pair.provider.Campus,
		ProposedSlot:  slot,
		Status:        constants.InvitationPending,
		CreatedAtUnix: time.Now().Unix(),
	}
	created, err := repository.TryCreateInvitation(inv, constants.MaxActiveInvitations)
	if errors.Is(err, repository.ErrDuplicate) {
		return model.Invitation{}, ErrDuplicateInvite
	}
	if errors.Is(err, repository.ErrActiveLimit) {
		return model.Invitation{}, ErrInvitationLimit
	}
	return created, nil
}

// AcceptResult 返回接受结果与终态提示。
type AcceptResult struct {
	Invitation model.Invitation `json:"invitation"`
	Accepted   bool             `json:"accepted"`
	Message    string           `json:"message"`
}

// AcceptInvitation 需求方确认接受邀请。
// 仅需求本人可操作；重复或并发确认只成功一次；其余同需求邀请立即失效。
func AcceptInvitation(id int, viewer string) (AcceptResult, error) {
	inv, acceptedNow, existed := repository.AcceptInvitation(id, viewer)
	if !existed {
		return AcceptResult{}, ErrInvitationMissing
	}
	if inv.ToUser != viewer {
		return AcceptResult{}, ErrForbidden
	}
	msg := "邀请已接受，其余有效邀请已立即失效"
	if !acceptedNow {
		msg = "该邀请已是终态（" + inv.Status + "），重复确认不会再次生效"
	}
	return AcceptResult{Invitation: inv, Accepted: acceptedNow, Message: msg}, nil
}
