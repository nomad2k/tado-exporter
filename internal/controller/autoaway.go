package controller

import (
	"github.com/clambin/tado-exporter/pkg/tado"
	log "github.com/sirupsen/logrus"
	"time"
)

// AutoAwayInfo contains the user we are tracking, and what zone to set to which temperature
// when ActivationTime occurs
type AutoAwayInfo struct {
	MobileDevice      *tado.MobileDevice
	Home              bool
	ActivationTime    time.Time
	ZoneID            int
	TargetTemperature float64
}

// runAutoAway runOverlayLimit checks if mobileDevices have come/left home and performs
// configured autoAway rules
func (controller *Controller) runAutoAway() error {
	if controller.Rules.AutoAway == nil {
		return nil
	}

	var (
		err     error
		actions []action
	)

	// update mobiles & zones for each autoaway entry
	if err = controller.updateAutoAwayInfo(); err == nil {
		// get actions for each autoaway setting
		if actions, err = controller.getAutoAwayActions(); err == nil {
			for _, action := range actions {
				// execute each action
				if err = controller.runAction(action); err != nil {
					break
				}
			}
		}
	}

	log.WithField("err", err).Debug("runAutoAway")
	return err
}

// updateAutoAwayInfo updates the mobile device & zone information for each autoAway rule.
// On exit, the map controller.AutoAwayInfo contains the up to date mobileDevice information
// for any mobile device mentioned in any autoAway rule.
func (controller *Controller) updateAutoAwayInfo() error {
	var (
		err           error
		mobileDevices []tado.MobileDevice
		zones         []tado.Zone
	)

	// If the map doesn't exist, create it
	if controller.AutoAwayInfo == nil {
		controller.AutoAwayInfo = make(map[int]AutoAwayInfo)
	}

	// get info we will need
	if mobileDevices, err = controller.GetMobileDevices(); err == nil {
		zones, err = controller.GetZones()
	}

	if err == nil {
		// for each autoaway setting, add/update a record for the mobileDevice
		for _, autoAway := range *controller.Rules.AutoAway {
			var (
				mobileDevice *tado.MobileDevice
				zone         *tado.Zone
			)
			// Rules file can contain either mobileDevice/zone ID or Name. Retrieve the ID for each of these
			// and discard any that aren't valid
			if mobileDevice = getMobileDevice(mobileDevices, autoAway.MobileDeviceID, autoAway.MobileDeviceName); mobileDevice == nil {
				log.WithFields(log.Fields{
					"deviceID":   autoAway.MobileDeviceID,
					"deviceName": autoAway.MobileDeviceName,
				}).Warning("skipping unknown mobile device in AutoAway rule")
				continue
			}
			if zone = getZone(zones, autoAway.ZoneID, autoAway.ZoneName); zone == nil {
				log.WithFields(log.Fields{
					"zoneID":   autoAway.ZoneID,
					"zoneName": autoAway.ZoneName,
				}).Warning("skipping unknown zone in AutoAway rule")
				continue
			}

			// Add/update the entry in the AutoAwayInfo map
			if entry, ok := controller.AutoAwayInfo[mobileDevice.ID]; ok == false {
				// We don't already have a record. Create it
				controller.AutoAwayInfo[mobileDevice.ID] = AutoAwayInfo{
					MobileDevice:      mobileDevice,
					Home:              mobileDevice.Location.AtHome,
					ActivationTime:    time.Now().Add(autoAway.WaitTime),
					ZoneID:            zone.ID,
					TargetTemperature: autoAway.TargetTemperature,
				}
			} else {
				// If we already have it, update it
				entry.MobileDevice = mobileDevice
			}
		}
	}

	return err
}

// getAutoAwayActions scans the AutoAwayInfo map and returns all required actions, i.e. any zones that
// need to be put in/out of Overlay mode.
func (controller *Controller) getAutoAwayActions() ([]action, error) {
	var (
		err     error
		actions = make([]action, 0)
	)

	for id, autoAway := range controller.AutoAwayInfo {
		log.WithFields(log.Fields{
			"mobileDeviceID":   autoAway.MobileDevice.ID,
			"mobileDeviceName": autoAway.MobileDevice.Name,
			"new_home":         autoAway.MobileDevice.Location.AtHome,
			"old_home":         autoAway.Home,
		}).Debug("autoAwayInfo")

		// if the mobile phone is now home but was away
		if autoAway.MobileDevice.Location.AtHome && !autoAway.Home {
			// mark the phone at home
			autoAway.Home = true
			controller.AutoAwayInfo[id] = autoAway
			// add action to disable the overlay
			actions = append(actions, action{
				Overlay: false,
				ZoneID:  autoAway.ZoneID,
			})
			log.WithFields(log.Fields{
				"MobileDeviceID": id,
				"ZoneID":         autoAway.ZoneID,
			}).Info("User returned home. Removing overlay")
		} else
		// if the mobile phone is away
		if !autoAway.MobileDevice.Location.AtHome {
			// if the phone was home, mark the phone away & record the time
			if autoAway.Home {
				autoAway.Home = false
				controller.AutoAwayInfo[id] = autoAway
			}
			// if the phone's been away for the required time
			if time.Now().After(autoAway.ActivationTime) {
				// add action to set the overlay
				actions = append(actions, action{
					Overlay:           true,
					ZoneID:            autoAway.ZoneID,
					TargetTemperature: autoAway.TargetTemperature,
				})
				log.WithFields(log.Fields{
					"MobileDeviceID":    id,
					"ZoneID":            autoAway.ZoneID,
					"TargetTemperature": autoAway.TargetTemperature,
				}).Info("User left. Setting overlay")
			}
		}
	}
	return actions, err
}
