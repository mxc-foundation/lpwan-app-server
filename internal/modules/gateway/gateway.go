package gateway

import (
	"context"
	"fmt"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/gws"
	"strings"

	"github.com/pkg/errors"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"

	psPb "github.com/mxc-foundation/lpwan-app-server/api/ps-serves-appserver"
	pscli "github.com/mxc-foundation/lpwan-app-server/internal/clients/psconn"
	"github.com/mxc-foundation/lpwan-app-server/internal/types"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/serverinfo"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

type controller struct {
	GatewayModuleInterface
	st                 *store.Handler
	ps                 pscli.ProvisioningServerStruct
	bindPortOldGateway string
	bindPortNewGateway string
}

var ctrl *controller

func Get() GatewayModuleInterface {
	return ctrl
}

type GatewayModuleInterface interface {
	Tx(ctx context.Context, f func(context.Context, *store.Handler) error) error
}

func (c *controller) Tx(ctx context.Context, f func(context.Context, *store.Handler) error) error {
	return c.st.Tx(ctx, f)
}

func Setup(h *store.Handler) error {
	ctrl = &controller{}
	SetupStore(h)
	return SetupFirmware()
}

func SetupStore(h *store.Handler) {
	ctrl.st = h
}

func SetupFirmware() error {
	ctrl.ps = pscli.ProvisioningServerStruct{
		Server:         ctrl.ps.Server,
		CACert:         ctrl.ps.CACert,
		TLSCert:        ctrl.ps.TLSCert,
		TLSKey:         ctrl.ps.TLSKey,
		UpdateSchedule: ctrl.ps.UpdateSchedule,
	}

	gwAPI := gws.GetSettings()
	if strArray := strings.Split(gwAPI.OldGateway.Bind, ":"); len(strArray) != 2 {
		return errors.New(fmt.Sprintf("Invalid API Bind settings for OldGateway: %s", gwAPI.OldGateway.Bind))
	} else {
		ctrl.bindPortOldGateway = strArray[1]
	}

	if strArray := strings.Split(gwAPI.NewGateway.Bind, ":"); len(strArray) != 2 {
		return errors.New(fmt.Sprintf("Invalid API Bind settings for NewGateway: %s", gwAPI.NewGateway.Bind))
	} else {
		ctrl.bindPortNewGateway = strArray[1]
	}

	return ctrl.updateFirmwareFromProvisioningServer(context.Background())
}

func (c *controller) updateFirmwareFromProvisioningServer(ctx context.Context) error {
	log.WithFields(log.Fields{
		"provisioning-server": c.ps.Server,
		"caCert":              c.ps.CACert,
		"tlsCert":             c.ps.TLSCert,
		"tlsKey":              c.ps.TLSKey,
		"schedule":            c.ps.UpdateSchedule,
	}).Info("Start schedule to update gateway firmware...")

	supernodeAddr := serverinfo.GetSettings().ServerAddr

	cron := cron.New()
	err := cron.AddFunc(c.ps.UpdateSchedule, func() {
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
