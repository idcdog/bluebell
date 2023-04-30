package redis

import (
	"bluebell/models"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func getIDsFromKey(key string, page, size int64) ([]string, error) {
	// 2、确认查询的索引起始点
	start := (page - 1) * size
	stop := start + size - 1
	// 3、zrevrange
	return rdb.ZRevRange(key, start, stop).Result()
}

func GetPostIDInOrder(p *models.ParamPostList) ([]string, error) {
	// 1、根据用户请求中携带的order参数确定要查询的redis key
	key := getRedisKey(KeyPostTimeZSet)

	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	return getIDsFromKey(key, p.Page, p.Size)
}

// GetCommunityPostIDInOrder 按照社区查询返回帖子id列表
func GetCommunityPostIDInOrder(p *models.ParamPostList) ([]string, error) {
	// 使用zinterstore吧分区的帖子set和帖子分数的zset生成一个新的zset
	// 针对新的zset按照之前的逻辑获取数据

	orderKey := getRedisKey(KeyPostTimeZSet)

	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}
	cKey := getRedisKey(KeyPostCommunitySetPF + strconv.FormatInt(p.CommunityID, 10))
	// 利用缓存key减少zinterstore操作的次数
	key := orderKey + strconv.FormatInt(p.CommunityID, 10)
	if rdb.Exists(key).Val() < 1 {
		// 不存在， 需要计算
		pipeline := rdb.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, cKey, orderKey)
		pipeline.Expire(key, 60*time.Second)
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	return getIDsFromKey(key, p.Page, p.Size)
}

// GetPostVoteData 根据ids查询每一个帖子的投票数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	// 如下方式也可以用， 但是ids比较大的时候， 要查询多次， 性能较差
	//data = make([]int64, 0, len(ids))
	//for i := 0; i < len(ids); i++ {
	//	key := getRedisKey(KeyPostVotedPF + ids[i])
	//	// 查找key中分数是1的元素的数量， 即统计每篇帖子赞成票数量
	//	v := rdb.ZCount(key, "1", "1").Val()
	//	data = append(data, v)
	//}

	// 使用pipeline一次性发送多条数据
	pipeline := rdb.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedPF + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(ids))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}

	return
}
