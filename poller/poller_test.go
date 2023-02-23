package poller_test

import (
	"context"
	"github.com/clambin/tado"
	"github.com/clambin/tado-exporter/poller"
	tado2 "github.com/clambin/tado-exporter/tado"
	"github.com/clambin/tado-exporter/tado/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func prepareMockAPI(api *mocks.API) {
	api.
		On("GetMobileDevices", mock.Anything).
		Return([]tado.MobileDevice{
			{
				ID:       1,
				Name:     "foo",
				Settings: tado.MobileDeviceSettings{GeoTrackingEnabled: true},
				Location: tado.MobileDeviceLocation{AtHome: true},
			},
			{
				ID:       2,
				Name:     "bar",
				Settings: tado.MobileDeviceSettings{GeoTrackingEnabled: true},
				Location: tado.MobileDeviceLocation{AtHome: false},
			}}, nil).
		Once()
	api.
		On("GetWeatherInfo", mock.Anything).
		Return(tado.WeatherInfo{
			OutsideTemperature: tado.Temperature{Celsius: 3.4},
			SolarIntensity:     tado.Percentage{Percentage: 13.3},
			WeatherState:       tado.Value{Value: "CLOUDY_MOSTLY"},
		}, nil).
		Once()
	api.On("GetZones", mock.Anything).
		Return(tado.Zones{
			{ID: 1, Name: "foo"},
			{ID: 2, Name: "bar"},
		}, nil).
		Once()
	api.
		On("GetZoneInfo", mock.Anything, 1).
		Return(tado.ZoneInfo{
			Setting: tado.ZonePowerSetting{
				Power:       "ON",
				Temperature: tado.Temperature{Celsius: 18.5},
			},
		}, nil).
		Once()
	api.
		On("GetZoneInfo", mock.Anything, 2).
		Return(tado.ZoneInfo{
			Setting: tado.ZonePowerSetting{
				Power: "OFF",
			},
			Overlay: tado.ZoneInfoOverlay{
				Type:        "MANUAL",
				Termination: tado.ZoneInfoOverlayTermination{Type: "MANUAL"},
			},
		}, nil).
		Once()
	api.
		On("GetHomeState", mock.Anything).
		Return(tado.HomeState{Presence: "HOME"}, nil).
		Once()
}

func TestPoller_Run(t *testing.T) {
	api := mocks.NewAPI(t)

	p := poller.New(api)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	prepareMockAPI(api)

	go p.Run(ctx, 10*time.Millisecond)
	ch := p.Register()
	update := <-ch

	require.Len(t, update.UserInfo, 2)
	assert.Equal(t, "foo", update.UserInfo[1].Name)
	device := update.UserInfo[1]
	assert.Equal(t, tado.DeviceHome, (&device).IsHome())
	assert.Equal(t, "bar", update.UserInfo[2].Name)
	device = update.UserInfo[2]
	assert.Equal(t, tado.DeviceAway, (&device).IsHome())

	assert.Equal(t, "CLOUDY_MOSTLY", update.WeatherInfo.WeatherState.Value)
	assert.Equal(t, 3.4, update.WeatherInfo.OutsideTemperature.Celsius)
	assert.Equal(t, 13.3, update.WeatherInfo.SolarIntensity.Percentage)

	require.Len(t, update.Zones, 2)
	assert.Equal(t, "foo", update.Zones[1].Name)
	assert.Equal(t, "bar", update.Zones[2].Name)

	require.Len(t, update.ZoneInfo, 2)
	info := update.ZoneInfo[1]
	assert.Equal(t, tado2.ZoneStateAuto, tado2.GetZoneState(info))
	info = update.ZoneInfo[2]
	assert.Equal(t, tado2.ZoneStateOff, tado2.GetZoneState(info))

	assert.True(t, update.Home)

	p.Unregister(ch)
}

func TestServer_Poll(t *testing.T) {
	api := mocks.NewAPI(t)
	prepareMockAPI(api)

	p := poller.New(api)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go p.Run(ctx, time.Minute)

	ch := p.Register()
	p.Refresh()
	update := <-ch

	require.Len(t, update.UserInfo, 2)

	p.Unregister(ch)
}

func TestServer_Refresh(t *testing.T) {
	api := mocks.NewAPI(t)
	prepareMockAPI(api)

	p := poller.New(api)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := p.Register()
	go p.Run(ctx, 10*time.Millisecond)

	update := <-ch
	require.Len(t, update.UserInfo, 2)

	prepareMockAPI(api)
	p.Refresh()
	update = <-ch
	require.Len(t, update.UserInfo, 2)
}
