package sentry_test

import (
	"os"
	"testing"

	sentrygo "github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"

	"github.com/simiancreative/simiango/logger"
	"github.com/simiancreative/simiango/monitoring/sentry"
)

func TestEnable(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func()
		expectedErr bool
		assertions  func(t *testing.T)
	}{
		{
			name: "Missing DSN",
			setupFunc: func() {
				os.Unsetenv("SENTRY_DSN")
			},
			assertions: func(t *testing.T) {
				hub := sentrygo.CurrentHub().Client()
				assert.Nil(t, hub)
			},
		},
		{
			name: "Init error",
			setupFunc: func() {
				os.Setenv("SENTRY_DSN", "invalid dsn")
			},
			assertions: func(t *testing.T) {
				hub := sentrygo.CurrentHub().Client()
				assert.Nil(t, hub)
			},
		},
		{
			name: "Valid case",
			setupFunc: func() {
				os.Setenv("SENTRY_DSN", "https://examplePublicKey@o0.ingest.sentry.io/0")
				os.Setenv("APP_VERSION", "1.0.0")
				os.Setenv("APP_ENV", "test")
			},
			assertions: func(t *testing.T) {
				hub := sentrygo.CurrentHub()
				assert.NotNil(t, hub)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.Enable()

			// Setup the test
			tt.setupFunc()

			// Call the function under test
			sentry.Enable()

			tt.assertions(t)
		})
	}
}

func TestRecoverAndThrow(t *testing.T) {
	// This test will panic, so we need to recover in the test
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	defer sentry.RecoverAndThrow()

	panic("test panic")
}
