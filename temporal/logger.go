package temporal

import (
	"github.com/simiancreative/simiango/logger"
	"github.com/sirupsen/logrus"
)

// CustomLogger adapts the Simiango logger to Temporal's logging interface.
type CustomLogger struct {
	simiangoLogger *logrus.Logger
}

// NewCustomLogger creates a new instance of CustomLogger.
func NewCustomLogger() *CustomLogger {
	return &CustomLogger{
		simiangoLogger: logger.New(),
	}
}

// Debug logs a debug message.
func (c *CustomLogger) Debug(msg string, keyvals ...interface{}) {
	c.simiangoLogger.Debug(msg, keyvals)
}

// Info logs an info message.
func (c *CustomLogger) Info(msg string, keyvals ...interface{}) {
	c.simiangoLogger.Info(msg, keyvals)
}

// Warn logs a warning message.
func (c *CustomLogger) Warn(msg string, keyvals ...interface{}) {
	c.simiangoLogger.Warn(msg, keyvals)
}

// Error logs an error message.
func (c *CustomLogger) Error(msg string, keyvals ...interface{}) {
	c.simiangoLogger.Error(msg, keyvals)
}
