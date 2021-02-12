package devprovision

import (
	"bytes"
	"context"
	cryptorand "crypto/rand"
	"encoding/binary"
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
const mockTxFreq = 471100000

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

func prepareAuthMessage(session *deviceSession, provisionIDHash []byte, verifyCode []byte, devNonce []byte) []byte {
	authpayload := make([]byte, 52)
	copy(authpayload[0:], provisionIDHash[:])
	copy(authpayload[32:], verifyCode[:])
	copy(authpayload[48:], devNonce[:])

	encpayload := session.encryptAuthPayload(authpayload, true)
	payload := []byte{0x11}
	payload = append(payload, session.rDevEui...)
	payload = append(payload, encpayload...)
	return payload
}

func extractAuthAccepted(macpayload []byte) (bool, []byte, []byte) {
	retok := false
	rdeveui := make([]byte, 8)
	encpayload := make([]byte, 32)

	if len(macpayload) == 41 {
		copy(rdeveui[:], macpayload[1:])
		copy(encpayload[:], macpayload[9:])
		retok = true
	}
	return retok, rdeveui, encpayload
}

func changeGwContext(req *as.HandleProprietaryUplinkRequest) {
	for _, rx := range req.RxInfo {
		randbuf := make([]byte, 4)
		_, err := cryptorand.Read(randbuf[:])
		if err != nil {
			log.Error("crypto.rand() failed. Fallback to Pseudorandom")
		}
		copy(rx.Context[:], randbuf)
	}
}

//
func TestMain(m *testing.M) {
	// Set log level
	log.SetLevel(logrus.DebugLevel)

	// Setup mock funcs
	funcGetNow = mockGetNow
	funcGen128Rand = mockGen128Rand
	funcFindDeviceBySnHash = mockFindDeviceBySnHash
	funcSaveDevice = mockSaveDevice

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

func TestDeviceSessionHandling(t *testing.T) {
	ctx := context.Background()
	if setupUnitTest(&handlerMock{}) != nil {
		t.Error("Uint test setup failed.")
	}

	// Prepare keys
	devicepublickey := make([]byte, ecdh.K233PubKeySize)
	fillByteArray(devicepublickey, 0x01)

	//
	request := as.HandleProprietaryUplinkRequest{
		MacPayload: prepareHelloMessage(rDevEui, devicepublickey),
		Mic:        []byte{0x00, 0x00, 0x00, 0x00},
		TxInfo:     &gwV3.UplinkTXInfo{},
		RxInfo:     mockRxInfo,
	}
	request.Mic = calProprietaryMic(request.MacPayload)

	// Send 1st frame
	processed, err := HandleReceivedFrame(ctx, &request)
	if err != nil {
		t.Fatalf("HandleReceivedFrame failed. %s", err)
	}
	if !processed {
		t.Fatalf("Request not processed.")
	}
	if len(deviceSessionList) != 1 {
		t.Fatalf("Number of active device session is wrong, expected 1, got %d", len(deviceSessionList))
	}

	// Send 2nd frame, same rDevEui
	processed, err = HandleReceivedFrame(ctx, &request)
	if err != nil {
		t.Fatalf("HandleReceivedFrame failed. %s", err)
	}
	if !processed {
		t.Fatalf("Request not processed.")
	}
	if len(deviceSessionList) != 1 {
		t.Fatalf("Expected only 1 device session, got %d", len(deviceSessionList))
	}

	// Send 3nd frame, different rDevEui
	otherdeveui := []byte{0x81, 0x82, 0x83, 0xff, 0xfe, 0x84, 0x85, 0x80}
	request.MacPayload = prepareHelloMessage(otherdeveui, devicepublickey)
	request.Mic = calProprietaryMic(request.MacPayload)

	processed, err = HandleReceivedFrame(ctx, &request)
	if err != nil {
		t.Fatalf("HandleReceivedFrame failed. %s", err)
	}
	if !processed {
		t.Fatalf("Request not processed.")
	}
	if len(deviceSessionList) != 2 {
		t.Fatalf("Expected there is 2 device session, got %d", len(deviceSessionList))
	}
}

func TestDeviceSessionExpire(t *testing.T) {
	ctx := context.Background()
	if setupUnitTest(&handlerMock{}) != nil {
		t.Error("Uint test setup failed.")
	}

	// Prepare keys
	devicepublickey := make([]byte, ecdh.K233PubKeySize)
	fillByteArray(devicepublickey, 0x01)

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
		request.Mic = calProprietaryMic(request.MacPayload)
		_, err := HandleReceivedFrame(ctx, &request)
		if err != nil {
			t.Fatalf("HandleReceivedFrame failed. %s", err)
		}

		timestamp := mockNowQueue[0]
		mockNowQueue = []time.Time{timestamp.Add(time.Second)}
	}
	if len(deviceSessionList) != maxNumberOfDevSession {
		t.Fatalf("Expected number device session is %d, got %d", maxNumberOfDevSession, len(deviceSessionList))
	}

	// Queue one more
	request.MacPayload = prepareHelloMessage(rDevEui, devicepublickey)
	request.Mic = calProprietaryMic(request.MacPayload)
	_, err := HandleReceivedFrame(ctx, &request)
	if err != nil {
		t.Fatalf("HandleReceivedFrame failed. %s", err)
	}
	if len(deviceSessionList) != maxNumberOfDevSession {
		t.Fatalf("Expected number device session is %d, got %d", maxNumberOfDevSession, len(deviceSessionList))
	}

	// Just before expire
	mockNowQueue = []time.Time{timestart.Add(deviceSessionLifeTime)}
	clearExpiredDevSession()
	if len(deviceSessionList) != 10 {
		t.Fatalf("Expected number device session is %d, got %d", 10, len(deviceSessionList))
	}

	// 1st expired
	mockNowQueue[0] = mockNowQueue[0].Add(time.Second)
	clearExpiredDevSession()
	if len(deviceSessionList) != 9 {
		t.Fatalf("Expected number device session is %d, got %d", 9, len(deviceSessionList))
	}

	// 4 more
	mockNowQueue[0] = mockNowQueue[0].Add(time.Second * 4)
	clearExpiredDevSession()
	if len(deviceSessionList) != 5 {
		t.Fatalf("Expected number device session is %d, got %d", 5, len(deviceSessionList))
	}

	// All expired
	mockNowQueue[0] = timestart.Add(deviceSessionLifeTime * 2)
	clearExpiredDevSession()
	if len(deviceSessionList) != 0 {
		t.Fatalf("Expected number device session is %d, got %d", 0, len(deviceSessionList))
	}
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
	fillByteArray(privkeyrand, 0x01)

	_, devicepublickey := ecdhK233.GenerateKeys(privkeyrand)
	if devicepublickey == nil {
		t.Error("ecdhK233.GenerateKeys() failed.")
	}

	//
	request := as.HandleProprietaryUplinkRequest{
		MacPayload: prepareHelloMessage(rDevEui, devicepublickey),
		Mic:        []byte{0x00, 0x00, 0x00, 0x00},
		TxInfo:     &gwV3.UplinkTXInfo{Frequency: mockTxFreq},
		RxInfo:     mockRxInfo,
	}
	request.Mic = calProprietaryMic(request.MacPayload)

	mockRandValue = -1
	downserverpubkey := make([][]uint8, 2)
	downservernonce := make([][]uint8, 2)
	for i := 0; i < 2; i++ {
		// Handle uplink request
		mockData.request = nil
		changeGwContext(&request)
		processed, err := HandleReceivedFrame(ctx, &request)
		if err != nil {
			t.Errorf("HandleReceivedFrame failed. %s", err)
		}
		if !processed {
			t.Error("Request not processed.")
		}
		if mockData.request == nil {
			t.Fatal("Mock frame not found.")
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
		if mockData.request.UplinkFreq != mockTxFreq {
			t.Errorf("UplinkFreq not set.")
		}
		if mockData.request.DownlinkFreq != 0 {
			t.Errorf("DownlinkFreq should not set.")
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

func TestHandleReceivedFrameWrongHello(t *testing.T) {
	ctx := context.Background()

	if setupUnitTest(&handlerMock{}) != nil {
		t.Error("Uint test setup failed.")
	}

	// Wrong MIC
	devicepublickey := make([]byte, 64)
	request := as.HandleProprietaryUplinkRequest{
		MacPayload: prepareHelloMessage(rDevEui, devicepublickey),
		Mic:        []byte{0x00, 0x00, 0x00, 0x00},
		TxInfo:     &gwV3.UplinkTXInfo{},
		RxInfo:     mockRxInfo,
	}
	changeGwContext(&request)
	processed, err := HandleReceivedFrame(ctx, &request)
	if err != nil {
		t.Errorf("HandleReceivedFrame failed. %s", err)
	}
	if processed {
		t.Error("Wrong MIC HELLO message should not processed, but it did.")
	}

	// Wrong size
	request.MacPayload = append(request.MacPayload, []byte{0}...)
	request.Mic = calProprietaryMic(request.MacPayload)

	changeGwContext(&request)
	processed, err = HandleReceivedFrame(ctx, &request)
	if err != nil {
		t.Errorf("HandleReceivedFrame failed. %s", err)
	}
	if processed {
		t.Error("Wrong size HELLO should not processed, but it did.")
	}

}

func TestHandleReceivedFrameValidAuth(t *testing.T) {
	provkey := []byte{0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03}
	servernonce := []byte{0x01, 0x02, 0x03, 0x04}
	devicenonce := []byte{0x55, 0xaa, 0x55, 0xaa}
	privisionid := "SERIALNUMBEROOOOOOOO"
	privisionidhash := []byte{
		0x34, 0xdf, 0xcb, 0x3d, 0xde, 0x1a, 0x09, 0xfd,
		0x34, 0x0f, 0xaf, 0xad, 0xa1, 0xe4, 0x31, 0xe8,
		0x40, 0x28, 0xfc, 0x53, 0xc3, 0x28, 0xd3, 0x59,
		0xa8, 0x82, 0x46, 0x13, 0xb8, 0x6d, 0x56, 0x8e}
	sharedkey := []byte{0x57, 0x57, 0x3A, 0x81, 0xE2, 0x7E, 0x48, 0x26, 0xFA, 0x8E, 0x18, 0x70, 0xCD,
		0x6B, 0x66, 0x40, 0xF3, 0x90, 0x5D, 0x98, 0x40, 0xF4, 0x12, 0xFA, 0xAE, 0x74,
		0x0B, 0x12, 0xE0, 0x01, 0x00, 0x00, 0xC4, 0xD8, 0x27, 0xA9, 0x37, 0x49, 0xEE,
		0x44, 0xEA, 0x1B, 0xAC, 0x1C, 0x18, 0x8C, 0x03, 0xAA, 0x6B, 0x02, 0xDA, 0x1C,
		0x68, 0xE9, 0xE8, 0xE6, 0xCA, 0xB9, 0xD1, 0xED, 0x91, 0x01, 0x00, 0x00}
	expectedpayload := []byte{
		0xf0, 0x54, 0xef, 0xcb, 0xad, 0x62, 0xbf, 0x87, 0xd9, 0x56, 0x53, 0x92, 0xf9, 0x12, 0x5a, 0xfb,
		0x2c, 0x7c, 0xa6, 0xa3, 0xcf, 0xe9, 0x31, 0x94, 0x6c, 0xa9, 0x28, 0x00, 0x09, 0xc5, 0xa3, 0xdc}

	ctx := context.Background()

	if setupUnitTest(&handlerMock{}) != nil {
		t.Error("Uint test setup failed.")
	}

	// Create a device sssion
	mockRandValue = 1
	session := makeDeviceSession()
	copy(session.rDevEui[:], rDevEui[:])
	copy(session.provKey[:], provkey[:])
	copy(session.sharedKey[:], sharedkey[:])
	copy(session.serverNonce[:], servernonce[:])
	copy(session.devNonce[:], devicenonce[:])
	session.deriveKeys()
	sessionid := binary.BigEndian.Uint64(rDevEui)
	deviceSessionList[sessionid] = session

	verifycode := session.calVerifyCode(privisionid, true)

	request := as.HandleProprietaryUplinkRequest{
		MacPayload: prepareAuthMessage(&session, privisionidhash, verifycode, devicenonce),
		Mic:        []byte{0x00, 0x00, 0x00, 0x00},
		TxInfo:     &gwV3.UplinkTXInfo{Frequency: mockTxFreq},
		RxInfo:     mockRxInfo,
	}
	request.Mic = calProprietaryMic(request.MacPayload)

	//
	changeGwContext(&request)
	processed, err := HandleReceivedFrame(ctx, &request)
	if err != nil {
		t.Errorf("HandleReceivedFrame failed. %s", err)
	}
	if !processed {
		t.Error("Request not processed.")
	}
	if mockData.request == nil {
		t.Fatal("Mock frame not found.")
	}

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
	if mockData.request.UplinkFreq != mockTxFreq {
		t.Errorf("UplinkFreq not set.")
	}
	if mockData.request.DownlinkFreq != 0 {
		t.Errorf("DownlinkFreq should not set.")
	}

	macpayload := mockData.request.MacPayload
	if len(macpayload) != 41 {
		t.Errorf("Auth Accept response should be 41 bytes, but got %d", len(macpayload))
	}
	if macpayload[0] != DownRespAuthAccept {
		t.Errorf("Wrong type value in Auth Accept response, expected %X but got %X", DownRespHello, macpayload[0])
	}

	_, downdeveui, encpayload := extractAuthAccepted(macpayload)
	if !bytes.Equal(rDevEui, downdeveui) {
		t.Error("Mismatch rDevEui at Hello response.")
	}
	if !bytes.Equal(encpayload, expectedpayload) {
		t.Error("Wrong encrypted payload.")
	}
}

func TestHandleReceivedFrameWrongAuth(t *testing.T) {
	provkey := []byte{0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03}
	servernonce := []byte{0x01, 0x02, 0x03, 0x04}
	devicenonce := []byte{0x55, 0xaa, 0x55, 0xaa}
	privisionid := "SERIALNUMBEROOOOOOOO"
	privisionidhash := []byte{
		0x34, 0xdf, 0xcb, 0x3d, 0xde, 0x1a, 0x09, 0xfd,
		0x34, 0x0f, 0xaf, 0xad, 0xa1, 0xe4, 0x31, 0xe8,
		0x40, 0x28, 0xfc, 0x53, 0xc3, 0x28, 0xd3, 0x59,
		0xa8, 0x82, 0x46, 0x13, 0xb8, 0x6d, 0x56, 0x8e}
	sharedkey := []byte{0x57, 0x57, 0x3A, 0x81, 0xE2, 0x7E, 0x48, 0x26, 0xFA, 0x8E, 0x18, 0x70, 0xCD,
		0x6B, 0x66, 0x40, 0xF3, 0x90, 0x5D, 0x98, 0x40, 0xF4, 0x12, 0xFA, 0xAE, 0x74,
		0x0B, 0x12, 0xE0, 0x01, 0x00, 0x00, 0xC4, 0xD8, 0x27, 0xA9, 0x37, 0x49, 0xEE,
		0x44, 0xEA, 0x1B, 0xAC, 0x1C, 0x18, 0x8C, 0x03, 0xAA, 0x6B, 0x02, 0xDA, 0x1C,
		0x68, 0xE9, 0xE8, 0xE6, 0xCA, 0xB9, 0xD1, 0xED, 0x91, 0x01, 0x00, 0x00}

	ctx := context.Background()

	if setupUnitTest(&handlerMock{}) != nil {
		t.Error("Uint test setup failed.")
	}

	// Create a device sssion
	session := makeDeviceSession()
	copy(session.rDevEui[:], rDevEui[:])
	copy(session.provKey[:], provkey[:])
	copy(session.sharedKey[:], sharedkey[:])
	copy(session.serverNonce[:], servernonce[:])
	copy(session.devNonce[:], devicenonce[:])
	session.deriveKeys()
	sessionid := binary.BigEndian.Uint64(rDevEui)
	deviceSessionList[sessionid] = session

	verifycode := session.calVerifyCode(privisionid, true)

	request := as.HandleProprietaryUplinkRequest{
		MacPayload: prepareAuthMessage(&session, privisionidhash, verifycode, devicenonce),
		Mic:        []byte{0x00, 0x00, 0x00, 0x00},
		TxInfo:     &gwV3.UplinkTXInfo{},
		RxInfo:     mockRxInfo,
	}

	// Wrong MIC
	request.Mic = []byte{0x01, 0x02, 0x03, 0x04}
	changeGwContext(&request)
	processed, err := HandleReceivedFrame(ctx, &request)
	if err != nil {
		t.Errorf("HandleReceivedFrame failed. %s", err)
	}
	if processed {
		t.Error("Wrong size AUTH should not processed, but it did.")
	}

	// Wrong Size
	request.MacPayload = append(request.MacPayload, []byte{0}...)
	request.Mic = calProprietaryMic(request.MacPayload)
	changeGwContext(&request)
	processed, err = HandleReceivedFrame(ctx, &request)
	if err != nil {
		t.Errorf("HandleReceivedFrame failed. %s", err)
	}
	if processed {
		t.Error("Wrong size AUTH should not processed, but it did.")
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
	request.Mic = calProprietaryMic(request.MacPayload)

	//
	changeGwContext(&request)
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
	request.Mic = calProprietaryMic(request.MacPayload)

	//
	changeGwContext(&request)
	processed, err := HandleReceivedFrame(ctx, &request)
	if err == nil {
		t.Error("Expected fail but it passed.")
	}
	if processed {
		t.Error("Expected the frame is not processed but it did.")
	}
}

func TestDeriveKeys(t *testing.T) {
	sharedkey := []byte{0x57, 0x57, 0x3A, 0x81, 0xE2, 0x7E, 0x48, 0x26, 0xFA, 0x8E, 0x18, 0x70, 0xCD,
		0x6B, 0x66, 0x40, 0xF3, 0x90, 0x5D, 0x98, 0x40, 0xF4, 0x12, 0xFA, 0xAE, 0x74,
		0x0B, 0x12, 0xE0, 0x01, 0x00, 0x00, 0xC4, 0xD8, 0x27, 0xA9, 0x37, 0x49, 0xEE,
		0x44, 0xEA, 0x1B, 0xAC, 0x1C, 0x18, 0x8C, 0x03, 0xAA, 0x6B, 0x02, 0xDA, 0x1C,
		0x68, 0xE9, 0xE8, 0xE6, 0xCA, 0xB9, 0xD1, 0xED, 0x91, 0x01, 0x00, 0x00}
	expectedappkey := []byte{0xFC, 0x3B, 0xDD, 0x59, 0x22, 0x87, 0xD9, 0x73, 0x48, 0xC0, 0x0B, 0xAC, 0x46, 0xB3, 0x05, 0x79}
	expectednwkkey := []byte{0x5B, 0x87, 0x83, 0xAF, 0x06, 0xFF, 0xB3, 0x62, 0x9D, 0x03, 0x77, 0x9B, 0xF3, 0x4E, 0x12, 0x89}
	expectedprovkey := []byte{0x29, 0x53, 0x01, 0x98, 0x2D, 0x35, 0xC7, 0x2F, 0x71, 0x42, 0xB9, 0xDD, 0x07, 0xFE, 0x1D, 0xEF}

	session := makeDeviceSession()

	copy(session.rDevEui[:], rDevEui[:])
	copy(session.sharedKey[:], sharedkey[:])

	session.deriveKeys()

	if !bytes.Equal(session.appKey, expectedappkey) {
		t.Error("Wrong appKey")
	}
	if !bytes.Equal(session.nwkKey, expectednwkkey) {
		t.Error("Wrong nwkKey")
	}
	if !bytes.Equal(session.provKey, expectedprovkey) {
		t.Error("Wrong provKey")
	}

}

func TestEncryptAuthPayload(t *testing.T) {
	payload := make([]byte, 52)
	provkey := []byte{0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03}
	expecteduplink := []byte{0x15, 0xD3, 0x9A, 0x14, 0xA2, 0x7C, 0x09, 0x9F, 0x4D, 0x14, 0x16, 0x3C, 0xFC, 0xA3, 0x37, 0x65,
		0x05, 0xC3, 0x8A, 0x04, 0xB2, 0x6C, 0x19, 0x8F, 0x5D, 0x04, 0x06, 0x2C, 0xEC, 0xB3, 0x27, 0x75,
		0x35, 0xF3, 0xBA, 0x34, 0x82, 0x5C, 0x29, 0xBF, 0x6D, 0x34, 0x36, 0x1C, 0xDC, 0x83, 0x17, 0x45,
		0x25, 0xE3, 0xAA, 0x24}
	expecteddownlink := []byte{0x5D, 0x08, 0x7E, 0x03, 0x28, 0xD7, 0x57, 0x64, 0xCE, 0x63, 0x13, 0x99, 0x07, 0x0C, 0x96, 0x1D,
		0x4D, 0x18, 0x6E, 0x13, 0x38, 0xC7, 0x47, 0x74, 0xDE, 0x73, 0x03, 0x89, 0x17, 0x1C, 0x86, 0x0D,
		0x7D, 0x28, 0x5E, 0x23, 0x08, 0xF7, 0x77, 0x44, 0xEE, 0x43, 0x33, 0xB9, 0x27, 0x2C, 0xB6, 0x3D,
		0x6D, 0x38, 0x4E, 0x33}

	session := makeDeviceSession()
	copy(session.rDevEui[:], rDevEui[:])
	copy(session.provKey[:], provkey[:])

	// Uplink, Dir = 0
	for i := range payload {
		payload[i] = uint8(i)
	}
	encpayload := session.encryptAuthPayload(payload, true)
	if !bytes.Equal(encpayload, expecteduplink) {
		t.Error("Wrong uplink encryption output")
	}
	encpayload = session.encryptAuthPayload(encpayload, true)
	if !bytes.Equal(encpayload, payload) {
		t.Error("Wrong uplink decryption output")
	}

	// Downlink, Dir = 1
	for i := range payload {
		payload[i] = uint8(i)
	}
	encpayload = session.encryptAuthPayload(payload, false)
	if !bytes.Equal(encpayload, expecteddownlink) {
		t.Error("Wrong downlink encryption output")
	}
	encpayload = session.encryptAuthPayload(encpayload, false)
	if !bytes.Equal(encpayload, payload) {
		t.Error("Wrong downlink decryption output")
	}

}

func TestCalVerifyCode(t *testing.T) {
	privisionid := "SERIALNUMBEROOOOOOOO"
	servernonce := []byte{0x01, 0x02, 0x03, 0x04}
	devicenonce := []byte{0x55, 0xaa, 0x55, 0xaa}
	expectedserveroutput := []byte{0x2E, 0x69, 0xBB, 0x5E, 0xD7, 0x8B, 0x5E, 0xE8, 0x0C, 0x6A, 0x8A, 0xDC, 0x81, 0x91, 0xDD, 0xF8}
	expecteddeviceoutput := []byte{0xF8, 0x4A, 0xE2, 0x97, 0x9C, 0x56, 0x49, 0x03, 0xB4, 0xFE, 0x7A, 0x93, 0xF1, 0xD6, 0xF8, 0x26}

	session := makeDeviceSession()
	copy(session.serverNonce[:], servernonce[:])
	copy(session.devNonce[:], devicenonce[:])

	// Use serverNonce
	output := session.calVerifyCode(privisionid, true)
	if !bytes.Equal(output, expectedserveroutput) {
		t.Error("Wrong verify code for server")
	}

	// Use deviceNonce
	output = session.calVerifyCode(privisionid, false)
	if !bytes.Equal(output, expecteddeviceoutput) {
		t.Error("Wrong verify code for device")
	}
}

func TestCalProprietaryMic(t *testing.T) {
	macpayload := []byte{
		0x01, 0xD9, 0x7A, 0xC7, 0x25, 0xF1, 0xDC, 0x99, 0x31, 0x09, 0x8E, 0x34, 0x9B, 0x2B, 0x06, 0xE8,
		0x6A, 0x5D, 0x41, 0x21, 0xDA, 0xA5, 0xC2, 0x18, 0xEC, 0x77, 0xCA, 0x21, 0x9A, 0x4C, 0x94, 0x7F,
		0x33, 0x60, 0xA7, 0x00, 0xA7, 0x2B, 0x00, 0x00, 0x00, 0x9E, 0xD7, 0x18, 0x84, 0x12, 0xF0, 0xE4,
		0x33, 0x18, 0x18, 0x02, 0x26, 0x51, 0x89, 0xDD, 0x17, 0x98, 0x13, 0xF2, 0xAB, 0x4E, 0x92, 0x2A,
		0xF2, 0x69, 0x28, 0x2F, 0xB9, 0x2D, 0x00, 0x00, 0x00, 0x01}
	expectedmic := []byte{0xB7, 0xC2, 0xCB, 0xB9}

	calmic := calProprietaryMic(macpayload[:])
	if !bytes.Equal(calmic, expectedmic) {
		t.Error("Wrong MIC code")
	}
}
