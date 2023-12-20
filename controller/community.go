package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"web_app/logic"
)

func CommunityHandler(c *gin.Context){
	//查询到所有社区（community_id,community_name)  以列表的形式返回（切片）
	data,err := logic.GetCommunityList()
	if err!=nil {
		zap.L().Error("logic.GetCommunityList() failed",zap.Error(err))
		ResponseError(c,CodeServerBusy) //不轻易把服务端报错暴露给外面，详细错误记录在日志
		return
	}
	ResponseSuccess(c,data)
}


// CommunityDetailHandler 社区分类详情
func CommunityDetailHandler(c *gin.Context){
	//获取社区id
	communityID := c.Param("id") //key要和controller里的路径参数名对应
	id,err := strconv.ParseInt(communityID,10,64)
	if err != nil {
		ResponseError(c,CodeInvalidParam) //参数错误
		return
	}
	//根据id获取社区详情
	data,err := logic.GetCommunityDetail(id)
	if err!=nil {
		zap.L().Error("logic.GetCommunityList() failed",zap.Error(err))
		ResponseError(c,CodeServerBusy) //不轻易把服务端报错暴露给外面，详细错误记录在日志
		return
	}
	ResponseSuccess(c,data)
}



