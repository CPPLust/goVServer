package main

import (
	"goVServer/api/defs"

	"github.com/gin-gonic/gin"
)

func sendErrorResponse(c *gin.Context, errResp defs.ErrorResponse) {
	c.JSON(errResp.HttpSC, errResp.Error)
	c.Abort() // 阻止后续处理
}

// 发送正常响应
func sendNormalResponse(c *gin.Context, resp string, sc int) {
	c.String(sc, resp)
}
