package controller

import (
	"bluebell/dao/redis"
	"bluebell/logic"
	"bluebell/models"
	"errors"

	"go.uber.org/zap"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

func PostVoteHandler(c *gin.Context) {
	// 参数校验
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		errData := removeTopStruct(errs.Translate(trans)) // 翻译并去掉错误提示中的结构体标识
		ResponseErrorWithMsg(c, CodeInvalidParam, errData)
		return
	}

	userID, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	// 处理
	// 返回响应
	if err := logic.VoteForPost(userID, p); err != nil {
		zap.L().Error("vote failed", zap.Int64("userID", userID), zap.Int64("postID", p.PostID), zap.Int8("direction", p.Direction), zap.Error(err))
		if errors.Is(err, redis.ErrVoteRepeated) {
			ResponseError(c, CodeVoteRepeated)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}
