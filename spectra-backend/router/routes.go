package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"spectra-backend/config"
	"spectra-backend/handlers"
	"spectra-backend/repository"
	"spectra-backend/services"
)

func SetupRoutes(router *gin.Engine, cfg *config.Config, logger *zap.Logger) {
	// 首页和健康检查路由
	HomeRoutes(router, logger)

	// 初始化数据库连接
	repo, err := repository.NewClickHouseRepository(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to initialize repository", zap.Error(err))
	}

	// 初始化服务
	logService := services.NewLogService(repo)

	// 初始化处理器
	logHandler := handlers.NewLogHandler(logService, logger)

	// API 路由组
	api := router.Group("/api")
	{
		// 错误日志相关路由
		api.POST("/error-logs", logHandler.RecordErrorLog)
		api.GET("/error-logs", logHandler.GetErrorLogs)

		// 性能指标相关路由
		api.POST("/performance-metrics", logHandler.RecordPerformanceMetric)
		api.GET("/performance-metrics", logHandler.GetPerformanceMetrics)

		// 用户行为相关路由
		api.POST("/user-actions", logHandler.RecordUserAction)
		api.GET("/user-actions", logHandler.GetUserActions)

		// 自定义事件相关路由
		api.POST("/custom-events", logHandler.RecordCustomEvent)
		api.GET("/custom-events", logHandler.GetCustomEvents)

		// 页面停留时长相关路由
		api.POST("/page-stays", logHandler.RecordPageStay)
		api.GET("/page-stays/average", logHandler.GetAveragePageStay)
	}
}

func HomeRoutes(router *gin.Engine, logger *zap.Logger) {
	router.GET("/", func(c *gin.Context) {
		logger.Info("Homepage accessed")
		c.HTML(200, "index.html", nil)
	})

	router.GET("/ping", func(c *gin.Context) {
		logger.Info("Ping API called")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
