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
	Msgs []jetstream.Msg
	Err  error
}

// NewMockMessageBatch creates a new instance of MockMessageBatch.
func NewMockMessageBatch(messages []jetstream.Msg, err error) *MockMessageBatch {
	return &MockMessageBatch{
		Msgs: messages,
		Err:  err,
	}
}

// Messages returns a channel of jetstream.Msg.
func (m *MockMessageBatch) Messages() <-chan jetstream.Msg {
	msgChan := make(chan jetstream.Msg, len(m.Msgs))
	for _, msg := range m.Msgs {
		msgChan <- msg
	}
	close(msgChan)
	return msgChan
}

// Error returns an error.
func (m *MockMessageBatch) Error() error {
	return m.Err
}
