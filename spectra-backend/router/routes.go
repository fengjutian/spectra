package router

import (
	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

func HomeRoutes(router *gin.Engine, logger *zap.Logger) {

	router.GET("/", func(c *gin.Context) {
		logger.Info("Ping API called")
		c.HTML(200, "index.html", nil)
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
