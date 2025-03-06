package natsjscon_test

import (
	"context"
	"errors"
	"testing"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/simiancreative/simiango/messaging/natsjscon"
	"github.com/stretchr/testify/assert"
)

// Mock implementations
type MockLogger struct{}

func (m *MockLogger) Debug(args ...any) {}
func (m *MockLogger) Error(args ...any) {}

type MockConnector struct {
	shouldFail bool
}

func (m *MockConnector) Connect() error {
	if m.shouldFail {
		return errors.New("connection failed")
	}
	return nil
}

type MockStrategy struct{}

func (m *MockStrategy) Setup(ctx context.Context) error {
	return nil
}

func (m *MockStrategy) Consume(ctx context.Context, workerID int) ([]jetstream.Msg, error) {
	return nil, nil
}

type MockDLQHandler struct{}

func (m *MockDLQHandler) Handle(msg jetstream.Msg) error {
	return nil
}

func TestNewConsumer(t *testing.T) {
	config := natsjscon.ConsumerConfig{}
	consumer := natsjscon.NewConsumer(config)
	assert.NotNil(t, consumer)
}

func TestSetLogger(t *testing.T) {
	consumer := natsjscon.NewConsumer(natsjscon.ConsumerConfig{})
	logger := &MockLogger{}
	c := consumer.SetLogger(logger)
	assert.NotNil(t, c)
}

func TestSetConnector(t *testing.T) {
	consumer := natsjscon.NewConsumer(natsjscon.ConsumerConfig{})
	connector := &MockConnector{}
	c := consumer.SetConnector(connector)
	assert.Equal(t, connector, consumer.cm)
}

//
// func TestSetStrategy(t *testing.T) {
// 	consumer := natsjscon.NewConsumer(natsjscon.ConsumerConfig{})
// 	strategy := &MockStrategy{}
// 	consumer.SetStrategy(strategy)
// 	assert.Equal(t, strategy, consumer.strategy)
// }
//
// func TestSetDLQHandler(t *testing.T) {
// 	consumer := natsjscon.NewConsumer(natsjscon.ConsumerConfig{})
// 	handler := &MockDLQHandler{}
// 	consumer.SetDLQHandler(handler)
// 	assert.Equal(t, handler, consumer.dlqHandler)
// }
//
// func TestSetProcessor(t *testing.T) {
// 	consumer := natsjscon.NewConsumer(natsjscon.ConsumerConfig{})
// 	processor := func(ctx context.Context, msgs []jetstream.Msg) map[jetstream.Msg]natsjscon.ProcessStatus {
// 		return nil
// 	}
// 	consumer.SetProcessor(processor)
// 	assert.Equal(t, processor, consumer.processor)
// }
//
// func TestSetup(t *testing.T) {
// 	consumer := natsjscon.NewConsumer(natsjscon.ConsumerConfig{})
// 	consumer.SetLogger(&MockLogger{})
// 	consumer.SetConnector(&MockConnector{})
// 	consumer.SetStrategy(&MockStrategy{})
// 	consumer.SetDLQHandler(&MockDLQHandler{})
// 	processor := func(ctx context.Context, msgs []jetstream.Msg) map[jetstream.Msg]natsjscon.ProcessStatus {
// 		return nil
// 	}
// 	consumer.Setup(processor)
// 	assert.NotNil(t, consumer.logger)
// 	assert.NotNil(t, consumer.cm)
// 	assert.NotNil(t, consumer.strategy)
// 	assert.NotNil(t, consumer.dlqHandler)
// 	assert.Equal(t, processor, consumer.processor)
// }
//
// func TestIsRunning(t *testing.T) {
// 	consumer := natsjscon.NewConsumer(natsjscon.ConsumerConfig{})
// 	assert.False(t, consumer.IsRunning())
// 	consumer.running = true
// 	assert.True(t, consumer.IsRunning())
// }
//
// func TestStop(t *testing.T) {
// 	consumer := natsjscon.NewConsumer(natsjscon.ConsumerConfig{})
// 	consumer.SetLogger(&MockLogger{})
// 	consumer.running = true
// 	ctx, cancel := context.WithCancel(context.Background())
// 	consumer.ctx = ctx
// 	consumer.cancel = cancel
// 	err := consumer.Stop()
// 	assert.NoError(t, err)
// 	assert.False(t, consumer.running)
// }
//
// func TestValidate(t *testing.T) {
// 	consumer := natsjscon.NewConsumer(natsjscon.ConsumerConfig{})
// 	consumer.SetLogger(&MockLogger{})
// 	consumer.SetConnector(&MockConnector{})
// 	consumer.SetStrategy(&MockStrategy{})
// 	consumer.SetDLQHandler(&MockDLQHandler{})
// 	processor := func(ctx context.Context, msgs []jetstream.Msg) map[jetstream.Msg]natsjscon.ProcessStatus {
// 		return nil
// 	}
// 	consumer.SetProcessor(processor)
//
// 	err := consumer.Validate()
// 	assert.Error(t, err)
//
// 	consumer.config.StreamName = "test-stream"
// 	consumer.config.ConsumerName = "test-consumer"
// 	consumer.config.Subject = "test-subject"
//
// 	err = consumer.Validate()
// 	assert.NoError(t, err)
// }
