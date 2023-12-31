// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// Server is an autogenerated mock type for the Server type
type Server struct {
	mock.Mock
}

// RegisterHandler provides a mock function with given fields: path, method, handler
func (_m *Server) RegisterHandler(path string, method string, handler func(http.ResponseWriter, *http.Request)) {
	_m.Called(path, method, handler)
}

// Start provides a mock function with given fields: address
func (_m *Server) Start(address string) error {
	ret := _m.Called(address)

	if len(ret) == 0 {
		panic("no return value specified for Start")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(address)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewServer creates a new instance of Server. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *Server {
	mock := &Server{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
