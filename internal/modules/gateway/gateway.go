package gateway

import (
	"context"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"

	"github.com/brocaar/lorawan"
)

type GatewayStore interface {
	AddNewDefaultGatewayConfig(defaultConfig *DefaultGatewayConfig) error
	UpdateDefaultGatewayConfig(defaultConfig *DefaultGatewayConfig) error
	GetDefaultGatewayConfig(defaultConfig *DefaultGatewayConfig) error

	AddGatewayFirmware(gwFw *GatewayFirmware) (model string, err error)
	GetGatewayFirmware(model string, forUpdate bool) (gwFw GatewayFirmware, err error)
	GetGatewayFirmwareList() (list []GatewayFirmware, err error)
	UpdateGatewayFirmware(gwFw *GatewayFirmware) (model string, err error)
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
	CheckCreateGatewayAccess(username string, organizationID, userID int64) (bool, error)
	CheckListGatewayAccess(username string, organizationID, userID int64) (bool, error)

	CheckReadGatewayAccess(username string, mac lorawan.EUI64, userID int64) (bool, error)
	CheckUpdateDeleteGatewayAccess(username string, mac lorawan.EUI64, userID int64) (bool, error)

	CheckReadOrganizationNetworkServerAccess(username string, organizationID, networkserverID, userID int64) (bool, error)
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
