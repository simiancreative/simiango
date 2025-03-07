package natsjscon

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/simiancreative/simiango/messaging/natsjscon"
	"github.com/stretchr/testify/mock"
)

type MockStrategy struct {
	mock.Mock
}

func (m *MockStrategy) Setup(ctx context.Context) error {
	c := m.Called(ctx)
	return c.Error(0)
}

func (m *MockStrategy) Consume(ctx context.Context, workerID int) ([]jetstream.Msg, error) {
	c := m.Called(ctx, workerID)
	return c.Get(0).([]jetstream.Msg), c.Error(1)
}

type ProcessorMock struct {
	mock.Mock
}

func (p *ProcessorMock) Process(
	ctx context.Context,
	msgs []jetstream.Msg,
) map[jetstream.Msg]natsjscon.ProcessStatus {
	args := p.Called(ctx, msgs)
	return args.Get(0).(map[jetstream.Msg]natsjscon.ProcessStatus)
}
