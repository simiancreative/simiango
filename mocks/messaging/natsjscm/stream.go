package natsjscm

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/mock"
)

// MockStream is a mock implementation of the Stream interface
type MockStream struct {
	mock.Mock
}

// CachedInfo provides a mock function with given fields:
func (_m *MockStream) CachedInfo() *jetstream.StreamInfo {
	ret := _m.Called()

	var r0 *jetstream.StreamInfo
	if rf, ok := ret.Get(0).(func() *jetstream.StreamInfo); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*jetstream.StreamInfo)
		}
	}

	return r0
}

// Consumer provides a mock function with given fields: ctx, name
func (_m *MockStream) Consumer(ctx context.Context, name string) (jetstream.Consumer, error) {
	ret := _m.Called(ctx, name)

	var r0 jetstream.Consumer
	if rf, ok := ret.Get(0).(func(context.Context, string) jetstream.Consumer); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(jetstream.Consumer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ConsumerNames provides a mock function with given fields: ctx
func (_m *MockStream) ConsumerNames(ctx context.Context) jetstream.ConsumerNameLister {
	ret := _m.Called(ctx)

	var r0 jetstream.ConsumerNameLister
	if rf, ok := ret.Get(0).(func(context.Context) jetstream.ConsumerNameLister); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(jetstream.ConsumerNameLister)
		}
	}

	return r0
}

// CreateConsumer provides a mock function with given fields: ctx, cfg
func (_m *MockStream) CreateConsumer(
	ctx context.Context,
	cfg jetstream.ConsumerConfig,
) (jetstream.Consumer, error) {
	ret := _m.Called(ctx, cfg)

	var r0 jetstream.Consumer
	if rf, ok := ret.Get(0).(func(context.Context, jetstream.ConsumerConfig) jetstream.Consumer); ok {
		r0 = rf(ctx, cfg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(jetstream.Consumer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, jetstream.ConsumerConfig) error); ok {
		r1 = rf(ctx, cfg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateOrUpdateConsumer provides a mock function with given fields: ctx, cfg
func (_m *MockStream) CreateOrUpdateConsumer(
	ctx context.Context,
	cfg jetstream.ConsumerConfig,
) (jetstream.Consumer, error) {
	ret := _m.Called(ctx, cfg)

	var r0 jetstream.Consumer
	if rf, ok := ret.Get(0).(func(context.Context, jetstream.ConsumerConfig) jetstream.Consumer); ok {
		r0 = rf(ctx, cfg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(jetstream.Consumer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, jetstream.ConsumerConfig) error); ok {
		r1 = rf(ctx, cfg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteConsumer provides a mock function with given fields: ctx, name
func (_m *MockStream) DeleteConsumer(ctx context.Context, name string) error {
	ret := _m.Called(ctx, name)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteMsg provides a mock function with given fields: ctx, seq
func (_m *MockStream) DeleteMsg(ctx context.Context, seq uint64) error {
	ret := _m.Called(ctx, seq)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) error); ok {
		r0 = rf(ctx, seq)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetLastMsgForSubject provides a mock function with given fields: ctx, subject
func (_m *MockStream) GetLastMsgForSubject(
	ctx context.Context,
	subject string,
) (*jetstream.RawStreamMsg, error) {
	ret := _m.Called(ctx, subject)

	var r0 *jetstream.RawStreamMsg
	if rf, ok := ret.Get(0).(func(context.Context, string) *jetstream.RawStreamMsg); ok {
		r0 = rf(ctx, subject)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*jetstream.RawStreamMsg)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, subject)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMsg provides a mock function with given fields: ctx, seq, opts
func (_m *MockStream) GetMsg(
	ctx context.Context,
	seq uint64,
	opts ...jetstream.GetMsgOpt,
) (*jetstream.RawStreamMsg, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, seq)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *jetstream.RawStreamMsg
	if rf, ok := ret.Get(0).(func(context.Context, uint64, ...jetstream.GetMsgOpt) *jetstream.RawStreamMsg); ok {
		r0 = rf(ctx, seq, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*jetstream.RawStreamMsg)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint64, ...jetstream.GetMsgOpt) error); ok {
		r1 = rf(ctx, seq, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Info provides a mock function with given fields: ctx, opts
func (_m *MockStream) Info(
	ctx context.Context,
	opts ...jetstream.StreamInfoOpt,
) (*jetstream.StreamInfo, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *jetstream.StreamInfo
	if rf, ok := ret.Get(0).(func(context.Context, ...jetstream.StreamInfoOpt) *jetstream.StreamInfo); ok {
		r0 = rf(ctx, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*jetstream.StreamInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, ...jetstream.StreamInfoOpt) error); ok {
		r1 = rf(ctx, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListConsumers provides a mock function with given fields: ctx
func (_m *MockStream) ListConsumers(ctx context.Context) jetstream.ConsumerInfoLister {
	ret := _m.Called(ctx)

	var r0 jetstream.ConsumerInfoLister
	if rf, ok := ret.Get(0).(func(context.Context) jetstream.ConsumerInfoLister); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(jetstream.ConsumerInfoLister)
		}
	}

	return r0
}

// OrderedConsumer provides a mock function with given fields: ctx, cfg
func (_m *MockStream) OrderedConsumer(
	ctx context.Context,
	cfg jetstream.OrderedConsumerConfig,
) (jetstream.Consumer, error) {
	ret := _m.Called(ctx, cfg)

	var r0 jetstream.Consumer
	if rf, ok := ret.Get(0).(func(context.Context, jetstream.OrderedConsumerConfig) jetstream.Consumer); ok {
		r0 = rf(ctx, cfg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(jetstream.Consumer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, jetstream.OrderedConsumerConfig) error); ok {
		r1 = rf(ctx, cfg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Purge provides a mock function with given fields: ctx, opts
func (_m *MockStream) Purge(ctx context.Context, opts ...jetstream.StreamPurgeOpt) error {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ...jetstream.StreamPurgeOpt) error); ok {
		r0 = rf(ctx, opts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SecureDeleteMsg provides a mock function with given fields: ctx, seq
func (_m *MockStream) SecureDeleteMsg(ctx context.Context, seq uint64) error {
	ret := _m.Called(ctx, seq)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) error); ok {
		r0 = rf(ctx, seq)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateConsumer provides a mock function with given fields: ctx, cfg
func (_m *MockStream) UpdateConsumer(
	ctx context.Context,
	cfg jetstream.ConsumerConfig,
) (jetstream.Consumer, error) {
	ret := _m.Called(ctx, cfg)

	var r0 jetstream.Consumer
	if rf, ok := ret.Get(0).(func(context.Context, jetstream.ConsumerConfig) jetstream.Consumer); ok {
		r0 = rf(ctx, cfg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(jetstream.Consumer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, jetstream.ConsumerConfig) error); ok {
		r1 = rf(ctx, cfg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
