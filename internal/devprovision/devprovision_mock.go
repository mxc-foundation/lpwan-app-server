package devprovision

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os/user"
	"time"

	"github.com/brocaar/lorawan"
	log "github.com/sirupsen/logrus"

	nsextra "github.com/mxc-foundation/lpwan-app-server/api/ns-extra"
	gwd "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway/data"
	nsd "github.com/mxc-foundation/lpwan-app-server/internal/networkserver_portal/data"
)

// Mock functions for ctrl.handler
type handlerMock struct {
}

func (h *handlerMock) GetGateway(ctx context.Context, mac lorawan.EUI64, forUpdate bool) (gwd.Gateway, error) {
	return gwd.Gateway{Name: "MockGateway"}, nil
}

func (h *handlerMock) GetNetworkServer(ctx context.Context, id int64) (nsd.NetworkServer, error) {
	return nsd.NetworkServer{Name: "MockNetworkServer", Server: "0.0.0.0"}, nil
}

// Mock for LoRa frame sending
type mockDataType struct {
	request *nsextra.SendDelayedProprietaryPayloadRequest
}

var mockData mockDataType

func sendToNsMock(n nsd.NetworkServer, req *nsextra.SendDelayedProprietaryPayloadRequest) error {
	mockData.request = req

	return nil
}

// Mock of get current time
var mockNowQueue []time.Time

func mockGetNow() time.Time {
	if len(mockNowQueue) == 0 {
		return time.Time{}
	}
	poptime := mockNowQueue[0]
	if len(mockNowQueue) > 1 {
		mockNowQueue = mockNowQueue[1:]
	}
	return poptime
}

// Mock of random buf generation. Set -1 to use pesudorandom
var mockRandValue int

func mockGen128Rand() []byte {
	softrand := softRand{}
	randbuf := make([]byte, 128)
	for i := range randbuf {
		if mockRandValue < 0 {
			randbuf[i] = uint8(softrand.Get())
		} else {
			randbuf[i] = uint8(mockRandValue)
		}
	}

	return randbuf
}

// Mock device list
var mockDeviceList = []deviceInfo{
	{SerialNumber: "SERIALNUMBEROOOOOOOO", SerialNumberHash: "34dfcb3dde1a09fd340fafada1e431e84028fc53c328d359a8824613b86d568e",
		ManufacturerID: 1, Model: "LoRaWatch", FixedDevEUI: true, DevEUI: "2462abfffeddc710", TimeCreated: time.Now()},
}

func mockFindDeviceBySnHash(serialnumberhash []byte) (bool, deviceInfo) {
	strhash := hex.EncodeToString(serialnumberhash)

	for i := range mockDeviceList {
		if mockDeviceList[i].SerialNumberHash == strhash {
			return true, mockDeviceList[i]
		}
	}
	return false, deviceInfo{}
}

func mockSaveDevice(deviceinfo deviceInfo) error {
	for i := range mockDeviceList {
		if mockDeviceList[i].SerialNumberHash == deviceinfo.SerialNumberHash {
			mockDeviceList[i] = deviceinfo
			break
		}
	}

	targetfile := "devicelist.json"
	user, err := user.Current()
	if err != nil {
		log.Errorf("Error to get current user. %s", err.Error())
	} else {
		targetfile = user.HomeDir + "/devicelist.json"
	}
	log.Debugf("Save device list to %s", targetfile)
	outputbuf, _ := json.MarshalIndent(mockDeviceList, "", "  ")
	_ = ioutil.WriteFile(targetfile, outputbuf, 0600)
	return nil
}
