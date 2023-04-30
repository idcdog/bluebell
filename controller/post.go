package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"strconv"

	"github.com/spf13/viper"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func CreatePostHandler(c *gin.Context) {
	// 获取参数并校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("创建post参数错误", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 获取用户id
	userID, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID

	// 创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost error", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 获取帖子详情
func GetPostDetailHandler(c *gin.Context) {
	// 获取参数
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("帖子id无效", zap.Error(err))
		ResponseError(c, CodeInvalidID)
		return
	}

	// 根据id取出帖子数据
	data, err := logic.GetPostById(pid)
	if err != nil {
		zap.L().Error("logic.GetPostById failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 返回数据
	ResponseSuccess(c, data)

}

// GetPostListHandler2 升级版获取帖子列表接口
// 根据前端传入参数动态获取列表
// 按照创建时间排序， 或 按照分数排序
// 1、获取请求的query string参数
// 2、去redis中查询id列表
// 3、根据id去数据库中查询帖子的详细信息
func GetPostListHandler2(c *gin.Context) {
	// 1、获取参数
	// 2、从redis查询id列表
	// 3、根据id从数据库中查询帖子详细信息

	// 获取分页参数
	page := viper.GetInt64("page")
	size := viper.GetInt64("size")
	mps := viper.GetInt64("max_page_size")
	if size > mps {
		size = mps
	}
	p := &models.ParamPostList{
		Page:  page,
		Size:  size,
		Order: models.OrderTime,
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	data, err := logic.GetPostListNew(p)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	if len(data) == 0 {
		ResponseSuccess(c, nil)
		return
	}
	ResponseSuccess(c, data)
}

// 获取接口列表
func GetPostListHandler(c *gin.Context) {
	// 获取分页参数
	page, size := getPageInfo(c)
	// 获取数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	if len(data) == 0 {
		ResponseSuccess(c, nil)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}
