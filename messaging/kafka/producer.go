package kafka

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	kafka "github.com/segmentio/kafka-go"
	"github.com/simiancreative/simiango/logger"
)

func getBatchTimeout() time.Duration {
	amount := 2000

	if amountStr, present := os.LookupEnv("KAFKA_BATCH_TIMEOUT"); present {
		amount, _ = strconv.Atoi(amountStr)
	}

	return time.Duration(amount) * time.Millisecond
}

func getBatchSize() int {
	size := 100

	if sizeStr, present := os.LookupEnv("KAFKA_BATCH_SIZE"); present {
		size, _ = strconv.Atoi(sizeStr)
	}

	return size
}

func getKafkaWriter(kafkaURL, topic string) *kafka.Writer {

	return &kafka.Writer{
		Addr:      kafka.TCP(kafkaURL),
		Topic:     topic,
		Balancer:  &kafka.LeastBytes{},
		BatchSize: getBatchSize(),
	}
}

func closeWriter(writer *kafka.Writer) {
	logger.Printf("closing writer")

	if err := writer.Close(); err != nil {
		logger.Fatal("failed to close writer:", logger.Fields{"err": err})
	}
}

func Writer(kafkaURL, topic string) *kafka.Writer {
	return getKafkaWriter(kafkaURL, topic)
}

func NewProducer(kafkaURL, topic string, in <-chan kafka.Message, done <-chan bool) {
	writer := getKafkaWriter(kafkaURL, topic)

	go func() {
		defer closeWriter(writer)

		batches := BatchMessages(in, getBatchSize(), getBatchTimeout())

		defer close(batches)

		for messages := range batches {
			select {
			case <-done:
				logger.Printf("Kafka: closing handler")
				return
			default:
				writeMessage(writer, messages)
			}
		}
	}()
}

func writeMessage(writer *kafka.Writer, messages []kafka.Message) {
	err := writer.WriteMessages(
		context.Background(),
		messages...,
	)

	if err != nil {
		logger.Error("Kafka Producer: Failed to write messages", logger.Fields{"err": err})
		return
	}

	logger.Info("Kafka Producer: wrote messages to kafka", logger.Fields{
		"topic":    writer.Topic,
		"count":    len(messages),
		"messages": fmt.Sprintf("%+v", messages),
	})
}
