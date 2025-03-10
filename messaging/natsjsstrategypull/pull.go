package natsjsstrategypull

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/simiancreative/simiango/circuitbreaker"
	"github.com/simiancreative/simiango/logger"
	"github.com/simiancreative/simiango/messaging/natsjscm"
	"github.com/simiancreative/simiango/messaging/natsjscon"
)

type Logger interface {
	Debugf(format string, args ...interface{})
	Debug(args ...any)
	Error(args ...any)
}

// PullStrategyConfig holds configuration specific to pull-based consumption
type Config struct {
	// ConsumerName is the name for this consumer
	ConsumerName string

	// BatchSize for processing
	BatchSize int

	// PollTimeout is how long to wait when polling for new messages
	PollTimeout time.Duration

	// AckWait is how long the server waits for an ack before redelivery
	AckWait time.Duration

	// MaxAckPending is the maximum number of pending acks
	MaxAckPending int

	// MaxRetries is the maximum number of retries
	MaxRetries int

	// RetentionPolicy is the retention policy for the stream
	RetentionPolicy jetstream.RetentionPolicy

	// Breaker is the configuration for the circuit breaker
	Breaker circuitbreaker.Breaker
}

// PullStrategy implements pull-based message consumption
type PullStrategy struct {
	config   Config
	cm       natsjscm.Connector
	cb       circuitbreaker.Breaker
	consumer jetstream.Consumer

	streamName string
	subject    string
	logger     Logger
}

// NewPullStrategy creates a new pull-based consumption strategy
func New(config Config) (*PullStrategy, error) {
	if config.ConsumerName == "" {
		return nil, errors.New("consumer name is required")
	}

	// Set defaults for unspecified configuration
	if config.BatchSize <= 0 {
		config.BatchSize = 100
	}

	if config.PollTimeout <= 0 {
		config.PollTimeout = 1 * time.Second
	}

	if config.AckWait <= 0 {
		config.AckWait = 30 * time.Second
	}

	if config.MaxAckPending <= 0 {
		config.MaxAckPending = 1000
	}

	if config.MaxRetries <= 0 {
		config.MaxRetries = 3
	}

	return &PullStrategy{
		config: config,
	}, nil
}

// Setup creates or updates the JetStream consumer
func (p *PullStrategy) Setup(ctx context.Context) error {
	p.streamName = ctx.Value(natsjscon.CtxKey("stream-name")).(string)
	p.subject = ctx.Value(natsjscon.CtxKey("subject")).(string)
	p.logger = ctx.Value(natsjscon.CtxKey("logger")).(Logger)
	p.cm = ctx.Value(natsjscon.CtxKey("connection-manager")).(natsjscm.Connector)
	p.cb = p.config.Breaker

	err := p.ensureStream(ctx)
	if err != nil {
		return fmt.Errorf("failed to ensure stream: %w", err)
	}

	js := p.cm.GetJetStream()
	if js == nil {
		return errors.New("JetStream connection not available")
	}

	// Get stream
	stream, err := js.Stream(ctx, p.streamName)
	if err != nil {
		return fmt.Errorf("failed to get stream: %w", err)
	}

	// Configure durable pull consumer
	consumerConfig := jetstream.ConsumerConfig{
		Name:          p.config.ConsumerName,
		Durable:       p.config.ConsumerName,
		AckPolicy:     jetstream.AckExplicitPolicy,
		AckWait:       p.config.AckWait,
		MaxAckPending: p.config.MaxAckPending,
		FilterSubject: p.subject,
		DeliverPolicy: jetstream.DeliverAllPolicy,
		MaxDeliver:    p.config.MaxRetries + 1, // Include first delivery
	}

	// Create or update the consumer
	consumer, err := stream.CreateOrUpdateConsumer(ctx, consumerConfig)
	if err != nil {
		return fmt.Errorf("failed to create strategy consumer: %w", err)
	}

	p.consumer = consumer
	return nil
}

