package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/simiancreative/simiango/service"
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

func TestHandleErrorOK(t *testing.T) {
	// Create a service.ResultError
	err := &service.ResultError{Status: 500, Message: "Internal Server Error"}

	// Call the function with the error
	resultErr := handleError(err)

	// Assert that the returned error is the same as the input error
	assert.Equal(t, err, resultErr)
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

func TestHandleRecovery(t *testing.T) {
	// Create a new gin context with a recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create a context map
	context := make(map[string]interface{})
	context["error"] = "test error"

	// Call the function with the test context and context map
	handleRecovery(c, context)

	// Assert that the HTTP status code was set correctly
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Assert that the JSON response was set correctly
	expectedBody := `{"message":"Internal Server Error","reasons":[{"error":"test error"}]}`
	assert.Equal(t, expectedBody, w.Body.String())
}

func TestHandleErrorResp(t *testing.T) {
	// Create a new gin context with a recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create a service.ResultError
	err := &service.ResultError{Status: 500, Message: "Internal Server Error"}

	// Call the function with the error and context
	resultErr := handleErrorResp(err, c)

	// Assert that the returned error is the same as the input error
	assert.Equal(t, err, resultErr)
}

func TestHandleErrorRespNil(t *testing.T) {
	// Create a new gin context with a recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Call the function with a nil error
	resultErr := handleErrorResp(nil, c)

	// Assert that the returned error is nil
	assert.Nil(t, resultErr)
}

func TestHandleErrorRespStatusLessThan500(t *testing.T) {
	// Create a new gin context with a recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create a service.ResultError with a status less than 500
	err := &service.ResultError{Status: 400, Message: "Bad Request"}

	// Call the function with the error and context
	resultErr := handleErrorResp(err, c)

	// Assert that the returned error is the same as the input error
	assert.Equal(t, err, resultErr)
}
