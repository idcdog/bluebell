package mysql

import (
	"bluebell/models"
	"database/sql"

	"go.uber.org/zap"
)

func GetCommunityList() (data []*models.Community, err error) {
	sqlStr := "select community_id,community_name from community"
	if err := db.Select(&data, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("select community is null")
			err = nil
		}
	}
	return
}

func GetCommunityDetail(id int64) (communityDetail *models.CommunityDetail, err error) {
	sqlStr := "select " +
		"community_id,community_name,introduction,create_time " +
		"from community " +
		"where community_id=?"
	communityDetail = new(models.CommunityDetail)
	if err = db.Get(communityDetail, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("select community is null")
			err = ErrorInvalidID
		}
	}
	return
}
