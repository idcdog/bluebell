package redis

import (
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// 判断投票限制
// 更新分数
// 记录用户为该帖子投过票

const (
	oneWeekInSeconds = 7 * 24 * 3600
)

var (
	scorePerVote      float64 = 432
	ErrVoteTimeExpire         = errors.New("投票时间已过")
	ErrVoteRepeated           = errors.New("已经投过票了")
)

// 投票功能
// 相关算法阅读 http://www.ruanyifeng.com/blog/algorithm/

// 投一票加432分数，86400/200 -> 需要200赞成票才能给帖子续一天

// direction=1时候， 分两种
// 1、之前没有投票，现在投赞成票， 差值绝对值 1 +432
// 2、之前投反对票，现在改投反对票， 差值绝对值 2 +432*2

// direction=0时候， 分两种
// 1、之前投赞成票， 现在取消，差值绝对值1 -432
// 2、之前投反对票，现在取消，差值绝对值1 +432

// direction=-1时候， 分两种
// 1、之前没有投票，现在投反对票，差值绝对值1 -432
// 2、之前投赞成票，现在改投反对票，差值绝对值2 -432*2

// 每个帖子自发表之日起一个星期内允许投票，超过后就不允许投票了
// 到期后
// 1、将redis中保存的赞成票以及反对票数保存到mysql中
// 2、删除redis中对应的key KeyPostVotedPF

func VoteForPost(userID, postID string, value float64) (err error) {
	// 1、从redis中拿到帖子发布时间
	postTime := rdb.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	// 2、更新帖子分数
	// 先查询当前用户给当前帖子投票的记录
	// 1表示赞成、0表示没投票、-1表示投的反对票
	ov := rdb.ZScore(getRedisKey(KeyPostVotedPF+postID), userID).Val()
	// 如果当前值等于数据库记录值， 就不做任何动作
	if ov == value {
		return ErrVoteRepeated
	}
	var dir float64
	if value > ov {
		dir = 1
	} else {
		dir = -1
	}
	diff := math.Abs(ov - value) // 计算两次投票差值
	pipeline := rdb.TxPipeline()
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), dir*diff*scorePerVote, postID)

	// 3、记录用户为该帖子投票的数据
	if value == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedPF+postID), userID)
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedPF+postID), redis.Z{
			Score:  value,
			Member: userID,
		})
	}
	_, err = pipeline.Exec()

	return
}

func CreatePost(postID, communityID int64) error {
	pipeline := rdb.TxPipeline()

	// 帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 创建帖子的时候， 将帖子id加入到community set中
	cKey := getRedisKey(KeyPostCommunitySetPF + strconv.FormatInt(communityID, 10))
	pipeline.SAdd(cKey, postID)

	_, err := pipeline.Exec()
	return err
}
