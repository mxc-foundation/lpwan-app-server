package mining

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/brocaar/lorawan"
	log "github.com/sirupsen/logrus"

	api "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	mxprotocolconn "github.com/mxc-foundation/lpwan-app-server/internal/clients/mxprotocol-server"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

// Config contains mining configuration
type Config struct {
	// If mining is enabled or not
	Enabled bool `mapstructure:"enabled"`
	// If we haven't got heartbeat for HeartbeatOfflineLimit seconds, we
	// consider gateway to be offline
	HeartbeatOfflineLimit int64 `mapstructure:"heartbeat_offline_limit"`
	// Gateway must be online for at leasts GwOnlineLimit seconds to receive mining reward
	GwOnlineLimit int64 `mapstructure:"gw_online_limit"`
	// Period is the length of the mining period in seconds
	Period int64 `mapstructure:"period"`
}

// controller regularly checks what gateways should be paid for mining and
// sends request to m2m to pay them
type controller struct {
	s         Config
	m2mClient api.MiningServiceClient
	st        *store.Handler
}

var ctrl *controller

func SettingsSetup(s Config) error {
	ctrl = &controller{
		s: s,
	}
	return nil
}
func GetSettings() Config {
	return ctrl.s
}

func Setup(h *store.Handler) (err error) {
	if !ctrl.s.Enabled {
		return nil
	}

	log.Info("mining cron task begin...")

	ctrl.st = h
	ctrl.m2mClient, err = mxprotocolconn.GetMiningServiceClient()
	if err != nil {
		return errors.Wrap(err, "get m2m mining service client error")
	}

	go func() {
		period := time.Duration(ctrl.s.Period) * time.Second
		for {
			nextRun := time.Now().Add(period).Truncate(period)
			time.Sleep(time.Until(nextRun))
			if err := ctrl.submitMining(context.Background()); err != nil {
				log.WithError(err).Error("couldn't submit mining")
			}
		}
	}()

	return nil
}

func (ctrl *controller) submitMining(ctx context.Context) error {
	current_time := time.Now().Unix()
	log.Infof("processing mining")

	// get the gateway list that should receive the mining tokens
	miningGws, err := ctrl.st.GetGatewayMiningList(
		ctx, current_time, ctrl.s.GwOnlineLimit,
	)
	if err != nil {
		log.WithError(err).Error("Cannot get mining gateway list from DB.")
		return err
	}

	if len(miningGws) == 0 {
		return nil
	}

	var macs []string

	// update the first heartbeat = 0
	for _, v := range miningGws {
		err := ctrl.st.UpdateFirstHeartbeatToZero(ctx, v)
		if err != nil {
			log.WithError(err).Error("tokenMining/update first heartbeat to zero error")
		}
		mac := lorawan.EUI64.String(v)
		macs = append(macs, mac)
	}

	// if error, resend after one minute
	for {
		if err := ctrl.sendMining(ctx, macs); err != nil {
			log.WithError(err).Error("send mining request to m2m error")
			time.Sleep(60 * time.Second)
			continue
		}
		break
	}

	return nil
}

func (ctrl *controller) sendMining(ctx context.Context, macs []string) error {
	_, err := ctrl.m2mClient.Mining(ctx, &api.MiningRequest{
		GatewayMac:    macs,
		PeriodSeconds: ctrl.s.Period,
	})

	return err
}
