package natsjscon

import (
	"context"
	"testing"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/simiancreative/simiango/messaging/natsjscon"
	"github.com/simiancreative/simiango/mocks/logger"
	"github.com/simiancreative/simiango/mocks/messaging/natsjscm"
	"github.com/simiancreative/simiango/mocks/messaging/natsjsdlq"
	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
)

func NewDependencies() *Dependencies {
	msg := &natsjscm.MockJetStreamMsg{}
	msg.On("Metadata").Return(new(jetstream.MsgMetadata), nil)
	msg.On("Ack").Return(nil)
	msgs := []jetstream.Msg{msg, msg, msg, msg}

	batch := natsjscm.NewMockMessageBatch(msgs, nil)

	processed := make(map[jetstream.Msg]natsjscon.ProcessStatus)
	for _, msg := range msgs {
		processed[msg] = natsjscon.Success
	}

	d := &Dependencies{
		Logger:     &logger.MockLogger{},
		Jetstream:  &natsjscm.MockJetStream{},
		Connector:  &natsjscm.MockConnectionManager{},
		Stream:     &natsjscm.MockStream{},
		Consumer:   &natsjscm.MockConsumer{},
		Batch:      batch,
		DLQHandler: &natsjsdlq.MockDLQHandler{},
		Strategy:   &MockStrategy{},
		Processor:  &ProcessorMock{},
	}

	d.Logger.On("Debug", mock.Anything).Maybe()
	d.Logger.On("Debugf", mock.Anything, mock.Anything).Maybe()
	d.Logger.On("Error", mock.Anything).Maybe()

	d.Connector.On("Connect").Return(nil).Maybe()
	d.Connector.On("IsConnected").Return(true).Maybe()
	d.Connector.On("EnsureStream", mock.Anything, mock.Anything).Return(d.Jetstream, nil).Maybe()
	d.Connector.On("GetJetStream").Return(d.Jetstream).Maybe()
	d.Connector.SetJetStream(d.Jetstream)

	d.Jetstream.On("Stream", mock.Anything, mock.Anything).
		Return(d.Stream, nil).
		Maybe()

	d.Stream.On("CreateOrUpdateConsumer", mock.Anything, mock.Anything).Return(d.Consumer, nil)
	d.Consumer.On("Fetch", mock.Anything, mock.Anything).Return(batch, nil)

	d.Processor.On("Process", mock.Anything, mock.Anything).Return(processed).Maybe()

	// Setup default strategy mocks to avoid test failures
	d.Strategy.On("Setup", mock.Anything).Return(nil).Maybe()
	d.Strategy.On("Consume", mock.Anything, mock.Anything).Return(msgs, nil).Maybe()

	return d
}

type Dependencies struct {
	Logger     *logger.MockLogger
	Jetstream  *natsjscm.MockJetStream
	Connector  *natsjscm.MockConnectionManager
	Stream     *natsjscm.MockStream
	Consumer   *natsjscm.MockConsumer
	Batch      *natsjscm.MockMessageBatch
	DLQHandler *natsjsdlq.MockDLQHandler
	Strategy   *MockStrategy
	Processor  *ProcessorMock
}

func (d *Dependencies) NewConsumer(
	t *testing.T,
	configs ...natsjscon.ConsumerConfig,
) *natsjscon.Consumer {
	if len(configs) == 0 {
		configs = append(configs, natsjscon.ConsumerConfig{
			StreamName:   "test",
			ConsumerName: "test-consumer",
			Subject:      "test.subject.>",
		})
	}

	config := configs[0]

	consumer := natsjscon.
		NewConsumer(config).
		SetLogger(d.Logger).
		SetConnector(d.Connector).
		SetDLQHandler(d.DLQHandler).
		SetStrategy(d.Strategy).
		SetProcessor(d.Processor.Process)

	assert.NotNil(t, consumer)

	return consumer
}

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
