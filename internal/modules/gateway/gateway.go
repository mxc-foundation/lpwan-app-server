package gateway

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"

	"github.com/brocaar/lorawan"

	psPb "github.com/mxc-foundation/lpwan-app-server/api/ps-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/provisionserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
	"github.com/mxc-foundation/lpwan-app-server/internal/types"
)

type GatewayStore interface {
	AddNewDefaultGatewayConfig(ctx context.Context, defaultConfig *DefaultGatewayConfig) error
	UpdateDefaultGatewayConfig(ctx context.Context, defaultConfig *DefaultGatewayConfig) error
	GetDefaultGatewayConfig(ctx context.Context, defaultConfig *DefaultGatewayConfig) error

	AddGatewayFirmware(ctx context.Context, gwFw *GatewayFirmware) (model string, err error)
	GetGatewayFirmware(ctx context.Context, model string, forUpdate bool) (gwFw GatewayFirmware, err error)
	GetGatewayFirmwareList(ctx context.Context, ) (list []GatewayFirmware, err error)
	UpdateGatewayFirmware(ctx context.Context, gwFw *GatewayFirmware) (model string, err error)
	UpdateGatewayConfigByGwId(ctx context.Context, config string, mac lorawan.EUI64) error
	CreateGateway(ctx context.Context, gw *Gateway) error
	UpdateGateway(ctx context.Context, gw *Gateway) error
	UpdateFirstHeartbeat(ctx context.Context, mac lorawan.EUI64, time int64) error
	UpdateLastHeartbeat(ctx context.Context, mac lorawan.EUI64, time int64) error
	SetAutoUpdateFirmware(ctx context.Context, mac lorawan.EUI64, autoUpdateFirmware bool) error
	DeleteGateway(ctx context.Context, mac lorawan.EUI64) error
	GetGateway(ctx context.Context, mac lorawan.EUI64, forUpdate bool) (Gateway, error)
	GetGatewayCount(ctx context.Context, search string) (int, error)
	GetGateways(ctx context.Context, limit, offset int32, search string) ([]Gateway, error)
	GetGatewayConfigByGwId(ctx context.Context, mac lorawan.EUI64) (string, error)
	GetFirstHeartbeat(ctx context.Context, mac lorawan.EUI64) (int64, error)
	UpdateFirstHeartbeatToZero(ctx context.Context, mac lorawan.EUI64) error
	GetLastHeartbeat(ctx context.Context, mac lorawan.EUI64) (int64, error)
	GetGatewayMiningList(ctx context.Context, time, limit int64) ([]lorawan.EUI64, error)
	GetGatewaysLoc(ctx context.Context, limit int) ([]GatewayLocation, error)
	GetGatewaysForMACs(ctx context.Context, macs []lorawan.EUI64) (map[lorawan.EUI64]Gateway, error)
	GetGatewayCountForOrganizationID(ctx context.Context, organizationID int64, search string) (int, error)
	GetGatewaysForOrganizationID(ctx context.Context, organizationID int64, limit, offset int, search string) ([]Gateway, error)
	GetGatewayCountForUser(ctx context.Context, username string, search string) (int, error)
	GetGatewaysForUser(ctx context.Context, username string, limit, offset int, search string) ([]Gateway, error)
	CreateGatewayPing(ctx context.Context, ping *GatewayPing) error
	GetGatewayPing(ctx context.Context, id int64) (GatewayPing, error)
	CreateGatewayPingRX(ctx context.Context, rx *GatewayPingRX) error
	DeleteAllGatewaysForOrganizationID(ctx context.Context, organizationID int64) error
	GetAllGatewayMacList(ctx context.Context) ([]string, error)
	GetGatewayPingRXForPingID(ctx context.Context, pingID int64) ([]GatewayPingRX, error)
	GetLastGatewayPingAndRX(ctx context.Context, mac lorawan.EUI64) (GatewayPing, []GatewayPingRX, error)

	// validator
	CheckCreateGatewayAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error)
	CheckListGatewayAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error)

	CheckReadGatewayAccess(ctx context.Context, username string, mac lorawan.EUI64, userID int64) (bool, error)
	CheckUpdateDeleteGatewayAccess(ctx context.Context, username string, mac lorawan.EUI64, userID int64) (bool, error)

	CheckReadOrganizationNetworkServerAccess(ctx context.Context, username string, organizationID, networkserverID, userID int64) (bool, error)
}

type Controller struct {
	St            GatewayStore
	Validator     Validator
	SupernodeAddr string
}

var Service *Controller

func Setup() error {
	Service.St = store.New(storage.DB().DB)
	return nil
}

func (c *Controller) UpdateFirmwareFromProvisioningServer(conf config.Config) error {
	log.WithFields(log.Fields{
		"provisioning-server": conf.ProvisionServer.ProvisionServer,
		"caCert":              conf.ProvisionServer.CACert,
		"tlsCert":             conf.ProvisionServer.TLSCert,
		"tlsKey":              conf.ProvisionServer.TLSKey,
		"schedule":            conf.ProvisionServer.UpdateSchedule,
	}).Info("Start schedule to update gateway firmware...")
	Service.SupernodeAddr = os.Getenv("APPSERVER")
	if strings.HasPrefix(Service.SupernodeAddr, "https://") {
		Service.SupernodeAddr = strings.Replace(Service.SupernodeAddr, "https://", "", -1)
	}
	if strings.HasPrefix(Service.SupernodeAddr, "http://") {
		Service.SupernodeAddr = strings.Replace(Service.SupernodeAddr, "http://", "", -1)
	}
	if strings.HasSuffix(Service.SupernodeAddr, ":8080") {
		Service.SupernodeAddr = strings.Replace(Service.SupernodeAddr, ":8080", "", -1)
	}
	Service.SupernodeAddr = strings.Replace(Service.SupernodeAddr, "/", "", -1)

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
		gwFwList, err := c.St.GetGatewayFirmwareList()
		if err != nil {
			log.WithError(err).Errorf("Failed to get gateway firmware list.")
			return
		}

		// send update
		psClient, err := provisionserver.CreateClientWithCert(conf.ProvisionServer.ProvisionServer,
			conf.ProvisionServer.CACert,
			conf.ProvisionServer.TLSCert,
			conf.ProvisionServer.TLSKey)
		if err != nil {
			log.WithError(err).Errorf("Create Provisioning server client error")
			return
		}

		for _, v := range gwFwList {
			res, err := psClient.GetUpdate(context.Background(), &psPb.GetUpdateRequest{
				Model:          v.Model,
				SuperNodeAddr:  Service.SupernodeAddr,
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

			gatewayFw := GatewayFirmware{
				Model:        v.Model,
				ResourceLink: res.ResourceLink,
				FirmwareHash: md5sum,
			}

			model, _ := c.St.UpdateGatewayFirmware(&gatewayFw)
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
