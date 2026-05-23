package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mall/app/Support/logger"
	"mall/bootstrap"
	"mall/config"
	"mall/routes"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 1. 启动应用 (加载配置、日志、数据库、Redis)
	bootstrap.Boot(".env")
	defer logger.Sync()

	// 2. 设置 Gin 模式
	if !config.App().Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// 3. 创建路由引擎
	router := gin.New()
	routes.Register(router)
	routes.RegisterAdmin(router)

	// 4. 启动 HTTP 服务 (支持优雅关闭)
	addr := ":" + config.App().Port
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		logger.Info("HTTP 服务启动", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP 服务启动失败: %v", err)
		}
	}()

	// 5. 等待中断信号, 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭服务...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("服务关闭异常", zap.Error(err))
	}

	logger.Info("服务已安全退出")
}
