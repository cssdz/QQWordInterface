package main

import (
	"context"
	"fmt"
	"go_web/dao/mysql"
	"go_web/dao/redis"
	"go_web/logger"
	"go_web/logic"
	"go_web/routes"
	"go_web/settings"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/spf13/viper"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) //内置log
	// 1. 加载配置文件
	if err := settings.Init(); err != nil {
		log.Panic("init settings failed, err:", err.Error())
	}
	log.Printf("[1] %s V%s Initialize... (%s)\n",
		viper.GetString("app.name"),
		viper.GetString("app.version"),
		viper.GetString("app.mode"),
	)

	// 2. 初始化日志
	if err := logger.Init(); err != nil {
		log.Panic("init logger failed, err:", err.Error())
	}
	defer zap.L().Sync()
	zap.L().Debug("Server starting...")
	log.Printf("[2] Logger...\n")

	// 3. 初始化MySql
	if err := mysql.Init(); err != nil {
		log.Panic("init mysql failed, err:", err.Error())
	}
	defer mysql.Close()
	log.Printf("[3] MySQL...DB>%s \n", viper.GetString("mysql.db"))

	// 4. 初始化Redis
	if err := redis.Init(); err != nil {
		log.Panic("init redis failed, err:%", err.Error())
	}
	defer redis.Close()
	log.Printf("[4] Redis...DB>%d \n", viper.GetInt("redis.db"))

	// 5. 注册路由
	r := routes.Setup()
	log.Printf("[5]routes.Setup() \n")

	// 6. 启动服务
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}
	log.Printf("[6] HTTP server start...%d \n", viper.GetInt("app.port"))
	zap.L().Info("HTTP server start...")

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
			zap.L().Error("Http ListenAndServe Failed", zap.Error(err))
		}
	}()

	// 7. 设置定时器
	go logic.ClearAchievement()

	// 8. 优雅关机
	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	log.Println("[7] Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
		zap.L().Error("Http srv.Shutdown", zap.Error(err))
	}
	log.Println("Server exited...")
	zap.L().Debug("Server exited...")
}
