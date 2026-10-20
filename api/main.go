package main

import (
	"goVServer/api/defs"

	"github.com/gin-gonic/gin"
)

func validateUserSessionMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		if !validateUserSession(c.Request) {
			// 使用错误响应函数
			sendErrorResponse(c, defs.ErrorNotAuthUser)
			return
		}

		c.Next()
	}
}
func RegisterRoutes() *gin.Engine {
	router := gin.Default()
	//公开的
	public := router.Group("/api")
	{
		public.POST("/user", CreateUser)
	}

	//需要验证的
	protected := router.Group("/protected")
	{
		protected.Use(validateUserSessionMiddleware())
		protected.POST("/user/:user_name", Login)
	}
	return router
}

func main() {
	server := gin.Default()
	//创建用户
	server.POST("/user", CreateUser)
	//登录
	server.POST("/user/:user_name", Login)

	server.Run(":8089")
}
