package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
{
	"code":10001, //程序中的错误码
	"msg":xx, //提示信息
	"data": {} //数据
}
*/

type ResponseData struct {
	Code ResCode `json:"code"`
	Msg interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"` //当字段是空时，就不在json序列化时展示整个字段
}

func ResponseError(c *gin.Context,code ResCode){
	c.  JSON(http.StatusOK,&ResponseData{
		Code: code,
		Msg: code.Msg(),
		Data: nil,
	})
}

// ResponseErrorWithMsg 自定义错误
func ResponseErrorWithMsg(c *gin.Context,code ResCode,msg interface{}){
	c.  JSON(http.StatusOK,&ResponseData{
		Code: code,
		Msg: msg,
		Data: nil,
	})
}

func ResponseSuccess(c *gin.Context,data interface{}){
	c.JSON(http.StatusOK,&ResponseData{
		Code: CodeSuccess,
		Msg: CodeSuccess.Msg(),
		Data: data,
	})
}
