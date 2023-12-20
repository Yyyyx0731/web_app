package router

import (
	"github.com/gin-gonic/gin"
	"time"

	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"net/http"
	"web_app/controller"
	"web_app/logger"
	"web_app/middlewares"

	_ "web_app/docs" // 千万不要忘了导入把你上一步生成的docs
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) //gin设置成发布模式
	}


	r := gin.New()
																		//每两秒钟添加一个
	r.Use(logger.GinLogger(), logger.GinRecovery(true),middlewares.RateLimitMiddleware(2*time.Second,1))

	//进行默认初始化和注册对应的路由 注册一个针对 swagger 的路由
	r.GET("/swagger/*any",ginSwagger.WrapHandler(swaggerFiles.Handler))

	//路由的分发

	v1 := r.Group("/api/v1")

	//注册业务路由
	v1.POST("/signup",controller.SignUpHandler)
	//登录业务路路由
	v1.POST("/login",controller.LoginHandler)

	v1.Use(middlewares.JWTAuthMiddleware() )//把token认证全放在 JWTAuthMiddleware() 中间件里

	{
		v1.GET("/community",controller.CommunityHandler)
		v1.GET("/community/:id",controller.CommunityDetailHandler)

		v1.POST("/post",controller.CreatePostHandler)
		v1.GET("/post/:id",controller.GetPostDetailHandler)
		v1.GET("/posts",controller.GetPostListHandler)
		//根据时间或分数获取帖子列表
		v1.GET("/posts2",controller.GetPostListHandler2)

		//投票
		v1.POST("/vote",controller.PostVoteController)
	}



	//设定请求url不存在的返回值
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK,gin.H{
			"msg":"404",
		})
	})
	//r.Run(":8080")
	return r
}


