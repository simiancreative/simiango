package natsjsdlq_test

import (
	"errors"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/simiancreative/simiango/messaging/natsjsdlq"
	"github.com/stretchr/testify/assert"
)

// MockJetStreamContext is a test double for nats.JetStreamContext
type MockJetStreamContext struct {
	AddStreamFunc  func(*nats.StreamConfig, ...nats.JSOpt) (*nats.StreamInfo, error)
	PublishMsgFunc func(*nats.Msg, ...nats.PubOpt) (*nats.PubAck, error)
	publishCalls   []*nats.Msg
}

func (m *MockJetStreamContext) AddStream(
	cfg *nats.StreamConfig,
	opts ...nats.JSOpt,
) (*nats.StreamInfo, error) {
	if m.AddStreamFunc != nil {
		return m.AddStreamFunc(cfg, opts...)
	}
	return &nats.StreamInfo{Config: *cfg}, nil
}

func (m *MockJetStreamContext) PublishMsg(
	msg *nats.Msg,
	opts ...nats.PubOpt,
) (*nats.PubAck, error) {
	m.publishCalls = append(m.publishCalls, msg)
	if m.PublishMsgFunc != nil {
		return m.PublishMsgFunc(msg, opts...)
	}
	return &nats.PubAck{}, nil
}

func TestNewHandler(t *testing.T) {
	tests := []struct {
		name    string
		deps    natsjsdlq.Dependencies
		config  natsjsdlq.Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid configuration",
			deps: natsjsdlq.Dependencies{
				JetStream: &MockJetStreamContext{},
			},
			config: natsjsdlq.Config{
				StreamName:    "test_dlq",
				Subject:       "test.dlq",
				MaxDeliveries: 3,
				Storage:       nats.FileStorage,
			},
			wantErr: false,
		},
		{
			name: "missing jetstream",
			deps: natsjsdlq.Dependencies{
				JetStream: nil,
			},
			config: natsjsdlq.Config{
				StreamName: "test_dlq",
				Subject:    "test.dlq",
			},
			wantErr: true,
			errMsg:  "JetStream context is required",
		},
		{
			name: "missing required config",
			deps: natsjsdlq.Dependencies{
				JetStream: &MockJetStreamContext{},
			},
			config:  natsjsdlq.Config{},
			wantErr: true,
			errMsg:  "stream name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, err := natsjsdlq.NewHandler(tt.deps, tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, handler)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, handler)
			}
		})
	}
}

func TestHandlerPublishMessage(t *testing.T) {
	tests := []struct {
		name           string
		msg            *nats.Msg
		reason         string
		publishErr     error
		wantErr        bool
		validateHeader func(*testing.T, nats.Header)
	}{
		{
			name: "successful publish",
			msg: &nats.Msg{
				Subject: "original.subject",
				Data:    []byte("test data"),
				Header:  nats.Header{"Original-Key": []string{"value"}},
			},
			reason: "test failure",
			validateHeader: func(t *testing.T, h nats.Header) {
				assert.Equal(t, "test failure", h.Get("DLQ-Reason"))
				assert.Equal(t, "original.subject", h.Get("Original-Subject"))
				assert.Equal(t, "value", h.Get("Original-Key"))
				assert.NotEmpty(t, h.Get("DLQ-Timestamp"))
			},
		},
		{
			name: "publish with no headers",
			msg: &nats.Msg{
				Subject: "original.subject",
				Data:    []byte("test data"),
			},
			reason: "test failure",
			validateHeader: func(t *testing.T, h nats.Header) {
				assert.Equal(t, "test failure", h.Get("DLQ-Reason"))
				assert.Equal(t, "original.subject", h.Get("Original-Subject"))
			},
		},
		{
			name: "publish error",
			msg: &nats.Msg{
				Subject: "original.subject",
				Data:    []byte("test data"),
			},
			reason:     "test failure",
			publishErr: errors.New("publish failed"),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockJetStreamContext{
				PublishMsgFunc: func(msg *nats.Msg, opts ...nats.PubOpt) (*nats.PubAck, error) {
					if tt.publishErr != nil {
						return nil, tt.publishErr
					}
					return &nats.PubAck{}, nil
				},
			}

			var errorCaught error
			handler, err := natsjsdlq.NewHandler(
				natsjsdlq.Dependencies{JetStream: mock},
				natsjsdlq.Config{
					StreamName:    "test_dlq",
					Subject:       "test.dlq",
					MaxDeliveries: 3,
					ErrorHandler: func(err error) {
						errorCaught = err
					},
				},
			)
			assert.NoError(t, err)

			err = handler.PublishMessage(tt.msg, tt.reason)
			if tt.wantErr {
				assert.Error(t, err)
				assert.NotNil(t, errorCaught)
			} else {
				assert.NoError(t, err)
				assert.Nil(t, errorCaught)
				assert.Len(t, mock.publishCalls, 1)
				if tt.validateHeader != nil {
					tt.validateHeader(t, mock.publishCalls[0].Header)
				}
			}
		})
	}
}

type MockMsg struct {
	*nats.Msg
	metadata      *nats.MsgMetadata
	metadataError error
}

func (m *MockMsg) Metadata() (*nats.MsgMetadata, error) {
	if m.metadataError != nil {
		return nil, m.metadataError
	}
	return m.metadata, nil
}

func TestHandlerShouldDLQ(t *testing.T) {
	tests := []struct {
		name          string
		msg           func() *MockMsg
		maxDeliveries int
		want          bool
		wantErr       bool
	}{
		{
			name: "should dlq when deliveries exceeded",
			msg: func() *MockMsg {
				return &MockMsg{
					Msg: &nats.Msg{
						Subject: "test.subject",
						Data:    []byte("test data"),
					},
					metadata: &nats.MsgMetadata{
						NumDelivered: 4,
					},
				}
			},
			maxDeliveries: 3,
			want:          true,
		},
		{
			name: "should not dlq when under max deliveries",
			msg: func() *MockMsg {
				return &MockMsg{
					Msg: &nats.Msg{
						Subject: "test.subject",
						Data:    []byte("test data"),
					},
					metadata: &nats.MsgMetadata{
						NumDelivered: 2,
					},
				}
			},
			maxDeliveries: 3,
			want:          false,
		},
		{
			name: "should not dlq on metadata error",
			msg: func() *MockMsg {
				return &MockMsg{
					Msg: &nats.Msg{
						Subject: "test.subject",
						Data:    []byte("test data"),
					},
					metadataError: errors.New("metadata error"),
				}
			},
			maxDeliveries: 3,
			want:          false,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a message with mock metadata function
			msg := tt.msg()

			var errorCaught error
			handler, err := natsjsdlq.NewHandler(natsjsdlq.Dependencies{
				JetStream: &MockJetStreamContext{},
			}, natsjsdlq.Config{
				StreamName:    "test_dlq",
				Subject:       "test.dlq",
				MaxDeliveries: tt.maxDeliveries,
				ErrorHandler: func(err error) {
					errorCaught = err
				},
			})
			assert.NoError(t, err)

			got := handler.ShouldDLQ(msg)
			assert.Equal(t, tt.want, got)

			assert.Equal(t, tt.wantErr, errorCaught != nil)
		})
	}
}
