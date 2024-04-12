package server

import (
	"os"

	"github.com/gin-gonic/gin"
	servertiming "github.com/p768lwy3/gin-server-timing"

	"github.com/simiancreative/simiango/logger"
	"github.com/simiancreative/simiango/service"
)

var (
	router   *gin.Engine
	services service.Collection
)

func New() {
	gin.SetMode(gin.ReleaseMode)
	router = gin.New()
	router.Use(
		servertiming.Middleware(),
		JSONLogMiddleware,
		Recovery,
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

	if err := router.Run(); err != nil {
		logger.Fatal(
			"Server: failed to start",
			logger.Fields{"err": err},
		)
	}
}

func GetRouter() *gin.Engine {
	return router
}
