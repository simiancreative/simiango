package kafka

import (
	"github.com/simiancreative/simiango/logger"
	"github.com/simiancreative/simiango/service"
)

type kafkaService struct {
	Hi string `json:"hi"`
}

func (s *kafkaService) Result() (service.Messages, error) {
	logger.Printf("Kafka Service: Result called - %v", s)

	messages := make(service.Messages, 1)
	messages[0] = service.Message{
		// partition key used to target a specific partition. Usually a record
		// identifier user, device, ...
		Key: "processed-example",

		// message value
		Value: struct {
			field string
		}{
			field: s.Hi,
		},
	}

	logger.Printf("Kafka Service: Result complete - %v", messages)

	return messages, nil
}
