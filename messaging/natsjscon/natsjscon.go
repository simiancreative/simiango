package natsjscon

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/simiancreative/simiango/logger"
	"github.com/simiancreative/simiango/messaging/natsjscm"
	"github.com/simiancreative/simiango/messaging/natsjsdlq"
)

// Logger is the interface for logging
type Logger interface {
	Debug(args ...any)
	Debugf(format string, args ...any)
	Error(args ...any)
}

// ProcessStatus represents the result of message processing
type ProcessStatus int

const (
	// Success indicates successful processing
	Success ProcessStatus = iota
	// Failure indicates processing failure (will be retried)
	Failure
	// TerminalFailure indicates a non-recoverable failure (sent to DLQ)
	TerminalFailure
)

var ResultNames = map[ProcessStatus]string{
	Success:         "SUCCESS",
	Failure:         "FAILURE",
	TerminalFailure: "TERMINAL-FAILURE",
}

// Processor processes messages
type Processor func(ctx context.Context, msgs []jetstream.Msg) map[jetstream.Msg]ProcessStatus

// ConsumerConfig holds general configuration for a consumer
type ConsumerConfig struct {
	// StreamName is the name of the stream to consume from
	StreamName string

	// ConsumerName is the name for this consumer
	ConsumerName string

	// Subject is the subject to subscribe to
	Subject string

	// MaxRetries is the maximum number of retries before sending to DLQ
	MaxRetries int

	// ProcessTimeout is the timeout for processing a message/batch
	ProcessTimeout time.Duration

	// WorkerCount is the number of concurrent workers
	WorkerCount int
}

// ConsumptionStrategy defines the interface for different message consumption strategies
type ConsumptionStrategy interface {
	// Setup prepares the consumption strategy
	Setup(ctx context.Context) error

	// Consume retrieves messages for processing
	Consume(ctx context.Context, workerID int) ([]jetstream.Msg, error)
}

// Consumer manages consuming and processing messages
type Consumer struct {
	// Configuration
	config ConsumerConfig

	// Dependencies
	logger     Logger
	strategy   ConsumptionStrategy
	processor  Processor
	cm         natsjscm.Connector
	dlqHandler natsjsdlq.Handler

	// State
	running bool
	wg      sync.WaitGroup
	mu      sync.RWMutex
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewConsumer creates a new consumer with the specified strategy
func NewConsumer(config ConsumerConfig) *Consumer {
	return &Consumer{config: config}
}

func (c *Consumer) debug(args ...any) {
	if c.logger == nil {
		return
	}

	c.logger.Debug(args...)
}

// SetLogger sets the logger
func (c *Consumer) SetLogger(logger Logger) *Consumer {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.logger = logger
	c.debug("logger set")
	return c
}

// SetConnector sets the connection manager
func (c *Consumer) SetConnector(cm natsjscm.Connector) *Consumer {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cm = cm
	c.debug("connection manager set")
	return c
}

// SetStrategy sets the consumption strategy
func (c *Consumer) SetStrategy(strategy ConsumptionStrategy) *Consumer {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.strategy = strategy
	c.debug("strategy set")
	return c
}

// SetDLQHandler sets the dead letter queue handler
func (c *Consumer) SetDLQHandler(handler natsjsdlq.Handler) *Consumer {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.dlqHandler = handler
	c.debug("DLQ handler set")
	return c
}

// SetProcessor sets the message processor
func (c *Consumer) SetProcessor(processor Processor) *Consumer {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.processor = processor
	c.debug("processor set")
	return c
}

// Stop stops the consumer
func (c *Consumer) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running {
		return
	}

	c.logger.Debug("stopping consumer", logger.Fields{
		"stream":   c.config.StreamName,
		"consumer": c.config.ConsumerName,
	})

	c.cancel()
	c.running = false
	c.wg.Wait()
}

// validate checks if the consumer is properly configured
func (c *Consumer) validate() error {
	if c.processor == nil {
		return errors.New("processor is required")
	}

	if c.cm == nil {
		return errors.New("connection manager is required")
	}

	if c.strategy == nil {
		return errors.New("consumption strategy is required")
	}

	if c.config.StreamName == "" {
		return errors.New("stream name is required")
	}

	if c.config.ConsumerName == "" {
		return errors.New("consumer name is required")
	}

	if c.config.Subject == "" {
		return errors.New("subject is required")
	}

	return nil
}

func (c *Consumer) setup() error {
	c.debug("starting consumer setup")

	if c.logger == nil {
		c.debug("setting logger")
		c.SetLogger(logger.New())
	}

	// Set defaults
	if c.config.MaxRetries <= 0 {
		c.debug("setting default max retries to 3")
		c.config.MaxRetries = 3
	}

	if c.config.ProcessTimeout <= 0 {
		c.debug("setting default process timeout to 30s")
		c.config.ProcessTimeout = 30 * time.Second
	}

	if c.config.WorkerCount <= 0 {
		c.debug("setting default worker count to 1")
		c.config.WorkerCount = 1
	}

	c.debug("consumer setup complete, connecting to NATS")

	// Ensure connection to NATS
	return c.cm.Connect()
}
