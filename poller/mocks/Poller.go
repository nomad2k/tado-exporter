// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	context "context"

	poller "github.com/clambin/tado-exporter/poller"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// Poller is an autogenerated mock type for the Poller type
type Poller struct {
	mock.Mock
}

// Refresh provides a mock function with given fields:
func (_m *Poller) Refresh() {
	_m.Called()
}

// Register provides a mock function with given fields:
func (_m *Poller) Register() chan *poller.Update {
	ret := _m.Called()

	var r0 chan *poller.Update
	if rf, ok := ret.Get(0).(func() chan *poller.Update); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan *poller.Update)
		}
	}

	return r0
}

// Run provides a mock function with given fields: ctx, interval
func (_m *Poller) Run(ctx context.Context, interval time.Duration) {
	_m.Called(ctx, interval)
}

// Unregister provides a mock function with given fields: ch
func (_m *Poller) Unregister(ch chan *poller.Update) {
	_m.Called(ch)
}

type mockConstructorTestingTNewPoller interface {
	mock.TestingT
	Cleanup(func())
}

// NewPoller creates a new instance of Poller. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPoller(t mockConstructorTestingTNewPoller) *Poller {
	mock := &Poller{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
