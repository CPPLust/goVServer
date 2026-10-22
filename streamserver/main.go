package main

import (
	"fmt"
	"goVServer/api/defs"

	"github.com/gin-gonic/gin"
)

func connListMiddleware(cl *ConnLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !cl.GetConn() {
			sendErrorResponse(c, defs.ErrorServerReachLimit)
			return
		}
		defer cl.ReleaseConn()
		c.Next()
	}
}

func RegisterHandlers() *gin.Engine {
	fmt.Println("RegisterHandlers")
	router := gin.Default()

	protected := router.Group("/")
	cl := NewConnLimiter(2)
	protected.Use(connListMiddleware(cl))
	protected.GET("/videos/:vid-id", streamHanler)
	protected.POST("/upload/:vid-id", uploadHanler)
	protected.GET("/testpage", testPageHandler)

	return router
}

func main() {
	server := RegisterHandlers()
	server.Run(":9000")

}
