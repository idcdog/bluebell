package controller

import (
	"bluebell/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 跟社区分类相关

// CommunityHandler 查询社区列表
// @Summary 查询社区列表
// @Description 可按照社区按时间或者分数排序查询帖子列表接口
// @Tags 社区相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponseCommunitys
// @Router /community [get]
func CommunityHandler(c *gin.Context) {
	// 查询所有的社区， 以列表方式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) // 服务端报错不轻易对外暴露
	}
	ResponseSuccess(c, data)
}

// CommunityDetailHandler 查询社区分类详情
// @Summary 查询社区列表
// @Description 可按照社区按时间或者分数排序查询帖子列表接口
// @Tags 社区相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "社区id"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponseCommunityDetail
// @Router /community/{id} [get]
func CommunityDetailHandler(c *gin.Context) {
	// 获取社区id
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 查询所有的社区， 以列表方式返回
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeInvalidID) // 服务端报错不轻易对外暴露
		return
	}
	ResponseSuccess(c, data)
}
