package devprovision

import (
	"time"
)

type deviceInfo struct {
	ProvisionID     string    `json:"provisionId"`
	ProvisionIDHash string    `json:"provisionIdHash"`
	ManufacturerID  int       `json:"manufacturerId"`
	Model           string    `json:"model"`
	FixedDevEUI     bool      `json:"fixedDevEUI"`
	DevEUI          string    `json:"devEUI"`
	AppEUI          string    `json:"appEUI"`
	AppKey          string    `json:"appKey"`
	NwkKey          string    `json:"nwkKey"`
	Status          int       `json:"status"`
	TimeCreated     time.Time `json:"timeCreated"`
	Server          string    `json:"server"`
}
