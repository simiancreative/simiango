package server

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/simiancreative/simiango/logger"
	"github.com/simiancreative/simiango/service"
)

func parseHeaders(request *http.Request) service.RawHeaders {
	var parsedHeaders = service.RawHeaders{}

	for key, values := range request.Header {
		param := service.ParamItem{Key: key}
		param.SetValues(values)
		parsedHeaders = append(parsedHeaders, param)
	}

	return parsedHeaders
}

func rawBody(source io.ReadCloser) []byte {
	reqBody, _ := ioutil.ReadAll(source)
	return []byte(reqBody)
}

func parseParams(params gin.Params, url *url.URL) service.RawParams {
	var parsedParams = service.RawParams{}

	for _, element := range params {
		parsedParams = append(parsedParams, service.ParamItem{
			Key:   element.Key,
			Value: element.Value,
		})
	}

	for k, v := range url.Query() {
		parsedParams = append(parsedParams, service.ParamItem{
			Key:    k,
			Value:  v[0],
			Values: v,
		})
	}

	return parsedParams
}

func handleError(err error) *service.ResultError {
	resultErr, ok := err.(*service.ResultError)
	if !ok {
		logger.Error("Service Build Failed", logger.Fields{"err": err})
		resultErr = service.ToResultError(err, "service failed", 500)
		return resultErr
	}

	return err.(*service.ResultError)
}

func handleAfter(config service.Config, req service.Req) {
	if config.After == nil {
		return
	}

	config.After(config, req)
}
