package natsjscm

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/mock"
)

// MockConsumer is a stretcher-based mock for the jetstream.Consumer interface
type MockConsumer struct {
	mock.Mock
}

// Fetch mocks the Consumer.Fetch method
func (m *MockConsumer) Fetch(
	batch int,
	opts ...jetstream.FetchOpt,
) (jetstream.MessageBatch, error) {
	args := []interface{}{batch}
	for _, opt := range opts {
		args = append(args, opt)
	}
	result := m.Called(args...)

	var messageBatch jetstream.MessageBatch
	if result.Get(0) != nil {
		messageBatch = result.Get(0).(jetstream.MessageBatch)
	}

	return messageBatch, result.Error(1)
}

// FetchBytes mocks the Consumer.FetchBytes method
func (m *MockConsumer) FetchBytes(
	maxBytes int,
	opts ...jetstream.FetchOpt,
) (jetstream.MessageBatch, error) {
	args := []interface{}{maxBytes}
	for _, opt := range opts {
		args = append(args, opt)
	}
	result := m.Called(args...)

	var messageBatch jetstream.MessageBatch
	if result.Get(0) != nil {
		messageBatch = result.Get(0).(jetstream.MessageBatch)
	}

	return messageBatch, result.Error(1)
}

// FetchNoWait mocks the Consumer.FetchNoWait method
func (m *MockConsumer) FetchNoWait(batch int) (jetstream.MessageBatch, error) {
	result := m.Called(batch)

	var messageBatch jetstream.MessageBatch
	if result.Get(0) != nil {
		messageBatch = result.Get(0).(jetstream.MessageBatch)
	}

	return messageBatch, result.Error(1)
}

// Consume mocks the Consumer.Consume method
func (m *MockConsumer) Consume(
	handler jetstream.MessageHandler,
	opts ...jetstream.PullConsumeOpt,
) (jetstream.ConsumeContext, error) {
	args := []interface{}{handler}
	for _, opt := range opts {
		args = append(args, opt)
	}
	result := m.Called(args...)

	var consumeContext jetstream.ConsumeContext
	if result.Get(0) != nil {
		consumeContext = result.Get(0).(jetstream.ConsumeContext)
	}

	return consumeContext, result.Error(1)
}

// Messages mocks the Consumer.Messages method
func (m *MockConsumer) Messages(
	opts ...jetstream.PullMessagesOpt,
) (jetstream.MessagesContext, error) {
	args := []interface{}{}
	for _, opt := range opts {
		args = append(args, opt)
	}
	result := m.Called(args...)

	var messagesContext jetstream.MessagesContext
	if result.Get(0) != nil {
		messagesContext = result.Get(0).(jetstream.MessagesContext)
	}

	return messagesContext, result.Error(1)
}

// Next mocks the Consumer.Next method
func (m *MockConsumer) Next(opts ...jetstream.FetchOpt) (jetstream.Msg, error) {
	args := []interface{}{}
	for _, opt := range opts {
		args = append(args, opt)
	}
	result := m.Called(args...)

	var msg jetstream.Msg
	if result.Get(0) != nil {
		msg = result.Get(0).(jetstream.Msg)
	}

	return msg, result.Error(1)
}

// Info mocks the Consumer.Info method
func (m *MockConsumer) Info(ctx context.Context) (*jetstream.ConsumerInfo, error) {
	result := m.Called(ctx)

	var consumerInfo *jetstream.ConsumerInfo
	if result.Get(0) != nil {
		consumerInfo = result.Get(0).(*jetstream.ConsumerInfo)
	}

	return consumerInfo, result.Error(1)
}

// CachedInfo mocks the Consumer.CachedInfo method
func (m *MockConsumer) CachedInfo() *jetstream.ConsumerInfo {
	result := m.Called()

	var consumerInfo *jetstream.ConsumerInfo
	if result.Get(0) != nil {
		consumerInfo = result.Get(0).(*jetstream.ConsumerInfo)
	}

	return consumerInfo
}

// Ensure the mock implements the interface
var _ jetstream.Consumer = (*MockConsumer)(nil)
