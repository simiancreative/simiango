package natsjscm

// type MockMessageBatch interface {
// 	Messages() <-chan jetstream.Msg
// 	Error() error
// }

import (
	"github.com/nats-io/nats.go/jetstream"
)

// MockMessageBatch is a mock implementation of the MockMessageBatch interface.
type MockMessageBatch struct {
	messages chan jetstream.Msg
	err      error
}

// NewMockMessageBatch creates a new instance of MockMessageBatch.
func NewMockMessageBatch(messages []jetstream.Msg, err error) *MockMessageBatch {
	msgChan := make(chan jetstream.Msg, len(messages))
	for _, msg := range messages {
		msgChan <- msg
	}
	close(msgChan)
	return &MockMessageBatch{
		messages: msgChan,
		err:      err,
	}
}

// Messages returns a channel of jetstream.Msg.
func (m *MockMessageBatch) Messages() <-chan jetstream.Msg {
	return m.messages
}

// Error returns an error.
func (m *MockMessageBatch) Error() error {
	return m.err
}
