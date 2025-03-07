package natsjspub

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/mock"
)

type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) Publish(ctx context.Context, msg *nats.Msg) (*jetstream.PubAck, error) {
	args := m.Called(ctx, msg)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jetstream.PubAck), args.Error(1)
}

type MockMsg struct {
	mock.Mock
}

func (m *MockMsg) Metadata() (*nats.MsgMetadata, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*nats.MsgMetadata), args.Error(1)
}
