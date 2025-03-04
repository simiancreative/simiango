package natsjsdlq_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/simiancreative/simiango/messaging/natsjscm"
	"github.com/simiancreative/simiango/messaging/natsjsdlq"
	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
)

// Mock implementations
type MockConnectionManager struct {
	mock.Mock
}

func (m *MockConnectionManager) Connect() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockConnectionManager) GetConnection() *nats.Conn {
	args := m.Called()
	return args.Get(0).(*nats.Conn)
}

func (m *MockConnectionManager) GetJetStream() natsjscm.JetStream {
	args := m.Called()
	return args.Get(0).(jetstream.JetStream)
}

func (m *MockConnectionManager) EnsureStream(
	ctx context.Context,
	config jetstream.StreamConfig,
) (natsjscm.JetStream, error) {
	args := m.Called(ctx, config)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(jetstream.JetStream), args.Error(1)
}

func (m *MockConnectionManager) Disconnect() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockConnectionManager) IsConnected() bool {
	args := m.Called()
	return args.Bool(0)
}

type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) Publish(ctx context.Context, msg *nats.Msg) (*jetstream.PubAck, error) {
	args := m.Called(ctx, msg)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jetstream.PubAck), args.Error(1)
}

type MockMsg struct {
	mock.Mock
}

func (m *MockMsg) Metadata() (*nats.MsgMetadata, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*nats.MsgMetadata), args.Error(1)
}

// Test Suite
type DLQHandlerTestSuite struct {
	mockCM  *MockConnectionManager
	mockPub *MockPublisher
	ctx     context.Context
}

func (suite *DLQHandlerTestSuite) SetupTest() {
	suite.mockCM = new(MockConnectionManager)
	suite.mockPub = new(MockPublisher)
	suite.ctx = context.Background()
}

func TestNewHandlerValidation(t *testing.T) {
	suite := new(DLQHandlerTestSuite)
	suite.SetupTest()

	testCases := []struct {
		name        string
		deps        natsjsdlq.Dependencies
		config      natsjsdlq.Config
		expectedErr string
	}{
		{
			name: "Valid configuration",
			deps: natsjsdlq.Dependencies{
				ConnectionManager: suite.mockCM,
				Publisher:         suite.mockPub,
			},
			config: natsjsdlq.Config{
				StreamName:    "test-dlq",
				Subject:       "test.dlq",
				MaxDeliveries: 3,
				Storage:       jetstream.FileStorage,
				Context:       suite.ctx,
			},
			expectedErr: "",
		},
		{
			name: "Missing connection manager",
			deps: natsjsdlq.Dependencies{
				ConnectionManager: nil,
				Publisher:         suite.mockPub,
			},
			config: natsjsdlq.Config{
				StreamName:    "test-dlq",
				Subject:       "test.dlq",
				MaxDeliveries: 3,
			},
			expectedErr: "invalid DLQ configuration: connection manager is required",
		},
		{
			name: "Missing publisher",
			deps: natsjsdlq.Dependencies{
				ConnectionManager: suite.mockCM,
				Publisher:         nil,
			},
			config: natsjsdlq.Config{
				StreamName:    "test-dlq",
				Subject:       "test.dlq",
				MaxDeliveries: 3,
			},
			expectedErr: "invalid DLQ configuration: publisher is required",
		},
		{
			name: "Empty stream name",
			deps: natsjsdlq.Dependencies{
				ConnectionManager: suite.mockCM,
				Publisher:         suite.mockPub,
			},
			config: natsjsdlq.Config{
				StreamName:    "",
				Subject:       "test.dlq",
				MaxDeliveries: 3,
			},
			expectedErr: "invalid DLQ configuration: stream name is required",
		},
		{
			name: "Empty subject",
			deps: natsjsdlq.Dependencies{
				ConnectionManager: suite.mockCM,
				Publisher:         suite.mockPub,
			},
			config: natsjsdlq.Config{
				StreamName:    "test-dlq",
				Subject:       "",
				MaxDeliveries: 3,
			},
			expectedErr: "invalid DLQ configuration: subject is required",
		},
		{
			name: "Invalid max deliveries",
			deps: natsjsdlq.Dependencies{
				ConnectionManager: suite.mockCM,
				Publisher:         suite.mockPub,
			},
			config: natsjsdlq.Config{
				StreamName:    "test-dlq",
				Subject:       "test.dlq",
				MaxDeliveries: 0,
			},
			expectedErr: "invalid DLQ configuration: max deliveries must be greater than 0",
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedErr == "" {
				// Mock successful stream creation
				streamConfig := jetstream.StreamConfig{
					Name:      tc.config.StreamName,
					Subjects:  []string{tc.config.Subject},
					Storage:   tc.config.Storage,
					Retention: jetstream.WorkQueuePolicy,
				}
				suite.mockCM.On("EnsureStream", mock.Anything, streamConfig).Return(nil, nil).Once()
			}

			handler, err := natsjsdlq.NewHandler(tc.deps, tc.config)

			if tc.expectedErr != "" {
				assert.EqualError(t, err, tc.expectedErr)
				assert.Nil(t, handler)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, handler)
				suite.mockCM.AssertExpectations(t)
			}
		})
	}
}

