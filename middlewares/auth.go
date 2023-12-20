package middlewares

import (
	"github.com/gin-gonic/gin"
	"strings"
	"web_app/controller"
	"web_app/pkg/jwt"
)

// JWTAuthMiddleware 基于JWT的认证中间件（检验token的中间件）
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI

		// 这里的具体实现方式要依据你的实际业务情况决定
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// Authorization:Bearer xxxxxxxx.xxx.xxx
		//例：判断请求头里是不是带Authorization这个认证的token
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" { //如果没有带
			//请求头里没有 Authorization 说明没有登录
			controller.ResponseError(c,controller.CodeNeedLogin)
			//c.JSON(http.StatusOK, gin.H{
			//	"code": 2003,
			//	"msg":  "请求头中auth为空",
			//})
			c.Abort() //直接返回响应
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") { //如果格式不对也返回 属于无效token
			controller.ResponseError(c,controller.CodeInvalidToken)
			//c.JSON(http.StatusOK, gin.H{
			//	"code": 2004,
			//	"msg":  "请求头中auth格式有误",
			//})
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			controller.ResponseError(c,controller.CodeInvalidToken)
			//c.JSON(http.StatusOK, gin.H{
			//	"code": 2005,
			//	"msg":  "无效的Token",
			//})
			c.Abort()
			return
		}
		//如果走到这里，说明当前用户带了有效且解析成功的token
		// 将当前请求的userID信息保存到请求的上下文c上
		c.Set(controller.CtxUserIDKey, mc.UserID)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}

