package kafka

import (
	"github.com/simiancreative/simiango/logger"
	"github.com/simiancreative/simiango/service"
)

type kafkaService struct {
	Time int64 `json:"time"`
}

func (s *kafkaService) Result() (service.Messages, error) {
	logger.Printf("Kafka Service: Result called - %+v", s)

	messages := make(service.Messages, 1)
	messages[0] = service.Message{
		// partition key used to target a specific partition. Usually a record
		// identifier user, device, ...
		Key: "processed-example",

		// message value
		Value: struct {
			// fields need to be exported
			Field int64
		}{
			Field: s.Time,
		},
	}

	logger.Printf("Kafka Service: Result complete - %+v", messages)

	return messages, nil
}
