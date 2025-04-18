package tadotools

import (
	"context"
	"fmt"
	"github.com/clambin/tado-exporter/internal/oapi"
	"github.com/clambin/tado/v2"
	"net/http"
	"time"
)

type TadoClient interface {
	SetZoneOverlayWithResponse(ctx context.Context, homeId tado.HomeId, zoneId tado.ZoneId, body tado.SetZoneOverlayJSONRequestBody, reqEditors ...tado.RequestEditorFn) (*tado.SetZoneOverlayResponse, error)
}

func SetOverlay(ctx context.Context, c TadoClient, homeId tado.HomeId, zoneId tado.ZoneId, temperature float32, duration time.Duration) error {
	// possibly set power to "off" if temp <= 5?
	req := tado.SetZoneOverlayJSONRequestBody{
		Setting: &tado.ZoneSetting{
			Type: oapi.VarP(tado.HEATING),
		},
		Termination: &tado.ZoneOverlayTermination{
			Type: oapi.VarP(tado.ZoneOverlayTerminationTypeMANUAL),
		},
	}
	if temperature < 5 {
		req.Setting.Power = oapi.VarP(tado.PowerOFF)
	} else {
		req.Setting.Power = oapi.VarP(tado.PowerON)
		req.Setting.Temperature = &tado.Temperature{Celsius: oapi.VarP(temperature)}
	}
	if duration > 0 {
		req.Termination.Type = oapi.VarP(tado.ZoneOverlayTerminationTypeTIMER)
		req.Termination.DurationInSeconds = oapi.VarP(int(duration.Seconds()))
	}
	resp, err := c.SetZoneOverlayWithResponse(ctx, homeId, zoneId, req)
	if err == nil && resp.StatusCode() != http.StatusOK {
		err = fmt.Errorf("SetZoneOverlayWithResponse: %s", resp.Status())
	}
	return err
}
