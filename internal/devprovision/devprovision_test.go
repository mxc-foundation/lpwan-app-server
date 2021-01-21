package devprovision

import (
	"bytes"
	"context"
	"encoding/hex"
	"os"
	"testing"
	"time"

	"github.com/brocaar/chirpstack-api/go/v3/as"
	gwV3 "github.com/brocaar/chirpstack-api/go/v3/gw"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	ecdh "github.com/mxc-foundation/lpwan-app-server/internal/devprovision/ecdh"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

// Mock data
const expectedRxInfoIdx = 1

var mockRxInfo = []*gwV3.UplinkRXInfo{
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
var rDevEui = []byte{0x81, 0x82, 0x83, 0xff, 0xfe, 0x84, 0x85, 0x86}

//
func prepareHelloMessage(rdeveui []byte, publickey []byte) []byte {
	payload := []byte{0x01}
	payload = append(payload, rdeveui...)
	payload = append(payload, publickey...)
	payload = append(payload, []byte{0x01}...)
	return payload
}

func prepareAuthMessage(rdeveui []byte, serialNumberHash []byte, verifyCode []byte, devNonce []byte) []byte {
	payload := []byte{0x11}
	payload = append(payload, rdeveui...)
	payload = append(payload, serialNumberHash...)
	payload = append(payload, verifyCode...)
	payload = append(payload, devNonce...)
	return payload
}

func isByteArrayAllZero(input []byte) bool {
	for i := range input {
		if input[i] != 0 {
			return false
		}
	}
	return true
}

func fillByteArray(input []byte, value uint8, len int) {
	for i := range input {
		input[i] = value
	}
}

//
func TestMain(m *testing.M) {
	// Set log level
	log.SetLevel(logrus.DebugLevel)

	// Setup mock funcs
	funcGetNow = mockGetNow
	funcGen128Rand = mockGen128Rand

	//reduce max number of session
	maxNumberOfDevSession = 10

	code := m.Run()
	os.Exit(code)
}

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
	if ctrl.sendToNsFunc == nil {
		t.Error("ctrl.sendToNsFunc not set.")
	}
}

func extractHelloResponse(macpayload []byte) (bool, []byte, []byte, []byte) {
	retok := false
	rdeveui := make([]byte, 8)
	serverpubkey := make([]byte, ecdh.K233PubKeySize)
	servernonce := make([]byte, 4)

	if len(macpayload) == 77 {
		copy(rdeveui[:], macpayload[1:])
		copy(serverpubkey[:], macpayload[9:])
		copy(servernonce[:], macpayload[73:])
		retok = true
	}
	return retok, rdeveui, serverpubkey, servernonce
}

func TestHandleReceivedFrameValidHello(t *testing.T) {
	ctx := context.Background()

	if setupUnitTest(&handlerMock{}) != nil {
		t.Error("Uint test setup failed.")
	}

	// Prepare keys
	var ecdhK233 ecdh.K233
	//	var devicepublickey []byte
	privkeyrand := make([]byte, ecdh.K233PrvKeySize)
	fillByteArray(privkeyrand, 0x01, len(privkeyrand))

	_, devicepublickey := ecdhK233.GenerateKeys(privkeyrand)
	if devicepublickey == nil {
		t.Error("ecdhK233.GenerateKeys() failed.")
	}

	//
	request := as.HandleProprietaryUplinkRequest{
		MacPayload: prepareHelloMessage(rDevEui, devicepublickey),
		Mic:        []byte{0x00, 0x00, 0x00, 0x00},
		TxInfo:     &gwV3.UplinkTXInfo{},
		RxInfo:     mockRxInfo,
	}

	mockRandValue = -1
	downserverpubkey := make([][]uint8, 2)
	downservernonce := make([][]uint8, 2)
	for i := 0; i < 2; i++ {
		// Handle uplink request
		processed, err := HandleReceivedFrame(ctx, &request)
		if err != nil {
			t.Errorf("HandleReceivedFrame failed. %s", err)
		}
		if !processed {
			t.Error("Request not processed.")
		}

		// Check downlink request sent
		expectedRxInfo := mockRxInfo[expectedRxInfoIdx]
		if !bytes.Equal(mockData.request.GatewayMacs[0], expectedRxInfo.GatewayId) {
			t.Error("Highest RSSI gateway not selected.")
		}
		if !bytes.Equal(mockData.request.Context, expectedRxInfo.Context) {
			t.Error("Context is mismatched.")
		}
		if !mockData.request.PolarizationInversion {
			t.Error("ipol is wrong, it should be true always")
		}
		if mockData.request.Delay.Seconds != 5 {
			t.Errorf("Delay should be 5, but got %d", mockData.request.Delay.Seconds)
		}

		macpayload := mockData.request.MacPayload

		if len(macpayload) != 77 {
			t.Errorf("Hello response should be 77 bytes, but got %d", len(macpayload))
		}
		if macpayload[0] != DownRespHello {
			t.Errorf("Wrong type value in Hello response, expected %X but got %X", DownRespHello, macpayload[0])
		}

		_, downdeveui, serverpukey, servernonce := extractHelloResponse(macpayload)
		downserverpubkey[i] = serverpukey
		downservernonce[i] = servernonce
		t.Logf("  rDevEui: %s", hex.EncodeToString(downdeveui))
		t.Logf("  serverPublicKey: %s", hex.EncodeToString(downserverpubkey[i]))
		t.Logf("  serverNonce: %s", hex.EncodeToString(downservernonce[i]))

		if !bytes.Equal(rDevEui, downdeveui) {
			t.Error("Mismatch rDevEui at Hello response.")
		}
		if isByteArrayAllZero(downserverpubkey[i]) {
			t.Error("Server key is not set at Hello response.")
		}
		if isByteArrayAllZero(downservernonce[i]) {
			t.Error("Server nonce is not set at Hello response.")
		}
	}
	// As rDevEui no changed, same device session should use => 2 downlink should have same server public key
	if !bytes.Equal(downserverpubkey[0], downserverpubkey[1]) {
		t.Error("Server key changed.")
	}
	if !bytes.Equal(downservernonce[0], downservernonce[1]) {
		t.Error("Server nonce changed.")
	}

}

func TestDeviceSession(t *testing.T) {
	ctx := context.Background()
	if setupUnitTest(&handlerMock{}) != nil {
		t.Error("Uint test setup failed.")
	}

	// Prepare keys
	devicepublickey := make([]byte, ecdh.K233PubKeySize)
	fillByteArray(devicepublickey, 0x01, len(devicepublickey))

	//
	request := as.HandleProprietaryUplinkRequest{
		MacPayload: prepareHelloMessage(rDevEui, devicepublickey),
		Mic:        []byte{0x00, 0x00, 0x00, 0x00},
		TxInfo:     &gwV3.UplinkTXInfo{},
		RxInfo:     mockRxInfo,
	}

	// Send 1st frame
	processed, err := HandleReceivedFrame(ctx, &request)
	if err != nil {
		t.Errorf("HandleReceivedFrame failed. %s", err)
	}
	if !processed {
		t.Error("Request not processed.")
	}
	if len(deviceSessionList) != 1 {
		t.Errorf("Number of active device session is wrong, expected 1, got %d", len(deviceSessionList))
	}

	// Send 2nd frame, same rDevEui
	processed, err = HandleReceivedFrame(ctx, &request)
	if err != nil {
		t.Errorf("HandleReceivedFrame failed. %s", err)
	}
	if !processed {
		t.Error("Request not processed.")
	}
	if len(deviceSessionList) != 1 {
		t.Errorf("Expected only 1 device session, got %d", len(deviceSessionList))
	}

	// Send 3nd frame, different rDevEui
	otherdeveui := []byte{0x81, 0x82, 0x83, 0xff, 0xfe, 0x84, 0x85, 0x80}
	request.MacPayload = prepareHelloMessage(otherdeveui, devicepublickey)
	processed, err = HandleReceivedFrame(ctx, &request)
	if err != nil {
		t.Errorf("HandleReceivedFrame failed. %s", err)
	}
	if !processed {
		t.Error("Request not processed.")
	}
	if len(deviceSessionList) != 2 {
		t.Errorf("Expected there is 2 device session, got %d", len(deviceSessionList))
	}
}

func TestDeviceSessionExpire(t *testing.T) {
	ctx := context.Background()
	if setupUnitTest(&handlerMock{}) != nil {
		t.Error("Uint test setup failed.")
	}

	// Prepare keys
	devicepublickey := make([]byte, ecdh.K233PubKeySize)
	fillByteArray(devicepublickey, 0x01, len(devicepublickey))

	//
	request := as.HandleProprietaryUplinkRequest{
		MacPayload: prepareHelloMessage(rDevEui, devicepublickey),
		Mic:        []byte{0x00, 0x00, 0x00, 0x00},
		TxInfo:     &gwV3.UplinkTXInfo{},
		RxInfo:     mockRxInfo,
	}

	//
	timestart := time.Now()
	mockNowQueue = []time.Time{timestart}

	// queue to max
	for i := 0; i < maxNumberOfDevSession; i++ {
		seqdeveui := []byte{0x00, 0x00, 0x83, 0xff, 0xfe, 0x84, 0x85, 0x86}
		seqdeveui[0] = uint8(i >> 8)
		seqdeveui[1] = uint8(i)
		request.MacPayload = prepareHelloMessage(seqdeveui, devicepublickey)
		_, err := HandleReceivedFrame(ctx, &request)
		if err != nil {
			t.Errorf("HandleReceivedFrame failed. %s", err)
		}

		timestamp := mockNowQueue[0]
		mockNowQueue = []time.Time{timestamp.Add(time.Second)}
	}
	if len(deviceSessionList) != maxNumberOfDevSession {
		t.Errorf("Expected number device session is %d, got %d", maxNumberOfDevSession, len(deviceSessionList))
	}

	// Queue one more
	request.MacPayload = prepareHelloMessage(rDevEui, devicepublickey)
	_, err := HandleReceivedFrame(ctx, &request)
	if err != nil {
		t.Errorf("HandleReceivedFrame failed. %s", err)
	}
	if len(deviceSessionList) != maxNumberOfDevSession {
		t.Errorf("Expected number device session is %d, got %d", maxNumberOfDevSession, len(deviceSessionList))
	}

	// Just before expire
	mockNowQueue = []time.Time{timestart.Add(deviceSessionLifeTime)}
	clearExpiredDevSession()
	if len(deviceSessionList) != 10 {
		t.Errorf("Expected number device session is %d, got %d", 10, len(deviceSessionList))
	}

	// 1st expired
	mockNowQueue[0] = mockNowQueue[0].Add(time.Second)
	clearExpiredDevSession()
	if len(deviceSessionList) != 9 {
		t.Errorf("Expected number device session is %d, got %d", 9, len(deviceSessionList))
	}

	// 4 more
	mockNowQueue[0] = mockNowQueue[0].Add(time.Second * 4)
	clearExpiredDevSession()
	if len(deviceSessionList) != 5 {
		t.Errorf("Expected number device session is %d, got %d", 5, len(deviceSessionList))
	}

	// All expired
	mockNowQueue[0] = timestart.Add(deviceSessionLifeTime * 2)
	clearExpiredDevSession()
	if len(deviceSessionList) != 0 {
		t.Errorf("Expected number device session is %d, got %d", 0, len(deviceSessionList))
	}
}

func TestHandleReceivedFrameValidAuth(t *testing.T) {
	ctx := context.Background()

	if setupUnitTest(&handlerMock{}) != nil {
		t.Error("Uint test setup failed.")
	}

	serialNumberHash := make([]byte, 32)
	verifyCode := make([]byte, 16)
	devNonce := make([]byte, 4)

	request := as.HandleProprietaryUplinkRequest{
		MacPayload: prepareAuthMessage(rDevEui, serialNumberHash, verifyCode, devNonce),
		Mic:        []byte{0x00, 0x00, 0x00, 0x00},
		TxInfo:     &gwV3.UplinkTXInfo{},
		RxInfo:     mockRxInfo,
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
	ctx := context.Background()

	if setupUnitTest(&handlerMock{}) != nil {
		t.Error("Uint test setup failed.")
	}

	request := as.HandleProprietaryUplinkRequest{
		MacPayload: []byte{0x00, 0x02, 0x03, 0x04},
		Mic:        []byte{0x00, 0x00, 0x00, 0x00},
		TxInfo:     &gwV3.UplinkTXInfo{},
		RxInfo:     mockRxInfo,
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
	ctx := context.Background()

	if setupUnitTest(&handlerMock{}) != nil {
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
