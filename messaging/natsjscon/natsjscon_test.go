package natsjscon_test

import (
	"context"
	"testing"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/simiancreative/simiango/messaging/natsjscon"
	"github.com/simiancreative/simiango/mocks/messaging/natsjscm"
	conmock "github.com/simiancreative/simiango/mocks/messaging/natsjscon"
	"github.com/stretchr/testify/assert"
)

func TestNewConsumer(t *testing.T) {
	d := conmock.NewDependencies()
	c := d.NewConsumer(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := c.Start(ctx)
	assert.NoError(t, err)

	// wait
	time.Sleep(500 * time.Millisecond)

	// does not panic
	assert.NotPanics(t, func() {
		c.Stop()
	}, "Stop should not panic")
}

func TestNewConsumerWithInvalidConfig(t *testing.T) {
	d := conmock.NewDependencies()

	tests := []struct {
		name        string
		fn          func(c *natsjscon.Consumer)
		shouldError bool
	}{
		{
			name: "SetProcessor",
			fn: func(c *natsjscon.Consumer) {
				c.SetProcessor(nil)
			},
			shouldError: true,
		},
		{
			name: "SetConnector",
			fn: func(c *natsjscon.Consumer) {
				c.SetConnector(nil)
			},
			shouldError: true,
		},
		{
			name: "ValidateStream",
			fn: func(c *natsjscon.Consumer) {
				*c = *d.NewConsumer(t, natsjscon.ConsumerConfig{
					StreamName:   "",
					ConsumerName: "test-consumer",
					Subject:      "test.subject.>",
				})
			},
			shouldError: true,
		},
		{
			name: "ValidateConsumerName",
			fn: func(c *natsjscon.Consumer) {
				*c = *d.NewConsumer(t, natsjscon.ConsumerConfig{
					StreamName:   "test-stream",
					ConsumerName: "",
					Subject:      "test.subject.>",
				})
			},
			shouldError: true,
		},
		{
			name: "ValidateSubject",
			fn: func(c *natsjscon.Consumer) {
				*c = *d.NewConsumer(t, natsjscon.ConsumerConfig{
					StreamName:   "test-stream",
					ConsumerName: "test-consumer",
					Subject:      "",
				})
			},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		c := d.NewConsumer(t)

		t.Run(tt.name, func(t *testing.T) {
			tt.fn(c)

			if tt.shouldError {
				assert.Error(t, c.Start(context.Background()))
			}

			if !tt.shouldError {
				assert.NoError(t, c.Start(context.Background()))
			}
		})
	}
}

func TestHandleResult(t *testing.T) {
	d := conmock.NewDependencies()
	c := d.NewConsumer(t)

	badAck := natsjscm.NewMockJetStreamMsg(jetstream.ErrMsgAlreadyAckd)

	batch := natsjscm.NewMockMessageBatch(
		[]jetstream.Msg{badAck},
		nil,
	)

	d.Strategy.Reset(batch)

	err := c.Start(context.Background())
	assert.NoError(t, err)

	assert.Eventually(t, func() bool {
		return d.Processor.AssertExpectations(t) && badAck.AssertExpectations(t)
	}, 500*time.Millisecond, 100*time.Millisecond)
}
