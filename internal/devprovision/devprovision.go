package devprovision

import (
	"bytes"
	"context"
	cryptorand "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"sync"
	"time"

	//	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/brocaar/chirpstack-api/go/v3/as"
	gwV3 "github.com/brocaar/chirpstack-api/go/v3/gw"
	"github.com/brocaar/lorawan"
	duration "github.com/golang/protobuf/ptypes/duration"

	nsextra "github.com/mxc-foundation/lpwan-app-server/api/ns-extra"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/networkserverextra"
	"github.com/mxc-foundation/lpwan-app-server/internal/devprovision/ecdh"
	gwd "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway/data"
	nsd "github.com/mxc-foundation/lpwan-app-server/internal/networkserver_portal/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"
)

// LoRa Frame Message and Response Type
//#define MAX_MESSAGE_SIZE 256
const (
	UpMessageHello     = 0x01
	UpMessageAuth      = 0x11
	DownRespHello      = 0x81
	DownRespAuthAccept = 0x91
	DownRespAuthReject = 0x92
)
const (
	sizeUpMessageHello = 74
	sizeUpMessageAuth  = 61
)

//
const moduleName = "devprovision"

// Proprietary Payload
type proprietaryPayload struct {
	MacPayload []byte
	GatewayMAC lorawan.EUI64
	Frequency  uint32
	DR         int
	Context    []byte
	Delay      *duration.Duration
	Mic        []byte
}

//
var ecdhK223 ecdh.K233
var deviceSessionList map[uint64]deviceSession
var mutexDeviceSessionList sync.RWMutex
var maxNumberOfDevSession = 5000
var deviceSessionLifeTime = time.Minute * 5

// Func to get current now, it will override at test
var funcGetNow = time.Now
var funcGen128Rand = gen128Rand

// Gen 128 bytes of random numbers
func gen128Rand() []byte {
	randbuf := make([]byte, 128)
	_, err := cryptorand.Read(randbuf[:])
	if err != nil {
		log.Error("crypto.rand() failed. Fallback to Pseudorandom")
		// Fallback to Pseudorandom
		softrand := softRand{}
		for i := range randbuf {
			randbuf[i] = uint8(softrand.Get())
		}
	}
	return randbuf
}

//
func init() {
	mgr.RegisterModuleSetup(moduleName, Setup)
}

// Function pointer for send payload to ns
type sendToNsFunc func(n nsd.NetworkServer, req *nsextra.SendDelayedProprietaryPayloadRequest) error

var ctrl struct {
	handler      *store.Handler
	handlerMock  *handlerMock
	sendToNsFunc sendToNsFunc

	moduleUp bool
}

// Setup prepares device provisioning service module
func Setup(name string, h *store.Handler) error {
	if ctrl.moduleUp {
		return nil
	}
	defer func() {
		ctrl.moduleUp = true
	}()

	ctrl.handler = h
	ctrl.handlerMock = nil
	ctrl.sendToNsFunc = sendToNs
	clearDeviceSessionList()

	go CleanUpLoop()

	return nil
}

// Setup for unit test
func setupUnitTest(h *handlerMock) error {
	ctrl.handler = nil
	ctrl.handlerMock = h
	ctrl.sendToNsFunc = sendToNsMock
	clearDeviceSessionList()
	return nil
}

// CleanUpLoop is a never returning function, performing cleanup
func CleanUpLoop() {
	for {
		clearExpiredDevSession()
		time.Sleep(time.Second * 10)
	}
}

