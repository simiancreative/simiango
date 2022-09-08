// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	crud "github.com/simiancreative/simiango/data/sql/crud"
	mock "github.com/stretchr/testify/mock"
)

// FilterableFunc is an autogenerated mock type for the FilterableFunc type
type FilterableFunc struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0
func (_m *FilterableFunc) Execute(_a0 interface{}) (crud.WhereKind, interface{}) {
	ret := _m.Called(_a0)

	var r0 crud.WhereKind
	if rf, ok := ret.Get(0).(func(interface{}) crud.WhereKind); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(crud.WhereKind)
	}

	var r1 interface{}
	if rf, ok := ret.Get(1).(func(interface{}) interface{}); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(interface{})
		}
	}

	return r0, r1
}

type mockConstructorTestingTNewFilterableFunc interface {
	mock.TestingT
	Cleanup(func())
}

// NewFilterableFunc creates a new instance of FilterableFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewFilterableFunc(t mockConstructorTestingTNewFilterableFunc) *FilterableFunc {
	mock := &FilterableFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}