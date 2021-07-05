package devprovision

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"sync"
	"time"

	"github.com/jacobsa/crypto/cmac"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/brocaar/chirpstack-api/go/v3/as"
	gwV3 "github.com/brocaar/chirpstack-api/go/v3/gw"
	"github.com/brocaar/lorawan"
	duration "github.com/golang/protobuf/ptypes/duration"

	nsPb "github.com/mxc-foundation/lpwan-app-server/api/networkserver"
	psPb "github.com/mxc-foundation/lpwan-app-server/api/ps-serves-appserver"
	nsd "github.com/mxc-foundation/lpwan-app-server/internal/api/external/ns"
	"github.com/mxc-foundation/lpwan-app-server/internal/devprovision/ecdh"
	gwd "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/nscli"
	"github.com/mxc-foundation/lpwan-app-server/internal/pscli"
)

// LoRa Frame Message and Response Type
//#define MAX_MESSAGE_SIZE 256
const (
	upMessageHello     = 0x01
	upMessageAuth      = 0x11
	downRespHello      = 0x81
	downRespAuthAccept = 0x91
	downRespAuthReject = 0x92
)
const (
	sizeUpMessageHello = 74
	sizeUpMessageAuth  = 61
)

// Proprietary Payload
type proprietaryPayload struct {
	MacPayload      []byte
	GatewayMAC      lorawan.EUI64
	DownlinkFreq    uint32
	UplinkFreq      uint32
	UplinkBandwidth uint32
	UplinkSf        uint32
	Context         []byte
	Delay           *duration.Duration
	Mic             []byte
}

type deviceSessionList struct {
	maxNumberOfDevSession  int //5000
	sessionlist            map[uint64]deviceSession
	mutexSessionList       sync.RWMutex
	deviceSessionLifeCycle time.Duration //time.Minute * 5
}

type nsClient interface {
	GetNetworkServerExtraServiceClient(networkServerID int64) (nsPb.NetworkServerExtraServiceClient, error)
}

type controller struct {
	psCli          psPb.DeviceProvisionClient
	nsCli          nsClient
	devSessionList deviceSessionList
}

var ctrl *controller

func (c *controller) sendToNs(networkServerID int64, req *nsPb.SendDelayedProprietaryPayloadRequest) error {
	nsClient, err := c.nsCli.GetNetworkServerExtraServiceClient(networkServerID)
	if err != nil {
		return err
	}
	_, err = nsClient.SendDelayedProprietaryPayload(context.Background(), req)
	if err != nil {
		return err
	}

	return nil
}

// Start prepares device provisioning service module
func Start(psCli *pscli.Client, nsCli *nscli.Client) error {
	ctrl = &controller{
		nsCli: nsCli,
		psCli: psCli.GetDeviceProvisionServiceClient(),
		devSessionList: deviceSessionList{
			sessionlist:            make(map[uint64]deviceSession),
			mutexSessionList:       sync.RWMutex{},
			maxNumberOfDevSession:  5000,
			deviceSessionLifeCycle: time.Minute * 5,
		},
	}

	ctrl.devSessionList.clearDeviceSessionList()

	go ctrl.cleanUpLoop()

	return nil
}

// cleanUpLoop is a never returning function, performing cleanup
func (c *controller) cleanUpLoop() {
	for {
		c.devSessionList.clearExpiredDevSession()
		time.Sleep(time.Second * 10)
	}
}

func (c *controller) processMessage(ctx context.Context, nID int64, req *as.HandleProprietaryUplinkRequest,
	targetgateway *gwV3.UplinkRXInfo) (bool, error) {
	processed := false
	messageType := req.MacPayload[0]
	messageSize := len(req.MacPayload)

	if (messageType == upMessageHello) && (messageSize == sizeUpMessageHello) {
		err := ctrl.handleHello(ctx, nID, req, targetgateway)
		if err != nil {
			return false, errors.Wrap(err, "process HELLO msg error")
		}
		processed = true
	} else if (messageType == upMessageAuth) && (messageSize == sizeUpMessageAuth) {
		err := ctrl.handleAuth(ctx, nID, req, targetgateway)
		if err != nil {
			return false, errors.Wrap(err, "process AUTH msg error")
		}
		processed = true
	} else {
		log.Debug("Unknown Message.")
	}

	return processed, nil
}

