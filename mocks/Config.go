// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Config is an autogenerated mock type for the Config type
type Config struct {
	mock.Mock
}

// AppMode provides a mock function with given fields:
func (_m *Config) AppMode() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for AppMode")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// DatabaseHost provides a mock function with given fields:
func (_m *Config) DatabaseHost() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for DatabaseHost")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// DatabaseName provides a mock function with given fields:
func (_m *Config) DatabaseName() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for DatabaseName")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// DatabasePassword provides a mock function with given fields:
func (_m *Config) DatabasePassword() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for DatabasePassword")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// DatabasePort provides a mock function with given fields:
func (_m *Config) DatabasePort() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for DatabasePort")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// DatabaseUser provides a mock function with given fields:
func (_m *Config) DatabaseUser() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for DatabaseUser")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Mux provides a mock function with given fields:
func (_m *Config) Mux() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Mux")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewRelicAppName provides a mock function with given fields:
func (_m *Config) NewRelicAppName() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for NewRelicAppName")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewRelicLicense provides a mock function with given fields:
func (_m *Config) NewRelicLicense() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for NewRelicLicense")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// ServerPort provides a mock function with given fields:
func (_m *Config) ServerPort() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ServerPort")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// TokenExpiry provides a mock function with given fields:
func (_m *Config) TokenExpiry() int {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for TokenExpiry")
	}

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// TokenSecret provides a mock function with given fields:
func (_m *Config) TokenSecret() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for TokenSecret")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewConfig creates a new instance of Config. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewConfig(t interface {
	mock.TestingT
	Cleanup(func())
}) *Config {
	mock := &Config{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
