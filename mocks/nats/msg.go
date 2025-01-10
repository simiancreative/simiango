package nats

import (
	"context"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

var NewJetStreamMsg = func(subject string, data []byte) *jetStreamMsg {
	return &jetStreamMsg{
		msg: &nats.Msg{
			Subject: subject,
			Data:    data,
		},
	}
}

type jetStreamMsg struct {
	msg  *nats.Msg
	ackd bool
	js   *jetstream.JetStream
	sync.Mutex
}

func (m *jetStreamMsg) Metadata() (*jetstream.MsgMetadata, error) {
	return nil, nil
}

// Data returns the message body.
func (m *jetStreamMsg) Data() []byte {
	return m.msg.Data
}

// Headers returns a map of headers for a message.
func (m *jetStreamMsg) Headers() nats.Header {
	return m.msg.Header
}

// Subject returns a subject on which a message is published.
func (m *jetStreamMsg) Subject() string {
	return m.msg.Subject
}

// Reply returns a reply subject for a JetStream message.
func (m *jetStreamMsg) Reply() string {
	return m.msg.Reply
}

func (m *jetStreamMsg) Ack() error {
	return nil
}

func (m *jetStreamMsg) DoubleAck(_ context.Context) error {
	return nil
}

func (m *jetStreamMsg) Nak() error {
	return nil
}

func (m *jetStreamMsg) NakWithDelay(_ time.Duration) error {
	return nil
}

func (m *jetStreamMsg) InProgress() error {
	return nil
}

func (m *jetStreamMsg) Term() error {
	return nil
}

func (m *jetStreamMsg) TermWithReason(_ string) error {
	return nil
}

type ackOpts struct {
	nakDelay   time.Duration
	termReason string
}

func (m *jetStreamMsg) ackReply(
	_ context.Context,
	_ []byte,
	_ bool,
	_ ackOpts,
) error {
	return nil
}

func (m *jetStreamMsg) checkReply() error {
	return nil
}
