package kafka

import (
	"context"
	"strings"

	kafka "github.com/segmentio/kafka-go"
	"github.com/simiancreative/simiango/logger"
)

func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Logger:      Logger{},
		StartOffset: kafka.LastOffset,
		Brokers:     brokers,
		GroupID:     groupID,
		Topic:       topic,
		MinBytes:    1,    // 1B
		MaxBytes:    10e6, // 10MB
	})
}

func NewConsumer(kafkaURL, topic, groupID string, done <-chan bool) <-chan kafka.Message {
	reader := getKafkaReader(kafkaURL, topic, groupID)

	stop := false
	messages := make(chan kafka.Message)

	go func() {
		defer close(messages)
		defer reader.Close()

		select {
		case <-done:
			stop = true
			logger.Printf("Kafka: closing consumer")
			return
		}
	}()

	go func() {
		for {
			if stop {
				return
			}

			readMessages(reader, messages)
		}
	}()

	return messages
}

func readMessages(reader *kafka.Reader, messages chan kafka.Message) {
	m, err := reader.ReadMessage(context.Background())
	if err != nil {
		logger.Error("Kafka: Error reading message", logger.Fields{"err": err.Error()})
		return
	}

	messages <- m

	logger.Printf("Kafka: read message (%+v)", m)
}
