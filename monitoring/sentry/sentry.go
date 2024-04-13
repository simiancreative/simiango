package sentry

import (
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"

	"github.com/simiancreative/simiango/logger"
	"github.com/simiancreative/simiango/service"
)

func Enable() {
	dsn, hasDsn := os.LookupEnv("SENTRY_DSN")
	if !hasDsn {
		logger.Warnf("SENTRY_DSN is not set")
		return
	}

	_, debug := os.LookupEnv("SENTRY_DEBUG")

	options := sentry.ClientOptions{
		Debug:            debug,
		Release:          os.Getenv("APP_VERSION"),
		Dsn:              dsn,
		Environment:      os.Getenv("APP_ENV"),
		TracesSampleRate: 0.1,
		EnableTracing:    true,
		AttachStacktrace: true,
	}

	err := sentry.Init(options)
	if err != nil {
		logger.Warnf("Sentry initialization failed: %v", err)
	}
}

func RecoverAndThrow() {
	err := recover()
	if err != nil {
		sentry.CurrentHub().Recover(err)
		sentry.Flush(time.Second * 5)
		panic(err)
	}
}

func HandleError(c *gin.Context, err *service.ResultError) *service.ResultError {
	hub := sentrygin.GetHubFromContext(c)

	if hub == nil {
		return err
	}

	hub.WithScope(func(scope *sentry.Scope) {
		scopeFunc, ok := c.Get("sentryScopeFunc")
		if ok {
			scopeFunc.(func(*gin.Context, *sentry.Scope))(c, scope)
		}

		hub.CaptureException(err)
	})

	return err
}
