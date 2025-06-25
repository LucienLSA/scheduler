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
		rG.GET("/task/:id", api.GetTask)
		rG.POST("/task", api.AddTask)
		rG.PUT("/task/:id", api.EditTask)
		rG.DELETE("/task/:id", api.DeleteTask)

		rG.GET("/tags", api.ListTag)
		rG.GET("/specs", api.ListSpec)

		rG.GET("/execute/:id", api.ExecuteTask)

		rG.GET("/records", api.ListRecord)

		rG.GET("/health", api.Health)
		rG.GET("/shutdown", api.Shutdown)
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
