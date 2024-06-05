package logger

import (
	"github.com/sirupsen/logrus/hooks/test"
)

// Mock creates a new logger and hook for testing purposes. use the hook to
// assert log messages.
func Mock() *test.Hook {
	logger, hook := test.NewNullLogger()

	global = logger

	return hook
}
