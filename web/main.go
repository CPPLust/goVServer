package main

import (
	"github.com/gin-gonic/gin"
)

func RegisterHandler() *gin.Engine {
	router := gin.Default()

	router.Static("/statics", "./templates")

	router.GET("/", homeHandler)
	router.POST("/", homeHandler)

	router.GET("/userhome", userHomeHandler)
	router.POST("/userhome", userHomeHandler)

	router.POST("/api", apiHomeHandler)

	return router
}

func main() {
	server := RegisterHandler()
	server.Run(":8080")
}
