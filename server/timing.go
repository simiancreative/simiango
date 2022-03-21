package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/go-server-timing"
)

func TimingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		timing := servertiming.Header{}
		c.Set("timing", &timing)

		metric := timing.NewMetric("req").Start()

		c.Next()

		metric.Stop()
	}
}
