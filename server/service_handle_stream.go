package server

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/simiancreative/simiango/logger"
	"github.com/simiancreative/simiango/service"
)

func handleStreamingServiceResult(result interface{}, c *gin.Context) {
	streamResult, ok := result.(service.StreamResult)

	if ok {
		c.Header("Content-Type", streamResult.Type)
		c.Header("Content-Length", streamResult.Length)
		c.Stream(streamResult.Writer)
		return
	}

	logger.Error("Streaming Service Failed", logger.Fields{})
	err := service.ToResultError(
		fmt.Errorf("service does not return a stream result"),
		"service failed",
		500,
	)

	c.JSON(err.GetStatus(), err.GetDetails())
}
