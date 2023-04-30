package controller

import (
	"errors"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

const CtxUserIDKey = "user_id"

var ErrorUserNotLogin = errors.New("用户未登陆")

// GetCurrentUser 获取当前登陆的用户id
func GetCurrentUser(c *gin.Context) (userID int64, err error) {
	uid, ok := c.Get(CtxUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}

func getPageInfo(c *gin.Context) (page, size int64) {
	pageStr := c.Query("page")
	sizeStr := c.Query("size")
	var (
		err error
	)
	page, err = strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		zap.L().Warn("page param invalid, use default 1", zap.String("page", pageStr))
		page = 1
	}
	size, err = strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		zap.L().Warn("limit param invalid, use default 10", zap.String("size", sizeStr))
		size = 10
	}
	return
}
