package controller

type ResCode int64

const (
	CodeUnknow          = 1
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy
	CodeInvalidToken
	CodeNeedLogin
	CodeInvalidID
	CodeVoteRepeated
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "参数错误",
	CodeUserExist:       "用户名已经存在",
	CodeUserNotExist:    "用户名不存在",
	CodeInvalidPassword: "用户名或者密码错误",
	CodeServerBusy:      "服务繁忙",
	CodeUnknow:          "未知错误",
	CodeInvalidToken:    "无效token",
	CodeNeedLogin:       "需要登陆",
	CodeInvalidID:       "无效的ID",
	CodeVoteRepeated:    "重复投票",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeUnknow]
	}
	return msg
}