// HandleReceivedFrame handles a ping received by one or multiple gateways.
func HandleReceivedFrame(ctx context.Context, req *as.HandleProprietaryUplinkRequest) (bool, error) {
	var mic lorawan.MIC
	copy(mic[:], req.Mic)

	log.Debugf("Rx MacPayload:\n%s", hex.Dump(req.MacPayload))
	log.Debugf("          MIC: %s", hex.EncodeToString(req.Mic))

	// Find max RSSI gw
	var maxRssiRx *gwV3.UplinkRXInfo = nil
	for _, rx := range req.RxInfo {
		if maxRssiRx == nil {
			maxRssiRx = rx
		} else if rx.Rssi > maxRssiRx.Rssi {
			maxRssiRx = rx
		}
	}
	if maxRssiRx == nil {
		return false, errors.Errorf("No gateway found.")
	}

	log.Debugf("  MAC:%s, RSSI: %d, Context: %s", hex.EncodeToString(maxRssiRx.GatewayId), maxRssiRx.Rssi,
		hex.EncodeToString(maxRssiRx.Context))

	// Get Gateway
	var mac lorawan.EUI64
	copy(mac[:], maxRssiRx.GatewayId)

	var gw gwd.Gateway
	var n nsd.NetworkServer
	var err error
	if ctrl.handlerMock != nil {
		gw, err = ctrl.handlerMock.GetGateway(ctx, mac, false)
	} else {
		gw, err = ctrl.handler.GetGateway(ctx, mac, false)
	}
	if err != nil {
		return false, errors.Wrap(err, "get gateway error")
	}

	if ctrl.handlerMock != nil {
		n, err = ctrl.handlerMock.GetNetworkServer(ctx, gw.NetworkServerID)
	} else {
		n, err = ctrl.handler.GetNetworkServer(ctx, gw.NetworkServerID)
	}
	if err != nil {
		return false, errors.Wrap(err, "get network-server error")
	}
	log.Debugf("  NetworkServer: %s", n.Server)

	//

	// Check Message Type
	processed := false
	messageType := req.MacPayload[0]
	messageSize := len(req.MacPayload)

	if (messageType == UpMessageHello) && (messageSize == sizeUpMessageHello) {
		err := handleHello(n, req, maxRssiRx)
		if err != nil {
			return false, errors.Wrap(err, "send proprietary error")
		}
		processed = true
	} else if (messageType == UpMessageAuth) && (messageSize == sizeUpMessageAuth) {
		err := handleAuth(n, req, maxRssiRx)
		if err != nil {
			return false, errors.Wrap(err, "send proprietary error")
		}
		processed = true
	}

	return processed, nil
}

func sendToNs(n nsd.NetworkServer, req *nsextra.SendDelayedProprietaryPayloadRequest) error {
	nsClient, err := networkserverextra.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err != nil {
		return err
	}
	_, err = nsClient.SendDelayedProprietaryPayload(context.Background(), req)
	if err != nil {
		return err
	}

	return nil
}

func sendProprietary(n nsd.NetworkServer, payload proprietaryPayload) error {
	req := nsextra.SendDelayedProprietaryPayloadRequest{
		MacPayload:            payload.MacPayload,
		GatewayMacs:           [][]byte{payload.GatewayMAC[:]},
		PolarizationInversion: true,
		Frequency:             uint32(payload.Frequency),
		Dr:                    uint32(payload.DR),
		Context:               payload.Context,
		Delay:                 payload.Delay,
		Mic:                   payload.Mic,
	}

	if ctrl.sendToNsFunc != nil {
		err := ctrl.sendToNsFunc(n, &req)
		if err != nil {
			return errors.Wrap(err, "send proprietary payload error")
		}
		log.WithFields(log.Fields{
			"gateway_mac": payload.GatewayMAC,
			"freq":        payload.Frequency,
		}).Infof("gateway proprietary payload sent to network-server %s", n.Server)
	} else {
		return errors.Errorf("ctrl.sendToNsFunc() not set.")
	}

	return nil
}

func makeHelloResponse(session deviceSession) []byte {
	payload := []byte{DownRespHello}
	payload = append(payload, session.rDevEui...)
	payload = append(payload, session.serverPublicKey...)
	payload = append(payload, session.serverNonce...)
	return payload
}

