package natsjspub

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/simiancreative/simiango/circuitbreaker"
	"github.com/simiancreative/simiango/logger"
	"github.com/simiancreative/simiango/messaging/natsjscm"
)

type Logger interface {
	Debugf(format string, args ...interface{})
}

type Publisher interface {
	Publish(ctx context.Context, msg *nats.Msg) (*jetstream.PubAck, error)
}

type JsonPublisher interface {
	PublishJSON(
		ctx context.Context,
		data interface{},
		headers ...map[string]string,
	) (*jetstream.PubAck, error)
}

// Config holds publisher configuration
type Config struct {
	// Stream name to publish to
	StreamName string

	// Subject to publish on
	Subject string

	// CircuitBreaker configuration (optional)
	CircuitBreaker *circuitbreaker.Config

	// Publish timeout (default 5s)
	Timeout time.Duration

	// Message retention policy (optional, default is WorkQueuePolicy)
	RetentionPolicy jetstream.RetentionPolicy
}

// Dependencies for the publisher
type Dependencies struct {
	// ConnectionManager for NATS
	ConnectionManager *natsjscm.ConnectionManager
}

// Publisher is a JetStream publisher with circuit breaker capabilities
type PublishManager struct {
	config Config
	cm     *natsjscm.ConnectionManager
	cb     *circuitbreaker.CircuitBreaker
	log    Logger
}

// NewPublisher creates a new JetStream publisher
func NewPublisher(deps Dependencies, config Config) (*PublishManager, error) {
	// Validation
	if deps.ConnectionManager == nil {
		return nil, fmt.Errorf("connection manager is required")
	}

	if config.StreamName == "" {
		return nil, fmt.Errorf("stream name is required")
	}

	if config.Subject == "" {
		return nil, fmt.Errorf("subject is required")
	}

	if config.Timeout <= 0 {
		config.Timeout = 5 * time.Second
	}

	pub := &PublishManager{
		config: config,
		cm:     deps.ConnectionManager,
		log:    logger.New(),
	}

	// Initialize circuit breaker if configured
	if config.CircuitBreaker != nil {
		cb, err := circuitbreaker.New(*config.CircuitBreaker)
		if err != nil {
			return nil, fmt.Errorf("failed to create circuit breaker: %w", err)
		}
		pub.cb = cb
	}

	// Ensure the stream exists
	err := pub.ensureStream(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to ensure stream: %w", err)
	}

	return pub, nil
}

// ensureStream makes sure the configured stream exists
func (p *PublishManager) ensureStream(ctx context.Context) error {
	// Get JetStream connection
	if !p.cm.IsConnected() {
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
		Name:      p.config.StreamName,
		Subjects:  []string{p.config.Subject},
		Retention: retentionPolicy,
	}

	// Ensure stream exists
	_, err := p.cm.EnsureStream(ctx, streamConfig)
	return err
}

// Publish publishes a message to the configured subject
func (p *PublishManager) Publish(
	ctx context.Context,
	msg *nats.Msg,
) (*jetstream.PubAck, error) {
	// Check if circuit breaker allows the request
	if p.cb != nil && !p.cb.Allow() {
		return nil, fmt.Errorf("circuit breaker is open")
	}

	// Record attempt start if circuit breaker is configured
	var cbRecorded bool
	if p.cb != nil {
		cbRecorded = p.cb.RecordStart()
		if !cbRecorded {
			return nil, fmt.Errorf("circuit breaker rejected request")
		}
	}

	// Get JetStream connection
	js := p.cm.GetJetStream()
	if js == nil {
		// Record failure if circuit breaker is configured
		if p.cb != nil && cbRecorded {
			p.cb.RecordResult(false)
		}
		return nil, fmt.Errorf("jetstream connection not available")
	}

	// Apply timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, p.config.Timeout)
	defer cancel()

	// Publish to JetStream
	ack, err := js.PublishMsg(ctxWithTimeout, msg)

	// Record result in circuit breaker
	if p.cb != nil && cbRecorded {
		p.cb.RecordResult(err == nil)
	}

	if err != nil {
		p.log.Debugf("failed to publish message: %s", err)
	}

	return ack, err
}

// PublishJSON publishes a JSON-serializable object to the configured subject
func (p *PublishManager) PublishJSON(
	ctx context.Context,
	data interface{},
	headers ...map[string]string,
) (*jetstream.PubAck, error) {
	msg := &nats.Msg{}

	// Marshal data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	for _, header := range headers {
		for k, v := range header {
			msg.Header.Add(k, v)
		}
	}

	// Add content-type header if not already present
	if val := msg.Header.Get("content-type"); val == "" {
		msg.Header.Add("content-type", "application/json")
	}

	msg.Data = jsonData
	msg.Subject = p.config.Subject

	return p.Publish(ctx, msg)
}

// Close cleans up resources
func (p *PublishManager) Close() error {
	if p.cb != nil {
		p.cb.Reset()
	}
	return nil
}
