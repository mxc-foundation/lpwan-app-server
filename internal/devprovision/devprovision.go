package devprovision

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"

	rs "github.com/mxc-foundation/lpwan-app-server/internal/modules/redis"

	"github.com/gofrs/uuid"
	//	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/brocaar/chirpstack-api/go/v3/as"
	gwV3 "github.com/brocaar/chirpstack-api/go/v3/gw"
	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/brocaar/lorawan"
	duration "github.com/golang/protobuf/ptypes/duration"

	nsextra "github.com/mxc-foundation/lpwan-app-server/api/ns-extra"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/networkserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/networkserverextra"
	"github.com/mxc-foundation/lpwan-app-server/internal/logging"
	gwd "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway/data"
	nsd "github.com/mxc-foundation/lpwan-app-server/internal/networkserver_portal/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
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

// ProprietaryPayload - Proprietary Payload
type proprietaryPayload struct {
	MacPayload []byte
	GatewayMAC lorawan.EUI64
	Frequency  uint32
	DR         int
	Context    []byte
	Delay      *duration.Duration
}

func init() {
	mgr.RegisterModuleSetup(moduleName, Setup)
}

const moduleName = "devprovision"

var ctrl struct {
	handler *store.Handler

	moduleUp bool
}

// Setup prepares device provisioning service module
func Setup(name string, h *store.Handler) error {
	if ctrl.moduleUp == true {
		return nil
	}
	defer func() {
		ctrl.moduleUp = true
	}()

	ctrl.handler = h

	go SendPingLoop()

	return nil
}

// SendPingLoop is a never returning function sending the gateway pings.
func SendPingLoop() {
	for {
		ctxID, err := uuid.NewV4()
		if err != nil {
			log.WithError(err).Error("new uuid error")
		}

		ctx := context.Background()
		ctx = context.WithValue(ctx, logging.ContextIDKey, ctxID)

		if err := sendGatewayPing(ctx, ctrl.handler); err != nil {
			log.Errorf("send gateway ping error: %s", err)
		}
		time.Sleep(time.Second)
	}
}

// HandleReceivedPing handles a ping received by one or multiple gateways.
func HandleReceivedPing(ctx context.Context, req *as.HandleProprietaryUplinkRequest) (bool, error) {
	var processed bool = false
	var mic lorawan.MIC
	copy(mic[:], req.Mic)

	log.Infof("MacPayload:\n%s", hex.Dump(req.MacPayload))

	// Find max RSSI gw
	var maxRssiRx *gwV3.UplinkRXInfo = nil
	for _, rx := range req.RxInfo {
		if maxRssiRx == nil {
			maxRssiRx = rx
		} else if rx.Rssi > maxRssiRx.Rssi {
			maxRssiRx = rx
		}
	}
	log.Infof("  MAC:%s, RSSI: %d, Context: %s", hex.EncodeToString(maxRssiRx.GatewayId), maxRssiRx.Rssi,
		hex.EncodeToString(maxRssiRx.Context))

	// Get Gateway
	var mac lorawan.EUI64
	copy(mac[:], maxRssiRx.GatewayId)

	gw, err := ctrl.handler.GetGateway(ctx, mac, false)
	if err != nil {
		return processed, errors.Wrap(err, "get gateway error")
	}

	n, err := ctrl.handler.GetNetworkServer(ctx, gw.NetworkServerID)
	if err != nil {
		return processed, errors.Wrap(err, "get network-server error")
	}
	log.Infof("  NetworkServer: %s", n.Server)

	//
	var upFreqChannel uint32 = (req.TxInfo.Frequency - 470300000) / 200000
	var downFreq uint32 = 500300000 + ((upFreqChannel % 48) * 200000)

	// Check Message Type
	var messageType byte = req.MacPayload[0]
	if messageType == UpMessageHello {
		log.Info("  HELLO Message.")

		payload := proprietaryPayload{
			MacPayload: []byte("HELLO"),
			GatewayMAC: mac,
			Frequency:  downFreq,
			DR:         3,
			Delay:      &duration.Duration{Seconds: 5, Nanos: 0},
			Context:    maxRssiRx.Context,
		}

		err = sendProprietary(n, payload)
		if err != nil {
			return processed, errors.Wrap(err, "send proprietary error")
		}

		processed = true
	} else if messageType == UpMessageAuth {
		log.Info("  AUTH Message.")
		processed = true
	}

	//	err := ctrl.handler.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
	// for _, rx := range req.RxInfo {
	// 	var mac lorawan.EUI64
	// 	copy(mac[:], rx.GatewayId)

	// 	// ignore pings received by the sending gateway
	// 	if ping.GatewayMAC == mac {
	// 		continue
	// 	}

	// 	var receivedAt *time.Time
	// 	if rx.Time != nil {
	// 		ts, err := ptypes.Timestamp(rx.Time)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		receivedAt = &ts
	// 	}

	// 	pingRX := gwd.GatewayPingRX{
	// 		PingID:     id,
	// 		GatewayMAC: mac,
	// 		ReceivedAt: receivedAt,
	// 		RSSI:       int(rx.Rssi),
	// 		LoRaSNR:    rx.LoraSnr,
	// 	}

	// 	if rx.Location != nil {
	// 		pingRX.Location = gwd.GPSPoint{
	// 			Latitude:  rx.Location.Latitude,
	// 			Longitude: rx.Location.Longitude,
	// 		}
	// 		pingRX.Altitude = rx.Location.Altitude
	// 	}

	// 	err := handler.CreateGatewayPingRX(ctx, &pingRX)
	// 	if err != nil {
	// 		return errors.Wrap(err, "create gateway ping rx error")
	// 	}
	// }
	//return false, nil
	//	})
	//	if err != nil {
	//		return processed, errors.Wrap(err, "transaction error")
	//	}

	return processed, nil
}

