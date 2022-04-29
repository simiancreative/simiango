package kafka

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	kafka "github.com/segmentio/kafka-go"
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
		Logger:    kl,
		Addr:      kafka.TCP(kafkaURL),
		Topic:     topic,
		Balancer:  &kafka.LeastBytes{},
		BatchSize: getBatchSize(),
	}
}

func closeWriter(writer *kafka.Writer) {
	kl.Printf("closing writer")

	if err := writer.Close(); err != nil {
		kl.Fatal("failed to close writer:", fields{"err": err})
	}
}

func Writer(kafkaURL, topic string) *kafka.Writer {
	return getKafkaWriter(kafkaURL, topic)
}

func NewProducer(kafkaURL, topic string, in <-chan kafka.Message) {
	writer := getKafkaWriter(kafkaURL, topic)

	go func() {
		defer closeWriter(writer)

		batches := BatchMessages(in, getBatchSize(), getBatchTimeout())

		for messages := range batches {
			writeMessage(writer, messages)
		}
	}()
}

func writeMessage(writer *kafka.Writer, messages []kafka.Message) {
	err := writer.WriteMessages(
		context.Background(),
		messages...,
	)

	if err != nil {
		kl.Error("Kafka Producer: Failed to write messages", fields{"err": err})
		return
	}

	kl.Info("Kafka Producer: wrote messages to kafka", fields{
		"topic":    writer.Topic,
		"count":    len(messages),
		"messages": fmt.Sprintf("%+v", messages),
	})
}
