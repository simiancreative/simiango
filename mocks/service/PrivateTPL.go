// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	meta "github.com/simiancreative/simiango/meta"
	mock "github.com/stretchr/testify/mock"

	service "github.com/simiancreative/simiango/service"
)

// PrivateTPL is an autogenerated mock type for the PrivateTPL type
type PrivateTPL struct {
	mock.Mock
}

// Auth provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *PrivateTPL) Auth(_a0 meta.RequestId, _a1 service.RawHeaders, _a2 service.RawBody, _a3 service.RawParams) error {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 error
	if rf, ok := ret.Get(0).(func(meta.RequestId, service.RawHeaders, service.RawBody, service.RawParams) error); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Result provides a mock function with given fields:
func (_m *PrivateTPL) Result() (interface{}, error) {
	ret := _m.Called()

	var r0 interface{}
	if rf, ok := ret.Get(0).(func() interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
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

type mockConstructorTestingTNewPrivateTPL interface {
	mock.TestingT
	Cleanup(func())
}

// NewPrivateTPL creates a new instance of PrivateTPL. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPrivateTPL(t mockConstructorTestingTNewPrivateTPL) *PrivateTPL {
	mock := &PrivateTPL{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