// Store defines db API used by device provision service
type Store interface {
	GetGateway(ctx context.Context, mac lorawan.EUI64, forUpdate bool) (gwd.Gateway, error)
	GetNetworkServer(ctx context.Context, id int64) (nsd.NetworkServer, error)
}

// HandleReceivedFrame handles a ping received by one or multiple gateways.
func HandleReceivedFrame(ctx context.Context, req *as.HandleProprietaryUplinkRequest, h Store) (bool, error) {
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
	gw, err = h.GetGateway(ctx, mac, false)
	if err != nil {
		return false, errors.Wrap(err, "get gateway error")
	}
	n, err = h.GetNetworkServer(ctx, gw.NetworkServerID)
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
	return ctrl.processMessage(ctx, n.ID, req, maxRssiRx)
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

func (c *controller) sendProprietary(networkServerID int64, payload proprietaryPayload) error {
	req := nsPb.SendDelayedProprietaryPayloadRequest{
		MacPayload:            payload.MacPayload,
		GatewayMacs:           [][]byte{payload.GatewayMAC[:]},
		PolarizationInversion: true,
		UplinkFreq:            payload.UplinkFreq,
		DownlinkFreq:          payload.DownlinkFreq,
		UplinkBandwidth:       payload.UplinkBandwidth,
		UplinkSf:              payload.UplinkSf,
		Context:               payload.Context,
		Delay:                 payload.Delay,
		Mic:                   calProprietaryMic(payload.MacPayload),
	}
	log.Debugf("  sendProprietary() MIC: %s", hex.EncodeToString(req.Mic))

	err := c.sendToNs(networkServerID, &req)
	if err != nil {
		return errors.Wrap(err, "send proprietary payload error")
	}
	log.WithFields(log.Fields{
		"gateway_mac": payload.GatewayMAC,
		"up_freq":     payload.UplinkFreq,
		"up_bw":       payload.UplinkBandwidth,
		"up_sf":       payload.UplinkSf,
		"down_freq":   payload.DownlinkFreq,
	}).Infof("gateway proprietary payload sent to network server %d", networkServerID)

	return nil
}

func makeHelloResponse(session deviceSession) []byte {
	payload := []byte{downRespHello}
	payload = append(payload, session.rDevEui...)
	payload = append(payload, session.serverPublicKey...)
	payload = append(payload, session.serverNonce...)
	return payload
}

func (c *controller) handleHello(ctx context.Context, nID int64, req *as.HandleProprietaryUplinkRequest,
	targetgateway *gwV3.UplinkRXInfo) error {
	log.Debug("  HELLO Message.")

	var err error
	var frameversion byte

	rdeveui := make([]byte, 8)
	copy(rdeveui[0:], req.MacPayload[1:])
	sessionid := binary.BigEndian.Uint64(rdeveui)
	log.Debugf("  sessionid=%X", sessionid)
	frameversion = req.MacPayload[73]

	ok, currentsession := c.devSessionList.searchDeviceSession(sessionid)
	if !ok {
		rdeveui := make([]byte, 8)
		devicepublickey := make([]byte, ecdh.K233PubKeySize)

		log.Debugf("  Creating new session")
		copy(rdeveui[0:], req.MacPayload[1:])
		copy(devicepublickey[0:], req.MacPayload[9:])
		ok, currentsession = c.devSessionList.createDeviceSession(sessionid, rdeveui, devicepublickey)
		if !ok {
			// Create session failed. drop this frame. return true to mark is processed.
			return nil
		}
	}

	// Drop if already sent to the same Gateway context
	ok, currentsession = c.devSessionList.checkDeviceSession(sessionid, targetgateway.Context)
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
		MacPayload:      makeHelloResponse(currentsession),
		GatewayMAC:      mac,
		UplinkFreq:      req.TxInfo.Frequency,
		UplinkBandwidth: req.TxInfo.GetLoraModulationInfo().GetBandwidth(),
		UplinkSf:        req.TxInfo.GetLoraModulationInfo().SpreadingFactor,
		DownlinkFreq:    0,
		Delay:           &duration.Duration{Seconds: 5, Nanos: 0},
		Context:         targetgateway.Context,
		Mic:             []byte{0x00, 0x00, 0x00, 0x00},
	}
	// log.Debugf("Tx MacPayload:\n%s", hex.Dump(payload.MacPayload))

	err = c.sendProprietary(nID, payload)
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

	payload := []byte{downRespAuthAccept}
	payload = append(payload, session.rDevEui...)
	payload = append(payload, encpayload...)

	return payload
}

