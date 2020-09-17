package storage

import (
	"context"
	"time"

	"github.com/brocaar/lorawan"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/brocaar/chirpstack-api/go/v3/ns"

	nscli "github.com/mxc-foundation/lpwan-app-server/internal/clients/networkserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/logging"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

// Device defines a LoRaWAN device.
type Device store.Device

// DeviceListItem defines the Device as list item.
type DeviceListItem store.DeviceListItem

// Validate validates the device data.
func (d Device) Validate() error {
	return nil
}

// DeviceKeys defines the keys for a LoRaWAN device.
type DeviceKeys store.DeviceKeys

// DevicesActiveInactive holds the active and inactive counts.
type DevicesActiveInactive store.DevicesActiveInactive

// DevicesDataRates holds the device counts by data-rate.
type DevicesDataRates store.DevicesDataRates

// CreateDevice creates the given device.
func CreateDevice(ctx context.Context, handler *store.Handler, d *Device) error {
	return handler.CreateDevice(ctx, (*store.Device)(d))
}

// GetDevice returns the device matching the given DevEUI.
// When forUpdate is set to true, then db must be a db transaction.
// When localOnly is set to true, no call to the network-server is made to
// retrieve additional device data.
func GetDevice(ctx context.Context, handler *store.Handler, devEUI lorawan.EUI64, forUpdate, localOnly bool) (Device, error) {
	d, err := handler.GetDevice(ctx, devEUI, forUpdate)
	if err != nil {
		return Device(d), err
	}

	if localOnly {
		return Device(d), nil
	}

	n, err := handler.GetNetworkServerForDevEUI(ctx, d.DevEUI)
	if err != nil {
		return Device(d), errors.Wrap(err, "get network-server error")
	}

	nsStruct := nscli.NSStruct{
		Server:  n.Server,
		CACert:  n.CACert,
		TLSCert: n.TLSCert,
		TLSKey:  n.TLSKey,
	}
	nsClient, err := nsStruct.GetNetworkServiceClient()
	if err != nil {
		return Device(d), errors.Wrap(err, "get network-server client error")
	}

	resp, err := nsClient.GetDevice(ctx, &ns.GetDeviceRequest{
		DevEui: d.DevEUI[:],
	})
	if err != nil {
		return Device(d), err
	}

	if resp.Device != nil {
		d.SkipFCntCheck = resp.Device.SkipFCntCheck
		d.ReferenceAltitude = resp.Device.ReferenceAltitude
		d.IsDisabled = resp.Device.IsDisabled
	}

	return Device(d), nil
}

// DeviceFilters provide filters that can be used to filter on devices.
// Note that empty values are not used as filter.
type DeviceFilters store.DeviceFilters

// SQL returns the SQL filter.
func (f DeviceFilters) SQL() string {
	df := store.DeviceFilters(f)
	return df.SQL()
}

// GetDeviceCount returns the number of devices.
func GetDeviceCount(ctx context.Context, handler *store.Handler, filters DeviceFilters) (int, error) {
	return handler.GetDeviceCount(ctx, store.DeviceFilters(filters))
}

// GetDevices returns a slice of devices.
func GetDevices(ctx context.Context, handler *store.Handler, filters DeviceFilters) ([]DeviceListItem, error) {
	res, err := handler.GetDevices(ctx, store.DeviceFilters(filters))
	if err != nil {
		return nil, err
	}

	var devices []DeviceListItem
	for _, v := range res {
		deviceItem := DeviceListItem(v)
		devices = append(devices, deviceItem)
	}

	return devices, nil
}

// UpdateDevice updates the given device.
// When localOnly is set, it will not update the device on the network-server.
func UpdateDevice(ctx context.Context, handler *store.Handler, d *Device, localOnly bool) error {
	err := handler.UpdateDevice(ctx, (*store.Device)(d))
	if err != nil {
		return err
	}

	// update the device on the network-server
	if localOnly {
		return nil
	}

	app, err := handler.GetApplication(ctx, d.ApplicationID)
	if err != nil {
		return errors.Wrap(err, "get application error")
	}

	n, err := handler.GetNetworkServerForDevEUI(ctx, d.DevEUI)
	if err != nil {
		return errors.Wrap(err, "get network-server error")
	}

	nsStruct := nscli.NSStruct{
		Server:  n.Server,
		CACert:  n.CACert,
		TLSCert: n.TLSCert,
		TLSKey:  n.TLSKey,
	}
	nsClient, err := nsStruct.GetNetworkServiceClient()
	if err != nil {
		return errors.Wrap(err, "get network-server client error")
	}

	rpID := ctrl.applicationServerID
	_, err = nsClient.UpdateDevice(ctx, &ns.UpdateDeviceRequest{
		Device: &ns.Device{
			DevEui:            d.DevEUI[:],
			DeviceProfileId:   d.DeviceProfileID.Bytes(),
			ServiceProfileId:  app.ServiceProfileID.Bytes(),
			RoutingProfileId:  rpID.Bytes(),
			SkipFCntCheck:     d.SkipFCntCheck,
			ReferenceAltitude: d.ReferenceAltitude,
			IsDisabled:        d.IsDisabled,
		},
	})
	if err != nil {
		return errors.Wrap(err, "update device error")
	}

	log.WithFields(log.Fields{
		"dev_eui": d.DevEUI,
		"ctx_id":  ctx.Value(logging.ContextIDKey),
	}).Info("device updated")

	return nil
}

// UpdateDeviceLastSeenAndDR updates the device last-seen timestamp and data-rate.
func UpdateDeviceLastSeenAndDR(ctx context.Context, handler *store.Handler, devEUI lorawan.EUI64, ts time.Time, dr int) error {
	return handler.UpdateDeviceLastSeenAndDR(ctx, devEUI, ts, dr)
}

// UpdateDeviceActivation updates the device address and the AppSKey.
func UpdateDeviceActivation(ctx context.Context, handler *store.Handler, devEUI lorawan.EUI64, devAddr lorawan.DevAddr, appSKey lorawan.AES128Key) error {
	return handler.UpdateDeviceActivation(ctx, devEUI, devAddr, appSKey)
}

// DeleteDevice deletes the device matching the given DevEUI.
func DeleteDevice(ctx context.Context, handler *store.Handler, devEUI lorawan.EUI64) error {
	return handler.DeleteDevice(ctx, devEUI)
}

// CreateDeviceKeys creates the keys for the given device.
func CreateDeviceKeys(ctx context.Context, handler *store.Handler, dc *DeviceKeys) error {
	return handler.CreateDeviceKeys(ctx, (*store.DeviceKeys)(dc))
}

// GetDeviceKeys returns the device-keys for the given DevEUI.
func GetDeviceKeys(ctx context.Context, handler *store.Handler, devEUI lorawan.EUI64) (DeviceKeys, error) {
	devKeys, err := handler.GetDeviceKeys(ctx, devEUI)
	return DeviceKeys(devKeys), err
}

// UpdateDeviceKeys updates the given device-keys.
func UpdateDeviceKeys(ctx context.Context, handler *store.Handler, dc *DeviceKeys) error {
	return handler.UpdateDeviceKeys(ctx, (*store.DeviceKeys)(dc))
}

// DeleteDeviceKeys deletes the device-keys for the given DevEUI.
func DeleteDeviceKeys(ctx context.Context, handler *store.Handler, devEUI lorawan.EUI64) error {
	return handler.DeleteDeviceKeys(ctx, devEUI)
}

// DeleteAllDevicesForApplicationID deletes all devices given an application id.
func DeleteAllDevicesForApplicationID(ctx context.Context, handler *store.Handler, applicationID int64) error {
	return handler.DeleteAllDevicesForApplicationID(ctx, applicationID)
}

// EnqueueDownlinkPayload adds the downlink payload to the network-server
// device-queue.
func EnqueueDownlinkPayload(ctx context.Context, handler *store.Handler, devEUI lorawan.EUI64, confirmed bool, fPort uint8, data []byte) (uint32, error) {
	return handler.EnqueueDownlinkPayload(ctx, devEUI, confirmed, fPort, data)
}

// GetDevicesActiveInactive returns the active / inactive devices.
func GetDevicesActiveInactive(ctx context.Context, handler *store.Handler, organizationID int64) (DevicesActiveInactive, error) {
	res, err := handler.GetDevicesActiveInactive(ctx, organizationID)
	return DevicesActiveInactive(res), err
}

// GetDevicesDataRates returns the device counts by data-rate.
func GetDevicesDataRates(ctx context.Context, handler *store.Handler, organizationID int64) (DevicesDataRates, error) {
	res, err := handler.GetDevicesDataRates(ctx, organizationID)
	return DevicesDataRates(res), err
}
