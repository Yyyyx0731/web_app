package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"web_app/logic"
	"web_app/models"
)



func CreatePostHandler(c *gin.Context){
	//1.获取参数即参数校验

	p := new(models.Post)
	if err := c.ShouldBindJSON(p);err != nil {
		zap.L().Debug("c.ShouldBindJSON(p) error",zap.Any("err",err))
		zap.L().Error("create post with invalid param")
		ResponseError(c,CodeInvalidParam)
		return
	}
	//从 c 获取当前发送请求的用户id
	userID,err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c,CodeNeedLogin) //获取不出来，让他重新登录
		return
	}
	p.AuthorID = userID
	//2.创建帖子
	if err:=logic.CreatePost(p);err != nil {
		zap.L().Error("logic.CreatePost(p) failed",zap.Error(err))
		ResponseError(c,CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c,nil)
}

// GetPostDetailHandler 获取帖子详情
func GetPostDetailHandler(c *gin.Context){
	//1.获取参数 从url帖子id
	pidStr := c.Param("id")
	pid,err := strconv.ParseInt(pidStr,10,64)
	if err != nil {
		zap.L().Error("get post detail with invalid param",zap.Error(err))
		ResponseError(c,CodeInvalidParam)
		return
	}
	//2.根据id获取帖子数据
	data,err := logic.GetPostByID(pid)
	if err!= nil {
		zap.L().Error("logic.GetPostByID(pid) failed",zap.Error(err))
		ResponseError(c,CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c,data)
}

// GetPostListHandler 获取帖子列表
func GetPostListHandler(c *gin.Context){
	//获取分页参数函数
	page,size := getPageInfo(c)
	//获取数据
	data,err := logic.GetPostList(page,size)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed",zap.Error(err))
		ResponseError(c,CodeServerBusy)
		return
	}
	//返回响应
	ResponseSuccess(c,data)
}

// GetPostListHandler2 升级版帖子列表接口
// @Summary 升级版帖子列表接口
// @Description 根据前端传来的参数动态获取帖子列表  按时间或按分数
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer JWT"
// @Param object query models.ParamPostList false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /posts2 [get]
func GetPostListHandler2(c *gin.Context){
	//1.获取分数
	//2.去redis查id列表
	//3.根据id去数据库查帖子详情

	//GET请求参数：/api/v1/posts2?page=1&size=10&order=time
	//初始化结构体时，指定初始参数
	p := &models.ParamPostList{
		Page:  models.DefaultPage,
		Size:  models.DefaultSize,
		Order: models.OrderTime, //默认按时间排序
	}
	//获取分页参数函数
	//c.ShouldBindJSON() 如果请求中携带的是json格式的数据，才能用这个方法获取到数据
	//c.ShouldBind() 根据请求类型，字段选择相应的方法获取数据
	//此处带的是query string参数
	if err := c.ShouldBindQuery(p);err != nil {
		zap.L().Error("GetPostListHandler2 with invalid params",zap.Error(err))
		ResponseError(c,CodeInvalidParam)
		return
	}

	data,err := logic.GetPostListNew(p) //更新：合二为一的查询帖子业务逻辑层
	if err != nil {
		zap.L().Error("logic.GetPostList2() failed",zap.Error(err))
		ResponseError(c,CodeServerBusy)
		return
	}
	//返回响应
	ResponseSuccess(c,data)
}


//// GetCommunityPostListHandler 根据社区查帖子列表
//func GetCommunityPostListHandler(c *gin.Context){
//	//初始化结构体时，指定初始参数
//	p := &models.ParamCommunityPostList{
//		ParamPostList:&models.ParamPostList{
//			Page:  models.DefaultPage,
//			Size:  models.DefaultSize,
//			Order: models.OrderTime, //默认按时间排序
//		},
//	}
//	//获取分页参数函数
//	//c.ShouldBindJSON() 如果请求中携带的是json格式的数据，才能用这个方法获取到数据
//	//c.ShouldBind() 根据请求类型，字段选择相应的方法获取数据
//	//此处带的是query string参数
//	if err := c.ShouldBindQuery(p);err != nil {
//		zap.L().Error("GetPostListHandler2 with invalid params",zap.Error(err))
//		ResponseError(c,CodeInvalidParam)
//		return
//	}
//	//获取数据
//	data,err := logic.GetCommunityPostList2(p)
//	if err != nil {
//		zap.L().Error("logic.GetPostList2() failed",zap.Error(err))
//		ResponseError(c,CodeServerBusy)
//		return
//	}
//	//返回响应
//	ResponseSuccess(c,data)
//}