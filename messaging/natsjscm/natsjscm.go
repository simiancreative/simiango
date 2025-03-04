package natsjscm

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/simiancreative/simiango/logger"
)

type Logger interface {
	Debugf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type Connector interface {
	Connect() error
	Disconnect() error
	IsConnected() bool
	GetConnection() *nats.Conn
	GetJetStream() JetStream
	EnsureStream(ctx context.Context, config jetstream.StreamConfig) (JetStream, error)
}

type JetStream interface {
	Stream(ctx context.Context, name string) (jetstream.Stream, error)
	CreateStream(ctx context.Context, config jetstream.StreamConfig) (jetstream.Stream, error)
	PublishMsg(context.Context, *nats.Msg, ...jetstream.PublishOpt) (*jetstream.PubAck, error)
}

// ConnectionConfig holds the configuration for NATS connection
type ConnectionConfig struct {
	URL             string
	Options         []nats.Option
	Logger          Logger // Optional logger
	ReconnectWait   time.Duration
	JetStreamDomain string // Optional JetStream domain
}

// ConnectionManager manages NATS connections
type ConnectionManager struct {
	config ConnectionConfig
	nc     *nats.Conn
	js     jetstream.JetStream
	mu     sync.RWMutex
	refs   int
	log    Logger
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager(config ConnectionConfig) (Connector, error) {
	if config.URL == "" {
		return nil, fmt.Errorf("NATS URL is required")
	}

	if config.Logger == nil {
		config.Logger = logger.New()
	}

	if config.ReconnectWait <= 0 {
		config.ReconnectWait = 5 * time.Second
	}

	// Add connection status callbacks
	options := append([]nats.Option{}, config.Options...)
	cm := &ConnectionManager{
		config: config,
		log:    config.Logger,
	}

	// Add connection event handlers
	options = append(options,
		nats.ReconnectHandler(func(_ *nats.Conn) {
			cm.handleReconnect()
		}),
	)

	cm.config.Options = options
	return cm, nil
}

func (cm *ConnectionManager) handleReconnect() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.log.Debugf("NATS connection reconnected")

	// Recreate JetStream context after reconnect
	if cm.nc == nil {
		return
	}

	js, err := cm.createJetStreamContext(cm.nc)
	if err != nil {
		cm.log.Errorf("Failed to recreate JetStream context %s", err)
		go cm.retryConnection()
		return
	}

	cm.js = js
	cm.log.Debugf("JetStream context recreated")
}

// createJetStreamContext creates a new JetStream context with the current configuration
func (cm *ConnectionManager) createJetStreamContext(nc *nats.Conn) (jetstream.JetStream, error) {
	// Create JetStream context
	return jetstream.New(nc)
}

// retryConnection attempts to reconnect to NATS periodically
func (cm *ConnectionManager) retryConnection() {
	for {
		time.Sleep(cm.config.ReconnectWait)

		cm.mu.Lock()
		if cm.nc != nil && cm.nc.IsConnected() {
			cm.mu.Unlock()
			return
		}

		cm.log.Debugf("Attempting to reconnect to NATS")
		nc, err := nats.Connect(cm.config.URL, cm.config.Options...)
		if err != nil {
			cm.log.Errorf("Failed to reconnect to NATS: %s", err)
			cm.mu.Unlock()
			continue
		}

		js, err := cm.createJetStreamContext(nc)
		if err != nil {
			cm.log.Errorf("Failed to recreate JetStream context: %s", err)
			nc.Close()
			cm.mu.Unlock()
			continue
		}

		cm.nc = nc
		cm.js = js
		cm.refs = 1
		cm.log.Debugf("Successfully reconnected to NATS")
		cm.mu.Unlock()
		return
	}
}

// Connect establishes a connection to NATS if not already connected
func (cm *ConnectionManager) Connect() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// If already connected, increment reference count
	if cm.nc != nil && cm.nc.IsConnected() {
		cm.refs++
		return nil
	}

	// Connect to NATS
	nc, err := nats.Connect(cm.config.URL, cm.config.Options...)
	if err != nil {
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}

	// Create JetStream context
	js, err := cm.createJetStreamContext(nc)
	if err != nil {
		nc.Close()
		return fmt.Errorf("failed to create JetStream context: %w", err)
	}

	cm.nc = nc
	cm.js = js
	cm.refs = 1

	return nil
}

// GetConnection returns the NATS connection
func (cm *ConnectionManager) GetConnection() *nats.Conn {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.nc
}

// GetJetStream returns the JetStream context
func (cm *ConnectionManager) GetJetStream() JetStream {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.js
}

// EnsureStream creates a stream if it doesn't exist
func (cm *ConnectionManager) EnsureStream(
	ctx context.Context,
	config jetstream.StreamConfig,
) (JetStream, error) {
	cm.mu.RLock()
	js := cm.js
	cm.mu.RUnlock()

	if js == nil {
		return nil, fmt.Errorf("not connected to JetStream")
	}

	// Try to get the stream first
	_, err := js.Stream(ctx, config.Name)

	// If stream exists, return it
	if err == nil {
		return js, nil
	}

	// If stream doesn't exist, create it
	if err == jetstream.ErrStreamNotFound {
		_, err := js.CreateStream(ctx, config)
		if err != nil {
			return nil, err
		}
		return js, nil
	}

	return nil, err
}

// Disconnect decrements the reference count and closes connection if no more references
func (cm *ConnectionManager) Disconnect() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.nc == nil {
		return nil
	}

	cm.refs--
	if cm.refs > 0 {
		return nil
	}

	cm.nc.Close()
	cm.nc = nil
	cm.js = nil

	return nil
}

// IsConnected checks if there is an active connection
func (cm *ConnectionManager) IsConnected() bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.nc != nil && cm.nc.IsConnected() && !cm.nc.IsClosed()
}
