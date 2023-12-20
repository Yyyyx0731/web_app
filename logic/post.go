package logic

import (
	"fmt"
	"go.uber.org/zap"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/models"
	"web_app/pkg/snowflake"
)

func CreatePost(p *models.Post) (err error){
	//1.生成postID
	p.ID = snowflake.GenID()
	//2.保存到数据库 //3.返回
	err =  mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.ID,p.CommunityID)
	return
}


func GetPostByID(pid int64) (data *models.ApiPostDetail,err error){
	//data = new(models.ApiPostDetail)
	post,err := mysql.GetPostByID(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostByID(pid) failed",zap.Int64("pid",pid),zap.Error(err))
		return
	}
	//根据作者id查作者信息
	user,err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserByID(post.AuthorID) failed",zap.Int64("author_id",post.AuthorID),zap.Error(err))
		return
	}
	//根据社区id查社区详情
	community,err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",zap.Int64("community_id",post.AuthorID),zap.Error(err))
		return
	}
	//接口数据拼接
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}
	return
}

// GetPostList 获取帖子列表
func GetPostList(page,size int64)(data []*models.ApiPostDetail,err error){
	posts,err := mysql.GetPostList(page,size)
	if err!=nil {
		return nil, err
	}

	//初始化帖子列表切片data
	data = make([]*models.ApiPostDetail,0,len(posts))

	for _,post := range posts {
		//根据作者id查作者信息
		user,err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID(post.AuthorID) failed",zap.Int64("author_id",post.AuthorID),zap.Error(err))
			continue
		}
		//根据社区id查社区详情
		community,err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",zap.Int64("community_id",post.AuthorID),zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:           post,
			CommunityDetail: community,
		}
		data = append(data,postDetail)
	}
	return
}



// GetPostList2 获取帖子列表
func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail,err error){
	//2.去redis查id列表
	ids,err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	//3.根据id去数据库查帖子详情
	if len(ids)==0 { //如果id列表是空的
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}
	zap.L().Debug("GetPostList2",zap.Any("ids",ids))
	//按照给定的id顺序返回查询结果
	posts,err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	zap.L().Debug("GetPostList2",zap.Any("posts",posts))

	//根据帖子的ids提前查好每个帖子的投票数
	voteData,err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	//将每个帖子的作者信息和分区信息查询出来 和每个post一起拼接在data结构体切片
	for idx,post := range posts {
		//根据作者id 查作者信息
		user,err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID(post.AuthorID) failed",zap.Int64("author_id",post.AuthorID),zap.Error(err))
			continue
		}
		//根据社区id 查社区详情
		community,err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",zap.Int64("community_id",post.AuthorID),zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum: 		voteData[idx],
			Post:           post,
			CommunityDetail: community,
		}
		fmt.Println(idx,voteData[idx])
		data = append(data,postDetail)
	}
	return
}


func GetCommunityPostList2(p *models.ParamPostList)(data []*models.ApiPostDetail,err error){
	//2.去redis查id列表
	ids,err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}
	//3.根据id去数据库查帖子详情
	if len(ids)==0 { //如果id列表是空的
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}
	zap.L().Debug("GetPostList2",zap.Any("ids",ids))
	//按照给定的id顺序返回查询结果
	posts,err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	zap.L().Debug("GetPostList2",zap.Any("posts",posts))

	//根据帖子的ids提前查好每个帖子的投票数
	voteData,err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	//将每个帖子的作者信息和分区信息查询出来 和每个post一起拼接在data结构体切片
	for idx,post := range posts {
		//根据作者id 查作者信息
		user,err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID(post.AuthorID) failed",zap.Int64("author_id",post.AuthorID),zap.Error(err))
			continue
		}
		//根据社区id 查社区详情
		community,err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",zap.Int64("community_id",post.AuthorID),zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum: 		voteData[idx],
			Post:           post,
			CommunityDetail: community,
		}
		fmt.Println(idx,voteData[idx])
		data = append(data,postDetail)
	}
	return
}


// GetPostListNew 将两个查询帖子列表的接口合二为一
//传了社区id和没传
func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail,err error){
	if p.CommunityID == 0{ //没传社区id，默认查所有社区的帖子
		//获取数据
		data,err = GetPostList2(p)
	}else{ //根据社区id查
		//获取数据
		data,err = GetCommunityPostList2(p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed",zap.Error(err))
		return nil, err
	}
	return
}//12
