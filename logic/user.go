package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
)

// 存放业务逻辑代码

// SignUp 注册业务逻辑
func SignUp(p *models.ParamSignUp) (err error) {
	// 判断用户是否存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		return err
	}

	// 生成UID
	userID := snowflake.GenID()
	// 构造一个user实例
	u := &models.User{
		UserID:   userID,
		UserName: p.Username,
		Password: p.Password,
	}
	// 保存进入数据库
	err = mysql.InsertUser(u)
	return err
}

// Login 登陆业务逻辑
func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		UserName: p.Username,
		Password: p.Password,
	}
	// 传递的是指针， 可以拿到user.UserID
	if err = mysql.Login(user); err != nil {
		return
	}
	// 生成JWT Token
	aToken, _, err := jwt.GenToken(user.UserID)
	if err != nil {
		return
	}
	user.Token = aToken
	return
}
