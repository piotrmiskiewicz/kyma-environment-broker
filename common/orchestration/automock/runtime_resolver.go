// Code generated by mockery v2.14.0. DO NOT EDIT.

package automock

import (
	orchestration "github.com/kyma-project/control-plane/components/kyma-environment-broker/common/orchestration"
	mock "github.com/stretchr/testify/mock"
)

// RuntimeResolver is an autogenerated mock type for the RuntimeResolver type
type RuntimeResolver struct {
	mock.Mock
}

// Resolve provides a mock function with given fields: targets
func (_m *RuntimeResolver) Resolve(targets orchestration.TargetSpec) ([]orchestration.Runtime, error) {
	ret := _m.Called(targets)

	var r0 []orchestration.Runtime
	if rf, ok := ret.Get(0).(func(orchestration.TargetSpec) []orchestration.Runtime); ok {
		r0 = rf(targets)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]orchestration.Runtime)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(orchestration.TargetSpec) error); ok {
		r1 = rf(targets)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewRuntimeResolver interface {
	mock.TestingT
	Cleanup(func())
}

// NewRuntimeResolver creates a new instance of RuntimeResolver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRuntimeResolver(t mockConstructorTestingTNewRuntimeResolver) *RuntimeResolver {
	mock := &RuntimeResolver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
