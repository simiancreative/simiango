package natsjscm

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/mock"
)

type JetStreamMsg struct {
	mock.Mock
}

// Metadata returns [MsgMetadata] for a JetStream message.
func (m *JetStreamMsg) Metadata() (*jetstream.MsgMetadata, error) {
	called := m.Called()
	return called.Get(0).(*jetstream.MsgMetadata), called.Error(1)
}

// Data returns the message body.
func (m *JetStreamMsg) Data() []byte {
	called := m.Called()
	return called.Get(0).([]byte)
}

// Headers returns a map of headers for a message.
func (m *JetStreamMsg) Headers() nats.Header {
	called := m.Called()
	return called.Get(0).(nats.Header)
}

// Subject returns a subject on which a message is published.
func (m *JetStreamMsg) Subject() string {
	called := m.Called()
	return called.String(0)
}

// Reply returns a reply subject for a JetStream message.
func (m *JetStreamMsg) Reply() string {
	called := m.Called()
	return called.String(0)
}

// Ack acknowledges a message. This tells the server that the message was
// successfully processed and it can move on to the next message.
func (m *JetStreamMsg) Ack() error {
	call := m.Called()
	return call.Error(0)
}

// DoubleAck acknowledges a message and waits for ack reply from the server.
// While it impacts performance, it is useful for scenarios where
// message loss is not acceptable.
func (m *JetStreamMsg) DoubleAck(ctx context.Context) error {
	call := m.Called(ctx)
	return call.Error(0)
}

// Nak negatively acknowledges a message. This tells the server to
// redeliver the message.
func (m *JetStreamMsg) Nak() error {
	call := m.Called()
	return call.Error(0)
}

// NakWithDelay negatively acknowledges a message. This tells the server
// to redeliver the message after the given delay.
func (m *JetStreamMsg) NakWithDelay(delay time.Duration) error {
	called := m.Called(delay)
	return called.Error(0)
}

// InProgress tells the server that this message is being worked on. It
// resets the redelivery timer on the server.
func (m *JetStreamMsg) InProgress() error {
	called := m.Called()
	return called.Error(0)
}

// Term tells the server to not redeliver this message, regardless of
// the value of MaxDeliver.
func (m *JetStreamMsg) Term() error {
	called := m.Called()
	return called.Error(0)
}

// TermWithReason tells the server to not redeliver this message, regardless of
// the value of MaxDeliver. The provided reason will be included in JetStream
// advisory event sent by the server.
//
// Note: This will only work with JetStream servers >= 2.10.4.
// For older servers, TermWithReason will be ignored by the server and the message
// will not be terminated.
func (m *JetStreamMsg) TermWithReason(reason string) error {
	called := m.Called(reason)
	return called.Error(0)
}
