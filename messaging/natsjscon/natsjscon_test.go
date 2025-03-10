package natsjscon_test

import (
	"context"
	"testing"
	"time"

	"github.com/simiancreative/simiango/messaging/natsjscon"
	"github.com/simiancreative/simiango/mocks/logger"
	"github.com/simiancreative/simiango/mocks/messaging/natsjscm"
	conmock "github.com/simiancreative/simiango/mocks/messaging/natsjscon"
	"github.com/simiancreative/simiango/mocks/messaging/natsjsdlq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newConsumer(t *testing.T) *natsjscon.Consumer {
	logger := &logger.MockLogger{}
	jetstream := &natsjscm.MockJetStream{}
	connector := &natsjscm.MockConnectionManager{}
	config := natsjscon.ConsumerConfig{}
	dlqHandler := &natsjsdlq.MockDLQHandler{}
	strategy := &conmock.MockStrategy{}
	processor := &conmock.ProcessorMock{}

	connector.SetJetStream(jetstream)

	logger.On("Debug", mock.Anything)

	consumer := natsjscon.
		NewConsumer(config).
		SetLogger(logger).
		SetConnector(connector).
		SetDLQHandler(dlqHandler).
		SetStrategy(strategy).
		SetProcessor(processor.Process)

	assert.NotNil(t, consumer)

	return consumer
}

func TestNewConsumer(t *testing.T) {
	c := newConsumer(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := c.Start(ctx)
	assert.NoError(t, err)

	// wait
	time.Sleep(500 * time.Millisecond)
}
