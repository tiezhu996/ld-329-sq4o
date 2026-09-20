package service

import (
	"fmt"

	"cyskillswap/internal/constants"
	bizerr "cyskillswap/internal/errors"
	"cyskillswap/internal/model"
	"cyskillswap/internal/repository"
)

// CreateInvitation 需求方向技能方发起邀请。
// 发起方与受邀方信用分均须≥80；同一需求最多保留三项有效邀请；
// 同一组合的待确认邀请去重；需求已达成则关闭。限额与去重在仓库层原子完成。
func CreateInvitation(needID int, toUser string) (model.Invitation, error) {
	profile := repository.GetProfile()
	fromUser := profile.Name

	from, ok := repository.GetUser(fromUser)
	if !ok || from.CreditScore < constants.CreditBrowseThreshold {
		return model.Invitation{}, bizerr.New(bizerr.CodeCreditDenied,
			fmt.Sprintf("信用分低于%d，仅保留浏览权限，不能发起邀请", constants.CreditBrowseThreshold))
	}
	to, ok := repository.GetUser(toUser)
	if !ok || to.CreditScore < constants.CreditBrowseThreshold {
		return model.Invitation{}, bizerr.New(bizerr.CodeCreditDenied,
			fmt.Sprintf("对方信用分低于%d，不能向其发起邀请", constants.CreditBrowseThreshold))
	}

	var needTitle string
	for _, n := range repository.ListNeeds() {
		if n.ID == needID {
			needTitle = n.Title
			break
		}
	}
	if needTitle == "" {
		return model.Invitation{}, bizerr.New(bizerr.CodeInvitationNotFound, "需求不存在")
	}

	inv, reject := repository.TryAddInvitation(model.Invitation{
		NeedID: needID, NeedTitle: needTitle,
		FromUser: fromUser, ToUser: toUser,
		Status: constants.InvitationPending, Note: "等待技能方确认",
	}, constants.MaxActiveInvitationsPerNeed)
	switch reject {
	case "":
		return inv, nil
	case repository.RejectClosed:
		return model.Invitation{}, bizerr.New(bizerr.CodeNeedClosed, "该需求已达成，邀请通道关闭")
	case repository.RejectDuplicate:
		return inv, bizerr.New(bizerr.CodeDuplicateInvite, "已存在相同组合的待确认邀请")
	default:
		return model.Invitation{}, bizerr.New(bizerr.CodeInvitationLimit,
			fmt.Sprintf("该需求最多保留%d项有效邀请", constants.MaxActiveInvitationsPerNeed))
	}
}

// AcceptInvitation 确认邀请：重复或并发确认只成功一次，
// changed=false 表示该邀请此前已进入终态。
func AcceptInvitation(id int) (model.Invitation, bool, error) {
	inv, changed := repository.AcceptInvitation(id)
	if inv.ID == 0 {
		return inv, false, bizerr.New(bizerr.CodeInvitationNotFound, "邀请不存在")
	}
	return inv, changed, nil
}
