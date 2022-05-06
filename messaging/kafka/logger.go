package kafka

import (
	ll "github.com/simiancreative/simiango/logger"
)

var kl = logger{}

type logger struct {
	Fields ll.Fields
}

type fields map[string]interface{}

func (l logger) Printf(pattern string, replacements ...interface{}) {
	ll.Debugf("Kafka: "+pattern, replacements...)
}

func (l logger) Error(pattern string, fields map[string]interface{}) {
	ll.Error("Kafka: "+pattern, fields)
}

func (l logger) Panic(pattern string, fields map[string]interface{}) {
	ll.Panic("Kafka: "+pattern, fields)
}

func (l logger) Fatal(pattern string, fields map[string]interface{}) {
	ll.Fatal("Kafka: "+pattern, fields)
}

func (l logger) Info(pattern string, fields map[string]interface{}) {
	ll.Info("Kafka: "+pattern, fields)
}

func (l logger) Warn(pattern string, fields map[string]interface{}) {
	ll.Warn("Kafka: "+pattern, fields)
}
