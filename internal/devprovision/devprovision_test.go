package devprovision

import (
	"bytes"
	"context"
	"testing"

	"github.com/brocaar/chirpstack-api/go/v3/as"
	gwV3 "github.com/brocaar/chirpstack-api/go/v3/gw"

	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

func TestSetup(t *testing.T) {
	h := store.Handler{}
	err := Setup("UnitTest", &h)
	if err != nil {
		t.Error("Setup() failed.")
	}
	if ctrl.handler != &h {
		t.Error("ctrl.handler not correctly set.")
	}
	if ctrl.handlerMock != nil {
		t.Error("ctrl.handlerMock should not set.")
	}
	if ctrl.networkServerMock != nil {
		t.Error("ctrl.networkServerMock should not set.")
	}
}

func TestHandleReceivedFrameValidHello(t *testing.T) {
	ns := networkServerMock{}
	ctx := context.Background()

	if setupUnitTest(&handlerMock{}, &ns) != nil {
		t.Error("Uint test setup failed.")
	}

	rxInfo := []*gwV3.UplinkRXInfo{
		&gwV3.UplinkRXInfo{
			GatewayId: []byte{0x00, 0x00, 0x00, 0xff, 0xfe, 0x00, 0x00, 0x00},
			Rssi:      -11,
			Context:   []byte{'0', '0', '0', '0', '0', '0'},
		},
		&gwV3.UplinkRXInfo{
			GatewayId: []byte{0x00, 0x00, 0x00, 0xff, 0xfe, 0x00, 0x00, 0x01},
			Rssi:      -10,
			Context:   []byte{'0', '0', '0', '0', '0', '1'},
		},
	}

	request := as.HandleProprietaryUplinkRequest{
		MacPayload: []byte{0x01, 0x02, 0x03, 0x04},
		Mic:        []byte{0x00, 0x00, 0x00, 0x00},
		TxInfo:     &gwV3.UplinkTXInfo{},
		RxInfo:     rxInfo,
	}

	//
	processed, err := HandleReceivedFrame(ctx, &request)
	if err != nil {
		t.Errorf("HandleReceivedFrame failed. %s", err)
	}
	if !processed {
		t.Error("Request not processed.")
	}

	// Check request
	expectedRxInfo := rxInfo[1]
	if !bytes.Equal(ns.request.GatewayMacs[0], expectedRxInfo.GatewayId) {
		t.Error("Highest RSSI gateway not selected.")
	}
	if !bytes.Equal(ns.request.Context, expectedRxInfo.Context) {
		t.Error("Context is mismatched.")
	}
	if !ns.request.PolarizationInversion {
		t.Error("ipol is wrong, it should be true always")
	}
	if ns.request.Delay.Seconds != 5 {
		t.Errorf("Delay should be 5, but got %d", ns.request.Delay.Seconds)
	}
}

func TestHandleReceivedFrameValidAuth(t *testing.T) {
	ns := networkServerMock{}
	ctx := context.Background()

	if setupUnitTest(&handlerMock{}, &ns) != nil {
		t.Error("Uint test setup failed.")
	}

	rxInfo := []*gwV3.UplinkRXInfo{
		&gwV3.UplinkRXInfo{
			GatewayId: []byte{0x00, 0x00, 0x00, 0xff, 0xfe, 0x00, 0x00, 0x00},
			Rssi:      -11,
			Context:   []byte{'0', '0', '0', '0', '0', '0'},
		},
		&gwV3.UplinkRXInfo{
			GatewayId: []byte{0x00, 0x00, 0x00, 0xff, 0xfe, 0x00, 0x00, 0x01},
			Rssi:      -10,
			Context:   []byte{'0', '0', '0', '0', '0', '1'},
		},
	}

	request := as.HandleProprietaryUplinkRequest{
		MacPayload: []byte{0x11, 0x02, 0x03, 0x04},
		Mic:        []byte{0x00, 0x00, 0x00, 0x00},
		TxInfo:     &gwV3.UplinkTXInfo{},
		RxInfo:     rxInfo,
	}

	//
	processed, err := HandleReceivedFrame(ctx, &request)
	if err != nil {
		t.Errorf("HandleReceivedFrame failed. %s", err)
	}
	if !processed {
		t.Error("Request not processed.")
	}

}

func TestHandleReceivedFrameUnknownMsg(t *testing.T) {
	ns := networkServerMock{}
	ctx := context.Background()

	if setupUnitTest(&handlerMock{}, &ns) != nil {
		t.Error("Uint test setup failed.")
	}

	rxInfo := []*gwV3.UplinkRXInfo{
		&gwV3.UplinkRXInfo{
			GatewayId: []byte{0x00, 0x00, 0x00, 0xff, 0xfe, 0x00, 0x00, 0x00},
			Rssi:      -11,
			Context:   []byte{'0', '0', '0', '0', '0', '0'},
		},
		&gwV3.UplinkRXInfo{
			GatewayId: []byte{0x00, 0x00, 0x00, 0xff, 0xfe, 0x00, 0x00, 0x01},
			Rssi:      -10,
			Context:   []byte{'0', '0', '0', '0', '0', '1'},
		},
	}

	request := as.HandleProprietaryUplinkRequest{
		MacPayload: []byte{0x00, 0x02, 0x03, 0x04},
		Mic:        []byte{0x00, 0x00, 0x00, 0x00},
		TxInfo:     &gwV3.UplinkTXInfo{},
		RxInfo:     rxInfo,
	}

	//
	processed, err := HandleReceivedFrame(ctx, &request)
	if err != nil {
		t.Errorf("HandleReceivedFrame failed. %s", err)
	}
	if processed {
		t.Error("Expected the frame is not processed but it did.")
	}

}
func TestHandleReceivedFrameNoRxInfo(t *testing.T) {
	ns := networkServerMock{}
	ctx := context.Background()

	if setupUnitTest(&handlerMock{}, &ns) != nil {
		t.Error("Uint test setup failed.")
	}

	rxInfo := []*gwV3.UplinkRXInfo{}

	request := as.HandleProprietaryUplinkRequest{
		MacPayload: []byte{0x01, 0x02, 0x03, 0x04},
		Mic:        []byte{0x00, 0x00, 0x00, 0x00},
		TxInfo:     &gwV3.UplinkTXInfo{},
		RxInfo:     rxInfo,
	}

	//
	processed, err := HandleReceivedFrame(ctx, &request)
	if err == nil {
		t.Error("Expected fail but it passed.")
	}
	if processed {
		t.Error("Expected the frame is not processed but it did.")
	}
}
