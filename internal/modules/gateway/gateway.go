package gateway

import (
	"context"

	"github.com/pkg/errors"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"

	psPb "github.com/mxc-foundation/lpwan-app-server/api/ps-serves-appserver"
	pscli "github.com/mxc-foundation/lpwan-app-server/internal/clients/psconn"
	"github.com/mxc-foundation/lpwan-app-server/internal/types"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

type Config struct {
	ServerAddr string
}

type controller struct {
	st                 *store.Handler
	ps                 pscli.ProvisioningServerStruct
	bindPortOldGateway string
	bindPortNewGateway string
	s                  Config
}

var ctrl *controller

func SettingsSetup(s Config) error {
	ctrl.s = s
	return nil
}

func Setup(h *store.Handler) error {
	// must be called after setup gateway API
	if ctrl.bindPortNewGateway == "" || ctrl.bindPortOldGateway == "" {
		return errors.New("bindPortNewGateway and bindPortOldGateway not initiated")
	}

	ctrl.st = h
	ctrl.ps = pscli.GetSettings()
	return ctrl.updateFirmwareFromProvisioningServer(context.Background())
}

func SetupFirmware(bindOld, bindNew string) {
	ctrl = &controller{
		bindPortOldGateway: bindOld,
		bindPortNewGateway: bindNew,
	}
}

func (c *controller) updateFirmwareFromProvisioningServer(ctx context.Context) error {
	log.WithFields(log.Fields{
		"provisioning-server": ctrl.ps.Server,
		"caCert":              ctrl.ps.CACert,
		"tlsCert":             ctrl.ps.TLSCert,
		"tlsKey":              ctrl.ps.TLSKey,
		"schedule":            ctrl.ps.UpdateSchedule,
	}).Info("Start schedule to update gateway firmware...")

	supernodeAddr := ctrl.s.ServerAddr

	cron := cron.New()
	err := cron.AddFunc(ctrl.ps.UpdateSchedule, func() {
		log.Info("Check firmware update...")
		gwFwList, err := c.st.GetGatewayFirmwareList(ctx)
		if err != nil {
			log.WithError(err).Errorf("Failed to get gateway firmware list.")
			return
		}

		// send update
		psClient, err := pscli.CreateClientWithCert()
		if err != nil {
			log.WithError(err).Errorf("Create Provisioning server client error")
			return
		}

		for _, v := range gwFwList {
			res, err := psClient.GetUpdate(context.Background(), &psPb.GetUpdateRequest{
				Model:          v.Model,
				SuperNodeAddr:  supernodeAddr,
				PortOldGateway: c.bindPortOldGateway,
				PortNewGateway: c.bindPortNewGateway,
			})
			if err != nil {
				log.WithError(err).Errorf("Failed to get update for gateway model: %s", v.Model)
				continue
			}

			var md5sum types.MD5SUM
			if err := md5sum.UnmarshalText([]byte(res.FirmwareHash)); err != nil {
				log.WithError(err).Errorf("Failed to unmarshal firmware hash: %s", res.FirmwareHash)
				continue
			}

			gatewayFw := store.GatewayFirmware{
				Model:        v.Model,
				ResourceLink: res.ResourceLink,
				FirmwareHash: md5sum,
			}

			model, _ := c.st.UpdateGatewayFirmware(ctx, &gatewayFw)
			if model == "" {
				log.Warnf("No row updated for gateway_firmware at model=%s", v.Model)
			}

		}
	})
	if err != nil {
		log.Fatalf("Failed to set update schedule when set up provisioning server config: %s", err.Error())
	}

	go cron.Start()

	return nil
}