// ensureStream makes sure the configured stream exists
func (p *PullStrategy) ensureStream(ctx context.Context) error {
	p.logger.Debugf("ensuring stream is connected")

	// Get JetStream connection
	if !p.cm.IsConnected() {
		p.logger.Debugf("connecting to NATS")
		if err := p.cm.Connect(); err != nil {
			return fmt.Errorf("failed to connect to NATS: %w", err)
		}
	}

	// Set default retention policy if not configured
	retentionPolicy := p.config.RetentionPolicy
	if retentionPolicy == 0 {
		retentionPolicy = jetstream.WorkQueuePolicy
	}

	// Create stream config
	streamConfig := jetstream.StreamConfig{
		Name:      p.streamName,
		Subjects:  []string{fmt.Sprintf("%v.>", p.streamName)},
		Retention: retentionPolicy,
	}

	// Ensure stream exists
	p.logger.Debugf("ensuring stream %s exists", p.streamName)
	_, err := p.cm.EnsureStream(ctx, streamConfig)
	return err
}

// Consume pulls a batch of messages and converts them to our Message type
func (p *PullStrategy) Consume(ctx context.Context, workerID int) ([]jetstream.Msg, error) {
	// Check if we need to reconnect/reinitialize
	if err := p.ensureConsumerActive(ctx); err != nil {
		return nil, fmt.Errorf("failed to ensure consumer: %w", err)
	}

	if p.cb != nil && !p.cb.Allow() {
		p.logger.Debug("circuit breaker open", logger.Fields{
			"worker_id": workerID,
			"state":     p.cb.GetState(),
		})
		return nil, errors.New("circuit breaker open")
	}

	p.logger.Debug("pulling messages", logger.Fields{
		"worker_id":  workerID,
		"batch_size": p.config.BatchSize,
		"timeout":    p.config.PollTimeout,
	})

	// Record attempt start if circuit breaker is configured
	var cbRecorded bool
	if p.cb != nil {
		cbRecorded = p.cb.RecordStart()
		if !cbRecorded {
			return nil, fmt.Errorf("circuit breaker rejected request")
		}
	}

	// Fetch messages from stream
	jsMsgs, err := p.consumer.Fetch(
		p.config.BatchSize,
		jetstream.FetchMaxWait(p.config.PollTimeout),
	)

	// Record result in circuit breaker
	if p.cb != nil && cbRecorded {
		p.cb.RecordResult(err == nil)
	}

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			// This is normal for pull-based when no messages are available
			return []jetstream.Msg{}, nil
		}
		return nil, fmt.Errorf("failed to fetch messages: %w", err)
	}

	// Convert jetstream messages to our Message type
	messages := []jetstream.Msg{}
	for jsMsg := range jsMsgs.Messages() {
		messages = append(messages, jsMsg)
	}

	return messages, nil
}

// ensureConsumerActive checks the connection status and reinitializes the consumer if needed
func (p *PullStrategy) ensureConsumerActive(ctx context.Context) error {
	if p.consumer != nil && p.cm.IsConnected() {
		// Check if the JS context in the connection manager has changed
		currentJS := p.cm.GetJetStream()
		if currentJS != nil {
			return nil // All good, connection is active
		}
	}

	p.logger.Debug("reconnecting consumer")
	
	// Ensure connection and stream
	if err := p.ensureStream(ctx); err != nil {
		return fmt.Errorf("failed to ensure stream: %w", err)
	}

	// Get current JetStream instance
	js := p.cm.GetJetStream()
	if js == nil {
		return errors.New("JetStream connection not available")
	}

	// Get stream
	stream, err := js.Stream(ctx, p.streamName)
	if err != nil {
		return fmt.Errorf("failed to get stream: %w", err)
	}

	// Configure durable pull consumer
	consumerConfig := jetstream.ConsumerConfig{
		Name:          p.config.ConsumerName,
		Durable:       p.config.ConsumerName,
		AckPolicy:     jetstream.AckExplicitPolicy,
		AckWait:       p.config.AckWait,
		MaxAckPending: p.config.MaxAckPending,
		FilterSubject: p.subject,
		DeliverPolicy: jetstream.DeliverAllPolicy,
		MaxDeliver:    p.config.MaxRetries + 1, // Include first delivery
	}

	// Create or update the consumer
	consumer, err := stream.CreateOrUpdateConsumer(ctx, consumerConfig)
	if err != nil {
		return fmt.Errorf("failed to create/update consumer: %w", err)
	}

	p.consumer = consumer
	p.logger.Debug("consumer reconnected successfully")
	return nil
}
