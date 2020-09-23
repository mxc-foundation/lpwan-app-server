package gwping

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	rs "github.com/mxc-foundation/lpwan-app-server/internal/modules/redis"

	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/brocaar/chirpstack-api/go/v3/as"
	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/brocaar/lorawan"

	"github.com/mxc-foundation/lpwan-app-server/internal/backend/networkserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/logging"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

var ctrl struct {
	handler *store.Handler
}

func Setup(handler *store.Handler) {
	ctrl.handler = handler
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
func HandleReceivedPing(ctx context.Context, req *as.HandleProprietaryUplinkRequest) error {
	var mic lorawan.MIC
	copy(mic[:], req.Mic)

	id, err := GetPingLookup(mic)
	if err != nil {
		return errors.Wrap(err, "get ping lookup error")
	}

	if err = DeletePingLookup(mic); err != nil {
		log.Errorf("delete ping lookup error: %s", err)
	}

	ping, err := ctrl.handler.GetGatewayPing(ctx, id)
	if err != nil {
		return errors.Wrap(err, "get gateway ping error")
	}

	err = ctrl.handler.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		for _, rx := range req.RxInfo {
			var mac lorawan.EUI64
			copy(mac[:], rx.GatewayId)

			// ignore pings received by the sending gateway
			if ping.GatewayMAC == mac {
				continue
			}

			var receivedAt *time.Time
			if rx.Time != nil {
				ts, err := ptypes.Timestamp(rx.Time)
				if err != nil {
					return err
				}
				receivedAt = &ts
			}

			pingRX := store.GatewayPingRX{
				PingID:     id,
				GatewayMAC: mac,
				ReceivedAt: receivedAt,
				RSSI:       int(rx.Rssi),
				LoRaSNR:    rx.LoraSnr,
			}

			if rx.Location != nil {
				pingRX.Location = store.GPSPoint{
					Latitude:  rx.Location.Latitude,
					Longitude: rx.Location.Longitude,
				}
				pingRX.Altitude = rx.Location.Altitude
			}

			err := handler.CreateGatewayPingRX(ctx, &pingRX)
			if err != nil {
				return errors.Wrap(err, "create gateway ping rx error")
			}
		}
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "transaction error")
	}

	return nil
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

		ping := store.GatewayPing{
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

func sendPing(mic lorawan.MIC, n store.NetworkServer, ping store.GatewayPing) error {
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
