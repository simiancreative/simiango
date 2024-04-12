package server

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/simiancreative/simiango/logger"
	"github.com/simiancreative/simiango/meta"
)

func getClientIP(c *gin.Context) string {
	// first check the X-Forwarded-For header
	requester := c.Request.Header.Get("X-Forwarded-For")
	// if empty, check the Real-IP header
	if len(requester) == 0 {
		requester = c.Request.Header.Get("X-Real-IP")
	}
	// if the requester is still empty, use the hard-coded address from the socket
	if len(requester) == 0 {
		requester = c.Request.RemoteAddr
	}

	// if requester is a comma delimited list, take the first one
	// (this happens when proxied via elastic load balancer then again through nginx)
	if strings.Contains(requester, ",") {
		requester = strings.Split(requester, ",")[0]
	}

	return requester
}

func getAuth(c *gin.Context) string {
	token, exists := c.Get("Authorization")
	if exists {
		return token.(string)
	}
	return ""
}

func JSONLogMiddleware(c *gin.Context) {
	// Start timer
	start := time.Now()

	// Process Request
	c.Next()

	// Stop timer
	duration := meta.GetDurationInMillseconds(start)

	entry := logger.GetEntry(logger.Fields{
		"client_ip":  getClientIP(c),
		"duration":   duration,
		"method":     c.Request.Method,
		"path":       c.Request.RequestURI,
		"status":     c.Writer.Status(),
		"auth":       getAuth(c),
		"referrer":   c.Request.Referer(),
		"request_id": c.Writer.Header().Get("X-Request-ID"),
	})

	if c.Writer.Status() >= 500 {
		entry.Error(c.Errors.String())
	} else {
		entry.Info("")
	}
}
