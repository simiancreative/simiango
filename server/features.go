package server

import (
	"fmt"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/mandrigin/gin-spa/spa"

	"github.com/simiancreative/simiango/logger"
)

// SetSPA hosts a Single Page App directory at the specified url
func SetSPA(url string, dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		panic(fmt.Sprintf("spa dir does not exist: %v", dir))
	}

	logger.Debug("Server: adding spa route", logger.Fields{
		"path": url,
		"dir":  dir,
	})

	router.Use(spa.Middleware(url, dir))
}

// SetCORS applies CORS config to the gin router
func SetCORS() {
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}

// AddSentry adds the Sentry middleware to the gin router
func AddSentry(scopeFunc func(*gin.Context, *sentry.Scope)) {
	router.Use(sentrygin.New(sentrygin.Options{
		Repanic:         true,
		WaitForDelivery: false,
	}))

	router.Use(sentrySetup(scopeFunc))
}

func sentrySetup(scopeFunc func(*gin.Context, *sentry.Scope)) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("sentryScopeFunc", scopeFunc)
		c.Next()
	}
}

func healthGET(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "UP",
	})
}

func EnableHealthCheck() {
	router.Handle("GET", "/status", healthGET)
}

func EnableCustomHealthCheck(path string, check gin.HandlerFunc) {
	router.Handle("GET", path, check)
}

func AddPprof() {
	pprof.Register(router)
}