//
func handleHello(nserver nsd.NetworkServer, req *as.HandleProprietaryUplinkRequest, targetgateway *gwV3.UplinkRXInfo) error {
	log.Debug("  HELLO Message.")

	var upFreqChannel uint32 = (req.TxInfo.Frequency - 470300000) / 200000
	var downFreq uint32 = 500300000 + ((upFreqChannel % 48) * 200000)
	var err error
	var frameversion byte

	//
	rdeveui := make([]byte, 8)
	copy(rdeveui[0:], req.MacPayload[1:])
	sessionid := binary.BigEndian.Uint64(rdeveui)
	log.Debugf("  sessionid=%X", sessionid)
	frameversion = req.MacPayload[73]

	//
	ok, currentsession := searchDeviceSession(sessionid)
	if !ok {
		rdeveui := make([]byte, 8)
		devicepublickey := make([]byte, ecdh.K233PubKeySize)

		log.Debugf("  Creating new session")
		copy(rdeveui[0:], req.MacPayload[1:])
		copy(devicepublickey[0:], req.MacPayload[9:])
		ok, currentsession = createDeviceSession(sessionid, rdeveui, devicepublickey)
		if !ok {
			// Create session failed. drop this frame. return true to mark is processed.
			return nil
		}
	}

	log.Debugf("  rDevEui: %s", hex.EncodeToString(currentsession.rDevEui))
	log.Debugf("  devicePublicKey: %s", hex.EncodeToString(currentsession.devicePublicKey))
	log.Debugf("  serverPrivateKey: %s", hex.EncodeToString(currentsession.serverPrivateKey))
	log.Debugf("  serverPublicKey: %s", hex.EncodeToString(currentsession.serverPublicKey))
	log.Debugf("  serverNonce: %s", hex.EncodeToString(currentsession.serverNonce))
	log.Debugf("  sharedKey: %s", hex.EncodeToString(currentsession.sharedKey))
	log.Debugf("  version: %d", frameversion)

	//
	var mac lorawan.EUI64
	copy(mac[:], targetgateway.GatewayId)

	payload := proprietaryPayload{
		MacPayload: makeHelloResponse(currentsession),
		GatewayMAC: mac,
		Frequency:  downFreq,
		DR:         3,
		Delay:      &duration.Duration{Seconds: 5, Nanos: 0},
		Context:    targetgateway.Context,
		Mic:        []byte{0x01, 0x02, 0x03, 0x04},
	}
	log.Debugf("Tx MacPayload:\n%s", hex.Dump(payload.MacPayload))
	log.Debugf("          MIC: %s", hex.EncodeToString(payload.Mic))

	err = sendProprietary(nserver, payload)
	if err != nil {
		return err
	}

	return nil
}

func makeAuthAccept(session deviceSession) []byte {
	authpayload := make([]byte, 32)

	encpayload := session.encryptAuthPayload(authpayload, true)

	payload := []byte{DownRespAuthAccept}
	payload = append(payload, session.rDevEui...)
	payload = append(payload, encpayload...)

	return payload
}

