package devprovision

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/user"
	"time"

	"google.golang.org/grpc"

	"github.com/brocaar/lorawan"
	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"

	nsextra "github.com/mxc-foundation/lpwan-app-server/api/networkserver"
	psPb "github.com/mxc-foundation/lpwan-app-server/api/ps-serves-appserver"
	gwd "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway/data"
	nsd "github.com/mxc-foundation/lpwan-app-server/internal/networkserver_portal/data"
)

// Mock for LoRa frame sending
type mockDataType struct {
	request *nsextra.SendDelayedProprietaryPayloadRequest
}

var mockData mockDataType

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
var timeCreated = time.Now()
var mockDeviceList = []deviceInfo{
	{ProvisionID: "SERIALNUMBEROOOOOOOO", ProvisionIDHash: "34dfcb3dde1a09fd340fafada1e431e84028fc53c328d359a8824613b86d568e",
		ManufacturerID: 1, Model: "LoRaWatch", FixedDevEUI: true,
		DevEUI:      []byte{0x24, 0x62, 0xab, 0xff, 0xfe, 0xdd, 0xc7, 0x10},
		AppEUI:      make([]byte, 8),
		AppKey:      make([]byte, 16),
		NwkKey:      make([]byte, 16),
		TimeCreated: &timeCreated},
	{ProvisionID: "TESTPIDOOOOOOOOOOOOO", ProvisionIDHash: "c8c7564b46b91c91ef6c4f37bcca8cf7e81baac6eb869dcc62e5fafdd0242497",
		ManufacturerID: 1, Model: "LoRaWatch", FixedDevEUI: true,
		DevEUI:      []byte{0x11, 0x22, 0x33, 0xff, 0xfe, 0x44, 0x55, 0x66},
		AppEUI:      make([]byte, 8),
		AppKey:      make([]byte, 16),
		NwkKey:      make([]byte, 16),
		Server:      "sn-dev.local",
		TimeCreated: &timeCreated},
}

type deviceInfoForJSON struct {
	ProvisionID     string `json:"provisionId"`
	ProvisionIDHash string `json:"provisionIdHash"`
	DevEUI          string `json:"devEUI"`
	AppEUI          string `json:"appEUI"`
	AppKey          string `json:"appKey"`
	NwkKey          string `json:"nwkKey"`
}

func mockSaveDevice(ctx context.Context, device deviceInfo) error {
	for i := range mockDeviceList {
		if mockDeviceList[i].ProvisionIDHash == device.ProvisionIDHash {
			mockDeviceList[i] = device
			break
		}
	}

	var devlist []deviceInfoForJSON
	for i := range mockDeviceList {
		devlist = append(devlist, deviceInfoForJSON{
			ProvisionID:     mockDeviceList[i].ProvisionID,
			ProvisionIDHash: mockDeviceList[i].ProvisionIDHash,
			DevEUI:          hex.EncodeToString(mockDeviceList[i].DevEUI),
			AppEUI:          hex.EncodeToString(mockDeviceList[i].AppEUI),
			AppKey:          hex.EncodeToString(mockDeviceList[i].AppKey),
			NwkKey:          hex.EncodeToString(mockDeviceList[i].NwkKey),
		})
	}

	targetfile := "devicelist.json"
	user, err := user.Current()
	if err != nil {
		log.Errorf("Error to get current user. %s", err.Error())
	} else {
		targetfile = user.HomeDir + "/devicelist.json"
	}
	log.Debugf("Save device list to %s", targetfile)
	outputbuf, _ := json.MarshalIndent(devlist, "", "  ")
	_ = ioutil.WriteFile(targetfile, outputbuf, 0600)
	return nil
}

type testPsCli struct {
}

func (t testPsCli) GetManufacturerByID(ctx context.Context, in *psPb.GetMfgByIdRequest, opts ...grpc.CallOption) (*psPb.GetMfgResponse, error) {
	return &psPb.GetMfgResponse{}, nil
}

func (t testPsCli) GetManufacturerByName(ctx context.Context, in *psPb.GetMfgByNameRequest, opts ...grpc.CallOption) (*psPb.GetMfgResponse, error) {
	return &psPb.GetMfgResponse{}, nil
}

func (t testPsCli) CreateManufacturer(ctx context.Context, in *psPb.CreateMfgRequest, opts ...grpc.CallOption) (*psPb.CreateMfgResponse, error) {
	return &psPb.CreateMfgResponse{}, nil
}

func (t testPsCli) UpdateManufacturer(ctx context.Context, in *psPb.UpdateMfgRequest, opts ...grpc.CallOption) (*psPb.UpdateMfgResponse, error) {
	return &psPb.UpdateMfgResponse{}, nil
}

func (t testPsCli) IsDeviceExist(ctx context.Context, in *psPb.IsDeviceExistRequest, opts ...grpc.CallOption) (*psPb.IsDeviceExistResponse, error) {
	return &psPb.IsDeviceExistResponse{}, nil
}

