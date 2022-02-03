package kafka

import (
	"fmt"

	"github.com/simiancreative/simiango/service"
)

type kafkaService struct {
}

func (s *kafkaService) Result() (service.Messages, error) {
	fmt.Println("Result, messages")
	messages := make(service.Messages, 1)
	messages[0] = service.Message{
		Key: "example",
		Value: struct {
			field string
		}{
			field: "hello",
		},
	}
	return messages, nil
}
