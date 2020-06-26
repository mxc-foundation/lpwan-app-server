package gateway

import (
	"context"
	"github.com/brocaar/lorawan"

	"github.com/jmoiron/sqlx"

	pg "github.com/mxc-foundation/lpwan-app-server/internal/storage/postgresql"
)

type GatewayTable interface {
	AddGatewayFirmware(db sqlx.Queryer, gwFw *GatewayFirmware) (model string, err error)
	GetGatewayFirmware(db sqlx.Queryer, model string, forUpdate bool) (gwFw GatewayFirmware, err error)
	GetGatewayFirmwareList(db sqlx.Queryer) (list []GatewayFirmware, err error)
	UpdateGatewayFirmware(db sqlx.Queryer, gwFw *GatewayFirmware) (model string, err error)
	UpdateGatewayConfigByGwId(ctx context.Context, db sqlx.Ext, config string, mac lorawan.EUI64) error
	CreateGateway(ctx context.Context, db sqlx.Ext, gw *Gateway) error
	UpdateGateway(ctx context.Context, db sqlx.Ext, gw *Gateway) error
	UpdateFirstHeartbeat(ctx context.Context, db sqlx.Ext, mac lorawan.EUI64, time int64) error
	UpdateLastHeartbeat(ctx context.Context, db sqlx.Ext, mac lorawan.EUI64, time int64) error
	SetAutoUpdateFirmware(ctx context.Context, db sqlx.Ext, mac lorawan.EUI64, autoUpdateFirmware bool) error
	DeleteGateway(ctx context.Context, db sqlx.Ext, mac lorawan.EUI64) error
	GetGateway(ctx context.Context, db sqlx.Queryer, mac lorawan.EUI64, forUpdate bool) (Gateway, error)
	GetGatewayCount(ctx context.Context, db sqlx.Queryer, search string) (int, error)
	GetGateways(ctx context.Context, db sqlx.Queryer, limit, offset int, search string) ([]Gateway, error)
	GetGatewayConfigByGwId(ctx context.Context, db sqlx.Queryer, mac lorawan.EUI64) (string, error)
	GetFirstHeartbeat(ctx context.Context, db sqlx.Queryer, mac lorawan.EUI64) (int64, error)
	UpdateFirstHeartbeatToZero(ctx context.Context, db sqlx.Execer, mac lorawan.EUI64) error
	GetLastHeartbeat(ctx context.Context, db sqlx.Queryer, mac lorawan.EUI64) (int64, error)
	GetGatewayMiningList(ctx context.Context, db sqlx.Queryer, time, limit int64) ([]lorawan.EUI64, error)
	GetGatewaysLoc(ctx context.Context, db sqlx.Queryer, limit int) ([]GatewayLocation, error)
	GetGatewaysForMACs(ctx context.Context, db sqlx.Queryer, macs []lorawan.EUI64) (map[lorawan.EUI64]Gateway, error)
	GetGatewayCountForOrganizationID(ctx context.Context, db sqlx.Queryer, organizationID int64, search string) (int, error)
	GetGatewaysForOrganizationID(ctx context.Context, db sqlx.Queryer, organizationID int64, limit, offset int, search string) ([]Gateway, error)
	GetGatewayCountForUser(ctx context.Context, db sqlx.Queryer, username string, search string) (int, error)
	GetGatewaysForUser(ctx context.Context, db sqlx.Queryer, username string, limit, offset int, search string) ([]Gateway, error)
	CreateGatewayPing(ctx context.Context, db sqlx.Queryer, ping *GatewayPing) error
	GetGatewayPing(ctx context.Context, db sqlx.Queryer, id int64) (GatewayPing, error)
	CreateGatewayPingRX(ctx context.Context, db sqlx.Queryer, rx *GatewayPingRX) error
	DeleteAllGatewaysForOrganizationID(ctx context.Context, db sqlx.Ext, organizationID int64) error
	GetAllGatewayMacList(ctx context.Context, db sqlx.Ext) ([]string, error)
	GetGatewayPingRXForPingID(ctx context.Context, db sqlx.Queryer, pingID int64) ([]GatewayPingRX, error)
	GetLastGatewayPingAndRX(ctx context.Context, db sqlx.Queryer, mac lorawan.EUI64) (GatewayPing, []GatewayPingRX, error)
}

var GatewayDB = GatewayTable(&pg.GatewayTable)
