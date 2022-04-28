package kafka

import (
	"os"

	"github.com/simiancreative/simiango/logger"
	"github.com/simiancreative/simiango/meta"
)

func Start(done chan bool) {
	url := os.Getenv("KAFKA_BROKERS")
	readerTopic := os.Getenv("KAFKA_READER_TOPIC")
	groupID := os.Getenv("KAFKA_READER_GROUP")

	messages := NewConsumer(url, readerTopic, groupID, done)
	results := Handle(messages)

	meta.AddCleanup(func() {
		// wait for the channels to close.
		<-messages
		<-results

		logger.Printf("Kafka: cleanup complete")
	})

	if writerTopic, present := os.LookupEnv("KAFKA_WRITER_TOPIC"); present {
		NewProducer(url, writerTopic, results)
		logger.Printf("Kafka: Producer setup (topic: %v)", writerTopic)
	}
}
