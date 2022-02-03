package kafka

import (
	"fmt"
	"os"
)

func Start() {
	url := os.Getenv("KAFKA_BROKERS")
	readerTopic := os.Getenv("KAFKA_READER_TOPIC")
	groupID := os.Getenv("KAFKA_READER_GROUP")

	c1 := NewConsumer(url, readerTopic, groupID)

	c2 := Handle(c1)

	if writerTopic, present := os.LookupEnv("KAFKA_WRITER_TOPIC"); present {
		done := NewProducer(url, writerTopic, c2)
		fmt.Println(<-done)
	}
}
