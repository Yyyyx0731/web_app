package mysql

import (
	"github.com/jmoiron/sqlx"
	"strings"
	"web_app/models"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post)(err error){
	sqlStr := `insert into post (post_id,title,content,author_id,community_id)values (?,?,?,?,?)`
	_,err = db.Exec(sqlStr,p.ID,p.Title,p.Content,p.AuthorID,p.CommunityID)
	return
}

// GetPostByID 根据id查询单个帖子详情
func GetPostByID(pid int64)(post *models.Post,err error){
	post = new(models.Post)
	sqlStr := "select post_id,title,content,author_id,community_id,create_time from post where post_id=?"
	err = db.Get(post,sqlStr,pid)
	return
}

// GetPostList 查询帖子列表
func GetPostList(page,size int64) (posts []*models.Post,err error){
	//按创建时间递减（实现新帖在前）
	sqlStr := `select post_id,title,content,author_id,community_id,create_time 
				from post
				order by create_time desc 
				limit ?,? `
	posts = make([]*models.Post,0,2)
	err = db.Select(&posts,sqlStr,(page-1)*size,size)
	return
}

// GetPostListByIDs 根据给定的id列表(切片ids) 查询帖子
func GetPostListByIDs(ids []string) (postList []*models.Post,err error){
	//FIND_IN_SET(str,strList)是mysql内置函数，查询s列表中是否有str；返回记录或NULL.strList之间用逗号分隔
	//例：FIND_IN_SET(1,2,3,1,4,5) 返回:3(2的下标是1)
	//在这里 FIND_IN_SET 实现了按照指定顺序查询
	sqlStr := `select post_id,title,content,author_id,community_id,create_time
				from post 
				where post_id in (?) 
				order by FIND_IN_SET(post_id,?)`
	//注意区分post表里 id字段和post_id字段不是一回事，id是自增的，post_id才是帖子的id

	//FIND_IN_SET的strList是用逗号分隔的，所以下面需要手动用逗号分隔切片，拼成字符串
	//query是拼接好的查询语句   args是变量切片，接口类型
	query,args,err := sqlx.In(sqlStr,ids,strings.Join(ids,","))
	if err != nil {
		return nil, err
	}
	//Rebind重新绑定query语句
	query = db.Rebind(query)
	err = db.Select(&postList,query,args...) //切片要加三个点别忘了
	return
}