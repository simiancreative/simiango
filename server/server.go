package server

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/simiancreative/simiango/logger"
	"github.com/simiancreative/simiango/service"
)

var router *gin.Engine
var services service.Collection

func init() {
	gin.SetMode(gin.ReleaseMode)
	router = gin.New()
	router.Use(
		TimingMiddleware(),
		JSONLogMiddleware(),
		gin.Recovery(),
	)
}

func Start() {
	for _, service := range services {
		InitService(service)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Debug(
		"Server: initialized and starting",
		logger.Fields{"port": port},
	)

	router.Run()
}
