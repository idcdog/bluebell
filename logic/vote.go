package logic

import (
	"bluebell/dao/redis"
	"bluebell/models"
	"strconv"

	"go.uber.org/zap"
)

// VoteForPost 为帖子投票
func VoteForPost(userID int64, p *models.ParamVoteData) error {
	zap.L().Debug("vote info", zap.Int64("userID", userID), zap.Int64("postID", p.PostID), zap.Int8("direction", p.Direction))

	return redis.VoteForPost(
		strconv.FormatInt(userID, 10),
		strconv.FormatInt(p.PostID, 10),
		float64(p.Direction))
}
