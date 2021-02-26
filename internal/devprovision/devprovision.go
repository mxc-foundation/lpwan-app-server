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

	"github.com/jacobsa/crypto/cmac"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/brocaar/chirpstack-api/go/v3/as"
	gwV3 "github.com/brocaar/chirpstack-api/go/v3/gw"
	"github.com/brocaar/lorawan"
	duration "github.com/golang/protobuf/ptypes/duration"

	nsextra "github.com/mxc-foundation/lpwan-app-server/api/ns-extra"
	psPb "github.com/mxc-foundation/lpwan-app-server/api/ps-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/networkserverextra"
	pscli "github.com/mxc-foundation/lpwan-app-server/internal/clients/psconn"
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
	MacPayload   []byte
	GatewayMAC   lorawan.EUI64
	DownlinkFreq uint32
	UplinkFreq   uint32
	DR           int
	Context      []byte
	Delay        *duration.Duration
	Mic          []byte
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
var funcFindDeviceBySnHash = findDeviceBySnHash
var funcSaveDevice = saveDevice

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

func processMessage(ctx context.Context, n nsd.NetworkServer, req *as.HandleProprietaryUplinkRequest,
	targetgateway *gwV3.UplinkRXInfo) (bool, error) {
	processed := false
	messageType := req.MacPayload[0]
	messageSize := len(req.MacPayload)

	if (messageType == UpMessageHello) && (messageSize == sizeUpMessageHello) {
		err := handleHello(ctx, n, req, targetgateway)
		if err != nil {
			return false, errors.Wrap(err, "process HELLO msg error")
		}
		processed = true
	} else if (messageType == UpMessageAuth) && (messageSize == sizeUpMessageAuth) {
		err := handleAuth(ctx, n, req, targetgateway)
		if err != nil {
			return false, errors.Wrap(err, "process AUTH msg error")
		}
		processed = true
	} else {
		log.Debug("Unknown Message.")
	}

	return processed, nil
}

