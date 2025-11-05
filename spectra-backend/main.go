package main

import (
	"fmt"
	"log"
	"spectra-backend/config"
	"spectra-backend/middleware"
	"spectra-backend/router"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {

	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	logger := middleware.InitLogger()
	defer logger.Sync()

	r := gin.Default()

	r.Use(middleware.GinLogger(logger))

	// 静态文件和模板
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	router.HomeRoutes(r, logger)

	// 启动服务器
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("Starting server", zap.String("address", serverAddr))
	if err := r.Run(serverAddr); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
