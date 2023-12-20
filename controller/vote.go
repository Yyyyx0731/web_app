package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"web_app/logic"
	"web_app/models"
)



func PostVoteController(c *gin.Context){
	//参数校验
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p);err != nil{
		errs,ok := err.(validator.ValidationErrors) //类型断言
		if !ok {
			ResponseError(c,CodeInvalidParam)
			return
		}
		// removeTopStruct 可去掉错误提示结构体名称前缀
		errData := removeTopStruct(errs.Translate(trans))
		ResponseErrorWithMsg(c,CodeInvalidParam,errData)
		return
	}
	fmt.Println("dir:",p.Direction)
	//获取当前请求的用户id
	userID,err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c,CodeNeedLogin)
		return
	}
	//投票业务逻辑
	if err := logic.VoteForPost(userID,p);err != nil {
		zap.L().Error("logic.VoteForPost(userID,p) failed",zap.Error(err))
		ResponseError(c,CodeServerBusy)
		return
	}
	fmt.Println(p.Direction)
	ResponseSuccess(c,nil)
}
