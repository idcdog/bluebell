package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
)

func GetCommunityList() (data []*models.Community, err error) {
	// 查询数据库并返回
	return mysql.GetCommunityList()
}

func GetCommunityDetail(id int64) (data *models.CommunityDetail, err error) {
	return mysql.GetCommunityDetail(id)
}
