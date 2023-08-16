// Code generated by mockery v2.32.4. DO NOT EDIT.

package mocks

import (
	context "context"
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// TadoSetter is an autogenerated mock type for the TadoSetter type
type TadoSetter struct {
	mock.Mock
}

type TadoSetter_Expecter struct {
	mock *mock.Mock
}

func (_m *TadoSetter) EXPECT() *TadoSetter_Expecter {
	return &TadoSetter_Expecter{mock: &_m.Mock}
}

// DeleteZoneOverlay provides a mock function with given fields: _a0, _a1
func (_m *TadoSetter) DeleteZoneOverlay(_a0 context.Context, _a1 int) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TadoSetter_DeleteZoneOverlay_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteZoneOverlay'
type TadoSetter_DeleteZoneOverlay_Call struct {
	*mock.Call
}

// DeleteZoneOverlay is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 int
func (_e *TadoSetter_Expecter) DeleteZoneOverlay(_a0 interface{}, _a1 interface{}) *TadoSetter_DeleteZoneOverlay_Call {
	return &TadoSetter_DeleteZoneOverlay_Call{Call: _e.mock.On("DeleteZoneOverlay", _a0, _a1)}
}

func (_c *TadoSetter_DeleteZoneOverlay_Call) Run(run func(_a0 context.Context, _a1 int)) *TadoSetter_DeleteZoneOverlay_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int))
	})
	return _c
}

func (_c *TadoSetter_DeleteZoneOverlay_Call) Return(_a0 error) *TadoSetter_DeleteZoneOverlay_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TadoSetter_DeleteZoneOverlay_Call) RunAndReturn(run func(context.Context, int) error) *TadoSetter_DeleteZoneOverlay_Call {
	_c.Call.Return(run)
	return _c
}

// SetZoneTemporaryOverlay provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *TadoSetter) SetZoneTemporaryOverlay(_a0 context.Context, _a1 int, _a2 float64, _a3 time.Duration) error {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, float64, time.Duration) error); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TadoSetter_SetZoneTemporaryOverlay_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetZoneTemporaryOverlay'
type TadoSetter_SetZoneTemporaryOverlay_Call struct {
	*mock.Call
}

// SetZoneTemporaryOverlay is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 int
//   - _a2 float64
//   - _a3 time.Duration
func (_e *TadoSetter_Expecter) SetZoneTemporaryOverlay(_a0 interface{}, _a1 interface{}, _a2 interface{}, _a3 interface{}) *TadoSetter_SetZoneTemporaryOverlay_Call {
	return &TadoSetter_SetZoneTemporaryOverlay_Call{Call: _e.mock.On("SetZoneTemporaryOverlay", _a0, _a1, _a2, _a3)}
}

func (_c *TadoSetter_SetZoneTemporaryOverlay_Call) Run(run func(_a0 context.Context, _a1 int, _a2 float64, _a3 time.Duration)) *TadoSetter_SetZoneTemporaryOverlay_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int), args[2].(float64), args[3].(time.Duration))
	})
	return _c
}

func (_c *TadoSetter_SetZoneTemporaryOverlay_Call) Return(_a0 error) *TadoSetter_SetZoneTemporaryOverlay_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TadoSetter_SetZoneTemporaryOverlay_Call) RunAndReturn(run func(context.Context, int, float64, time.Duration) error) *TadoSetter_SetZoneTemporaryOverlay_Call {
	_c.Call.Return(run)
	return _c
}

// NewTadoSetter creates a new instance of TadoSetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTadoSetter(t interface {
	mock.TestingT
	Cleanup(func())
}) *TadoSetter {
	mock := &TadoSetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
