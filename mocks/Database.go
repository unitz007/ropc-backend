// Code generated by mockery v2.37.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Database is an autogenerated mock type for the Database type
type Database[T interface{}] struct {
	mock.Mock
}

// GetDatabaseConnection provides a mock function with given fields:
func (_m *Database[T]) GetDatabaseConnection() *T {
	ret := _m.Called()

	var r0 *T
	if rf, ok := ret.Get(0).(func() *T); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*T)
		}
	}

	return r0
}

// NewDatabase creates a new instance of Database. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDatabase[T interface{}](t interface {
	mock.TestingT
	Cleanup(func())
}) *Database[T] {
	mock := &Database[T]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
