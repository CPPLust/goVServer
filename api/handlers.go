package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"goVServer/api/dbops"
	"goVServer/api/defs"
	"goVServer/api/session"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	fmt.Println("CreateUser")
	res, _ := io.ReadAll(c.Request.Body)
	ubody := &defs.UserCredential{}

	fmt.Println(string(res))
	if err := json.Unmarshal(res, ubody); err != nil {
		sendErrorResponse(c, defs.ErrorRequestBodyParseFailed)
		return
	}
	fmt.Printf("recv name: %v, recv pwd:%v", ubody.Username, ubody.Pwd)

	if err := dbops.AddUserCredential(ubody.Username, ubody.Pwd); err != nil {
		sendErrorResponse(c, defs.ErrorDBError)
		return
	}

	id := session.GenerateNewSessionId(ubody.Username)
	su := &defs.SignedUp{Success: true, Session: id}

	if resp, err := json.Marshal(su); err != nil {
		sendErrorResponse(c, defs.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(c, string(resp), http.StatusCreated)
	}
}

func Login(c *gin.Context) {
	// 获取 URL 参数
	userName := c.Param("user_name")

	// 这里可以添加登录验证逻辑
	// 例如：验证密码、生成 token 等

	// 返回成功响应
	c.JSON(200, gin.H{
		"code":      200,
		"message":   "登录成功",
		"user_name": userName,
		"token":     "mock_token_12345", // 实际项目中应该生成真实的 JWT token
	})
}
