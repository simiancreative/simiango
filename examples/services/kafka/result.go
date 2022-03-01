package kafka

import (
	"fmt"

	"github.com/simiancreative/simiango/service"
)

type kafkaService struct {
	Hi string `json:"hi"`
}

func (s *kafkaService) Result() (service.Messages, error) {
	fmt.Println("Result, called", s)

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

	fmt.Println("Result, complete", messages)

	return messages, nil
}
