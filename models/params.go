package models


const(
	OrderTime = "time"
	OrderScore = "score"

	DefaultPage = 1
	DefaultSize = 10
)



// 定义请求的参数结构体

// ParamSignUp 注册请求参数
type ParamSignUp struct {
	//gin框架中validator库用tag binding做参数校验
	//required代表需要该字段，如果缺少则会在ShouldBindJSON报错
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// ParamLogin 登录请求参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}


// ParamVoteData 投票数据
type ParamVoteData struct {
	//UserID 从请求中获取当前用户id
	PostID string `json:"post_id" binding:"required"`
	Direction int8 `json:"direction,string" binding:"oneof=-1 0 1"` //赞成_1 或 反对_-1 取消_0
		//oneof：限定值只能在空格分隔的几个值中间，可以是很多类型
}


// ParamPostList 获取帖子列表query string参数
type ParamPostList struct {
	CommunityID int64 `json:"community_id" form:"community_id"` //可以为空
	Page int64 `json:"page" form:"page"` //页数
	Size int64 `json:"size" form:"size"` //一页几个
	Order string `json:"order" form:"order" example:"score"` //按时间还是按分数
}



