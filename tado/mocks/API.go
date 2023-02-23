// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	context "context"

	clambintado "github.com/clambin/tado"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// API is an autogenerated mock type for the API type
type API struct {
	mock.Mock
}

// DeleteZoneOverlay provides a mock function with given fields: _a0, _a1
func (_m *API) DeleteZoneOverlay(_a0 context.Context, _a1 int) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetHomeState provides a mock function with given fields: ctx
func (_m *API) GetHomeState(ctx context.Context) (clambintado.HomeState, error) {
	ret := _m.Called(ctx)

	var r0 clambintado.HomeState
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (clambintado.HomeState, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) clambintado.HomeState); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(clambintado.HomeState)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMobileDevices provides a mock function with given fields: _a0
func (_m *API) GetMobileDevices(_a0 context.Context) ([]clambintado.MobileDevice, error) {
	ret := _m.Called(_a0)

	var r0 []clambintado.MobileDevice
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]clambintado.MobileDevice, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []clambintado.MobileDevice); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]clambintado.MobileDevice)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWeatherInfo provides a mock function with given fields: _a0
func (_m *API) GetWeatherInfo(_a0 context.Context) (clambintado.WeatherInfo, error) {
	ret := _m.Called(_a0)

	var r0 clambintado.WeatherInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (clambintado.WeatherInfo, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(context.Context) clambintado.WeatherInfo); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(clambintado.WeatherInfo)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetZoneInfo provides a mock function with given fields: _a0, _a1
func (_m *API) GetZoneInfo(_a0 context.Context, _a1 int) (clambintado.ZoneInfo, error) {
	ret := _m.Called(_a0, _a1)

	var r0 clambintado.ZoneInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (clambintado.ZoneInfo, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) clambintado.ZoneInfo); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(clambintado.ZoneInfo)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetZones provides a mock function with given fields: _a0
func (_m *API) GetZones(_a0 context.Context) (clambintado.Zones, error) {
	ret := _m.Called(_a0)

	var r0 clambintado.Zones
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (clambintado.Zones, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(context.Context) clambintado.Zones); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(clambintado.Zones)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetZoneOverlay provides a mock function with given fields: _a0, _a1, _a2
func (_m *API) SetZoneOverlay(_a0 context.Context, _a1 int, _a2 float64) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, float64) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetZoneTemporaryOverlay provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *API) SetZoneTemporaryOverlay(_a0 context.Context, _a1 int, _a2 float64, _a3 time.Duration) error {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, float64, time.Duration) error); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewAPI interface {
	mock.TestingT
	Cleanup(func())
}

// NewAPI creates a new instance of API. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAPI(t mockConstructorTestingTNewAPI) *API {
	mock := &API{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
