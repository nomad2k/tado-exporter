package rules

import (
	"github.com/clambin/tado"
	"github.com/clambin/tado-exporter/poller"
	tado2 "github.com/clambin/tado-exporter/tado"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestAutoAwayRule_Evaluate(t *testing.T) {
	tests := []testCase{
		{
			name: "user goes away",
			update: &poller.Update{
				Zones:    map[int]tado.Zone{10: {ID: 10, Name: "living room"}},
				ZoneInfo: map[int]tado.ZoneInfo{10: {Setting: tado.ZonePowerSetting{Power: "ON"}}},
				UserInfo: map[int]tado.MobileDevice{100: {ID: 100, Name: "foo", Settings: tado.MobileDeviceSettings{GeoTrackingEnabled: true}, Location: tado.MobileDeviceLocation{AtHome: false}}},
			},
			action: NextState{ZoneID: 10, ZoneName: "living room", State: tado2.ZoneStateOff, Delay: time.Hour, ActionReason: "foo is away", CancelReason: "foo is home"},
		},
		{
			name: "user comes home",
			update: &poller.Update{
				Zones: map[int]tado.Zone{10: {ID: 10, Name: "living room"}},
				ZoneInfo: map[int]tado.ZoneInfo{10: {
					Setting: tado.ZonePowerSetting{Power: "OFF"},
					Overlay: tado.ZoneInfoOverlay{
						Type:        "MANUAL",
						Termination: tado.ZoneInfoOverlayTermination{Type: "MANUAL"},
					}}},
				UserInfo: map[int]tado.MobileDevice{100: {ID: 100, Name: "foo", Settings: tado.MobileDeviceSettings{GeoTrackingEnabled: true}, Location: tado.MobileDeviceLocation{AtHome: true}}},
			},
			action: NextState{ZoneID: 10, ZoneName: "living room", State: tado2.ZoneStateAuto, Delay: 0, ActionReason: "foo is home", CancelReason: "foo is away"},
		},
		{
			name: "user is home",
			update: &poller.Update{
				Zones:    map[int]tado.Zone{10: {ID: 10, Name: "living room"}},
				ZoneInfo: map[int]tado.ZoneInfo{10: {Setting: tado.ZonePowerSetting{Power: "ON"}}},
				UserInfo: map[int]tado.MobileDevice{100: {ID: 100, Name: "foo", Settings: tado.MobileDeviceSettings{GeoTrackingEnabled: true}, Location: tado.MobileDeviceLocation{AtHome: true}}},
			},
		},
		{
			// TODO: not quite sure why this rule was here, or why it even worked
			name: "non-geolocation user",
			update: &poller.Update{
				Zones: map[int]tado.Zone{10: {ID: 10, Name: "living room"}},
				ZoneInfo: map[int]tado.ZoneInfo{10: {
					Setting: tado.ZonePowerSetting{Power: "ON", Temperature: tado.Temperature{Celsius: 15.0}},
					Overlay: tado.ZoneInfoOverlay{
						Type:        "MANUAL",
						Termination: tado.ZoneInfoOverlayTermination{Type: "MANUAL"},
					}}},
				UserInfo: map[int]tado.MobileDevice{100: {ID: 100, Name: "foo", Settings: tado.MobileDeviceSettings{GeoTrackingEnabled: false}}},
			},
			//action: NextState{ZoneID: 10, ZoneName: "living room", State: tado.ZoneStateAuto, Delay: time.Hour, ActionReason: " are home", CancelReason: " are away"},
		},
	}

	r := &AutoAwayRule{
		zoneID:   10,
		zoneName: "living room",
		delay:    time.Hour,
		users:    []string{"foo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := r.Evaluate(tt.update)
			require.NoError(t, err)
			assert.Equal(t, tt.action, a)
		})
	}
}

func TestAutoAwayRule_Evaluate_MultipleUsers(t *testing.T) {
	tests := []testCase{
		{
			name: "one user goes away",
			update: &poller.Update{
				Zones:    map[int]tado.Zone{10: {ID: 10, Name: "living room"}},
				ZoneInfo: map[int]tado.ZoneInfo{10: {Setting: tado.ZonePowerSetting{Power: "ON"}}},
				UserInfo: map[int]tado.MobileDevice{
					100: {ID: 100, Name: "foo", Settings: tado.MobileDeviceSettings{GeoTrackingEnabled: true}, Location: tado.MobileDeviceLocation{AtHome: false}},
					110: {ID: 100, Name: "bar", Settings: tado.MobileDeviceSettings{GeoTrackingEnabled: true}, Location: tado.MobileDeviceLocation{AtHome: true}},
				},
			},
		},
		{
			name: "all users are away",
			update: &poller.Update{
				Zones:    map[int]tado.Zone{10: {ID: 10, Name: "living room"}},
				ZoneInfo: map[int]tado.ZoneInfo{10: {Setting: tado.ZonePowerSetting{Power: "ON"}}},
				UserInfo: map[int]tado.MobileDevice{
					100: {ID: 100, Name: "foo", Settings: tado.MobileDeviceSettings{GeoTrackingEnabled: true}, Location: tado.MobileDeviceLocation{AtHome: false}},
					110: {ID: 100, Name: "bar", Settings: tado.MobileDeviceSettings{GeoTrackingEnabled: true}, Location: tado.MobileDeviceLocation{AtHome: false}},
				},
			},
			action: NextState{ZoneID: 10, ZoneName: "living room", State: tado2.ZoneStateOff, Delay: time.Hour, ActionReason: "foo, bar are away", CancelReason: "foo, bar are home"},
		},
		{
			name: "one user is home",
			update: &poller.Update{
				Zones: map[int]tado.Zone{10: {ID: 10, Name: "living room"}},
				ZoneInfo: map[int]tado.ZoneInfo{10: {
					Setting: tado.ZonePowerSetting{Power: "OFF"},
					Overlay: tado.ZoneInfoOverlay{
						Type:        "MANUAL",
						Termination: tado.ZoneInfoOverlayTermination{Type: "MANUAL"},
					}}},
				UserInfo: map[int]tado.MobileDevice{
					100: {ID: 100, Name: "foo", Settings: tado.MobileDeviceSettings{GeoTrackingEnabled: true}, Location: tado.MobileDeviceLocation{AtHome: false}},
					110: {ID: 100, Name: "bar", Settings: tado.MobileDeviceSettings{GeoTrackingEnabled: true}, Location: tado.MobileDeviceLocation{AtHome: true}},
				},
			},
			action: NextState{ZoneID: 10, ZoneName: "living room", State: tado2.ZoneStateAuto, Delay: 0, ActionReason: "bar is home", CancelReason: "bar is away"},
		},
	}

	r := &AutoAwayRule{
		zoneID:   10,
		zoneName: "living room",
		delay:    time.Hour,
		users:    []string{"foo", "bar"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := r.Evaluate(tt.update)
			require.NoError(t, err)
			assert.Equal(t, tt.action, a)
		})
	}
}
