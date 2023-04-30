package controller

import "bluebell/models"

type _ResponseCommunitys struct {
	Code    ResCode             `json:"code"`    // 业务响应状态码
	Message string              `json:"message"` // 提示信息
	Data    []*models.Community `json:"data"`    // 数据
}
type _ResponseCommunityDetail struct {
	Code    ResCode                   `json:"code"`    // 业务响应状态码
	Message string                    `json:"message"` // 提示信息
	Data    []*models.CommunityDetail `json:"data"`    // 数据
}
