package logger

import "github.com/stretchr/testify/mock"

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debugf(fmt string, args ...any) {
	m.Called(args)
}

func (m *MockLogger) Errorf(fmt string, args ...any) {
	m.Called(args)
}

func (m *MockLogger) Debug(args ...any) {
	m.Called(args)
}

func (m *MockLogger) Error(args ...any) {
	m.Called(args)
}
