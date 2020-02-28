// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import mock "github.com/stretchr/testify/mock"

// DirectorClient is an autogenerated mock type for the DirectorClient type
type DirectorClient struct {
	mock.Mock
}

// GetConsoleURL provides a mock function with given fields: accountID, runtimeID
func (_m *DirectorClient) GetConsoleURL(accountID string, runtimeID string) (string, error) {
	ret := _m.Called(accountID, runtimeID)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(accountID, runtimeID)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(accountID, runtimeID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