func handleAuth(nserver nsd.NetworkServer, req *as.HandleProprietaryUplinkRequest, targetgateway *gwV3.UplinkRXInfo) error {
	log.Debug("  AUTH Message.")

	var upFreqChannel uint32 = (req.TxInfo.Frequency - 470300000) / 200000
	var downFreq uint32 = 500300000 + ((upFreqChannel % 48) * 200000)

	//
	rdeveui := make([]byte, 8)
	copy(rdeveui[0:], req.MacPayload[1:])
	sessionid := binary.BigEndian.Uint64(rdeveui)
	log.Debugf("  sessionid=%X", sessionid)

	//
	ok, currentsession := searchDeviceSession(sessionid)
	if !ok {
		log.Debugf("  Auth message without active session. Frame dropped.")
		return nil
	}

	authpayload := make([]byte, 52)
	copy(authpayload[:], req.MacPayload[9:])
	authpayload = currentsession.encryptAuthPayload(authpayload, true)

	serialnumberhash := make([]byte, 32)
	verifycode := make([]byte, 16)
	copy(serialnumberhash[:], authpayload[0:])
	copy(verifycode[:], authpayload[32:])
	copy(currentsession.devNonce[:], authpayload[48:])

	log.Debugf("  rDevEui: %s", hex.EncodeToString(currentsession.rDevEui))
	log.Debugf("  devNonce: %s", hex.EncodeToString(currentsession.devNonce))
	log.Debugf("  serialNumberHash: %s", hex.EncodeToString(serialnumberhash))
	log.Debugf("  verifycode: %s", hex.EncodeToString(verifycode))

	found, deviceinfo := findDeviceBySnHash(serialnumberhash)
	if !found {
		return errors.Errorf("Device %s not found.", hex.EncodeToString(serialnumberhash))
	}
	log.Debugf("  Device found. %s, mfgID=%d, server=%s", deviceinfo.serialNumber, deviceinfo.manufacturerID, deviceinfo.server)
	log.Debugf("  devEUI=%s, appEUI=%s, appKey=%s, nwkKey=%s", deviceinfo.devEUI, deviceinfo.appEUI, deviceinfo.appKey, deviceinfo.nwkKey)
	log.Debugf("  status=%d, model=%s, fixedDevEUI=%v, created=%v", deviceinfo.status, deviceinfo.model, deviceinfo.fixedDevEUI, deviceinfo.timeCreated)

	calverifycode := currentsession.calVerifyCode(deviceinfo.serialNumber, true)
	if !bytes.Equal(verifycode, calverifycode) {
		return errors.Errorf("Incorrect verify code at Auth message")
	}

	//
	var mac lorawan.EUI64
	copy(mac[:], targetgateway.GatewayId)

	payload := proprietaryPayload{
		MacPayload: makeAuthAccept(currentsession),
		GatewayMAC: mac,
		Frequency:  downFreq,
		DR:         3,
		Delay:      &duration.Duration{Seconds: 5, Nanos: 0},
		Context:    targetgateway.Context,
		Mic:        []byte{0x01, 0x02, 0x03, 0x04},
	}
	log.Debugf("Tx MacPayload:\n%s", hex.Dump(payload.MacPayload))
	log.Debugf("          MIC: %s", hex.EncodeToString(payload.Mic))

	err := sendProprietary(nserver, payload)
	if err != nil {
		return err
	}

	return nil
}

// Device session handling
func searchDeviceSession(sessionid uint64) (bool, deviceSession) {
	mutexDeviceSessionList.Lock()
	defer mutexDeviceSessionList.Unlock()
	currentsession, sessionfound := deviceSessionList[sessionid]
	if !sessionfound {
		return false, deviceSession{}
	}
	return true, currentsession
}

func createDeviceSession(sessionid uint64, rdeveui []byte, devicepublickey []byte) (bool, deviceSession) {
	mutexDeviceSessionList.Lock()
	defer mutexDeviceSessionList.Unlock()

	if len(deviceSessionList) >= maxNumberOfDevSession {
		log.Warnf("Maximum number (%d) of device provisioning session reached. Request dropped.", maxNumberOfDevSession)
		return false, deviceSession{}
	}

	// New session
	currentsession := makeDeviceSession()
	copy(currentsession.rDevEui[0:], rdeveui)
	copy(currentsession.devicePublicKey[0:], devicepublickey)

	currentsession.genServerKeys()
	currentsession.genSharedKey()
	deviceSessionList[sessionid] = currentsession

	return true, currentsession
}

func clearExpiredDevSession() {
	mutexDeviceSessionList.Lock()
	defer mutexDeviceSessionList.Unlock()

	now := funcGetNow()
	for key, session := range deviceSessionList {
		if now.After(session.expireTime) {
			delete(deviceSessionList, key)
		}
	}
}

func clearDeviceSessionList() {
	mutexDeviceSessionList.Lock()
	defer mutexDeviceSessionList.Unlock()
	deviceSessionList = make(map[uint64]deviceSession)
}

//
func fillByteArray(input []byte, value uint8) {
	for i := range input {
		input[i] = value
	}
}

func findDeviceBySnHash(serialnumberhash []byte) (bool, deviceInfo) {
	strhash := hex.EncodeToString(serialnumberhash)

	for i := range fakeDeviceList {
		if fakeDeviceList[i].serialNumberHash == strhash {
			return true, fakeDeviceList[i]
		}
	}
	return false, deviceInfo{}
}
