package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	c.String(http.StatusOK, "用户创建成功")



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
