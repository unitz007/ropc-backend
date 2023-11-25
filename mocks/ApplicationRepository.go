// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	model "ropc-backend/model"

	mock "github.com/stretchr/testify/mock"
)

// ApplicationRepository is an autogenerated mock type for the ApplicationRepository type
type ApplicationRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: client
func (_m *ApplicationRepository) Create(client *model.Application) error {
	ret := _m.Called(client)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Application) error); ok {
		r0 = rf(client)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: id
func (_m *ApplicationRepository) Delete(id uint) error {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(uint) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAll provides a mock function with given fields: userId
func (_m *ApplicationRepository) GetAll(userId uint) []model.Application {
	ret := _m.Called(userId)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []model.Application
	if rf, ok := ret.Get(0).(func(uint) []model.Application); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Application)
		}
	}

	return r0
}

// GetByClientId provides a mock function with given fields: clientId
func (_m *ApplicationRepository) GetByClientId(clientId string) (*model.Application, error) {
	ret := _m.Called(clientId)

	if len(ret) == 0 {
		panic("no return value specified for GetByClientId")
	}

	var r0 *model.Application
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*model.Application, error)); ok {
		return rf(clientId)
	}
	if rf, ok := ret.Get(0).(func(string) *model.Application); ok {
		r0 = rf(clientId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Application)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(clientId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByClientIdAndUserId provides a mock function with given fields: clientId, userId
func (_m *ApplicationRepository) GetByClientIdAndUserId(clientId string, userId uint) (*model.Application, error) {
	ret := _m.Called(clientId, userId)

	if len(ret) == 0 {
		panic("no return value specified for GetByClientIdAndUserId")
	}

	var r0 *model.Application
	var r1 error
	if rf, ok := ret.Get(0).(func(string, uint) (*model.Application, error)); ok {
		return rf(clientId, userId)
	}
	if rf, ok := ret.Get(0).(func(string, uint) *model.Application); ok {
		r0 = rf(clientId, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Application)
		}
	}

	if rf, ok := ret.Get(1).(func(string, uint) error); ok {
		r1 = rf(clientId, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByName provides a mock function with given fields: name
func (_m *ApplicationRepository) GetByName(name string) (*model.Application, error) {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for GetByName")
	}

	var r0 *model.Application
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*model.Application, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) *model.Application); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Application)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByNameAndUserId provides a mock function with given fields: name, userId
func (_m *ApplicationRepository) GetByNameAndUserId(name string, userId uint) (*model.Application, error) {
	ret := _m.Called(name, userId)

	if len(ret) == 0 {
		panic("no return value specified for GetByNameAndUserId")
	}

	var r0 *model.Application
	var r1 error
	if rf, ok := ret.Get(0).(func(string, uint) (*model.Application, error)); ok {
		return rf(name, userId)
	}
	if rf, ok := ret.Get(0).(func(string, uint) *model.Application); ok {
		r0 = rf(name, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Application)
		}
	}

	if rf, ok := ret.Get(1).(func(string, uint) error); ok {
		r1 = rf(name, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: app
func (_m *ApplicationRepository) Update(app *model.Application) (*model.Application, error) {
	ret := _m.Called(app)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *model.Application
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.Application) (*model.Application, error)); ok {
		return rf(app)
	}
	if rf, ok := ret.Get(0).(func(*model.Application) *model.Application); ok {
		r0 = rf(app)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Application)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.Application) error); ok {
		r1 = rf(app)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewApplicationRepository creates a new instance of ApplicationRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewApplicationRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *ApplicationRepository {
	mock := &ApplicationRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
