package gateway

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"

	psPb "github.com/mxc-foundation/lpwan-app-server/api/ps-serves-appserver"
	pscli "github.com/mxc-foundation/lpwan-app-server/internal/clients/psconn"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/types"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/serverinfo"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

type Pserver struct {
	ProvisionServer string
	CACert          string
	TLSCert         string
	TLSKey          string
	UpdateSchedule  string
}

type Controller struct {
	GatewayModuleInterface
	st                 *store.Handler
	ps                 Pserver
	bindPortOldGateway string
	bindPortNewGateway string
}

var service = &Controller{}

func Get() GatewayModuleInterface {
	return service
}

type GatewayModuleInterface interface {
	Handler() *store.Handler
}

func (c *Controller) Handler() *store.Handler {
	return c.st
}

func Setup(conf config.Config, h *store.Handler) error {
	service.st = h
	service.ps = Pserver{
		ProvisionServer: conf.ProvisionServer.ProvisionServer,
		CACert:          conf.ProvisionServer.CACert,
		TLSCert:         conf.ProvisionServer.TLSCert,
		TLSKey:          conf.ProvisionServer.TLSKey,
		UpdateSchedule:  conf.ProvisionServer.UpdateSchedule,
	}

	if strArray := strings.Split(conf.ApplicationServer.APIForGateway.OldGateway.Bind, ":"); len(strArray) != 2 {
		return errors.New(fmt.Sprintf("Invalid API Bind settings for OldGateway: %s", conf.ApplicationServer.APIForGateway.OldGateway.Bind))
	} else {
		service.bindPortOldGateway = strArray[1]
	}

	if strArray := strings.Split(conf.ApplicationServer.APIForGateway.NewGateway.Bind, ":"); len(strArray) != 2 {
		return errors.New(fmt.Sprintf("Invalid API Bind settings for NewGateway: %s", conf.ApplicationServer.APIForGateway.NewGateway.Bind))
	} else {
		service.bindPortNewGateway = strArray[1]
	}

	return service.UpdateFirmwareFromProvisioningServer(context.Background())
}

func (c *Controller) UpdateFirmwareFromProvisioningServer(ctx context.Context) error {
	log.WithFields(log.Fields{
		"provisioning-server": c.ps.ProvisionServer,
		"caCert":              c.ps.CACert,
		"tlsCert":             c.ps.TLSCert,
		"tlsKey":              c.ps.TLSKey,
		"schedule":            c.ps.UpdateSchedule,
	}).Info("Start schedule to update gateway firmware...")

	supernodeAddr := serverinfo.Service.SupernodeAddr

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
