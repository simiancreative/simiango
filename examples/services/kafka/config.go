package kafka

import (
	"github.com/simiancreative/simiango/messaging/kafka"
	"github.com/simiancreative/simiango/service"
)

var Config = service.Config{
	// Matched to KAFKA_HANDLER when selecting the runtime service that will
	// process messages
	Key:           "example-handler",
	BuildMessages: Build,
}

// dont forget to import your package in your main.go for initialization
// _ "path/to/project/stream"
func init() {
	kafka.AddService(Config)
}
