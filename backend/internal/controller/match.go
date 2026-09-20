package controller

import (
	"net/http"

	"cyskillswap/internal/service"
	"github.com/gin-gonic/gin"
)

// MatchBoard 返回匹配页数据：匹配依据、过滤原因、邀请终态，刷新后整体回读。
func MatchBoard(c *gin.Context) {
	viewer := c.Query("viewer")
	c.JSON(http.StatusOK, service.MatchBoard(viewer))
}
