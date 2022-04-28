package kafaingest

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/simiancreative/simiango/server"
	"github.com/simiancreative/simiango/service"

	kafkago "github.com/segmentio/kafka-go"
	"github.com/simiancreative/simiango/messaging/kafka"
)

var writer *kafkago.Writer

var Config = service.Config{
	Kind:   service.DIRECT,
	Method: "GET",
	Path:   "/kafka-event",
	Direct: direct,
}

func direct(req service.Req) (interface{}, error) {
	content := map[string]int64{"time": time.Now().UnixNano()}
	marshalled, _ := json.Marshal(content)
	msg := kafkago.Message{
		Key:   []byte("time"),
		Value: marshalled,
	}

	go func() {
		writer.WriteMessages(
			context.Background(),
			msg,
		)
	}()

	return msg, nil
}

// dont forget to import your package in your main.go for initialization
// _ "path/to/project/direct"
func init() {
	writer = kafka.Writer(os.Getenv("KAFKA_BROKERS"), os.Getenv("KAFKA_READER_TOPIC"))
	server.AddService(Config)
}
