package controller

import (
	"github.com/clambin/tado-exporter/pkg/tado"
	log "github.com/sirupsen/logrus"
	"time"
)

// runOverlayLimit checks for new overlays and expires any that have exceeded their limit
func (controller *Controller) runOverlayLimit() error {
	if controller.Rules.OverlayLimit == nil {
		return nil
	}

	if controller.Overlays == nil {
		controller.Overlays = make(map[int]time.Time)
	}

	var (
		err     error
		actions []action
	)

	if err = controller.updateOverlays(); err == nil {
		if actions, err = controller.expireOverlays(); err == nil {
			for _, action := range actions {
				if err = controller.runAction(action); err == nil {
					break
				}
			}
		}
	}

	log.WithFields(log.Fields{
		"err":      err,
		"overlays": len(controller.Overlays),
	}).Debug("runOverlayLimit")

	return err
}

// updateOverlays gets any overlays that are currently active and stores them in controller.Overlays.
// Any zones that are not in overlay are removed from controller.Overlays
func (controller *Controller) updateOverlays() error {
	var (
		err   error
		zones []tado.Zone
	)

	if zones, err = controller.GetZones(); err == nil {
		for _, overlayLimit := range *controller.Rules.OverlayLimit {
			var (
				zone     *tado.Zone
				zoneInfo *tado.ZoneInfo
			)
			if zone = getZone(zones, overlayLimit.ZoneID, overlayLimit.ZoneName); zone == nil {
				log.WithFields(log.Fields{
					"ZoneID":   overlayLimit.ZoneID,
					"ZoneName": overlayLimit.ZoneName,
				}).Warning("skipping unknown zone in OverlayLimit rule")
				continue
			}

			if zoneInfo, err = controller.GetZoneInfo(zone.ID); err == nil {
				if zoneInfo.Overlay.Type == "MANUAL" && zoneInfo.Overlay.Setting.Type == "HEATING" {
					// Zone in overlay. If we're not already tracking it, add it now
					if _, ok := controller.Overlays[zone.ID]; ok == false {
						expiry := time.Now().Add(overlayLimit.MaxTime)
						controller.Overlays[zone.ID] = expiry
						log.WithFields(log.Fields{
							"zoneID": zone.ID,
							"expiry": expiry,
						}).Info("new zone in overlay")
					}
				} else {
					// Zone is not in overlay. Remove it from the tracking map
					if _, ok := controller.Overlays[zone.ID]; ok == true {
						delete(controller.Overlays, zone.ID)
						log.WithField("zoneID", zone.ID).Info("zone no longer in overlay")
					}
				}
			}

		}
	}

	return err
}

// expireOverlays deletes any overlays that have expired
func (controller *Controller) expireOverlays() ([]action, error) {
	var (
		err     error
		actions = make([]action, 0)
	)
	for zoneID, expiryTimer := range controller.Overlays {
		if time.Now().After(expiryTimer) {
			actions = append(actions, action{
				Overlay: false,
				ZoneID:  zoneID,
			})
			log.WithField("zoneID", zoneID).Info("expiring overlay in zone")
			// Technically not needed (next run will do this automatically, but facilitates unit testing
			delete(controller.Overlays, zoneID)
		}
	}
	return actions, err
}
