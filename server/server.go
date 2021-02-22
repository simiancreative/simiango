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
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	"github.com/simiancreative/simiango/logger"
	"github.com/simiancreative/simiango/meta"
	"github.com/simiancreative/simiango/service"
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
	requestID := meta.Id()

	return func(c *gin.Context) {
		c.Header("X-Request-ID", string(requestID))

		logger.Debug("Server: handling route", logger.Fields{
			"request_id": requestID,
			"method":     config.Method,
			"path":       config.Path,
			"url":        c.Request.URL,
		})

		s, buildErr := buildService(requestID, config, c)
		if buildErr != nil {
			c.JSON(buildErr.GetStatus(), buildErr.GetDetails())
			return
		}

		result, err := serviceResult(s)

		if err == nil && result == nil {
			c.Writer.WriteHeader(http.StatusNoContent)
			return
		}

		if err == nil {
			c.JSON(http.StatusOK, result)
			return
		}

		logger.Error("Service Failed", logger.Fields{"err": err})
		c.JSON(err.GetStatus(), err.GetDetails())
	}
}

func serviceResult(s service.TPL) (interface{}, *service.ResultError) {
	result, err := s.Result()

	if err == nil && result == nil {
		return nil, nil
	}

	if err == nil {
		return result, nil
	}

	resultErr, ok := err.(*service.ResultError)
	if ok {
		return nil, resultErr
	}

	logger.Error("Service Failed", logger.Fields{"err": err})
	resultErr = service.ToResultError(err, "service failed", 500)
	return nil, resultErr
}

func buildService(requestID meta.RequestId, config service.Config, c *gin.Context) (service.TPL, *service.ResultError) {
	parsedHeaders := parseHeaders(c.Request)
	parsedBody := rawBody(c.Request.Body)
	parsedParams := parseParams(c.Params, c.Request.URL)
	s, err := config.Build(requestID, parsedHeaders, parsedBody, parsedParams)

	err = handleAuth(
		requestID,
		parsedHeaders,
		parsedBody,
		parsedParams,
		s,
		config,
	)

	if err == nil {
		return s, nil
	}

	resultErr, ok := err.(*service.ResultError)
	if !ok {
		logger.Error("Service Build Failed", logger.Fields{"err": err})
		resultErr = service.ToResultError(err, "service failed", 500)
		return nil, resultErr
	}

	return nil, err.(*service.ResultError)
}

func handleAuth(
	requestID meta.RequestId,
	parsedHeaders service.RawHeaders,
	parsedBody service.RawBody,
	parsedParams service.RawParams,
	s service.TPL,
	config service.Config,
) error {
	var err error

	if !config.IsPrivate {
		return nil
	}

	if _, ok := interface{}(s).(service.PrivateTPL); config.IsPrivate && !ok {
		return fmt.Errorf("Private Service requires the Auth method")
	}

	ps, _ := interface{}(s).(service.PrivateTPL)
	ok := ps.Auth(requestID, parsedHeaders, parsedBody, parsedParams)

	if !ok {
		err = fmt.Errorf("Authentication Failed")
		logger.Error("Authentication Failed", logger.Fields{"err": err})
		resultErr := service.ToResultError(err, "service auth failed", 401)
		return resultErr
	}

	return nil
}

func parseHeaders(Request *http.Request) service.RawHeaders {
	var parsedHeaders = service.RawHeaders{}

	for key, values := range Request.Header {
		parsedHeaders = append(parsedHeaders, service.RawHeader{
			Key:    key,
			Values: values,
		})
	}

	return parsedHeaders
}

func rawBody(source io.ReadCloser) []byte {
	buf := make([]byte, 1024)
	num, _ := source.Read(buf)
	reqBody := string(buf[0:num])
	return []byte(reqBody)
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
			Key:    k,
			Values: v,
		})
	}

	return parsedParams
}

func EnableSwagger() {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
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

func healthGET(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "UP",
	})
}

func EnableHealthCheck() {
	router.Handle("GET", "/status", healthGET)
}

func Start() {
	for _, service := range services {
		InitService(service)
	}

	router.Run()
}