// sendGatewayPing selects the next gateway to ping, creates the "ping"
// frame and sends this frame to the network-server for transmission.
func sendGatewayPing(ctx context.Context, handler *store.Handler) error {
	err := handler.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		gw, err := handler.GetGatewayForPing(ctx)
		if err != nil {
			return errors.Wrap(err, "get gateway for ping error")
		}
		if gw == nil {
			return nil
		}

		n, err := handler.GetNetworkServer(ctx, gw.NetworkServerID)
		if err != nil {
			return errors.Wrap(err, "get network-server error")
		}

		ping := gwd.GatewayPing{
			GatewayMAC: gw.MAC,
			Frequency:  n.GatewayDiscoveryTXFrequency,
			DR:         n.GatewayDiscoveryDR,
		}
		err = handler.CreateGatewayPing(ctx, &ping)
		if err != nil {
			return errors.Wrap(err, "create gateway ping error")
		}

		var mic lorawan.MIC
		if _, err = rand.Read(mic[:]); err != nil {
			return errors.Wrap(err, "read random bytes error")
		}

		err = CreatePingLookup(mic, ping.ID)
		if err != nil {
			return errors.Wrap(err, "store mic lookup error")
		}

		err = sendPing(mic, n, ping)
		if err != nil {
			return errors.Wrap(err, "send ping error")
		}

		gw.LastPingID = &ping.ID
		gw.LastPingSentAt = &ping.CreatedAt

		err = handler.UpdateGateway(ctx, gw)
		if err != nil {
			return errors.Wrap(err, "update gateway error")
		}

		return nil
	})

	return err
}

func sendProprietary(n nsd.NetworkServer, payload proprietaryPayload) error {
	nsClient, err := networkserverextra.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err != nil {
		return errors.Wrap(err, "get network-server client error")
	}
	_, err = nsClient.SendDelayedProprietaryPayload(context.Background(), &nsextra.SendDelayedProprietaryPayloadRequest{
		MacPayload:            payload.MacPayload,
		GatewayMacs:           [][]byte{payload.GatewayMAC[:]},
		PolarizationInversion: true,
		Frequency:             uint32(payload.Frequency),
		Dr:                    uint32(payload.DR),
		Context:               payload.Context,
		Delay:                 payload.Delay,
	})
	if err != nil {
		return errors.Wrap(err, "send proprietary payload error")
	}

	log.WithFields(log.Fields{
		"gateway_mac": payload.GatewayMAC,
		"freq":        payload.Frequency,
	}).Info("gateway proprietary payload sent to network-server")

	return nil
}

func sendPing(mic lorawan.MIC, n nsd.NetworkServer, ping gwd.GatewayPing) error {
	nsClient, err := networkserver.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err != nil {
		return errors.Wrap(err, "get network-server client error")
	}

	_, err = nsClient.SendProprietaryPayload(context.Background(), &ns.SendProprietaryPayloadRequest{
		Mic:                   mic[:],
		GatewayMacs:           [][]byte{ping.GatewayMAC[:]},
		PolarizationInversion: false,
		Frequency:             uint32(ping.Frequency),
		Dr:                    uint32(ping.DR),
	})
	if err != nil {
		return errors.Wrap(err, "send proprietary payload error")
	}

	log.WithFields(log.Fields{
		"gateway_mac": ping.GatewayMAC,
		"id":          ping.ID,
	}).Info("gateway ping sent to network-server")

	return nil
}

// CreatePingLookup creates an automatically expiring MIC to ping id lookup.
func CreatePingLookup(mic lorawan.MIC, id int64) error {
	keyWord := fmt.Sprintf("%s", mic)
	return rs.RedisClient().Set(fmt.Sprintf(rs.MicLookupTempl, keyWord), id, rs.MicLookupExpire).Err()
}

// GetPingLookup :
func GetPingLookup(mic lorawan.MIC) (int64, error) {
	keyWord := fmt.Sprintf("%s", mic)
	return rs.RedisClient().Get(fmt.Sprintf(rs.MicLookupTempl, keyWord)).Int64()
}

// DeletePingLookup :
func DeletePingLookup(mic lorawan.MIC) error {
	keyWord := fmt.Sprintf("%s", mic)
	return rs.RedisClient().Del(fmt.Sprintf(rs.MicLookupTempl, keyWord)).Err()
}
