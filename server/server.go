package server

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	servertiming "github.com/p768lwy3/gin-server-timing"

	"github.com/simiancreative/simiango/logger"
	"github.com/simiancreative/simiango/meta"
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
		meta.GinRecovery(recoveryHandler),
	)
}

func recoveryHandler(c *gin.Context, context map[string]interface{}) {
	c.JSON(http.StatusInternalServerError, service.ResultError{
		Status:  http.StatusInternalServerError,
		Message: "Internal Server Error",
		Reasons: []map[string]interface{}{context},
	})

	c.AbortWithStatus(http.StatusInternalServerError)
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
