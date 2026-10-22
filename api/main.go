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

	// 公开路由
	router.POST("/user", CreateUser)

	// 需要验证的路由
	protected := router.Group("/")
	protected.Use(validateUserSessionMiddleware())
	{
		protected.POST("/user/:user_name", Login)
	}

	return router
}

func main() {
	server := RegisterRoutes()
	server.Run(":8089")
}
