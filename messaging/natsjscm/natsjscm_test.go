package natsjscm_test

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"testing"
	"time"

	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go/jetstream"

	"github.com/nats-io/nats.go"
	"github.com/simiancreative/simiango/messaging/natsjscm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock Logger for testing
type mockLogger struct {
	debugMessages []string
	errorMessages []string
}

func newMockLogger() *mockLogger {
	return &mockLogger{
		debugMessages: []string{},
		errorMessages: []string{},
	}
}

func (m *mockLogger) Debugf(format string, args ...interface{}) {
	m.debugMessages = append(m.debugMessages, fmt.Sprintf(format, args...))
}

func (m *mockLogger) Errorf(format string, args ...interface{}) {
	m.errorMessages = append(m.errorMessages, fmt.Sprintf(format, args...))
}

func MockServer(args ...int) func() {
	port := 0
	if len(args) > 0 {
		port = args[0]
	}

	err := getPort(&port)
	if err != nil {
		panic(err)
	}

	os.Setenv("NATS_HOST", fmt.Sprintf("localhost:%v", port))
	os.Setenv("NATS_PORT", fmt.Sprintf("%v", port))

	opts := natsserver.DefaultTestOptions
	opts.Port = port
	opts.JetStream = true

	return natsserver.RunServer(&opts).Shutdown
}

func getPort(port *int) error {
	if *port != 0 {
		return nil
	}

	addr, _ := net.ResolveTCPAddr("tcp", "localhost:0")

	listener, _ := net.ListenTCP("tcp", addr)
	defer listener.Close()

	*port = listener.Addr().(*net.TCPAddr).Port

	return nil
}

// TestNewConnectionManager tests the creation of the ConnectionManager
func TestNewConnectionManager(t *testing.T) {
	tests := []struct {
		name        string
		config      natsjscm.ConnectionConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid configuration",
			config: natsjscm.ConnectionConfig{
				URL:           "nats://localhost:4222",
				ReconnectWait: 5 * time.Second,
			},
			expectError: false,
		},
		{
			name: "Missing URL",
			config: natsjscm.ConnectionConfig{
				URL:           "",
				ReconnectWait: 5 * time.Second,
			},
			expectError: true,
			errorMsg:    "NATS URL is required",
		},
		{
			name: "With custom logger",
			config: natsjscm.ConnectionConfig{
				URL:           "nats://localhost:4222",
				ReconnectWait: 5 * time.Second,
				Logger:        newMockLogger(),
			},
			expectError: false,
		},
		{
			name: "With custom JetStream domain",
			config: natsjscm.ConnectionConfig{
				URL:             "nats://localhost:4222",
				ReconnectWait:   5 * time.Second,
				JetStreamDomain: "test-domain",
			},
			expectError: false,
		},
		{
			name: "Zero reconnect wait should use default",
			config: natsjscm.ConnectionConfig{
				URL:           "nats://localhost:4222",
				ReconnectWait: 0,
			},
			expectError: false,
		},
		{
			name: "With NATS options",
			config: natsjscm.ConnectionConfig{
				URL:           "nats://localhost:4222",
				ReconnectWait: 5 * time.Second,
				Options:       []nats.Option{nats.Name("test-connection")},
			},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cm, err := natsjscm.NewConnectionManager(tc.config)

			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)
				require.NotNil(t, cm)
			}
		})
	}
}

