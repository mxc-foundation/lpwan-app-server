package mining

import (
	"context"
	"time"

	"github.com/brocaar/lorawan"
	log "github.com/sirupsen/logrus"

	api "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
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

// Store is interface to DB that stores information about gateways and heartbeat times
type Store interface {
	// GetGatewayMiningList returns the list of gateways that were online for
	// at least onlineLimit
	GetGatewayMiningList(ctx context.Context, time, onlineLimit int64) ([]lorawan.EUI64, error)
	UpdateFirstHeartbeatToZero(ctx context.Context, mac lorawan.EUI64) error
}

// Controller regularly checks what gateways should be paid for mining and
// sends request to m2m to pay them
type Controller struct {
	gwOnlineLimit int64
	period        int64
	m2mClient     api.MiningServiceClient
	store         Store
}

func Setup(conf Config, store Store, m2mClient api.MiningServiceClient) error {
	if !conf.Enabled {
		return nil
	}

	log.Info("mining cron task begin...")

	ctrl := &Controller{
		gwOnlineLimit: conf.GwOnlineLimit,
		period:        conf.Period,
		store:         store,
		m2mClient:     m2mClient,
	}

	go func() {
		period := time.Duration(ctrl.period) * time.Second
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

func (ctrl *Controller) submitMining(ctx context.Context) error {
	current_time := time.Now().Unix()
	log.Infof("processing mining")

	// get the gateway list that should receive the mining tokens
	miningGws, err := ctrl.store.GetGatewayMiningList(
		ctx, current_time, ctrl.gwOnlineLimit,
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
		err := ctrl.store.UpdateFirstHeartbeatToZero(ctx, v)
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

func (ctrl *Controller) sendMining(ctx context.Context, macs []string) error {
	_, err := ctrl.m2mClient.Mining(ctx, &api.MiningRequest{
		GatewayMac:    macs,
		PeriodSeconds: ctrl.period,
	})

	return err
}
