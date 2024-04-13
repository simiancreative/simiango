package sentry_test

import (
	"fmt"
	"testing"

	sentrygo "github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/simiancreative/simiango/monitoring/sentry"
	"github.com/simiancreative/simiango/service"
)

func TestHandleError(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() *gin.Context
		err       *service.ResultError
	}{
		{
			name: "With nil hub",
			setupFunc: func() *gin.Context {
				// Mock gin.Context
				return &gin.Context{}
			},
			err: service.ToResultError(
				fmt.Errorf("test error"),
				"test error",
				500,
			),
		},
		{
			name: "With scopeFunc",
			setupFunc: func() *gin.Context {
				// Mock gin.Context
				c := &gin.Context{}

				// Add a sentry.Hub to the gin.Context
				hub := sentrygo.CurrentHub().Clone()
				c.Set("sentry", hub)

				// Add a scopeFunc to the gin.Context
				scopeFunc := func(_ *gin.Context, _ *sentrygo.Scope) {}
				c.Set("sentryScopeFunc", scopeFunc)

				return c
			},
			err: service.ToResultError(
				fmt.Errorf("test error"),
				"test error",
				500,
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function under test
			result := sentry.HandleError(tt.setupFunc(), tt.err)

			// Assert that the result is the same as the input error
			assert.Equal(t, tt.err, result)
		})
	}
}
