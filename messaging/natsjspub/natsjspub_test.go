package natsjspub_test

import (
	"context"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/simiancreative/simiango/messaging/natsjspub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewPublisher(t *testing.T) {
	mockCM := new(MockConnectionManager)
	config := natsjspub.Config{
		StreamName: "test-stream",
		Subject:    "test-subject",
		Timeout:    5 * time.Second,
	}

	t.Run("success", func(t *testing.T) {
		mockCM.
			On("EnsureStream", mock.Anything, mock.Anything).
			Return(new(MockJetStream), nil)

		mockCM.
			On("IsConnected").
			Return(true)

		deps := natsjspub.Dependencies{
			Connector: mockCM,
		}

		pub, err := natsjspub.NewPublisher(deps, config)
		assert.NoError(t, err)
		assert.NotNil(t, pub)
	})

	t.Run("missing connection manager", func(t *testing.T) {
		deps := natsjspub.Dependencies{}

		pub, err := natsjspub.NewPublisher(deps, config)
		assert.Error(t, err)
		assert.Nil(t, pub)
	})

	t.Run("missing stream name", func(t *testing.T) {
		config.StreamName = ""
		deps := natsjspub.Dependencies{
			Connector: mockCM,
		}

		pub, err := natsjspub.NewPublisher(deps, config)
		assert.Error(t, err)
		assert.Nil(t, pub)
	})
}

type Mocks struct {
	cm *MockConnectionManager
	js *MockJetStream
	cb *MockCircuitBreaker
}

func createPublisher(t *testing.T) (natsjspub.PublishMulti, *Mocks) {
	mocks := &Mocks{
		cm: new(MockConnectionManager),
		js: new(MockJetStream),
		cb: new(MockCircuitBreaker),
	}

	mocks.cm.SetJetStream(mocks.js)

	config := natsjspub.Config{
		StreamName: "test-stream",
		Subject:    "test-subject",
		Timeout:    5 * time.Second,
	}

	// new publisher sets up the conection manager and stream
	mocks.cm.On("IsConnected").Return(true)
	mocks.cm.On("Connect").Return(nil)
	mocks.cm.On("EnsureStream", mock.Anything, mock.Anything).Return(mocks.js, nil)

	pub, err := natsjspub.NewPublisher(natsjspub.Dependencies{
		Connector: mocks.cm,
		Breaker:   mocks.cb,
	}, config)
	assert.NoError(t, err)

	return pub, mocks
}

func TestPublish(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		pub, mocks := createPublisher(t)
		mocks.cb.On("Allow").Return(true)
		mocks.cb.On("RecordStart").Return(true)
		mocks.cm.On("GetJetStream").Return(mocks.js)
		mocks.js.
			On("PublishMsg", mock.Anything, mock.Anything, mock.Anything).
			Return(&jetstream.PubAck{}, nil)
		mocks.cb.On("RecordResult", mock.Anything)

		msg := &nats.Msg{Subject: "test-subject", Data: []byte("test")}
		ack, err := pub.Publish(context.Background(), msg)
		assert.NoError(t, err)
		assert.NotNil(t, ack)
	})

	t.Run("circuit breaker open", func(t *testing.T) {
		pub, mocks := createPublisher(t)

		mocks.cb.On("Allow").Return(false)

		msg := &nats.Msg{Subject: "test-subject", Data: []byte("test")}
		ack, err := pub.Publish(context.Background(), msg)
		assert.Error(t, err)
		assert.Nil(t, ack)
	})

	t.Run("jetstream connection not available", func(t *testing.T) {
		pub, mocks := createPublisher(t)

		mocks.cb.On("Allow").Return(true)
		mocks.cb.On("RecordStart").Return(true)
		mocks.cm.On("GetJetStream").Return(nil)
		mocks.cb.On("RecordResult", mock.Anything)

		msg := &nats.Msg{Subject: "test-subject", Data: []byte("test")}
		ack, err := pub.Publish(context.Background(), msg)
		assert.Error(t, err)
		assert.Nil(t, ack)
	})
}

func TestPublishJSON(t *testing.T) {
	t.Run("json marshal error", func(t *testing.T) {
		pub, _ := createPublisher(t)
		data := make(chan int) // non-serializable type
		ack, err := pub.PublishJSON(context.Background(), data)
		assert.Error(t, err)
		assert.Nil(t, ack)
	})
}