func (t testPsCli) GenDevEUI(ctx context.Context, in *psPb.GenDevEuiRequest, opts ...grpc.CallOption) (*psPb.GenDevEuiResponse, error) {
	return &psPb.GenDevEuiResponse{}, nil
}

func (t testPsCli) GenProvisionID(ctx context.Context, in *psPb.GenProvisionIdRequest, opts ...grpc.CallOption) (*psPb.GenProvisionIdResponse, error) {
	return &psPb.GenProvisionIdResponse{}, nil
}

func (t testPsCli) CreateDevice(ctx context.Context, in *psPb.CreateDeviceRequest, opts ...grpc.CallOption) (*psPb.CreateDeviceResponse, error) {
	return &psPb.CreateDeviceResponse{}, nil
}

func (t testPsCli) GetDeviceByID(ctx context.Context, in *psPb.GetDeviceByIdRequest, opts ...grpc.CallOption) (*psPb.GetDeviceResponse, error) {
	return &psPb.GetDeviceResponse{}, nil
}

func (t testPsCli) GetDeviceByIDHash(ctx context.Context, in *psPb.GetDeviceByIdHashRequest, opts ...grpc.CallOption) (*psPb.GetDeviceResponse, error) {
	for i := range mockDeviceList {
		device := mockDeviceList[i]
		if device.ProvisionIDHash == in.ProvisionIdHash {
			return &psPb.GetDeviceResponse{
				ProvisionId:     device.ProvisionID,
				ProvisionIdHash: in.ProvisionIdHash,
				ManufacturerId:  device.ManufacturerID,
				Model:           device.Model,
				SerialNumber:    device.SerialNumber,
				FixedDevEUI:     device.FixedDevEUI,
				DevEUI:          device.DevEUI,
				AppKey:          device.AppKey,
				AppEUI:          device.AppEUI,
				NwkKey:          device.NwkKey,
				Status:          device.Status,
				Server:          device.Server,
			}, nil
		}
	}
	return nil, fmt.Errorf("no device found")
}

func (t testPsCli) UpdateDeviceInfo(ctx context.Context, in *psPb.UpdateDeviceInfoRequest, opts ...grpc.CallOption) (*psPb.UpdateDeviceResponse, error) {
	return &psPb.UpdateDeviceResponse{}, nil
}

func (t testPsCli) SetDeviceProvisioned(ctx context.Context, in *psPb.SetDeviceProvisionedRequest, opts ...grpc.CallOption) (*psPb.UpdateDeviceResponse, error) {
	for i := range mockDeviceList {
		device := mockDeviceList[i]
		if device.ProvisionID == in.ProvisionId {
			device.DevEUI = in.DevEUI
			device.AppEUI = in.AppEUI
			device.AppKey = in.AppKey
			device.NwkKey = in.NwkKey
			break
		}
	}

	var devlist []deviceInfoForJSON
	for i := range mockDeviceList {
		devlist = append(devlist, deviceInfoForJSON{
			ProvisionID:     mockDeviceList[i].ProvisionID,
			ProvisionIDHash: mockDeviceList[i].ProvisionIDHash,
			DevEUI:          hex.EncodeToString(mockDeviceList[i].DevEUI),
			AppEUI:          hex.EncodeToString(mockDeviceList[i].AppEUI),
			AppKey:          hex.EncodeToString(mockDeviceList[i].AppKey),
			NwkKey:          hex.EncodeToString(mockDeviceList[i].NwkKey),
		})
	}

	targetfile := "devicelist.json"
	user, err := user.Current()
	if err != nil {
		log.Errorf("Error to get current user. %s", err.Error())
	} else {
		targetfile = user.HomeDir + "/devicelist.json"
	}
	log.Debugf("Save device list to %s", targetfile)
	outputbuf, _ := json.MarshalIndent(devlist, "", "  ")
	_ = ioutil.WriteFile(targetfile, outputbuf, 0600)
	return &psPb.UpdateDeviceResponse{}, nil
}

func (t testPsCli) SetDeviceServer(ctx context.Context, in *psPb.SetDeviceServerRequest, opts ...grpc.CallOption) (*psPb.UpdateDeviceResponse, error) {
	return &psPb.UpdateDeviceResponse{}, nil
}

type testDb struct {
}

func (t testDb) GetGateway(ctx context.Context, mac lorawan.EUI64, forUpdate bool) (gwd.Gateway, error) {
	return gwd.Gateway{NetworkServerID: 1}, nil
}

func (t testDb) GetNetworkServer(ctx context.Context, id int64) (nsd.NetworkServer, error) {
	return nsd.NetworkServer{
		ID: 1,
	}, nil
}

type testNsCli struct {
}

func (t testNsCli) GetNetworkServerExtraServiceClient(networkServerID int64) (nsextra.NetworkServerExtraServiceClient, error) {
	return &networkServerExtraServiceClient{}, nil
}

type networkServerExtraServiceClient struct {
}

func (n networkServerExtraServiceClient) SendDelayedProprietaryPayload(ctx context.Context,
	in *nsextra.SendDelayedProprietaryPayloadRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	var response empty.Empty
	mockData.request = in
	return &response, nil
}