// TestConnectionLifecycle tests connect, disconnect and status methods
func TestConnectionLifecycle(t *testing.T) {
	defer MockServer()()

	tests := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{
			name: "Connect and check status",
			fn: func(t *testing.T) {
				cm, err := natsjscm.NewConnectionManager(natsjscm.ConnectionConfig{
					URL:           os.Getenv("NATS_HOST"),
					ReconnectWait: 100 * time.Millisecond,
				})
				require.NoError(t, err)

				// Initially not connected
				assert.False(t, cm.IsConnected())

				// Connect
				err = cm.Connect()
				require.NoError(t, err)
				assert.True(t, cm.IsConnected())

				// Get connection and JetStream
				nc := cm.GetConnection()
				require.NotNil(t, nc)

				js := cm.GetJetStream()
				require.NotNil(t, js)

				// Disconnect
				err = cm.Disconnect()
				require.NoError(t, err)
				assert.False(t, cm.IsConnected())
			},
		},
		{
			name: "Reference counting in connect/disconnect",
			fn: func(t *testing.T) {
				cm, err := natsjscm.NewConnectionManager(natsjscm.ConnectionConfig{
					URL:           os.Getenv("NATS_HOST"),
					ReconnectWait: 100 * time.Millisecond,
				})
				require.NoError(t, err)

				// First connection
				err = cm.Connect()
				require.NoError(t, err)
				assert.True(t, cm.IsConnected())

				// Second connection (should increment ref count)
				err = cm.Connect()
				require.NoError(t, err)

				// First disconnect (should decrement ref count but stay connected)
				err = cm.Disconnect()
				require.NoError(t, err)
				assert.True(t, cm.IsConnected())

				// Second disconnect (should close connection)
				err = cm.Disconnect()
				require.NoError(t, err)
				assert.False(t, cm.IsConnected())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, tc.fn)
	}
}

// TestEnsureStream tests stream management functionality
func TestEnsureStream(t *testing.T) {
	defer MockServer()()

	cm, err := natsjscm.NewConnectionManager(natsjscm.ConnectionConfig{
		URL:           os.Getenv("NATS_HOST"),
		ReconnectWait: 100 * time.Millisecond,
	})
	require.NoError(t, err)

	err = cm.Connect()
	require.NoError(t, err)
	defer func() {
		err := cm.Disconnect()
		require.NoError(t, err)
	}()

	ctx := context.Background()

	tests := []struct {
		name        string
		config      jetstream.StreamConfig
		expectError bool
	}{
		{
			name: "Create new stream",
			config: jetstream.StreamConfig{
				Name:     "test-stream-1",
				Subjects: []string{"test.subject.1"},
				Storage:  jetstream.MemoryStorage,
			},
			expectError: false,
		},
		{
			name: "Get existing stream",
			config: jetstream.StreamConfig{
				Name:     "test-stream-1", // Same as previous test
				Subjects: []string{"test.subject.1"},
				Storage:  jetstream.MemoryStorage,
			},
			expectError: false,
		},
		{
			name: "Create stream with invalid config",
			config: jetstream.StreamConfig{
				Name:     "", // Invalid - empty name
				Subjects: []string{"test.subject.2"},
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			js, err := cm.EnsureStream(ctx, tc.config)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, js)

				// Verify the stream exists
				streamInfo, err := js.Stream(ctx, tc.config.Name)
				require.NoError(t, err)
				require.NotNil(t, streamInfo)
				assert.Equal(t, tc.config.Name, streamInfo.CachedInfo().Config.Name)
			}
		})
	}
}

// TestConnectionEvents tests the connection event handlers
func TestConnectionEvents(t *testing.T) {
	shutdown := MockServer()

	// Create a mock logger to capture log messages
	logger := newMockLogger()

	// Create a connection manager with event handlers
	cm, err := natsjscm.NewConnectionManager(natsjscm.ConnectionConfig{
		URL:           os.Getenv("NATS_HOST"),
		ReconnectWait: 100 * time.Millisecond,
		Logger:        logger,
	})
	require.NoError(t, err)

	// Connect to the server
	err = cm.Connect()
	require.NoError(t, err)
	assert.True(t, cm.IsConnected())

	// Simulate server shutdown to trigger disconnect
	shutdown()

	assert.Eventually(t, func() bool {
		return !cm.IsConnected()
	}, 5000*time.Millisecond, 100*time.Millisecond, "Connection should be closed")

	// Restart the server to trigger reconnect
	port, _ := strconv.Atoi(os.Getenv("NATS_PORT"))
	defer MockServer(port)()

	// Verify reconnect handler was called
	assert.Eventually(t, func() bool {
		return cm.IsConnected()
	}, 5000*time.Millisecond, 100*time.Millisecond, "Connection should be reconnected")

	// Verify JetStream context is properly recreated
	js := cm.GetJetStream()
	require.NotNil(t, js)

	// Clean up
	err = cm.Disconnect()
	require.NoError(t, err)
	assert.False(t, cm.IsConnected())
}
