package controller

// 请求参数的校验	路由的转发

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"web_app/dao/mysql"
	"web_app/logic"
	"web_app/models"
)

// SignUpHandler 处理注册请求的函数
func SignUpHandler(c *gin.Context){
	//1.获取参数和参数校验
	//让gin框架能够自动从请求中把想要的数据绑定到一个结构体里
	//var p models.ParamSignUp //用此方法传&p
	p := new(models.ParamSignUp)//返回的直接是指针，下面的ShouldBindJSON()和logic.SignUp()可以直接传p
	//ShouldBindJSON只能识别简单的 例如是不是json格式，字段类型，所以还需要详细判断
	if err:=c.ShouldBindJSON(p);err!=nil{
		//请求参数有误，直接返回响应，并return
		zap.L().Error("SignUp with invalid param",zap.Error(err))
		//判断err是不是validator.ValidationErrors类型（有可能在上面json序列化的时候就错了）
		errs,ok := err.(validator.ValidationErrors)
		if !ok { //如果不是validator.ValidationErrors类型,则是在序列化时就有错误
			//c.JSON(http.StatusOK,gin.H{
			//	"msg":err.Error(),
			//})
			ResponseError(c,CodeInvalidParam)
			return
		}
		//c.JSON(http.StatusOK,gin.H{//是validator.ValidationErrors类型
		//	"msg":removeTopStruct(errs.Translate(trans)),//翻译
		//})
		ResponseErrorWithMsg(c,CodeInvalidParam,removeTopStruct(errs.Translate(trans)))
		return
	}
	fmt.Println(p)
	////参数详细校验（用结构体tag进行参数校验，就不需要在此单独写判断了）
	//if len(p.Username)==0||len(p.Password)==0||len(p.RePassword)==0||p.RePassword!=p.Password {
	//	//请求参数有误，直接返回响应，并return
	//	zap.L().Error("SignUp with invalid param")
	//	c.JSON(http.StatusOK,gin.H{
	//		"msg":"请求参数有误",
	//	})
	//	return
	//}

	//2.业务层 业务处理
	if err := logic.SignUp(p);err!=nil {
		zap.L().Error("logic.SignUp failed",zap.Error(err))
		if errors.Is(err,mysql.ErrorUserExist){ //判断是不是错误”用户已存在“
			ResponseError(c,CodeUserExist) //如果是
			return
		}
		//c.JSON(http.StatusOK,gin.H{
		//	"msg":"注册失败",
		//})
		ResponseError(c,CodeServerBusy)
		return
	}

	//3.返回响应
	//c.JSON(http.StatusOK,gin.H{
	//	"msg":"success",
	//})
	ResponseSuccess(c,nil)

}


func LoginHandler(c *gin.Context){
	//1.获取请求参数即参数校验
	p := new(models.ParamLogin)
	if err:= c.ShouldBindJSON(p);err!=nil {
		//请求参数有误，直接返回响应，并return
		zap.L().Error("Login with invalid param",zap.Error(err))
		//判断err是不是validator.ValidationErrors类型（有可能在上面json序列化的时候就错了）
		errs,ok := err.(validator.ValidationErrors)
		if !ok { //如果不是validator.ValidationErrors类型,则是在序列化时就有错误
			//c.JSON(http.StatusOK,gin.H{
			//	"msg":err.Error(),
			//})
			ResponseError(c,CodeInvalidParam)
			return
		}
		//c.JSON(http.StatusOK,gin.H{//是validator.ValidationErrors类型
		//	"msg":removeTopStruct(errs.Translate(trans)),//翻译
		//})
		ResponseErrorWithMsg(c,CodeInvalidParam,errs.Translate(trans))
		return
	}

	//2.业务逻辑处理
	//获取生成的token
	user,err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login failed",zap.String("username",p.Username),zap.Error(err))
		//c.JSON(http.StatusOK,gin.H{
		//	"msg":"用户名或密码错误",
		//})
		if errors.Is(err,mysql.ErrorUserNotExist){
			ResponseError(c,CodeUserNotExist)
			return
		}
		ResponseError(c,CodeInvalidPassword)
		return
	}

	//3.返回响应
	//c.JSON(http.StatusOK,gin.H{f
	//	"msg":"登录成功",
	//})
	//ResponseSuccess(c,token)
	ResponseSuccess(c,gin.H{
		"user_id":fmt.Sprintf("%d",user.UserID),
		"user_name":user.Username,
		"token":user.Token,
	})
}