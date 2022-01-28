package kafka

import (
	"fmt"

	kafka "github.com/segmentio/kafka-go"
	"github.com/simiancreative/simiango/logger"
	"github.com/simiancreative/simiango/meta"
	"github.com/simiancreative/simiango/service"
)

func handle(c <-chan kafka.Message) <-chan []kafka.Message {
	out := make(chan []kafka.Message)

	handler := func(message kafka.Message) {
		fmt.Printf("******* the kafka message received: %+v\n", message)
		readerConfig, err := findService("reader")
		if err != nil {
			logger.Error("error finding reader service", logger.Fields{
				"message_key":   string(message.Key),
				"message_value": string(message.Value),
				"err":           err.Error(),
			})
			return
		}
		fmt.Printf("******* the kafka message.Key: %s\n", message.Key)
		fmt.Printf("******* the kafka message.Value: %s\n", message.Value)
		requestID := meta.Id()
		readerService, err := buildService(requestID, readerConfig, message.Value)
		if err != nil {
			logger.Error("error building service", logger.Fields{"err": err.Error()})
			return
		}

		object, err := readerService.Result()
		if err != nil {
			logger.Error("error on exec result", logger.Fields{"err": err.Error()})
			return
		}

		writerConfig, err := findService("writer")
		if err != nil {
			logger.Warn("no writer service found", logger.Fields{
				"err": err.Error(),
			})
			return
		}

		body, ok := object.(service.RawBody)
		if !ok {
			logger.Error("Error casting interface object to RawBody", logger.Fields{
				"err": fmt.Sprintf("casted to type %T, value: %v", body, body),
			})
			return
		}
		writerService, err := buildService(requestID, writerConfig, body)
		if err != nil {
			logger.Error("error building service", logger.Fields{"err": err.Error()})
			return
		}

		object, err = writerService.Result()
		if err != nil {
			logger.Error("error on exec result", logger.Fields{"err": err.Error()})
			return
		}

		messages, ok := object.([]kafka.Message)
		if !ok {
			logger.Error("Error casting interface object to array of messages", logger.Fields{
				"err": fmt.Sprintf("casted to type %T, value: %v", messages, messages),
			})
			return
		}
		fmt.Printf("Got messages %+v\n", messages)
		out <- messages
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
