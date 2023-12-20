package redis

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"math"
	"strconv"
	"time"
)


/*
投票限制：
	每个帖子发表后一个星期不允许再投票
	1.到期后将redis中保存的赞成及反对票数储存到mysql表中
	2.到期后删除 KeyPostVotedZSetPF
*/

/*
投票的记者情况：
	direction=1时：
		1.之前没投，现在赞成  	新旧票差值的绝对值：1+		分数改动：+432
		2.之前反对，现在赞成 	新旧票差值的绝对值：2+				+432*2
	direction=0：
		1.之前赞成，现在取消 	新旧票差值的绝对值：1-				-432
		2.之前反对，现在取消 	新旧票差值的绝对值：1+				+432
	direction=-1：
		1.之前没投，现在反对 	新旧票差值的绝对值：1-				-432
		2.之前赞成，现在反对 	新旧票差值的绝对值：2-				-432*2
*/


const(
	oneWeekInSecond = 7*24*3600 //一个星期的秒数
	scorePerVote = 432 //每一票值多少分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated = errors.New("不允许重复投票")
)


// CreatePost 创建帖子
func CreatePost(postId,communityID int64) error {
	////帖子时间
	//_,err := client.ZAdd(getRedisKey(KeyPostTimeZSet),redis.Z{
	//	Score:  float64(time.Now().Unix()),
	//	Member: postId,
	//}).Result()
	//
	////帖子分数 投一票续432秒
	//_,err = client.ZAdd(getRedisKey(KeyPostScoreZSet),redis.Z{
	//	Score:  float64(time.Now().Unix()),
	//	Member: postId,
	//}).Result()

	//使用事务 下面创建的两件事，要么同时成功，要么同时失败
	pipeline := client.TxPipeline() //获取一个事务
	//帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet),redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postId,
	})
	//帖子分数 投一票续432秒
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet),redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postId,
	})
	//更新：把帖子id加到对应社区的set
	cKey := getRedisKey(KeyCommunitySetPF+strconv.Itoa(int(communityID)))
	pipeline.SAdd(cKey,postId)
	_,err := pipeline.Exec()//真正执行
	return err
}


// VoteForPost 用float64原因：redis涉及有序集合分数的都是用float
func VoteForPost(userID,postID string,value float64) error{
	fmt.Println("v:",value)
	//1.判断投票限制
	//根据postID 去redis取帖子发布时间
	postTime := client.ZScore(getRedisKey(KeyPostTimeZSet),postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSecond {
		//查过时间限制
		return ErrVoteTimeExpire
	}

	//2的部分和3需要放到一个pipeline事务中操作
	//2.更新帖子分数
	//根据userID 查当前用户给当前帖子之前的投票纪录
	//ov是旧分数
	ov := client.ZScore(getRedisKey(KeyPostVotedZSetPF+postID),userID).Val()
	//如果这一次投票的值和之前保存的值一致，就提示不允许重复投票
	if value == ov {
		return ErrVoteRepeated
	}
	var op float64 //dir是新的选择 决定正负
	if value > ov {
		op = 1
	}else{
		op = -1
	}
	diff := math.Abs(ov - value) //两次投票的差值 决定是1还是2倍
	////更新分数
	//_,err := client.ZIncrBy(getRedisKey(KeyPostScoreZSet),op*diff*scorePerVote,postID).Result()
	///////////////////////////
	////if err != nil {
	////	return ErrVoteTimeExpire
	////}
	//if ErrVoteTimeExpire!=nil {
	//	return err
	//}
	////3.记录用户为该帖子投票的数据
	////value 表示当前新票是赞成还是反对或是取消之前的投票
	//if value == 0 { //取消
	//	//删除一个用户
	//	_,err = client.ZRem(getRedisKey(KeyPostVotedZSetPF+postID),userID).Result()
	//}else{
	//	_,err = client.ZAdd(getRedisKey(KeyPostVotedZSetPF+postID),redis.Z{
	//		Score:  value,
	//		Member: userID,
	//	}).Result()
	//}

	pipeline := client.TxPipeline() //获取一个事务
	//更新分数
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet),op*diff*scorePerVote,postID)
	//3.记录用户为该帖子投票的数据
	//value 表示当前新票是赞成还是反对或是取消之前的投票
	if value == 0 { //取消
		//删除一个用户
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPF+postID),userID)
	}else{
		pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPF+postID),redis.Z{
			Score:  value,
			Member: userID,
		})
	}
	_,err := pipeline.Exec() //真正执行
	fmt.Println("v:",value)
	return err
}