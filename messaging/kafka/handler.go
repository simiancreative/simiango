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

func buildKafkaMessages(messages service.Messages, out chan<- kafka.Message) {
	if len(messages) == 0 {
		return
	}

	for _, message := range messages {
		marshalled, _ := json.Marshal(message.Value)
		out <- kafka.Message{
			Key:   []byte(message.Key),
			Value: marshalled,
		}
	}
}

func Handle(messages <-chan kafka.Message) <-chan kafka.Message {
	handlerName := os.Getenv("KAFKA_HANDLER")

	readerConfig, err := findService(handlerName)

	if err != nil {
		logger.Panic("error finding reader service", logger.Fields{
			"err": err.Error(),
		})
	}

	results := make(chan kafka.Message)

	handler := func(message kafka.Message) {
		requestID := meta.Id()

		defer meta.RescuePanic(requestID, message)

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

		buildKafkaMessages(messages, results)
	}

	go func() {
		defer close(results)

		for {
			if _, open := <-messages; !open {
				logger.Printf("Kafka: closing handler")
				return
			}
		}
	}()

	go func() {
		for message := range messages {
			handler(message)
		}
	}()

	return results
}

func findService(key string) (service.Config, error) {
	for _, config := range services {
		if config.Key == key {
			return config, nil
		}
	}

	return service.Config{}, fmt.Errorf("No service found")
}
