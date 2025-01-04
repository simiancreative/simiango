package backoff_test

import (
	"errors"
	"testing"
	"time"

	"github.com/simiancreative/simiango/backoff"
	"github.com/simiancreative/simiango/config"
	"github.com/stretchr/testify/assert"
)

func TestRetry(t *testing.T) {
	// need to move config to common
	config.SetupTest()

	t.Run("successful retry", func(t *testing.T) {
		retries := 0

		fn := func() (interface{}, error) {
			retries++

			if retries < 3 {
				return nil, errors.New("failed")
			}

			return "success", nil
		}

		config := backoff.Config{
			MaxRetries:   5,
			InitialDelay: time.Second,
			MaxDelay:     time.Second,
		}

		result, err := backoff.Retry(fn, config)

		assert.NoError(t, err)
		assert.Equal(t, "success", result)
		assert.Equal(t, 3, retries)
	})

	t.Run("max retries exceeded", func(t *testing.T) {
		fn := func() (interface{}, error) {
			return nil, errors.New("failed")
		}

		config := backoff.Config{
			MaxRetries:   2,
			InitialDelay: time.Second,
			MaxDelay:     time.Second,
		}

		start := time.Now()
		result, err := backoff.Retry(fn, config)
		elapsed := time.Since(start)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.GreaterOrEqual(t, elapsed, time.Second)
	})
}
