package kafka

import (
	"fmt"
	"os"

	"github.com/simiancreative/simiango/logger"
)

func Start() {
	url := os.Getenv("KAFKA_BROKERS")
	readerTopic := os.Getenv("KAFKA_READER_TOPIC")
	groupID := os.Getenv("KAFKA_READER_GROUP")

	c1 := NewConsumer(url, readerTopic, groupID)

	c2, handlerDone := Handle(c1)

	if writerTopic, present := os.LookupEnv("KAFKA_WRITER_TOPIC"); present {
		go NewProducer(url, writerTopic, c2)
		logger.Printf("Kafka: Producer setup (topic: %v)", writerTopic)
	}
	fmt.Println(<-handlerDone)
}
