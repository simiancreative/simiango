package kafka

import (
	"context"
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

func NewProducer(kafkaURL, topic string, in <-chan []kafka.Message) <-chan bool {
	writer := getKafkaWriter(kafkaURL, topic)
	done := make(chan bool)

	go func() {
		defer closeWriter(writer)

		for messages := range in {
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
