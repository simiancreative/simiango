package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	servertiming "github.com/p768lwy3/gin-server-timing"

	"github.com/simiancreative/simiango/service"
	"github.com/stretchr/testify/assert"
)

func TestDirectService(t *testing.T) {
	buf := new(bytes.Buffer)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("POST", "/", buf)
	c.Params = []gin.Param{{Key: "k", Value: "v"}}
	h := servertiming.Header{}
	servertiming.NewContext(c, &h)

	req := parseRequest(c)

	config := service.Config{
		Kind: service.DIRECT,
		Direct: func(r service.Req) (interface{}, error) {
			defer req.Timer.NewMetric("testing").Start().Stop()
			return nil, nil
		},
	}

	result, err := handleDirect(config, req)

	assert.Nil(t, result)
	assert.Nil(t, err)
	assert.Regexp(t, "testing;dur=", req.Timer.String())
}
