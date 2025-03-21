package natsjsdlq

import (
	"github.com/nats-io/nats.go"
	"github.com/simiancreative/simiango/messaging/natsjsdlq"
	"github.com/stretchr/testify/mock"
)

type MockDLQHandler struct {
	mock.Mock
}

func (m *MockDLQHandler) PublishMessage(msg *nats.Msg, reason string) error {
	args := m.Called(msg, reason)
	return args.Error(0)
}

func (m *MockDLQHandler) ShouldDLQ(msg natsjsdlq.Msg) bool {
	args := m.Called(msg)
	return args.Bool(0)
}
