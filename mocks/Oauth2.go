// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	model "ropc-backend/model"

	mock "github.com/stretchr/testify/mock"
)

// Oauth2 is an autogenerated mock type for the Oauth2 type
type Oauth2 struct {
	mock.Mock
}

// ClientCredentials provides a mock function with given fields: clientId, clientSecret
func (_m *Oauth2) ClientCredentials(clientId string, clientSecret string) (*model.AccessToken, error) {
	ret := _m.Called(clientId, clientSecret)

	if len(ret) == 0 {
		panic("no return value specified for ClientCredentials")
	}

	var r0 *model.AccessToken
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (*model.AccessToken, error)); ok {
		return rf(clientId, clientSecret)
	}
	if rf, ok := ret.Get(0).(func(string, string) *model.AccessToken); ok {
		r0 = rf(clientId, clientSecret)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AccessToken)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(clientId, clientSecret)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewOauth2 creates a new instance of Oauth2. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOauth2(t interface {
	mock.TestingT
	Cleanup(func())
}) *Oauth2 {
	mock := &Oauth2{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
