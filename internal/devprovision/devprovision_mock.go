package devprovision

import (
	"context"
	"time"

	"github.com/brocaar/lorawan"

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
