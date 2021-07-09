package logger

import (
	log "github.com/sirupsen/logrus"
	"os"
)

type Fields map[string]interface{}

func init() {
	logType, _ := os.LookupEnv("LOG_TYPE")
	logLevel, logLevelExists := os.LookupEnv("LOG_LEVEL")

	if logType != "line" {
		log.SetFormatter(&log.JSONFormatter{})
	}

	if logType == "line" {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
	}

	if !logLevelExists {
		log.SetLevel(log.TraceLevel)
		return
	}

	switch logLevel {
	case "trace":
		log.SetLevel(log.TraceLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	}
}

func Trace(message string, fields Fields) {
	logFields := log.Fields(fields)
	log.WithFields(logFields).Trace(message)
}

func Debug(message string, fields Fields) {
	logFields := log.Fields(fields)
	log.WithFields(logFields).Debug(message)
}

func Info(message string, fields Fields) {
	logFields := log.Fields(fields)
	log.WithFields(logFields).Info(message)
}

func Warn(message string, fields Fields) {
	logFields := log.Fields(fields)
	log.WithFields(logFields).Warn(message)
}

func Error(message string, fields Fields) {
	logFields := log.Fields(fields)
	log.WithFields(logFields).Error(message)
}

func Fatal(message string, fields Fields) {
	logFields := log.Fields(fields)
	log.WithFields(logFields).Fatal(message)
}

func Panic(message string, fields Fields) {
	logFields := log.Fields(fields)
	log.WithFields(logFields).Panic(message)
}
