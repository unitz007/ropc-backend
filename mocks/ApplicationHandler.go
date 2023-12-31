// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// ApplicationHandler is an autogenerated mock type for the ApplicationHandler type
type ApplicationHandler struct {
	mock.Mock
}

// CreateApplication provides a mock function with given fields: w, r
func (_m *ApplicationHandler) CreateApplication(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// DeleteApplication provides a mock function with given fields: response, request
func (_m *ApplicationHandler) DeleteApplication(response http.ResponseWriter, request *http.Request) {
	_m.Called(response, request)
}

// GenerateSecret provides a mock function with given fields: w, r
func (_m *ApplicationHandler) GenerateSecret(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// GetApplication provides a mock function with given fields: w, r
func (_m *ApplicationHandler) GetApplication(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// GetApplications provides a mock function with given fields: w, r
func (_m *ApplicationHandler) GetApplications(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// NewApplicationHandler creates a new instance of ApplicationHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewApplicationHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *ApplicationHandler {
	mock := &ApplicationHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
