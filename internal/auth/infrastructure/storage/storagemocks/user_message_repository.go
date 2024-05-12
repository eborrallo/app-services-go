// Code generated by mockery v2.39.1. DO NOT EDIT.

package storagemocks

import mock "github.com/stretchr/testify/mock"

// UserMessageRepository is an autogenerated mock type for the UserMessageRepository type
type UserMessageRepository struct {
	mock.Mock
}

// GetMessage provides a mock function with given fields: address
func (_m *UserMessageRepository) GetMessage(address string) (string, error) {
	ret := _m.Called(address)

	if len(ret) == 0 {
		panic("no return value specified for GetMessage")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(address)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(address)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveMessage provides a mock function with given fields: address, message
func (_m *UserMessageRepository) SaveMessage(address string, message string) {
	_m.Called(address, message)
}

// NewUserMessageRepository creates a new instance of UserMessageRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserMessageRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserMessageRepository {
	mock := &UserMessageRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}