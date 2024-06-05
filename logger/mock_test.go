package logger_test

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/simiancreative/simiango/logger"
)

// TestMock demonstrates how to use the Mock function to create a mock logger
// and hook for testing purposes.
func TestMock(t *testing.T) {
	hook := logger.Mock()

	logger.Info("Test message", logger.Fields{})

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.InfoLevel, hook.LastEntry().Level)
	assert.Equal(t, "Test message", hook.LastEntry().Message)
}
