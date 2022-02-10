package rabbit

import (
	"github.com/simiancreative/simiango/messaging/amqp"
)

type rabbitService struct {
	Hello string `json:"hello"`
}

func (s *rabbitService) Result() (interface{}, error) {
	return nil, amqp.Publish(amqp.Publisher{
		Exchange: "hello",
		Queue:    "sweetie",
		Type:     "direct",
		Data:     s,
	})
}