func TestShouldDLQ(t *testing.T) {
	suite := new(DLQHandlerTestSuite)
	suite.SetupTest()

	testCases := []struct {
		name           string
		maxDeliveries  int
		numDelivered   uint64
		metadataErr    error
		expectedResult bool
	}{
		{
			name:           "Should send to DLQ - equal to max",
			maxDeliveries:  3,
			numDelivered:   3,
			metadataErr:    nil,
			expectedResult: true,
		},
		{
			name:           "Should send to DLQ - greater than max",
			maxDeliveries:  3,
			numDelivered:   5,
			metadataErr:    nil,
			expectedResult: true,
		},
		{
			name:           "Should not send to DLQ - less than max",
			maxDeliveries:  3,
			numDelivered:   2,
			metadataErr:    nil,
			expectedResult: false,
		},
		{
			name:           "Metadata error",
			maxDeliveries:  3,
			numDelivered:   0,
			metadataErr:    errors.New("metadata error"),
			expectedResult: false,
		},
	}

	// Setup for all test cases
	validConfig := natsjsdlq.Config{
		StreamName:    "test-dlq",
		Subject:       "test.dlq",
		MaxDeliveries: 3,
		Context:       context.Background(),
	}

	deps := natsjsdlq.Dependencies{
		ConnectionManager: suite.mockCM,
		Publisher:         suite.mockPub,
	}

	// Mock stream creation once for all test cases
	streamConfig := jetstream.StreamConfig{
		Name:      validConfig.StreamName,
		Subjects:  []string{validConfig.Subject},
		Storage:   jetstream.FileStorage,
		Retention: jetstream.WorkQueuePolicy,
	}
	suite.mockCM.On("EnsureStream", mock.Anything, streamConfig).Return(nil, nil).Once()

	handler, err := natsjsdlq.NewHandler(deps, validConfig)
	assert.NoError(t, err)
	assert.NotNil(t, handler)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock message with the configured metadata
			mockMsg := new(MockMsg)
			metadata := &nats.MsgMetadata{
				NumDelivered: tc.numDelivered,
			}

			mockMsg.On("Metadata").Return(metadata, tc.metadataErr).Once()

			// Test the ShouldDLQ method
			result := handler.ShouldDLQ(mockMsg)
			assert.Equal(t, tc.expectedResult, result)
			mockMsg.AssertExpectations(t)
		})
	}
}

func TestPublishMessage(t *testing.T) {
	suite := new(DLQHandlerTestSuite)
	suite.SetupTest()

	testCases := []struct {
		name               string
		originalMsg        *nats.Msg
		reason             string
		publishError       error
		expectedDLQSubject string
		expectHeaders      map[string]string
	}{
		{
			name: "Successful DLQ publish with headers",
			originalMsg: &nats.Msg{
				Subject: "original.subject",
				Data:    []byte("test data"),
				Header: nats.Header{
					"Nats-Msg-Id": []string{"msg-123"},
					"Custom":      []string{"value"},
				},
			},
			reason:             "processing failed",
			publishError:       nil,
			expectedDLQSubject: "test.dlq",
			expectHeaders: map[string]string{
				"DLQ-Reason":          "processing failed",
				"Original-Subject":    "original.subject",
				"Original-Message-ID": "msg-123",
				"Custom":              "value",
			},
		},
		{
			name: "Successful DLQ publish without headers",
			originalMsg: &nats.Msg{
				Subject: "original.subject",
				Data:    []byte("test data"),
			},
			reason:             "processing failed",
			publishError:       nil,
			expectedDLQSubject: "test.dlq",
			expectHeaders: map[string]string{
				"DLQ-Reason":       "processing failed",
				"Original-Subject": "original.subject",
			},
		},
		{
			name: "Failed DLQ publish",
			originalMsg: &nats.Msg{
				Subject: "original.subject",
				Data:    []byte("test data"),
			},
			reason:             "processing failed",
			publishError:       errors.New("publish error"),
			expectedDLQSubject: "test.dlq",
			expectHeaders: map[string]string{
				"DLQ-Reason":       "processing failed",
				"Original-Subject": "original.subject",
			},
		},
	}

	// Setup for all test cases
	validConfig := natsjsdlq.Config{
		StreamName:    "test-dlq",
		Subject:       "test.dlq",
		MaxDeliveries: 3,
		Context:       context.Background(),
		ErrorHandler:  func(err error) {}, // No-op error handler
	}

	deps := natsjsdlq.Dependencies{
		ConnectionManager: suite.mockCM,
		Publisher:         suite.mockPub,
	}

	// Mock stream creation once for all test cases
	streamConfig := jetstream.StreamConfig{
		Name:      validConfig.StreamName,
		Subjects:  []string{validConfig.Subject},
		Storage:   jetstream.FileStorage,
		Retention: jetstream.WorkQueuePolicy,
	}
	suite.mockCM.On("EnsureStream", mock.Anything, streamConfig).Return(nil, nil).Once()

	handler, err := natsjsdlq.NewHandler(deps, validConfig)
	assert.NoError(t, err)
	assert.NotNil(t, handler)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup the publisher mock
			suite.mockPub.On("Publish", mock.Anything, mock.MatchedBy(func(msg *nats.Msg) bool {
				// Verify the DLQ message properties
				if msg.Subject != tc.expectedDLQSubject {
					return false
				}

				if len(msg.Data) != len(tc.originalMsg.Data) {
					return false
				}

				// Verify headers
				for key, expectedValue := range tc.expectHeaders {
					if msg.Header.Get(key) == "" {
						return false
					}

					if key == "DLQ-Timestamp" { // Skip timestamp check as it's dynamic
						continue
					}

					if expectedValue != msg.Header.Get(key) {
						return false
					}
				}

				// Verify timestamp header exists and is valid
				timestampStr := msg.Header.Get("DLQ-Timestamp")
				if timestampStr == "" {
					return false
				}
				_, err := time.Parse(time.RFC3339, timestampStr)
				return err == nil
			})).Return(nil, tc.publishError).Once()

			// Test the PublishMessage method
			err := handler.PublishMessage(tc.originalMsg, tc.reason)

			if tc.publishError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.publishError, err)
			} else {
				assert.NoError(t, err)
			}

			suite.mockPub.AssertExpectations(t)
		})
	}
}
