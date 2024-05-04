package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/simiancreative/simiango/service"
)

func errorMiddleware(c *gin.Context) {
	c.Next()

	if len(c.Errors) == 0 {
		return
	}

	resp := service.ResultError{
		Status:  http.StatusInternalServerError,
		Message: "Internal Server Error",
	}

	for _, err := range c.Errors {
		resp.Reasons = append(resp.Reasons, map[string]interface{}{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusInternalServerError, resp)
}