func makeAuthReject(session deviceSession) []byte {
	payload := []byte{downRespAuthReject}
	payload = append(payload, session.rDevEui...)
	return payload
}

func (c *controller) handleAuth(ctx context.Context, nID int64, req *as.HandleProprietaryUplinkRequest,
	targetgateway *gwV3.UplinkRXInfo) error {
	log.Debug("  AUTH Message.")

	//
	rdeveui := make([]byte, 8)
	copy(rdeveui[0:], req.MacPayload[1:])
	sessionid := binary.BigEndian.Uint64(rdeveui)
	log.Debugf("  sessionid=%X", sessionid)

	//
	ok, currentsession := c.devSessionList.searchDeviceSession(sessionid)
	if !ok {
		log.Debugf("  Auth message without active session. Frame dropped.")
		return nil
	}

	// Drop if already sent to the same Gateway context
	ok, currentsession = c.devSessionList.checkDeviceSession(sessionid, targetgateway.Context)
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

	authaccepted := true
	found, deviceinfo := findDeviceBySnHash(ctx, privisionidhash, c.psCli)
	if !found {
		return errors.Errorf("Device %s not found.", hex.EncodeToString(privisionidhash))
	} else if deviceinfo.Status == "DISABLED" {
		log.Errorf("Device %s disabled.", deviceinfo.ProvisionID)
		authaccepted = false
	}
	log.Debugf("  Device found. %s, mfgID=%d, server=%s", deviceinfo.ProvisionID, deviceinfo.ManufacturerID, deviceinfo.Server)
	log.Debugf("  devEUI=%s, appEUI=%s, appKey=%s, nwkKey=%s",
		hex.EncodeToString(deviceinfo.DevEUI), hex.EncodeToString(deviceinfo.AppEUI),
		hex.EncodeToString(deviceinfo.AppKey), hex.EncodeToString(deviceinfo.NwkKey))
	log.Debugf("  status=%v, model=%v, fixedDevEUI=%v, created=%v", deviceinfo.Status, deviceinfo.Model, deviceinfo.FixedDevEUI,
		deviceinfo.TimeCreated)
	if deviceinfo.Server != "" {
		log.Errorf("Device %s registered to %v, provisioning not allowed.", deviceinfo.ProvisionID, deviceinfo.Server)
		authaccepted = false
	}

	if authaccepted {
		calverifycode := currentsession.calVerifyCode(deviceinfo.ProvisionID, true)
		if !bytes.Equal(verifycode, calverifycode) {
			return errors.Errorf("Incorrect verify code at Auth message")
		}

		currentsession, deviceinfo, err := c.updateDevice(ctx, currentsession, deviceinfo)
		if err != nil {
			return errors.Wrap(err, "updateDevice error")
		}
		c.devSessionList.updateDeviceSession(sessionid, currentsession)

		err = saveDevice(ctx, deviceinfo, c.psCli)
		if err != nil {
			return errors.Wrap(err, "saveDevice error")
		}
	}

	//
	var mac lorawan.EUI64
	copy(mac[:], targetgateway.GatewayId)
	verifycode = currentsession.calVerifyCode(deviceinfo.ProvisionID, false)

	payload := proprietaryPayload{
		GatewayMAC:      mac,
		UplinkFreq:      req.TxInfo.Frequency,
		DownlinkFreq:    0,
		UplinkBandwidth: req.TxInfo.GetLoraModulationInfo().GetBandwidth(),
		UplinkSf:        req.TxInfo.GetLoraModulationInfo().SpreadingFactor,
		Delay:           &duration.Duration{Seconds: 5, Nanos: 0},
		Context:         targetgateway.Context,
		Mic:             []byte{0x00, 0x00, 0x00, 0x00},
	}
	if authaccepted {
		payload.MacPayload = makeAuthAccept(currentsession, verifycode)
	} else {
		payload.MacPayload = makeAuthReject(currentsession)
	}
	// log.Debugf("Tx MacPayload:\n%s", hex.Dump(payload.MacPayload))

	err := c.sendProprietary(nID, payload)
	if err != nil {
		return err
	}

	return nil
}

