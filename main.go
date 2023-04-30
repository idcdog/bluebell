package main

import (
	"bluebell/controller"
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/logger"
	"bluebell/pkg/snowflake"
	"bluebell/routes"
	"bluebell/settings"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// @title bluebell gin框架demo项目
// @version 1.0
// @description gin框架学习教程专用的demo项目
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	// 1. 加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("init config failed, err: %s\n", err.Error())
		return
	}
	// 2. 初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("init logger failed, err: %s\n", err.Error())
		return
	}
	zap.L().Debug("logger init success...")
	defer zap.L().Sync()
	// 3. 初始化MySQL连接
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, err: %s\n", err.Error())
		return
	}
	zap.L().Debug("mysql init success...")
	defer mysql.Close()
	// 4. 初始化Redis连接
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed, err: %s\n", err.Error())
		return
	}
	defer redis.Close()
	zap.L().Debug("redis init success...")
	// 初始化序列号生成器
	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		zap.L().Fatal("init snowflake failed", zap.Error(err))
	}
	// 初始化校验器
	if err := controller.InitTrans("zh"); err != nil {
		zap.L().Fatal("init validtor trans failed", zap.Error(err))
	}
	// 5. 注册路由
	r := routes.SetupRouter(settings.Conf.Mode)
	// 6. 优雅关机
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.Port),
		Handler: r,
	}
	go func() {
		// 开启一个goroutine处理请求
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen failed", zap.Error(err))
		}
	}()
	// 创建一个接收信号的通道
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务， 超过5秒则超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown error: ", zap.Error(err))
	}
	zap.L().Info("Server exited!")

}
