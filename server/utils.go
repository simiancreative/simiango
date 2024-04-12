package server

import (
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/simiancreative/simiango/monitoring/sentry"
	"github.com/simiancreative/simiango/service"
)

func parseHeaders(request *http.Request) service.RawHeaders {
	parsedHeaders := service.RawHeaders{}

	for key, values := range request.Header {
		param := service.ParamItem{Key: key}
		param.SetValues(values)
		parsedHeaders = append(parsedHeaders, param)
	}

	return parsedHeaders
}

func rawBody(source io.ReadCloser) []byte {
	reqBody, _ := io.ReadAll(source)
	return []byte(reqBody)
}

func parseParams(params gin.Params, url *url.URL) service.RawParams {
	parsedParams := service.RawParams{}

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

func handleErrorResp(err *service.ResultError, c *gin.Context) *service.ResultError {
	if err == nil {
		return nil
	}

	status := err.GetStatus()

	if status < 500 {
		return err
	}

	return sentry.HandleError(c, err)
}

func handleError(err error) *service.ResultError {
	resultErr, ok := err.(*service.ResultError)

	if !ok {
		resultErr = service.ToResultError(err, "service failed", 500)
		return resultErr
	}

	return resultErr
}

func handleAfter(config service.Config, req service.Req) {
	if config.After == nil {
		return
	}

	config.After(config, req)
}
