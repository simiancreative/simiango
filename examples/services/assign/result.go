package assign

import (
	"github.com/simiancreative/simiango/service"
)

type assignService struct {
	rawEvent []byte
	event    Event
}

func (s *assignService) Result() (service.Messages, error) {
	event := s.event
	assignable := AssignableStruct{}
	event.Body.Assign(assignable)

	messages := make(service.Messages, 2)
	return messages, nil
}
