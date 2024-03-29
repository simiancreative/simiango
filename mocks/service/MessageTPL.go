// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	service "github.com/simiancreative/simiango/service"
	mock "github.com/stretchr/testify/mock"
)

// MessageTPL is an autogenerated mock type for the MessageTPL type
type MessageTPL struct {
	mock.Mock
}

// Result provides a mock function with given fields:
func (_m *MessageTPL) Result() (service.Messages, error) {
	ret := _m.Called()

	var r0 service.Messages
	if rf, ok := ret.Get(0).(func() service.Messages); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.Messages)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMessageTPL interface {
	mock.TestingT
	Cleanup(func())
}

// NewMessageTPL creates a new instance of MessageTPL. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMessageTPL(t mockConstructorTestingTNewMessageTPL) *MessageTPL {
	mock := &MessageTPL{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
