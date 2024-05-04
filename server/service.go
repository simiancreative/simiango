package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	servertiming "github.com/p768lwy3/gin-server-timing"

	"github.com/simiancreative/simiango/service"
)

var kinds = map[service.Kind]func(
	service.Config,
	service.Req,
) (interface{}, *service.ResultError){
	service.DEFAULT: handleDefault,
	service.DIRECT:  handleDirect,
}

func handleService(config service.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		timing := servertiming.FromContext(c)
		timer := timing.NewMetric("total").Start()
		req := parseRequest(c)

		var result interface{}
		var err *service.ResultError

		handler, ok := kinds[config.Kind]
		if !ok {
			result, err = handleDefault(config, req)
		}

		if ok {
			result, err = handler(config, req)
		}

		timer.Stop()
		servertiming.WriteHeader(c)

		go handleAfter(config, req)

		if handleErrorResp(err, c) != nil {
			c.JSON(err.GetStatus(), err.GetDetails())
			return
		}

		if result == nil {
			c.Writer.WriteHeader(http.StatusNoContent)
			return
		}

		if config.IsStream {
			handleStreamingServiceResult(result, c)
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
