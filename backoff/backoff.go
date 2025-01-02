package backoff

import (
	"errors"
	"time"

	"github.com/simiancreative/simiango/logger"
)

// RetryFunc defines the signature of the function to be retried.
type RetryFunc func() (interface{}, error)

// Config holds the configuration for the backoff.
type Config struct {
	MaxRetries   int
	InitialDelay time.Duration
	MaxDelay     time.Duration
}

// Retry retries the given function using exponential backoff.
func Retry(fn RetryFunc, config Config) (interface{}, error) {
	var result interface{}
	var err error

	for i := 0; i < config.MaxRetries; i++ {
		result, err = fn()
		if err == nil {
			return result, nil
		}

		logger.Debugf("RETRYING: attempt %d", i+1)

		delay := config.InitialDelay * (1 << i)
		if delay > config.MaxDelay {
			delay = config.MaxDelay
		}

		time.Sleep(delay)
	}

	return nil, errors.New("failed after maximum retries")
}
