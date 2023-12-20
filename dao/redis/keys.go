package redis


//redis key 注意使用，命名空间的方式，方便查询和拆分（下例使用冒号）

const (
	Prefix = "web_app:"

	KeyPostTimeZSet = "post:time" // zset;帖子及发帖时间 （zset是redis的一种数据类型）
	KeyPostScoreZSet = "post:score" // zset;帖子及投票所得的分数
	KeyPostVotedZSetPF = "post:voted:" //zset;记录用户及投票类型;参数是post_id(PF是prefix，表示数据还需要拼接其他参数才完整，这里是post_id）

	KeyCommunitySetPF = "community:" //set;保存每个分区下帖子的id
)

// getRedisKey 给redis key加上前缀
func getRedisKey(key string) string {
	return Prefix+key
}

//4.38