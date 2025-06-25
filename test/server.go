package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 支持所有常用方法
	r.Any("/test", func(c *gin.Context) {
		var jsonData map[string]interface{}
		c.BindJSON(&jsonData)
		fmt.Printf("收到请求: %s %v\n", c.Request.Method, jsonData)
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"method":  c.Request.Method,
			"data":    jsonData,
			"message": "收到请求",
		})
	})

	// 备用接口
	r.Any("/backup/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "backup success",
			"message": "这是备用接口",
		})
	})

	r.Run(":5000") // 监听5000端口
}
