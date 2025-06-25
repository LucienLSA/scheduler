package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"scheduler/config"
	"scheduler/db/mysql"
	"scheduler/logs"
	"scheduler/router"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	// 1. 加载配置
	if err := config.Init(); err != nil {
		fmt.Printf("init settings failed, err:%v\n", err)
		return
	}
	fmt.Println("init config success!")
	sConf := config.Conf

	// 2. 初始化日志
	if err := logs.InitLogger(sConf.LogConfig, sConf.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	zap.L().Info("init logger success")
	// 延迟将缓存区的日志追加
	defer zap.L().Sync()

	// 3. 初始化MySQL
	if err := mysql.Init(sConf.DBConfig); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	zap.L().Info("init mysql success")
	zap.L().Info(fmt.Sprintf("http://localhost:%d%s/web/index.html", sConf.ServerConfig.Port, sConf.ServerConfig.ContextPath))

	r := router.Init(sConf.Mode)
	host := fmt.Sprintf("%s:%d", sConf.ServerConfig.Addr, sConf.ServerConfig.Port)
	zap.L().Info("Starting server on " + host)

	srv := &http.Server{
		Addr:    host,
		Handler: r,
	}
	go func() {
		// 开启一个goroutine启动服务
		zap.L().Info("Server is starting...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("Server failed to start", zap.Error(err))
		}
	}()
	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}
	zap.L().Info("Server exiting")
}
