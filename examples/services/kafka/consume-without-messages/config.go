package kafka

import (
	"github.com/simiancreative/simiango/messaging/kafka"
	"github.com/simiancreative/simiango/service"
)

// Send a message manually
//
// docker exec -it simiango-kafka-1 bash
// kafka-console-producer.sh --broker-list kafka:9092 --topic decoded
//
// at the prompt enter you message and hit return, the service will be invoked
// and service handling will be logged

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
