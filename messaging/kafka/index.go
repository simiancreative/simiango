package kafka

import (
	"context"
	"os"

	"github.com/simiancreative/simiango/meta"
)

func Start(done context.Context) {
	url := os.Getenv("KAFKA_BROKERS")
	readerTopic := os.Getenv("KAFKA_READER_TOPIC")
	groupID := os.Getenv("KAFKA_READER_GROUP")

	messages := NewConsumer(url, readerTopic, groupID, done)
	results := Handle(messages)

	meta.AddCleanup(func() {
		// wait for the channels to close.
		<-messages
		<-results

		kl.Printf("cleanup complete")
	})

	if writerTopic, present := os.LookupEnv("KAFKA_WRITER_TOPIC"); present {
		NewProducer(url, writerTopic, results)
		kl.Printf("Producer setup (topic: %v)", writerTopic)
	}
}
