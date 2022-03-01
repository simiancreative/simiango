package kafka

import (
	"encoding/json"
	"fmt"
	"os"

	kafka "github.com/segmentio/kafka-go"
	"github.com/simiancreative/simiango/logger"
	"github.com/simiancreative/simiango/meta"
	"github.com/simiancreative/simiango/service"
)

func buildKafkaMessages(messages service.Messages) []kafka.Message {
	kafkaMessages := []kafka.Message{}

	for _, message := range messages {
		marshalled, _ := json.Marshal(message.Value)
		kafkaMessages = append(kafkaMessages, kafka.Message{
			Key:   []byte(message.Key),
			Value: marshalled,
		})
	}

	return kafkaMessages
}

func Handle(c <-chan kafka.Message) <-chan []kafka.Message {
	handlerName := os.Getenv("KAFKA_HANDLER")

	readerConfig, err := findService(handlerName)

	if err != nil {
		logger.Panic("error finding reader service", logger.Fields{
			"err": err.Error(),
		})
	}

	out := make(chan []kafka.Message)

	handler := func(message kafka.Message) {
		requestID := meta.Id()

		service, err := buildService(requestID, readerConfig, message)
		if err != nil {
			logger.Error("error building service", logger.Fields{"err": err.Error()})
			return
		}

		messages, err := service.Result()
		if err != nil {
			logger.Error("error on exec result", logger.Fields{"err": err.Error()})
			return
		}

		if len(messages) > 0 {
			out <- buildKafkaMessages(messages)
		}
	}

	go func() {
		for message := range c {
			go handler(message)
		}
		close(out)
	}()

	return out
}

func findService(key string) (service.Config, error) {
	for _, config := range services {
		if config.Key == key {
			return config, nil
		}
	}

	return service.Config{}, fmt.Errorf("No service found")
}
