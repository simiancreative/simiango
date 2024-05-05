package sentry_test

import (
	"errors"
	"net/http/httptest"
	"testing"

	sentrygin "github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/simiancreative/simiango/monitoring/sentry"
)

func TestGinCaptureError(t *testing.T) {
	t.Run("sentryScopeFunc is not ok", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		hub := sentrygin.CurrentHub()
		c.Set("sentry", hub)

		err := errors.New("test error")
		returnedErr := sentry.GinCaptureError(c, err)

		assert.Equal(t, err, returnedErr)
	})
}

func TestScopeFunctionError(t *testing.T) {
	t.Run("scopeFunction is not ok", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		hub := sentrygin.CurrentHub()
		c.Set("sentry", hub)

		// Set a non-function value for "sentryScopeFunc"
		c.Set("sentryScopeFunc", "not a function")

		err := errors.New("test error")
		returnedErr := sentry.GinCaptureError(c, err)

		assert.Equal(t, err, returnedErr)
	})
}
