package api

import (
	"github.com/gin-gonic/gin"
)

// NewRouter 创建并配置 Gin 路由（管理后台 API）
func NewRouter() *gin.Engine {
	r := gin.Default()

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1
	v1 := r.Group("/api/v1")
	{
		// 账号管理
		v1.GET("/accounts", listAccounts)       // 账号列表
		v1.POST("/accounts", createAccount)     // 创建账号
		v1.PUT("/accounts/:id", updateAccount)  // 更新账号
		v1.DELETE("/accounts/:id", deleteAccount) // 删除账号

		// 角色管理
		v1.GET("/characters", listCharacters)
		v1.PUT("/characters/:id", updateCharacter)

		// 服务器状态
		v1.GET("/server/info", serverInfo)
	}

	return r
}

// TODO: 实现以下 API handler

func listAccounts(c *gin.Context) {
	c.JSON(200, gin.H{"accounts": []interface{}{}})
}

func createAccount(c *gin.Context) {
	c.JSON(201, gin.H{"message": "not implemented"})
}

func updateAccount(c *gin.Context) {
	c.JSON(200, gin.H{"message": "not implemented"})
}

func deleteAccount(c *gin.Context) {
	c.JSON(200, gin.H{"message": "not implemented"})
}

func listCharacters(c *gin.Context) {
	c.JSON(200, gin.H{"characters": []interface{}{}})
}

func updateCharacter(c *gin.Context) {
	c.JSON(200, gin.H{"message": "not implemented"})
}

func serverInfo(c *gin.Context) {
	c.JSON(200, gin.H{
		"name":       "BeiDou-Go",
		"version":    "GMS v0.83",
		"status":     "running",
		"online":     0,
	})
}
