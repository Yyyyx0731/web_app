package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)


// 在多个模块可能使用的变量，将其定义为常量 或者写成可配置项
const CtxUserIDKey = "userID" //ctx是context缩写


var ErrorUserNotLogin = errors.New("用户未登录")

// getCurrentUserID 获取当前登录的用户的id
func getCurrentUserID(c *gin.Context) (userID int64,err error) {
	uid,ok := c.Get(CtxUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID ,ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}


func getPageInfo(c *gin.Context)(int64,int64){
	//获取分页参数
	pageStr := c.Query("page")
	sizeStr := c.Query("size")
	var(
		size int64
		page int64
		err error
	)
	page,err = strconv.ParseInt(pageStr,10,64)
	if err != nil {
		page = 1
	}
	size,err = strconv.ParseInt(sizeStr,10,64)
	if err != nil {
		size = 10
	}
	return page,size
}
