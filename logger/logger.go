package logger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var global *logrus.Logger

type Fields map[string]interface{}

func Type() string {
	logType, _ := os.LookupEnv("LOG_TYPE")
	return logType
}

func Level() logrus.Level {
	logLevel, _ := os.LookupEnv("LOG_LEVEL")

	switch logLevel {
	case "trace":
		return logrus.TraceLevel
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	}

	return logrus.DebugLevel
}

func New() *logrus.Logger {
	inst := logrus.New()

	if Type() != "line" {
		inst.SetFormatter(&logrus.JSONFormatter{})
	}

	if Type() == "line" {
		inst.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	inst.SetLevel(Level())
	return inst
}

func init() {
	global = New()

	Trace(
		"Logger: ready",
		Fields{"type": Type(), "level": Level()},
	)
}

func Debugf(message string, replacements ...interface{}) {
	global.Debug(fmt.Sprintf(message, replacements...))
}

func Printf(message string, replacements ...interface{}) {
	global.Info(fmt.Sprintf(message, replacements...))
}

func Errorf(message string, replacements ...interface{}) {
	global.Error(fmt.Sprintf(message, replacements...))
}

func Trace(message string, fields Fields) {
	logFields := logrus.Fields(fields)
	global.WithFields(logFields).Trace(message)
}

func Debug(message string, fields Fields) {
	logFields := logrus.Fields(fields)
	global.WithFields(logFields).Debug(message)
}

func Info(message string, fields Fields) {
	logFields := logrus.Fields(fields)
	global.WithFields(logFields).Info(message)
}

func Warn(message string, fields Fields) {
	logFields := logrus.Fields(fields)
	global.WithFields(logFields).Warn(message)
}

func Error(message string, fields Fields) {
	logFields := logrus.Fields(fields)
	global.WithFields(logFields).Error(message)
}

func Fatal(message string, fields Fields) {
	logFields := logrus.Fields(fields)
	global.WithFields(logFields).Fatal(message)
}

func Panic(message string, fields Fields) {
	logFields := logrus.Fields(fields)
	global.WithFields(logFields).Panic(message)
}

func GetEntry(fields Fields) *logrus.Entry {
	logFields := logrus.Fields(fields)
	return global.WithFields(logFields)
}
