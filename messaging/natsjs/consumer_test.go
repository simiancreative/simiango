package natsjs_test

import (
	"testing"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/simiancreative/simiango/errors"
	"github.com/simiancreative/simiango/logger"
	"github.com/simiancreative/simiango/messaging/natsjs"
	"github.com/simiancreative/simiango/mocks/nats"
	"github.com/stretchr/testify/assert"
)

func testNewConsumer(t *testing.T) {
	defer nats.MockServer()()

	assert.NotPanics(t, func() {
		natsjs.SetTimeout(30)
		natsjs.Connect()

		consumer := natsjs.NewConsumer("test", "test", "test")
		assert.NotNil(t, consumer)
	})
}

func testNewConsumerFailure(t *testing.T) {
	assert.Panics(t, func() {
		natsjs.SetTimeout(1)
		natsjs.NewConsumer("test", "test", "test")
	})
}

func testNewConsumerStop(t *testing.T) {
	logger.Enable()
	defer nats.MockServer()()

	natsjs.SetTimeout(30)
	natsjs.Connect()

	stop, err := natsjs.
		NewConsumer("test", "test", "test.>").
		Consume(func(_ jetstream.Msg) error {
			return nil
		})

	assert.NoError(t, err)
	assert.NotNil(t, stop)
	assert.NotPanics(t, stop)
}

func testNewConsumerConsume(t *testing.T) {
	logger.Enable()
	defer nats.MockServer()()

	result := make(chan string)

	natsjs.SetTimeout(30)
	natsjs.Connect()

	stop, err := natsjs.
		NewConsumer("test", "test", "test.>").
		Consume(func(msg jetstream.Msg) error {
			logger.Debugf("Received message: %v", msg)
			result <- string(msg.Data())
			return nil
		})

	assert.NoError(t, err)
	assert.NotNil(t, stop)

	natsjs.New().
		NewMessage().
		SetSubject("test", "test.test").
		SetData("test data").
		Publish()

	select {
	case value := <-result:
		assert.Equal(t, value, `"test data"`)
	case <-time.After(500 * time.Millisecond):
		t.Errorf("Timeout waiting for async result")
	}
}

func testNewConsumerConsumeFailure(t *testing.T) {
	hook := logger.Mock()
	defer nats.MockServer()()

	natsjs.SetTimeout(30)
	natsjs.Connect()

	result := make(chan string)

	callback := func(msg jetstream.Msg) error {
		result <- string(msg.Data())
		return errors.New("test error")
	}

	natsjs.
		NewConsumer("test", "test", "test.>").
		Consume(callback)

	natsjs.New().
		NewMessage().
		SetSubject("test", "test.test").
		SetData("test data").
		Publish()

	select {
	case <-result:
		time.Sleep(100 * time.Millisecond)
	case <-time.After(500 * time.Millisecond):
		t.Errorf("Timeout waiting for async result")
	}

	assert.Equal(t, `Acking message: "test data"`, hook.LastEntry().Message)
}

func testNewConsumerNoConnection(t *testing.T) {
	natsjs.Close()

	assert.Panics(t, func() {
		natsjs.
			NewConsumer("test", "test", "test.>").
			Consume(func(msg jetstream.Msg) error {
				logger.Debugf("Received message: %v", msg)
				return errors.New("test error")
			})
	})
}
