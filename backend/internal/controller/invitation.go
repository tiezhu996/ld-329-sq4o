package controller

import (
	"errors"
	"net/http"
	"strconv"

	errorsx "cyskillswap/internal/errors"
	"cyskillswap/internal/service"
	"github.com/gin-gonic/gin"
)

// CreateInvitation POST /api/invitations
func CreateInvitation(c *gin.Context) {
	var in service.InviteInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, errorsx.Payload("请求参数不完整"))
		return
	}
	inv, err := service.CreateInvitation(in)
	if err != nil {
		writeInviteError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"invitation": inv})
}

// AcceptInvitation POST /api/invitations/:id/accept
// 幂等：重复或并发确认只成功一次。
func AcceptInvitation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorsx.Payload("邀请编号不合法"))
		return
	}
	viewer := c.Query("viewer")
	result, err := service.AcceptInvitation(id, viewer)
	if err != nil {
		writeInviteError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func writeInviteError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrReadOnly):
		c.JSON(http.StatusForbidden, errorsx.Of("READ_ONLY_ACCOUNT", "信用分低于 80，仅保留浏览权限，不能发起或确认邀请"))
	case errors.Is(err, service.ErrInvitationLimit):
		c.JSON(http.StatusConflict, errorsx.Of("INVITATION_LIMIT", "该需求最多保留 3 项有效邀请"))
	case errors.Is(err, service.ErrDuplicateInvite):
		c.JSON(http.StatusConflict, errorsx.Of("DUPLICATE_INVITATION", "你已对该需求发起过有效邀请，请勿重复提交"))
	case errors.Is(err, service.ErrInvitationClosed):
		c.JSON(http.StatusConflict, errorsx.Of("INVITATION_CLOSED", "邀请已是终态，无法再次操作"))
	case errors.Is(err, service.ErrInvitationMissing):
		c.JSON(http.StatusNotFound, errorsx.Of("INVITATION_NOT_FOUND", "邀请不存在或已被清理"))
	case errors.Is(err, service.ErrForbidden):
		c.JSON(http.StatusForbidden, errorsx.Of("FORBIDDEN", "该匹配不满足同校区或连续时隙约束，或你无权操作"))
	default:
		c.JSON(http.StatusInternalServerError, errorsx.Of("INTERNAL", "服务暂时不可用"))
	}
}