// HandleReceivedFrame handles a ping received by one or multiple gateways.
func HandleReceivedFrame(ctx context.Context, req *as.HandleProprietaryUplinkRequest) (bool, error) {
	var mic lorawan.MIC
	copy(mic[:], req.Mic)

	// log.Debugf("Rx MacPayload:\n%s", hex.Dump(req.MacPayload))
	// log.Debugf("          MIC: %s", hex.EncodeToString(req.Mic))

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

	// Check MIC
	calmic := calProprietaryMic(req.MacPayload)
	if !bytes.Equal(calmic, req.Mic) {
		// log.Debugf("MacPayload:\n%s", hex.Dump(req.MacPayload))
		log.Debugf("Wrong MIC calmic=%s, rxed mic=%s", hex.EncodeToString(calmic), hex.EncodeToString(req.Mic))
		return false, errors.Wrap(err, "Wrong MIC for MacPayload")
	}
	return processMessage(ctx, n, req, maxRssiRx)
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

func calProprietaryMic(macpayload []byte) []byte {
	micbuf := make([]byte, 4)
	hash, err := cmac.New(getFixedKey())
	if err != nil {
		return micbuf
	}

	if _, err = hash.Write([]byte{0xe0}); err != nil {
		return micbuf
	}
	if _, err = hash.Write(macpayload); err != nil {
		return micbuf
	}

	hb := hash.Sum([]byte{})

	copy(micbuf[0:], hb[:])

	return micbuf
}

func sendProprietary(n nsd.NetworkServer, payload proprietaryPayload) error {
	req := nsextra.SendDelayedProprietaryPayloadRequest{
		MacPayload:            payload.MacPayload,
		GatewayMacs:           [][]byte{payload.GatewayMAC[:]},
		PolarizationInversion: true,
		UplinkFreq:            payload.UplinkFreq,
		DownlinkFreq:          payload.DownlinkFreq,
		Dr:                    uint32(payload.DR),
		Context:               payload.Context,
		Delay:                 payload.Delay,
		Mic:                   calProprietaryMic(payload.MacPayload),
	}
	log.Debugf("  sendProprietary() MIC: %s", hex.EncodeToString(req.Mic))

	if ctrl.sendToNsFunc != nil {
		err := ctrl.sendToNsFunc(n, &req)
		if err != nil {
			return errors.Wrap(err, "send proprietary payload error")
		}
		log.WithFields(log.Fields{
			"gateway_mac": payload.GatewayMAC,
			"up_freq":     payload.UplinkFreq,
			"down_freq":   payload.DownlinkFreq,
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
func handleHello(ctx context.Context, nserver nsd.NetworkServer, req *as.HandleProprietaryUplinkRequest,
	targetgateway *gwV3.UplinkRXInfo) error {
	log.Debug("  HELLO Message.")

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

	// Drop if already sent to the same Gateway context
	ok, currentsession = checkDeviceSession(sessionid, targetgateway.Context)
	if !ok {
		return nil
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
		MacPayload:   makeHelloResponse(currentsession),
		GatewayMAC:   mac,
		UplinkFreq:   req.TxInfo.Frequency,
		DownlinkFreq: 0,
		DR:           3,
		Delay:        &duration.Duration{Seconds: 5, Nanos: 0},
		Context:      targetgateway.Context,
		Mic:          []byte{0x00, 0x00, 0x00, 0x00},
	}
	// log.Debugf("Tx MacPayload:\n%s", hex.Dump(payload.MacPayload))

	err = sendProprietary(nserver, payload)
	if err != nil {
		return err
	}

	return nil
}

func makeAuthAccept(session deviceSession, verifycode []byte) []byte {
	authpayload := make([]byte, 32)
	copy(authpayload[0:], session.assignedDevEui[:])
	copy(authpayload[8:], session.assignedAppEui[:])
	copy(authpayload[16:], verifycode[:])
	encpayload := session.encryptAuthPayload(authpayload, false)

	payload := []byte{DownRespAuthAccept}
	payload = append(payload, session.rDevEui...)
	payload = append(payload, encpayload...)

	return payload
}

func handleAuth(ctx context.Context, nserver nsd.NetworkServer, req *as.HandleProprietaryUplinkRequest,
	targetgateway *gwV3.UplinkRXInfo) error {
	log.Debug("  AUTH Message.")

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

	// Drop if already sent to the same Gateway context
	ok, currentsession = checkDeviceSession(sessionid, targetgateway.Context)
	if !ok {
		return nil
	}

	//
	authpayload := make([]byte, 52)
	copy(authpayload[:], req.MacPayload[9:])
	authpayload = currentsession.encryptAuthPayload(authpayload, true)

	privisionidhash := make([]byte, 32)
	verifycode := make([]byte, 16)
	copy(privisionidhash[:], authpayload[0:])
	copy(verifycode[:], authpayload[32:])
	copy(currentsession.devNonce[:], authpayload[48:])

	log.Debugf("  rDevEui: %s", hex.EncodeToString(currentsession.rDevEui))
	log.Debugf("  devNonce: %s", hex.EncodeToString(currentsession.devNonce))
	log.Debugf("  privisionidhash: %s", hex.EncodeToString(privisionidhash))
	log.Debugf("  verifycode: %s", hex.EncodeToString(verifycode))

	found, deviceinfo := funcFindDeviceBySnHash(ctx, privisionidhash)
	if !found {
		return errors.Errorf("Device %s not found.", hex.EncodeToString(privisionidhash))
	} else if deviceinfo.Status == "DISABLED" {
		return errors.Errorf("Device %s disabled.", deviceinfo.ProvisionID)
	}
	log.Debugf("  Device found. %s, mfgID=%d, server=%s", deviceinfo.ProvisionID, deviceinfo.ManufacturerID, deviceinfo.Server)
	log.Debugf("  devEUI=%s, appEUI=%s, appKey=%s, nwkKey=%s",
		hex.EncodeToString(deviceinfo.DevEUI), hex.EncodeToString(deviceinfo.AppEUI),
		hex.EncodeToString(deviceinfo.AppKey), hex.EncodeToString(deviceinfo.NwkKey))
	log.Debugf("  status=%v, model=%v, fixedDevEUI=%v, created=%v", deviceinfo.Status, deviceinfo.Model, deviceinfo.FixedDevEUI,
		deviceinfo.TimeCreated.Time)

	calverifycode := currentsession.calVerifyCode(deviceinfo.ProvisionID, true)
	if !bytes.Equal(verifycode, calverifycode) {
		return errors.Errorf("Incorrect verify code at Auth message")
	}

	currentsession, deviceinfo, err := updateDevice(ctx, currentsession, deviceinfo)
	if err != nil {
		return errors.Wrap(err, "updateDevice error")
	}
	updateDeviceSession(sessionid, currentsession)

	err = funcSaveDevice(ctx, deviceinfo)
	if err != nil {
		return errors.Wrap(err, "saveDevice error")
	}

	//
	var mac lorawan.EUI64
	copy(mac[:], targetgateway.GatewayId)
	verifycode = currentsession.calVerifyCode(deviceinfo.ProvisionID, false)

	payload := proprietaryPayload{
		MacPayload:   makeAuthAccept(currentsession, verifycode),
		GatewayMAC:   mac,
		UplinkFreq:   req.TxInfo.Frequency,
		DownlinkFreq: 0,
		DR:           3,
		Delay:        &duration.Duration{Seconds: 5, Nanos: 0},
		Context:      targetgateway.Context,
		Mic:          []byte{0x00, 0x00, 0x00, 0x00},
	}
	// log.Debugf("Tx MacPayload:\n%s", hex.Dump(payload.MacPayload))

	err = sendProprietary(nserver, payload)
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

func updateDeviceSession(sessionid uint64, newsession deviceSession) {
	mutexDeviceSessionList.Lock()
	defer mutexDeviceSessionList.Unlock()
	_, sessionfound := deviceSessionList[sessionid]
	if sessionfound {
		deviceSessionList[sessionid] = newsession
	}

}

func checkDeviceSession(sessionid uint64, gwcontext []byte) (bool, deviceSession) {
	mutexDeviceSessionList.Lock()
	defer mutexDeviceSessionList.Unlock()

	currentsession, sessionfound := deviceSessionList[sessionid]
	if !sessionfound {
		return false, deviceSession{}
	}

	if bytes.Equal(currentsession.lastGwContext, gwcontext) {
		// Same gateway context already handled, retrun false to cause the frame being drop
		return false, deviceSession{}
	}

	// Save gateway context
	currentsession.lastGwContext = make([]byte, len(gwcontext))
	copy(currentsession.lastGwContext[:], gwcontext)
	deviceSessionList[sessionid] = currentsession

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
	currentsession.deriveKeys()
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

func updateDevice(ctx context.Context, session deviceSession, deviceinfo deviceInfo) (deviceSession, deviceInfo, error) {

	if isByteArrayAllZero(session.assignedDevEui) || !bytes.Equal(session.assignedDevEui, deviceinfo.DevEUI) {
		// Session is new
		deveui := make([]byte, 8)
		appeui := make([]byte, 8)
		copy(deveui[:], deviceinfo.DevEUI)
		copy(appeui[:], deviceinfo.AppEUI)

		if !deviceinfo.FixedDevEUI || isByteArrayAllZero(deveui) {
			// Generate devEUI
			psClient, err := pscli.GetDevProClient()
			if err != nil {
				return session, deviceinfo, err
			}

			resp, err := psClient.GenDevEUI(ctx, &psPb.GenDevEuiRequest{})
			if err != nil {
				return session, deviceinfo, err
			}
			copy(deveui[:], resp.DevEUI[:])
			copy(deviceinfo.DevEUI[:], resp.DevEUI[:])
		}

		copy(session.assignedDevEui[:], deveui[:])
		copy(session.assignedAppEui[:], appeui[:])

		copy(deviceinfo.AppKey[:], session.appKey)
		copy(deviceinfo.NwkKey[:], session.nwkKey)
		copy(deviceinfo.AppEUI[:], appeui)
	}

	return session, deviceinfo, nil
}

//
func fillByteArray(input []byte, value uint8) {
	for i := range input {
		input[i] = value
	}
}

func isByteArrayAllZero(input []byte) bool {
	for i := range input {
		if input[i] != 0 {
			return false
		}
	}
	return true
}
