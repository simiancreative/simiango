package lib

import (
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

func SentryScope(c *gin.Context, s *sentry.Scope) {
	s.SetExtra("extra", "extra")
}
