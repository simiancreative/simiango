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

func newDependencies() *dependencies {
	d := &dependencies{
		logger:     &logger.MockLogger{},
		jetstream:  &natsjscm.MockJetStream{},
		connector:  &natsjscm.MockConnectionManager{},
		dlqHandler: &natsjsdlq.MockDLQHandler{},
		strategy:   &conmock.MockStrategy{},
		processor:  &conmock.ProcessorMock{},
	}

	msg := &natsjscm.JetStreamMsg{}
	msg.On("Metadata").Return(new(jetstream.MsgMetadata), nil)
	msg.On("Ack").Return(nil)
	msgs := []jetstream.Msg{msg, msg, msg, msg}

	processed := make(map[jetstream.Msg]natsjscon.ProcessStatus)
	for _, msg := range msgs {
		processed[msg] = natsjscon.Success
	}

	d.logger.On("Debug", mock.Anything).Maybe()
	d.logger.On("Debugf", mock.Anything, mock.Anything).Maybe()
	d.logger.On("Error", mock.Anything).Maybe()
	d.connector.On("Connect").Return(nil).Maybe()
	d.processor.On("Process", mock.Anything, mock.Anything).Return(processed).Maybe()

	// Setup default strategy mocks to avoid test failures
	d.strategy.On("Setup", mock.Anything).Return(nil).Maybe()

	d.strategy.On("Consume", mock.Anything, mock.Anything).Return(msgs, nil).Maybe()

	d.connector.SetJetStream(d.jetstream)

	return d
}

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
	d := newDependencies()
	c := d.newConsumer(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := c.Start(ctx)
	assert.NoError(t, err)

	// wait
	time.Sleep(500 * time.Millisecond)
}
