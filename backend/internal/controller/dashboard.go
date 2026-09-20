package controller

import (
	"net/http"
	"strconv"

	"cyskillswap/internal/constants"
	bizerr "cyskillswap/internal/errors"
	"cyskillswap/internal/service"
	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": constants.ServiceName})
}
func Overview(c *gin.Context)     { c.JSON(http.StatusOK, service.Overview()) }
func Skills(c *gin.Context)       { c.JSON(http.StatusOK, service.Skills()) }
func Needs(c *gin.Context)        { c.JSON(http.StatusOK, service.Needs()) }
func Matches(c *gin.Context)      { c.JSON(http.StatusOK, service.Matches()) }
func Appointments(c *gin.Context) { c.JSON(http.StatusOK, service.Appointments()) }
func Reviews(c *gin.Context)      { c.JSON(http.StatusOK, service.Reviews()) }
func Messages(c *gin.Context)     { c.JSON(http.StatusOK, service.Messages()) }
func Profile(c *gin.Context)      { c.JSON(http.StatusOK, service.Profile()) }

// MatchBoard 匹配页：推荐（匹配依据+终态）、过滤原因、邀请列表
func MatchBoard(c *gin.Context) { c.JSON(http.StatusOK, service.BuildMatchBoard()) }

type createInvitationRequest struct {
	NeedID int    `json:"needId"`
	ToUser string `json:"toUser"`
}

// CreateInvitation 发起邀请，受信用分与每需求三项有效邀请上限约束
func CreateInvitation(c *gin.Context) {
	var req createInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.NeedID == 0 || req.ToUser == "" {
		c.JSON(http.StatusBadRequest, bizerr.New("INVALID_REQUEST", "needId 与 toUser 必填"))
		return
	}
	inv, err := service.CreateInvitation(req.NeedID, req.ToUser)
	if err != nil {
		status := http.StatusConflict
		if be, ok := err.(bizerr.BusinessError); ok && be.Code == bizerr.CodeCreditDenied {
			status = http.StatusForbidden
		}
		c.JSON(status, err)
		return
	}
	c.JSON(http.StatusCreated, inv)
}

// AcceptInvitation 确认邀请，重复或并发确认只成功一次
func AcceptInvitation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, bizerr.New("INVALID_REQUEST", "邀请 ID 非法"))
		return
	}
	inv, changed, err := service.AcceptInvitation(id)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"invitation": inv, "changed": changed})
}