// Device session handling
func (l *deviceSessionList) searchDeviceSession(sessionid uint64) (bool, deviceSession) {
	l.mutexSessionList.Lock()
	currentsession, sessionfound := l.sessionlist[sessionid]
	l.mutexSessionList.Unlock()

	if !sessionfound {
		return false, deviceSession{}
	}
	return true, currentsession
}

func (l *deviceSessionList) updateDeviceSession(sessionid uint64, newsession deviceSession) {
	l.mutexSessionList.Lock()
	_, sessionfound := l.sessionlist[sessionid]
	l.mutexSessionList.Unlock()

	if sessionfound {
		l.sessionlist[sessionid] = newsession
	}
}

func (l *deviceSessionList) checkDeviceSession(sessionid uint64, gwcontext []byte) (bool, deviceSession) {
	l.mutexSessionList.Lock()
	defer l.mutexSessionList.Unlock()

	currentsession, sessionfound := l.sessionlist[sessionid]
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
	l.sessionlist[sessionid] = currentsession

	return true, currentsession
}

func (l *deviceSessionList) createDeviceSession(sessionid uint64, rdeveui []byte, devicepublickey []byte) (bool, deviceSession) {
	l.mutexSessionList.Lock()
	defer l.mutexSessionList.Unlock()

	if len(l.sessionlist) >= l.maxNumberOfDevSession {
		log.Warnf("Maximum number (%d) of device provisioning session reached. Request dropped.", l.maxNumberOfDevSession)
		return false, deviceSession{}
	}

	// New session
	currentsession := makeDeviceSession(l.deviceSessionLifeCycle)
	copy(currentsession.rDevEui[0:], rdeveui)
	copy(currentsession.devicePublicKey[0:], devicepublickey)

	currentsession.genServerKeys()
	currentsession.genSharedKey()
	currentsession.deriveKeys()
	l.sessionlist[sessionid] = currentsession

	return true, currentsession
}

func (l *deviceSessionList) clearExpiredDevSession() {
	l.mutexSessionList.Lock()
	now := time.Now()
	for key, session := range l.sessionlist {
		if now.After(session.expireTime) {
			delete(l.sessionlist, key)
		}
	}
	l.mutexSessionList.Unlock()
}

func (l *deviceSessionList) clearDeviceSessionList() {
	l.mutexSessionList.Lock()
	for k := range l.sessionlist {
		delete(l.sessionlist, k)
	}
	l.mutexSessionList.Unlock()
}

func (c *controller) updateDevice(ctx context.Context, session deviceSession, deviceinfo deviceInfo) (deviceSession, deviceInfo, error) {

	if isByteArrayAllZero(session.assignedDevEui) || !bytes.Equal(session.assignedDevEui, deviceinfo.DevEUI) {
		// Session is new
		deveui := make([]byte, 8)
		appeui := make([]byte, 8)
		copy(deveui[:], deviceinfo.DevEUI)
		copy(appeui[:], deviceinfo.AppEUI)

		if !deviceinfo.FixedDevEUI || isByteArrayAllZero(deveui) {
			// Generate devEUI
			resp, err := c.psCli.GenDevEUI(ctx, &psPb.GenDevEuiRequest{})
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
