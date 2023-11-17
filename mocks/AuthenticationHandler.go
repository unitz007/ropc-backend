// Code generated by mockery v2.37.1. DO NOT EDIT.

package mocks

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// AuthenticationHandler is an autogenerated mock type for the AuthenticationHandler type
type AuthenticationHandler struct {
	mock.Mock
}

// Authenticate provides a mock function with given fields: w, r
func (_m *AuthenticationHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// NewAuthenticationHandler creates a new instance of AuthenticationHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAuthenticationHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *AuthenticationHandler {
	mock := &AuthenticationHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
