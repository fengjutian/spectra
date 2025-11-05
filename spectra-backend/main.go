package main

import (
	"spectra-backend/middleware"

	"github.com/gin-gonic/gin"
)

func main() {

	// logger, _ := zap.NewProduction()
	// defer logger.Sync()

	// logger.Info("Server started",
	// 	zap.String("env", "production"),
	// 	zap.Int("port", 8080),
	// )

	router := gin.Default()

	logger := middleware.InitLogger()
	defer logger.Sync()

	router.Use(middleware.GinLogger(logger))

	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		logger.Info("Ping API called")
		c.HTML(200, "index.html", nil)
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.Run()
}
