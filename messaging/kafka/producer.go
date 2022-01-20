package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	kafka "github.com/segmentio/kafka-go"
)

func getKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func closeWriter(writer *kafka.Writer) {
	fmt.Println("closing writer")
	if err := writer.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}

func buildMessages(result Result) []kafka.Message {
	messages := []kafka.Message{}

	for _, content := range result.Content {
		marshalled, _ := json.Marshal(content)
		messages = append(messages, kafka.Message{
			Key:   []byte(result.Key),
			Value: marshalled,
		})
	}

	return messages
}

func NewProducer(kafkaURL, topic string, in <-chan Result) <-chan bool {
	writer := getKafkaWriter(kafkaURL, topic)
	done := make(chan bool)

	go func() {
		defer closeWriter(writer)

		for result := range in {
			messages := buildMessages(result)
			err := writer.WriteMessages(
				context.Background(),
				messages...,
			)
			if err != nil {
				log.Fatal("Failed to write messages:", err)
			}
		}

		done <- true
	}()
	return done
}
