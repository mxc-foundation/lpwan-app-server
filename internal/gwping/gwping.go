package gwping

import (
	"context"
	"crypto/rand"
	"time"

	gwmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway"

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
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
)

// SendPingLoop is a never returning function sending the gateway pings.
func SendPingLoop() {
	for {
		ctxID, err := uuid.NewV4()
		if err != nil {
			log.WithError(err).Error("new uuid error")
		}

		ctx := context.Background()
		ctx = context.WithValue(ctx, logging.ContextIDKey, ctxID)

		if err := sendGatewayPing(ctx); err != nil {
			log.Errorf("send gateway ping error: %s", err)
		}
		time.Sleep(time.Second)
	}
}

// HandleReceivedPing handles a ping received by one or multiple gateways.
func HandleReceivedPing(ctx context.Context, req *as.HandleProprietaryUplinkRequest) error {
	var mic lorawan.MIC
	copy(mic[:], req.Mic)

	gp := store.GatewayPing{}
	id, err := gp.GetPingLookup(mic)
	if err != nil {
		return errors.Wrap(err, "get ping lookup error")
	}

	if err = gp.DeletePingLookup(mic); err != nil {
		log.Errorf("delete ping lookup error: %s", err)
	}

	ping, err := storage.GetGatewayPing(ctx, storage.DB(), id)
	if err != nil {
		return errors.Wrap(err, "get gateway ping error")
	}

	err = storage.Transaction(func(ctx context.Context, handler *store.Handler) error {
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
func sendGatewayPing(ctx context.Context) error {
	err := gwmod.Get().Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
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

		err = ping.CreatePingLookup(mic)
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
