package natsjscon_test

import (
	"context"
	"testing"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/simiancreative/simiango/messaging/natsjscon"
	"github.com/simiancreative/simiango/mocks/logger"
	"github.com/simiancreative/simiango/mocks/messaging/natsjscm"
	conmock "github.com/simiancreative/simiango/mocks/messaging/natsjscon"
	"github.com/simiancreative/simiango/mocks/messaging/natsjsdlq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type dependencies struct {
	logger     *logger.MockLogger
	jetstream  *natsjscm.MockJetStream
	connector  *natsjscm.MockConnectionManager
	dlqHandler *natsjsdlq.MockDLQHandler
	strategy   *conmock.MockStrategy
	processor  *conmock.ProcessorMock
}

func (d *dependencies) newConsumer(t *testing.T) *natsjscon.Consumer {
	config := natsjscon.ConsumerConfig{
		StreamName:   "test",
		ConsumerName: "test-consumer",
		Subject:      "test.subject.>",
	}

	d.logger = &logger.MockLogger{}
	d.jetstream = &natsjscm.MockJetStream{}
	d.connector = &natsjscm.MockConnectionManager{}
	d.dlqHandler = &natsjsdlq.MockDLQHandler{}
	d.strategy = &conmock.MockStrategy{}
	d.processor = &conmock.ProcessorMock{}

	d.connector.SetJetStream(d.jetstream)

	d.logger.On("Debug", mock.Anything)
	d.connector.On("Connect", mock.Anything).Return(nil)

	consumer := natsjscon.
		NewConsumer(config).
		SetLogger(d.logger).
		SetConnector(d.connector).
		SetDLQHandler(d.dlqHandler).
		SetStrategy(d.strategy).
		SetProcessor(d.processor.Process)

	assert.NotNil(t, consumer)

	return consumer
}

func TestNewConsumer(t *testing.T) {
	d := new(dependencies)
	c := d.newConsumer(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	d.strategy.On("Setup", mock.Anything).Return(nil)
	d.strategy.
		On("Consume", mock.Anything, mock.Anything).
		Return([]jetstream.Msg{}, nil)

	err := c.Start(ctx)
	assert.NoError(t, err)

	// wait
	time.Sleep(500 * time.Millisecond)
}
