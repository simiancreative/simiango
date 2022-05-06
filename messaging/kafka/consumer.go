package kafka

import (
	"context"
	"strings"

	kafka "github.com/segmentio/kafka-go"
)

func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Logger:      kl,
		StartOffset: kafka.LastOffset,
		Brokers:     brokers,
		GroupID:     groupID,
		Topic:       topic,
		MinBytes:    1,    // 1B
		MaxBytes:    10e6, // 10MB
	})
}

func NewConsumer(kafkaURL, topic, groupID string, done context.Context) <-chan kafka.Message {
	reader := getKafkaReader(kafkaURL, topic, groupID)

	stop := false
	messages := make(chan kafka.Message)

	go func() {
		defer close(messages)
		defer reader.Close()

		select {
		case <-done.Done():
			stop = true
			kl.Printf("closing consumer")
			return
		}
	}()

	go func() {
		for !stop {
			readMessages(reader, messages)
		}
	}()

	return messages
}

func readMessages(reader *kafka.Reader, messages chan kafka.Message) {
	m, err := reader.ReadMessage(context.Background())
	if err != nil {
		kl.Error("Error reading message", fields{"err": err.Error()})
		return
	}

	messages <- m

	kl.Printf("read message (%+v)", MessageSimplified(m))
}
