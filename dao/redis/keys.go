package redis

// redis key使用命名空间方式， 方便查询和拆分

const (
	KeyPrefix             = "bluebell:"
	KeyPostTimeZSet       = "post:time"   // zset;帖子及发帖时间
	KeyPostScoreZSet      = "post:score"  // zset;帖子及投票分数
	KeyPostVotedPF        = "post:voted:" // zset;记录用户以及投票类型
	KeyPostCommunitySetPF = "community:"  // set 保存每个分区下帖子的id
)

// getRedisKey 为redis key添加前缀
func getRedisKey(key string) string {
	return KeyPrefix + key
}
