package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"time"
	"web_app/models"
)


func getIDsFromKey(key string,page,size int64) ([]string,error){
	start := (page - 1) * size
	end := start + size - 1
	//3. ZRecRange 查询指定数量的 从高到低，降序（ZRange是升序）
	return client.ZRevRange(key,start,end).Result()
}


func GetPostIDsInOrder(p *models.ParamPostList) ([]string , error) {
	//从redis获取id
	//1.根据用户请求中携带的order参数确定要查询的redis key
	key := getRedisKey(KeyPostTimeZSet) //默认按时间排序
	if p.Order == models.OrderScore { //如果是要按分数
		key = getRedisKey(KeyPostScoreZSet) //改成按分数
	}
	//2.确定查询的索引起始点和结尾 并返回
	return getIDsFromKey(key,p.Page,p.Size)
}


// GetPostVoteData 根据ids查询每篇帖子的投赞成票的数据
func GetPostVoteData(ids []string) (data []int64,err error) {
	//data = make([]int64,0,len(ids))
	//for _,id := range ids {
	//	key := getRedisKey(KeyPostVotedZSetPF+id)
	//	//查询范围min-max，这里都是1 其实就是查1的数量  注意参数都是string类型的
	//	v1 := client.ZCount(key,"1","1").Val()
	//	data = append(data,v1)
	//}

	//使用pineline一次发送多条命令，减少RTT
	pipeline := client.Pipeline()
	//拼接key
	for _,id := range ids {
		key := getRedisKey(KeyPostVotedZSetPF+id)
		pipeline.ZCount(key,"1","1")
	}
	cmders,err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	data = make([]int64,0,len(cmders))
	for _,cmder := range cmders {
		//类型转换 转换成redis.IntCmd类型，再转换成值
		v := cmder.(*redis.IntCmd).Val()
		data = append(data,v)
		fmt.Println("redis post data:",v)
	}
	return
}


// GetCommunityPostIDsInOrder 按社区查询ids
func GetCommunityPostIDsInOrder(p *models.ParamPostList)([]string,error){
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}

	//使用zinterstore，把分区的帖子set与帖子分数的zset 生成一个新的zset
	//针对新的zset 按之前的逻辑取数据

	//社区的key
	cKey := getRedisKey(KeyCommunitySetPF+strconv.Itoa(int(p.CommunityID)))
	//利用缓存key减少zinterstore执行的次数
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	//判断key存不存在
	if client.Exists(orderKey).Val() < 1 { //如果不存在
		//不存在，需要计算
		pipeline := client.Pipeline()
		pipeline.ZInterStore(key,redis.ZStore{ //交集存在key集合
			Aggregate: "MAX", //聚合时要执行的函数
		},cKey,orderKey) //zinterstore计算
		pipeline.Expire(key,60*time.Second) //设置超时时间
		_,err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	//存在的话直接根据key查询ids
	return getIDsFromKey(key,p.Page,p.Size)
}