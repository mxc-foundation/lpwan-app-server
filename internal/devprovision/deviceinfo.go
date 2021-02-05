package devprovision

import (
	"context"
	"database/sql"
	"encoding/hex"

	"github.com/apex/log"

	psPb "github.com/mxc-foundation/lpwan-app-server/api/ps-serves-appserver"
	pscli "github.com/mxc-foundation/lpwan-app-server/internal/clients/psconn"
)

type deviceInfo struct {
	ProvisionID     string       `json:"provisionId"`
	ProvisionIDHash string       `json:"provisionIdHash"`
	ManufacturerID  int64        `json:"manufacturerId"`
	Model           string       `json:"model"`
	SerialNumber    string       `json:"serialNumber"`
	FixedDevEUI     bool         `json:"fixedDevEUI"`
	DevEUI          []byte       `json:"devEUI"`
	AppEUI          []byte       `json:"appEUI"`
	AppKey          []byte       `json:"appKey"`
	NwkKey          []byte       `json:"nwkKey"`
	Status          string       `json:"status"`
	Server          string       `json:"server"`
	TimeCreated     sql.NullTime `json:"timeCreated"`
	TimeProvisioned sql.NullTime `json:"timeProvisioned"`
	TimeAddToServer sql.NullTime `json:"timeAddToServer"`
}

//
func findDeviceBySnHash(ctx context.Context, provisionIdhash []byte) (bool, deviceInfo) {
	psClient, err := pscli.GetDevProClient()
	if err != nil {
		log.WithError(err).Errorf("find device failed.")
		return false, deviceInfo{}
	}

	resp, err := psClient.GetDeviceByIDHash(ctx, &psPb.GetDeviceByIdHashRequest{
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
		retdevice.TimeCreated.Valid = true
		retdevice.TimeCreated.Time = resp.TimeCreated.AsTime()
	}
	if resp.TimeProvisioned != nil {
		retdevice.TimeProvisioned.Valid = true
		retdevice.TimeProvisioned.Time = resp.TimeProvisioned.AsTime()
	}
	if resp.TimeAddToServer != nil {
		retdevice.TimeAddToServer.Valid = true
		retdevice.TimeAddToServer.Time = resp.TimeAddToServer.AsTime()
	}
	return true, retdevice
}

//
func saveDevice(ctx context.Context, device deviceInfo) error {
	psClient, err := pscli.GetDevProClient()
	if err != nil {
		return err
	}

	_, err = psClient.SetDeviceProvisioned(ctx, &psPb.SetDeviceProvisionedRequest{
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
