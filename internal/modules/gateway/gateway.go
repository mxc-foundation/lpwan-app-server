package gateway

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"

	"github.com/pkg/errors"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"

	psPb "github.com/mxc-foundation/lpwan-app-server/api/ps-serves-appserver"
	pscli "github.com/mxc-foundation/lpwan-app-server/internal/clients/psconn"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/types"
)

type Controller struct {
	St            *store.Handler
	Validator     Validator
	SupernodeAddr string
}

var Service = &Controller{}

func Setup(s store.Store) error {
	Service.St, _ = store.New(s)
	return nil
}

func (c *Controller) UpdateFirmwareFromProvisioningServer(ctx context.Context, conf config.Config) error {
	log.WithFields(log.Fields{
		"provisioning-server": conf.ProvisionServer.ProvisionServer,
		"caCert":              conf.ProvisionServer.CACert,
		"tlsCert":             conf.ProvisionServer.TLSCert,
		"tlsKey":              conf.ProvisionServer.TLSKey,
		"schedule":            conf.ProvisionServer.UpdateSchedule,
	}).Info("Start schedule to update gateway firmware...")
	c.SupernodeAddr = os.Getenv("APPSERVER")
	if strings.HasPrefix(c.SupernodeAddr, "https://") {
		c.SupernodeAddr = strings.Replace(c.SupernodeAddr, "https://", "", -1)
	}
	if strings.HasPrefix(c.SupernodeAddr, "http://") {
		c.SupernodeAddr = strings.Replace(c.SupernodeAddr, "http://", "", -1)
	}
	if strings.HasSuffix(c.SupernodeAddr, ":8080") {
		c.SupernodeAddr = strings.Replace(c.SupernodeAddr, ":8080", "", -1)
	}
	c.SupernodeAddr = strings.Replace(c.SupernodeAddr, "/", "", -1)

	var bindPortOldGateway string
	var bindPortNewGateway string

	if strArray := strings.Split(conf.ApplicationServer.APIForGateway.OldGateway.Bind, ":"); len(strArray) != 2 {
		return errors.New(fmt.Sprintf("Invalid API Bind settings for OldGateway: %s", conf.ApplicationServer.APIForGateway.OldGateway.Bind))
	} else {
		bindPortOldGateway = strArray[1]
	}

	if strArray := strings.Split(conf.ApplicationServer.APIForGateway.NewGateway.Bind, ":"); len(strArray) != 2 {
		return errors.New(fmt.Sprintf("Invalid API Bind settings for NewGateway: %s", conf.ApplicationServer.APIForGateway.NewGateway.Bind))
	} else {
		bindPortNewGateway = strArray[1]
	}

	cron := cron.New()
	err := cron.AddFunc(conf.ProvisionServer.UpdateSchedule, func() {
		log.Info("Check firmware update...")
		gwFwList, err := c.St.GetGatewayFirmwareList(ctx)
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
				SuperNodeAddr:  c.SupernodeAddr,
				PortOldGateway: bindPortOldGateway,
				PortNewGateway: bindPortNewGateway,
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

			model, _ := c.St.UpdateGatewayFirmware(ctx, &gatewayFw)
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
