package kafka

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/simiancreative/simiango/logger"
)

func SetupDone() chan bool {
	done := make(chan bool)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-sigs
		logger.Printf("Kafka: SIGINT received (%+v) closing...", sig)

		done <- true
		close(done)
	}()

	return done
}

func Start() {
	url := os.Getenv("KAFKA_BROKERS")
	readerTopic := os.Getenv("KAFKA_READER_TOPIC")
	groupID := os.Getenv("KAFKA_READER_GROUP")

	done := SetupDone()

	messages := NewConsumer(url, readerTopic, groupID, done)
	results := Handle(messages, done)

	if writerTopic, present := os.LookupEnv("KAFKA_WRITER_TOPIC"); present {
		NewProducer(url, writerTopic, results, done)
		logger.Printf("Kafka: Producer setup (topic: %v)", writerTopic)
	}

	logger.Printf("Kafka: done: %+v", <-done)
}
