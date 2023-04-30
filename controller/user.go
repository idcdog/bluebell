package controller

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"errors"
	"strconv"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	// 1、获取请求参数及参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("Login with invalid param", zap.Error(err))
		// 判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}

		ResponseErrorWithMsg(c, CodeServerBusy, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 2、业务逻辑处理
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("用户登陆失败", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorInvalidPassword) {
			ResponseError(c, CodeInvalidPassword)
			return
		}
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	zap.L().Error("用户登陆成功", zap.String("username", p.Username))
	// 3、返回响应
	ResponseSuccess(c, gin.H{
		"user_id":   strconv.FormatInt(user.UserID, 10), // id 大于 1<<53-1, js展示会出现失真
		"user_name": user.UserName,
		"token":     user.Token,
	})
}

func SignUpHandler(c *gin.Context) {
	// 1、获取参数和参数校验
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		// 判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))

		return
	}

	// 2、业务处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("注册失败", zap.Error(err))

		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3、返回响应
	ResponseSuccess(c, "注册成功")
}

// 自定义结构体的校验方法
func SignUpParamStructLevelValidation(sl validator.StructLevel) {
	su := sl.Current().Interface().(models.ParamSignUp)
	if su.Password != su.RePassword {
		sl.ReportError(su.RePassword, "re_password", "RePassword", "eqfield", "password")
	}
}
