package kafka

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
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

func NewProducer(kafkaURL, topic string, in <-chan kafka.Message) <-chan bool {
	writer := getKafkaWriter(kafkaURL, topic)
	done := make(chan bool)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	go func() {
		sig := <-sigs
		logger.Printf("Producer, SIGINT received %v, closing...", sig)
		done <- true
		closeWriter(writer)
	}()

	go func() {
		defer closeWriter(writer)

		batches := BatchMessages(in, getBatchSize(), getBatchTimeout())
		for messages := range batches {
			err := writer.WriteMessages(
				context.Background(),
				messages...,
			)

			if err != nil {
				logger.Error("Kafka Producer: Failed to write messages", logger.Fields{"err": err})
				continue
			}

			logger.Info("Kafka Producer: wrote messages to kafka", logger.Fields{
				"topic":    topic,
				"count":    len(messages),
				"messages": fmt.Sprintf("%+v", messages),
			})
		}

		done <- true
	}()

	return done
}
