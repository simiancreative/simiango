package kafka

import (
	//_ "consumers/notifications/templates/access_method"
	"os"

	"github.com/simiancreative/simiango/messaging/kafka"
	"github.com/simiancreative/simiango/service"
)

var Config = service.Config{
	IsStream:      true, // please check
	IsPrivate:     false,
	Key:           os.Getenv("KAFKA_READER_TOPIC"), // please check
	BuildMessages: Build,
}

// dont forget to import your package in your main.go for initialization
// _ "path/to/project/stream"
func init() {
	kafka.AddService(Config)
}
