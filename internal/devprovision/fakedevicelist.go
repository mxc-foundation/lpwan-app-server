package devprovision

import (
	"time"
)

type deviceInfo struct {
	serialNumber     string
	serialNumberHash string
	manufacturerID   int
	model            string
	fixedDevEUI      bool
	devEUI           string
	appEUI           string
	appKey           string
	nwkKey           string
	status           int
	timeCreated      time.Time
	server           string
}

var fakeDeviceList = []deviceInfo{
	{serialNumber: "SERIALNUMBEROOOOOOOO", serialNumberHash: "34dfcb3dde1a09fd340fafada1e431e84028fc53c328d359a8824613b86d568e",
		manufacturerID: 1, model: "LoRaWatch", fixedDevEUI: true, devEUI: "2462abfffeddc710", timeCreated: time.Now()},
}
