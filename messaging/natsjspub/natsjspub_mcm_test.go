package natsjspub_test

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/simiancreative/simiango/messaging/natsjscm"
	"github.com/stretchr/testify/mock"
)

// Mock dependencies
type MockConnectionManager struct {
	mock.Mock
	jetstream natsjscm.JetStream
}

func (m *MockConnectionManager) SetJetStream(jetstream natsjscm.JetStream) {
	m.jetstream = jetstream
}

func (m *MockConnectionManager) Connect() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockConnectionManager) Disconnect() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockConnectionManager) IsConnected() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockConnectionManager) GetConnection() *nats.Conn {
	args := m.Called()
	return args.Get(0).(*nats.Conn)
}

func (m *MockConnectionManager) GetJetStream() natsjscm.JetStream {
	args := m.Called()
	if js, ok := args.Get(0).(natsjscm.JetStream); ok {
		return js
	}
	return nil
}

func (m *MockConnectionManager) EnsureStream(
	ctx context.Context,
	cfg jetstream.StreamConfig,
) (natsjscm.JetStream, error) {
	args := m.Called(ctx, cfg)
	return args.Get(0).(*MockJetStream), args.Error(1)
}

type MockJetStream struct {
	mock.Mock
}

func (m *MockJetStream) Stream(ctx context.Context, name string) (jetstream.Stream, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(jetstream.Stream), args.Error(1)
}

func (m *MockJetStream) CreateStream(
	ctx context.Context,
	cfg jetstream.StreamConfig,
) (jetstream.Stream, error) {
	args := m.Called(ctx, cfg)
	return args.Get(0).(jetstream.Stream), args.Error(1)
}

func (m *MockJetStream) PublishMsg(
	ctx context.Context,
	msg *nats.Msg,
	opts ...jetstream.PublishOpt,
) (*jetstream.PubAck, error) {
	args := m.Called(ctx, msg, opts)
	return args.Get(0).(*jetstream.PubAck), args.Error(1)
}
