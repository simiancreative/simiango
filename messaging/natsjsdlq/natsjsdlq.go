package natsjsdlq

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

type JetStreamContext interface {
	AddStream(cfg *nats.StreamConfig, opts ...nats.JSOpt) (*nats.StreamInfo, error)
	PublishMsg(msg *nats.Msg, opts ...nats.PubOpt) (*nats.PubAck, error)
}

type Msg interface {
	Metadata() (*nats.MsgMetadata, error)
}

// Config holds DLQ configuration
type Config struct {
	// StreamName for the DLQ
	StreamName string

	// Subject to publish dead letters to
	Subject string

	// MaxDeliveries before message is considered dead
	MaxDeliveries int

	// Storage type for the DLQ stream
	Storage nats.StorageType

	// Optional handler for DLQ errors
	ErrorHandler func(error)
}

type Dependencies struct {
	JetStream JetStreamContext
}

// Handler manages dead letter queue operations
type Handler struct {
	config Config
	js     JetStreamContext
}

// NewHandler creates a new DLQ handler
func NewHandler(deps Dependencies, config Config) (*Handler, error) {
	if err := validateConfig(deps, config); err != nil {
		return nil, fmt.Errorf("invalid DLQ configuration: %w", err)
	}

	handler := &Handler{
		config: config,
		js:     deps.JetStream,
	}

	if err := handler.setup(); err != nil {
		return nil, err
	}

	return handler, nil
}

func validateConfig(deps Dependencies, config Config) error {
	if deps.JetStream == nil {
		return fmt.Errorf("JetStream context is required")
	}

	if config.StreamName == "" {
		return fmt.Errorf("stream name is required")
	}

	if config.Subject == "" {
		return fmt.Errorf("subject is required")
	}

	if config.MaxDeliveries <= 0 {
		return fmt.Errorf("max deliveries must be greater than 0")
	}

	if config.Storage == 0 {
		config.Storage = nats.FileStorage
	}

	return nil
}

// setup ensures the DLQ stream exists
func (h *Handler) setup() error {
	streamConfig := &nats.StreamConfig{
		Name:      h.config.StreamName,
		Subjects:  []string{h.config.Subject},
		Storage:   h.config.Storage,
		Retention: nats.WorkQueuePolicy,
	}

	_, err := h.js.AddStream(streamConfig)
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		return fmt.Errorf("failed to create DLQ stream: %w", err)
	}

	return nil
}

// PublishMessage sends a message to the DLQ
func (h *Handler) PublishMessage(msg *nats.Msg, reason string) error {
	// Clone original message headers
	headers := nats.Header{}
	if msg.Header != nil {
		for k, v := range msg.Header {
			headers[k] = v
		}
	}

	// Add DLQ metadata
	headers.Set("DLQ-Reason", reason)
	headers.Set("DLQ-Timestamp", time.Now().UTC().Format(time.RFC3339))
	headers.Set("Original-Subject", msg.Subject)
	if msg.Header != nil {
		headers.Set("Original-Message-ID", msg.Header.Get("Nats-Msg-Id"))
	}

	dlqMsg := nats.NewMsg(h.config.Subject)
	dlqMsg.Header = headers
	dlqMsg.Data = msg.Data

	// Publish to DLQ
	_, err := h.js.PublishMsg(dlqMsg)
	if err != nil && h.config.ErrorHandler != nil {
		h.config.ErrorHandler(fmt.Errorf("failed to publish to DLQ: %w", err))
	}

	return err
}

// ShouldDLQ determines if a message should be sent to DLQ based on delivery count
func (h *Handler) ShouldDLQ(msg Msg) bool {
	metadata, err := msg.Metadata()
	if err != nil {
		if h.config.ErrorHandler != nil {
			h.config.ErrorHandler(fmt.Errorf("failed to get message metadata: %w", err))
		}
		return false
	}

	return metadata.NumDelivered >= uint64(h.config.MaxDeliveries)
}
