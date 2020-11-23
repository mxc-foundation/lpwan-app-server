package mining

import (
	"context"
	"time"

	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"

	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"

	api "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/mining/data"
	mxprotocolconn "github.com/mxc-foundation/lpwan-app-server/internal/mxp_portal"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
)

func init() {
	mgr.RegisterSettingsSetup(moduleName, SettingsSetup)
	mgr.RegisterModuleSetup(moduleName, Setup)
}

const moduleName = "mining"

// controller regularly checks what gateways should be paid for mining and
// sends request to m2m to pay them
type controller struct {
	s         Config
	m2mClient api.MiningServiceClient
	st        *store.Handler

	moduleUp bool
}

var ctrl *controller

// SettingsSetup initialize module settings on start
func SettingsSetup(name string, s config.Config) error {

	ctrl = &controller{
		s: s.ApplicationServer.MiningSetUp,
	}
	return nil
}
func GetSettings() Config {
	return ctrl.s
}

func Setup(name string, h *store.Handler) (err error) {
	if ctrl.moduleUp == true {
		return nil
	}
	defer func() {
		ctrl.moduleUp = true
	}()

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

	var gws []*api.GatewayMining

	// update the first heartbeat = 0
	for _, v := range miningGws {
		err := ctrl.st.UpdateFirstHeartbeatToZero(ctx, v.GatewayMac)
		if err != nil {
			log.WithError(err).Error("tokenMining/update first heartbeat to zero error")
		}
		gw := api.GatewayMining{
			GatewayMac: v.GatewayMac.String(),
			OwnerOrgId: v.OwnerOrgID,
			StcOrgId:   0,
		}
		if v.StcOrgID != nil {
			gw.StcOrgId = *v.StcOrgID
		}

		gws = append(gws, &gw)
	}

	// if error, resend after one minute
	for {
		if err := ctrl.sendMining(ctx, gws); err != nil {
			log.WithError(err).Error("send mining request to m2m error")
			time.Sleep(60 * time.Second)
			continue
		}
		break
	}

	return nil
}

func (ctrl *controller) sendMining(ctx context.Context, gws []*api.GatewayMining) error {
	_, err := ctrl.m2mClient.Mining(ctx, &api.MiningRequest{
		GatewayMining: gws,
		PeriodSeconds: ctrl.s.Period,
	})

	return err
}
