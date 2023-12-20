package logic

import (
	"go.uber.org/zap"
	"strconv"
	"web_app/dao/redis"
	"web_app/models"
)

//基于用户投票的相关算法：http://www.ruanyifeng.com/blog/algorithm/

/*
本项目使用简化版投票分数
	投一票加432分
 */

/*
投票限制：
	每个帖子发表后一个星期不允许再投票
	1.到期后将redis中保存的赞成及反对票数储存到mysql表中
	2.到期后删除 KeyPostVotedZSetPF
 */

/*
投票的记者情况：
	direction=1时：
		1.之前没投，现在赞成
		2.之前反对，现在赞成
	direction=0：
		1.之前赞成，现在取消
		2.之前反对，现在取消
	direction=-1：
		1.之前没投，现在反对
		2.之前赞成，现在反对
 */

// VoteForPost 给帖子投票
func VoteForPost(userID int64,p *models.ParamVoteData) error {
	zap.L().Debug("VoteForPost",
		zap.Int64("userID",userID),
		zap.String("postID",p.PostID),
		zap.Int8("direction", p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(userID)),p.PostID,float64(p.Direction))
}
