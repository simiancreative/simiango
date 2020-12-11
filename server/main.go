package server

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mandrigin/gin-spa/spa"

	"simian/context"
	"simian/logger"
	"simian/service"
)

var router *gin.Engine
var services service.Collection

func init() {
	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
}

func AddService(addedservice service.Config) {
	services = append(services, addedservice)
}

func AddServices(addedservices []service.Config) {
	services = append(services, addedservices...)
}

func InitService(config service.Config) {
	logger.Debug("Server: adding route", logger.Fields{
		"method": config.Method,
		"path":   config.Path,
	})

	router.Handle(config.Method, config.Path, handleService(config))
}

func handleService(config service.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		parsedBody := rawBody(c.Request.Body)
		parsedParams := parseParams(c.Params, c.Request.URL)

		s, err := config.Build(context.Id(), parsedBody, parsedParams)
		if err != nil {
			buildErr := service.ToResultError(err, "service failed to build", 500)
			c.JSON(buildErr.GetStatus(), buildErr.GetDetails())
			return
		}

		result, err := s.Result()

		if err == nil {
			c.JSON(http.StatusOK, result)
			return
		}

		resultErr, ok := err.(*service.ResultError)
		if !ok {
			resultErr = service.ToResultError(err, "service failed", 500)
			c.JSON(resultErr.GetStatus(), resultErr.GetDetails())
			return
		}

		c.JSON(resultErr.GetStatus(), resultErr.GetDetails())
	}
}

func parseParams(params gin.Params, url *url.URL) service.RawParams {
	var parsedParams = service.RawParams{}

	for _, element := range params {
		parsedParams = append(parsedParams, service.RawParam{
			Key:   element.Key,
			Value: element.Value,
		})
	}

	for k, v := range url.Query() {
		parsedParams = append(parsedParams, service.RawParam{
			Key:   k,
			Value: v,
		})
	}

	return parsedParams
}

func rawBody(source io.ReadCloser) []byte {
	buf := make([]byte, 1024)
	num, _ := source.Read(buf)
	reqBody := string(buf[0:num])
	return []byte(reqBody)
}

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

func Start() {
	for _, service := range services {
		InitService(service)
	}

	router.Run()
}
