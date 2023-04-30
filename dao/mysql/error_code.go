package mysql

import "errors"

var (
	ErrorUserExist       = errors.New("用户已存在")
	ErrorUserNotExist    = errors.New("用户不存在")
	ErrorInvalidPassword = errors.New("密码错误")
	ErrorOperationDB     = errors.New("数据库操作失败")
	ErrorInvalidID       = errors.New("ID不存在")
)
