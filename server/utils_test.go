package server

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/simiancreative/simiango/service"
	"github.com/stretchr/testify/assert"
)

func TestUtilsHeaders(t *testing.T) {
	req := http.Request{Header: http.Header{
		"Content-Type": []string{"application/json"},
	}}

	parsed := parseHeaders(&req)

	assert.IsType(t, parsed, service.RawHeaders{})
}

func TestUtilsParams(t *testing.T) {
	req := gin.Params{
		gin.Param{
			Key:   "Content-Type",
			Value: "application/json",
		},
	}

	u := url.URL{RawQuery: "hi=there"}

	parsed := parseParams(req, &u)

	assert.IsType(t, parsed, service.RawParams{})
}

func TestHandleError(t *testing.T) {
	err := fmt.Errorf("new test error")
	res := handleError(err)

	assert.IsType(t, res, &service.ResultError{})
}

func TestHandleAfterPass(t *testing.T) {
	config := service.Config{
		Kind: service.DIRECT,
		Direct: func(r service.Req) (interface{}, error) {
			return nil, nil
		},
	}

	req := service.Req{}

	handleAfter(config, req)
}

func TestHandleAfterCall(t *testing.T) {
	config := service.Config{
		Kind: service.DIRECT,
		Direct: func(r service.Req) (interface{}, error) {
			return nil, nil
		},
		After: func(config service.Config, req service.Req) {},
	}

	req := service.Req{}

	handleAfter(config, req)
}
