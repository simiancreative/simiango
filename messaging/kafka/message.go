package kafka

import (
	"fmt"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

type SimpleMessage struct {
	Topic         string
	Partition     int
	Offset        int64
	HighWaterMark int64
	Headers       []kafka.Header
	Time          time.Time
	Key           string
	Value         string
}

func MessageSimplified(message kafka.Message) SimpleMessage {
	return SimpleMessage{
		Topic:         message.Topic,
		Partition:     message.Partition,
		Offset:        message.Offset,
		HighWaterMark: message.HighWaterMark,
		Headers:       message.Headers,
		Time:          message.Time,
		Key:           string(message.Key),
		Value:         string(message.Value),
	}
}

func MessagesSimplified(messages []kafka.Message) []SimpleMessage {
	simplified := []SimpleMessage{}

	for _, message := range messages {
		simple := MessageSimplified(message)

		simplified = append(simplified, simple)
	}

	return simplified
}

func MessageAsString(message kafka.Message) string {
	return fmt.Sprintf("%+v", MessageSimplified(message))
}
