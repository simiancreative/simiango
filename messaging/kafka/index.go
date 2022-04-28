package kafka

import (
	"os"

	"github.com/simiancreative/simiango/logger"
	"github.com/simiancreative/simiango/meta"
)

func Start() {
	url := os.Getenv("KAFKA_BROKERS")
	readerTopic := os.Getenv("KAFKA_READER_TOPIC")
	groupID := os.Getenv("KAFKA_READER_GROUP")

	done := meta.CatchSig()

	messages := NewConsumer(url, readerTopic, groupID, done)
	results := Handle(messages, done)

	meta.AddCleanup(func() {
		// wait for the channels to close.
		<-messages
		<-results
	})

	if writerTopic, present := os.LookupEnv("KAFKA_WRITER_TOPIC"); present {
		NewProducer(url, writerTopic, results, done)
		logger.Printf("Kafka: Producer setup (topic: %v)", writerTopic)
	}
}
