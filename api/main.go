package main

import (
	"fmt"
	"goVServer/api/defs"

	"github.com/gin-gonic/gin"
)

func validateUserSessionMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		if !validateUserSession(c.Request) {
			// 使用错误响应函数
			fmt.Println("无效的认证")
			sendErrorResponse(c, defs.ErrorNotAuthUser)
			return
		}

		c.Next()
	}
}
func RegisterRoutes() *gin.Engine {
	router := gin.Default()

	// 公开路由
	router.POST("/user", CreateUser)
	router.POST("/user/:username", Login) // 登录无需会话验证

	router.GET("/user/:username", GetUserInfo)

	// 需要验证的路由
	protected := router.Group("/")
	protected.Use(validateUserSessionMiddleware())
	{
		protected.POST("/user/:username/videos", AddNewVideo)
		protected.GET("/user/:username/videos", ListAllVideos)
		//protected.DELETE("/user/:username/videos/:vid-id", DeleteVideo)
	}

	return router
}

func main() {
	server := RegisterRoutes()
	server.Run(":8000")
}
