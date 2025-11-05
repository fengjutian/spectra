package main

import (
	"spectra-backend/middleware"
	"spectra-backend/router"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	logger := middleware.InitLogger()
	defer logger.Sync()

	r.Use(middleware.GinLogger(logger))

	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	router.HomeRoutes(r, logger)

	r.Run()
}
