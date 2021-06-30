package devprovision

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/apex/log"

	psPb "github.com/mxc-foundation/lpwan-app-server/api/ps-serves-appserver"
)

type deviceInfo struct {
	ProvisionID     string     `json:"provisionId"`
	ProvisionIDHash string     `json:"provisionIdHash"`
	ManufacturerID  int64      `json:"manufacturerId"`
	Model           string     `json:"model"`
	SerialNumber    string     `json:"serialNumber"`
	FixedDevEUI     bool       `json:"fixedDevEUI"`
	DevEUI          []byte     `json:"devEUI"`
	AppEUI          []byte     `json:"appEUI"`
	AppKey          []byte     `json:"appKey"`
	NwkKey          []byte     `json:"nwkKey"`
	Status          string     `json:"status"`
	Server          string     `json:"server"`
	TimeCreated     *time.Time `json:"timeCreated"`
	TimeProvisioned *time.Time `json:"timeProvisioned"`
	TimeAddToServer *time.Time `json:"timeAddToServer"`
}

func findDeviceBySnHash(ctx context.Context, provisionIdhash []byte, psCli psPb.DeviceProvisionClient) (bool, deviceInfo) {
	resp, err := psCli.GetDeviceByIDHash(ctx, &psPb.GetDeviceByIdHashRequest{
		ProvisionIdHash: hex.EncodeToString(provisionIdhash)})
	if err != nil {
		log.WithError(err).Errorf("Failed to get device, hash: %s", hex.EncodeToString(provisionIdhash))
		return false, deviceInfo{}
	}

	retdevice := deviceInfo{
		ProvisionID:     resp.ProvisionId,
		ProvisionIDHash: resp.ProvisionIdHash,
		ManufacturerID:  resp.ManufacturerId,
		Model:           resp.Model,
		SerialNumber:    resp.SerialNumber,
		FixedDevEUI:     resp.FixedDevEUI,
		DevEUI:          make([]byte, 8),
		AppEUI:          make([]byte, 8),
		AppKey:          make([]byte, 16),
		NwkKey:          make([]byte, 16),
		Status:          resp.Status,
		Server:          resp.Server,
	}
	copy(retdevice.DevEUI[:], resp.DevEUI)
	copy(retdevice.AppEUI[:], resp.AppEUI)
	copy(retdevice.AppKey[:], resp.AppKey)
	copy(retdevice.NwkKey[:], resp.NwkKey)
	if resp.TimeCreated != nil {
		time := resp.TimeCreated.AsTime()
		retdevice.TimeCreated = &time
	}
	if resp.TimeProvisioned != nil {
		time := resp.TimeProvisioned.AsTime()
		retdevice.TimeProvisioned = &time
	}
	if resp.TimeAddToServer != nil {
		time := resp.TimeAddToServer.AsTime()
		retdevice.TimeAddToServer = &time
	}
	return true, retdevice
}

func saveDevice(ctx context.Context, device deviceInfo, psCli psPb.DeviceProvisionClient) error {
	_, err := psCli.SetDeviceProvisioned(ctx, &psPb.SetDeviceProvisionedRequest{
		ProvisionId: device.ProvisionID,
		DevEUI:      device.DevEUI,
		AppEUI:      device.AppEUI,
		AppKey:      device.AppKey,
		NwkKey:      device.NwkKey,
	})
	if err != nil {
		return err
	}

	return nil
}
