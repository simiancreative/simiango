package natsjspub_test

import (
	"github.com/sanity-io/litter"
	"github.com/simiancreative/simiango/circuitbreaker"
	"github.com/stretchr/testify/mock"
)

type MockCircuitBreaker struct {
	mock.Mock
}

func (m *MockCircuitBreaker) Allow() bool {
	args := m.Called()
	litter.Dump(args)
	return args.Bool(0)
}

func (m *MockCircuitBreaker) GetState() circuitbreaker.State {
	args := m.Called()
	return args.Get(0).(circuitbreaker.State)
}

func (m *MockCircuitBreaker) RecordStart() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockCircuitBreaker) RecordResult(success bool) {
	m.Called(success)
}

func (m *MockCircuitBreaker) Reset() {
	m.Called()
}
