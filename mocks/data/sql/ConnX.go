// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	databasesql "database/sql"

	mock "github.com/stretchr/testify/mock"

	sqlx "github.com/jmoiron/sqlx"
)

// ConnX is an autogenerated mock type for the ConnX type
type ConnX struct {
	mock.Mock
}

// Beginx provides a mock function with given fields:
func (_m *ConnX) Beginx() (*sqlx.Tx, error) {
	ret := _m.Called()

	var r0 *sqlx.Tx
	if rf, ok := ret.Get(0).(func() *sqlx.Tx); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sqlx.Tx)
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

// BindNamed provides a mock function with given fields: _a0, _a1
func (_m *ConnX) BindNamed(_a0 string, _a1 interface{}) (string, []interface{}, error) {
	ret := _m.Called(_a0, _a1)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, interface{}) string); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 []interface{}
	if rf, ok := ret.Get(1).(func(string, interface{}) []interface{}); ok {
		r1 = rf(_a0, _a1)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]interface{})
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string, interface{}) error); ok {
		r2 = rf(_a0, _a1)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Get provides a mock function with given fields: _a0, _a1, _a2
func (_m *ConnX) Get(_a0 interface{}, _a1 string, _a2 ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, _a0, _a1)
	_ca = append(_ca, _a2...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}, string, ...interface{}) error); ok {
		r0 = rf(_a0, _a1, _a2...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MapperFunc provides a mock function with given fields: _a0
func (_m *ConnX) MapperFunc(_a0 func(string) string) {
	_m.Called(_a0)
}

// MustBegin provides a mock function with given fields:
func (_m *ConnX) MustBegin() *sqlx.Tx {
	ret := _m.Called()

	var r0 *sqlx.Tx
	if rf, ok := ret.Get(0).(func() *sqlx.Tx); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sqlx.Tx)
		}
	}

	return r0
}

// MustExec provides a mock function with given fields: _a0, _a1
func (_m *ConnX) MustExec(_a0 string, _a1 ...interface{}) databasesql.Result {
	var _ca []interface{}
	_ca = append(_ca, _a0)
	_ca = append(_ca, _a1...)
	ret := _m.Called(_ca...)

	var r0 databasesql.Result
	if rf, ok := ret.Get(0).(func(string, ...interface{}) databasesql.Result); ok {
		r0 = rf(_a0, _a1...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(databasesql.Result)
		}
	}

	return r0
}

// NamedExec provides a mock function with given fields: _a0, _a1
func (_m *ConnX) NamedExec(_a0 string, _a1 interface{}) (databasesql.Result, error) {
	ret := _m.Called(_a0, _a1)

	var r0 databasesql.Result
	if rf, ok := ret.Get(0).(func(string, interface{}) databasesql.Result); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(databasesql.Result)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, interface{}) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NamedQuery provides a mock function with given fields: _a0, _a1
func (_m *ConnX) NamedQuery(_a0 string, _a1 interface{}) (*sqlx.Rows, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *sqlx.Rows
	if rf, ok := ret.Get(0).(func(string, interface{}) *sqlx.Rows); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sqlx.Rows)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, interface{}) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PrepareNamed provides a mock function with given fields: _a0
func (_m *ConnX) PrepareNamed(_a0 string) (*sqlx.NamedStmt, error) {
	ret := _m.Called(_a0)

	var r0 *sqlx.NamedStmt
	if rf, ok := ret.Get(0).(func(string) *sqlx.NamedStmt); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sqlx.NamedStmt)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Preparex provides a mock function with given fields: _a0
func (_m *ConnX) Preparex(_a0 string) (*sqlx.Stmt, error) {
	ret := _m.Called(_a0)

	var r0 *sqlx.Stmt
	if rf, ok := ret.Get(0).(func(string) *sqlx.Stmt); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sqlx.Stmt)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Query provides a mock function with given fields: query, args
func (_m *ConnX) Query(query string, args ...interface{}) (*databasesql.Rows, error) {
	var _ca []interface{}
	_ca = append(_ca, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 *databasesql.Rows
	if rf, ok := ret.Get(0).(func(string, ...interface{}) *databasesql.Rows); ok {
		r0 = rf(query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*databasesql.Rows)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, ...interface{}) error); ok {
		r1 = rf(query, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryRowx provides a mock function with given fields: query, args
func (_m *ConnX) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	var _ca []interface{}
	_ca = append(_ca, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 *sqlx.Row
	if rf, ok := ret.Get(0).(func(string, ...interface{}) *sqlx.Row); ok {
		r0 = rf(query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sqlx.Row)
		}
	}

	return r0
}

// Queryx provides a mock function with given fields: query, args
func (_m *ConnX) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	var _ca []interface{}
	_ca = append(_ca, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 *sqlx.Rows
	if rf, ok := ret.Get(0).(func(string, ...interface{}) *sqlx.Rows); ok {
		r0 = rf(query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sqlx.Rows)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, ...interface{}) error); ok {
		r1 = rf(query, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Rebind provides a mock function with given fields: _a0
func (_m *ConnX) Rebind(_a0 string) string {
	ret := _m.Called(_a0)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Select provides a mock function with given fields: _a0, _a1, _a2
func (_m *ConnX) Select(_a0 interface{}, _a1 string, _a2 ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, _a0, _a1)
	_ca = append(_ca, _a2...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}, string, ...interface{}) error); ok {
		r0 = rf(_a0, _a1, _a2...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Unsafe provides a mock function with given fields:
func (_m *ConnX) Unsafe() *sqlx.DB {
	ret := _m.Called()

	var r0 *sqlx.DB
	if rf, ok := ret.Get(0).(func() *sqlx.DB); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sqlx.DB)
		}
	}

	return r0
}