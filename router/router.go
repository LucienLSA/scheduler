package router

import (
	"net/http"
	"scheduler/api"
	"scheduler/config"
	"scheduler/logs"
	"scheduler/middlewares"

	"github.com/gin-gonic/gin"
)

func Init(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		// gin设置发布模式
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	// 处理异常
	r.NoMethod(HandleNotFound)
	r.NoRoute(HandleNotFound)
	// 使用自定义的ginlogger中间件
	sConf := config.Conf
	r.Use(logs.GinLogger(), logs.GinRecovery(true))
	r.Use(middlewares.LogMiddleware())
	rG := r.Group(sConf.ServerConfig.ContextPath)
	{
		rG.GET("/tasks", api.ListTask)
		rG.POST("/task", api.AddTask)
	}
	// 加载静态文件和html
	if sConf.ServerConfig.ConsoleEnable {
		rG.Static("/web", "web")
	}
	return r
}

// 未找到资源
func HandleNotFound(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "404",
	})
}
