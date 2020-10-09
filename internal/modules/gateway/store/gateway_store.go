package store

import (
	"context"

	"github.com/brocaar/lorawan"

	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"
)

func NewStore(pgst pgstore.PgStore) *gs {
	return &gs{
		pgc: pgst,
		pg:  pgst,
	}
}

type gs struct {
	pgc pgstore.GatewayDefaultConfigPgStore
	pg  pgstore.GatewayPgStore
}

type GatewayStore interface {
	AddNewDefaultGatewayConfig(ctx context.Context, defaultConfig *DefaultGatewayConfig) error
	UpdateDefaultGatewayConfig(ctx context.Context, defaultConfig *DefaultGatewayConfig) error
	GetDefaultGatewayConfig(ctx context.Context, defaultConfig *DefaultGatewayConfig) error

	AddGatewayFirmware(ctx context.Context, gwFw *GatewayFirmware) (model string, err error)
	GetGatewayFirmware(ctx context.Context, model string, forUpdate bool) (gwFw GatewayFirmware, err error)
	GetGatewayFirmwareList(ctx context.Context) (list []GatewayFirmware, err error)
	UpdateGatewayFirmware(ctx context.Context, gwFw *GatewayFirmware) (model string, err error)
	UpdateGatewayConfigByGwId(ctx context.Context, config string, mac lorawan.EUI64) error
	CreateGateway(ctx context.Context, gw *Gateway) error
	UpdateGateway(ctx context.Context, gw *Gateway) error
	UpdateFirstHeartbeat(ctx context.Context, mac lorawan.EUI64, time int64) error
	UpdateLastHeartbeat(ctx context.Context, mac lorawan.EUI64, time int64) error
	SetAutoUpdateFirmware(ctx context.Context, mac lorawan.EUI64, autoUpdateFirmware bool) error
	DeleteGateway(ctx context.Context, mac lorawan.EUI64) error
	GetGateway(ctx context.Context, mac lorawan.EUI64, forUpdate bool) (Gateway, error)
	GetGatewayCount(ctx context.Context, filters GatewayFilters) (int, error)
	GetGateways(ctx context.Context, filters GatewayFilters) ([]GatewayListItem, error)
	GetGatewayConfigByGwId(ctx context.Context, mac lorawan.EUI64) (string, error)
	GetFirstHeartbeat(ctx context.Context, mac lorawan.EUI64) (int64, error)
	UpdateFirstHeartbeatToZero(ctx context.Context, mac lorawan.EUI64) error
	GetLastHeartbeat(ctx context.Context, mac lorawan.EUI64) (int64, error)
	GetGatewayMiningList(ctx context.Context, time, limit int64) ([]lorawan.EUI64, error)
	GetGatewaysLoc(ctx context.Context, limit int) ([]GatewayLocation, error)
	GetGatewaysForMACs(ctx context.Context, macs []lorawan.EUI64) (map[lorawan.EUI64]Gateway, error)
	CreateGatewayPing(ctx context.Context, ping *GatewayPing) error
	GetGatewayPing(ctx context.Context, id int64) (GatewayPing, error)
	CreateGatewayPingRX(ctx context.Context, rx *GatewayPingRX) error
	DeleteAllGatewaysForOrganizationID(ctx context.Context, organizationID int64) error
	GetAllGatewayMacList(ctx context.Context) ([]string, error)
	GetGatewayPingRXForPingID(ctx context.Context, pingID int64) ([]GatewayPingRX, error)
	GetLastGatewayPingAndRX(ctx context.Context, mac lorawan.EUI64) (GatewayPing, []GatewayPingRX, error)
	GetGatewaysActiveInactive(ctx context.Context, organizationID int64) (GatewaysActiveInactive, error)

	GetGatewayForPing(ctx context.Context) (*Gateway, error)

	// validator
	CheckCreateGatewayAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error)
	CheckListGatewayAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error)

	CheckReadGatewayAccess(ctx context.Context, username string, mac lorawan.EUI64, userID int64) (bool, error)
	CheckUpdateDeleteGatewayAccess(ctx context.Context, username string, mac lorawan.EUI64, userID int64) (bool, error)

	CheckReadOrganizationNetworkServerAccess(ctx context.Context, username string, organizationID, networkserverID, userID int64) (bool, error)
}

func (h *gs) GetGatewayForPing(ctx context.Context) (*Gateway, error) {
	return h.pg.GetGatewayForPing(ctx)
}
func (h *gs) GetGatewaysActiveInactive(ctx context.Context, organizationID int64) (GatewaysActiveInactive, error) {
	return h.pg.GetGatewaysActiveInactive(ctx, organizationID)
}

func (h *gs) AddNewDefaultGatewayConfig(ctx context.Context, defaultConfig *DefaultGatewayConfig) error {
	return h.pgc.AddNewDefaultGatewayConfig(ctx, defaultConfig)
}
func (h *gs) UpdateDefaultGatewayConfig(ctx context.Context, defaultConfig *DefaultGatewayConfig) error {
	return h.pgc.UpdateDefaultGatewayConfig(ctx, defaultConfig)
}
func (h *gs) GetDefaultGatewayConfig(ctx context.Context, defaultConfig *DefaultGatewayConfig) error {
	return h.pgc.GetDefaultGatewayConfig(ctx, defaultConfig)
}

func (h *gs) AddGatewayFirmware(ctx context.Context, gwFw *GatewayFirmware) (model string, err error) {
	return h.pg.AddGatewayFirmware(ctx, gwFw)
}
func (h *gs) GetGatewayFirmware(ctx context.Context, model string, forUpdate bool) (gwFw GatewayFirmware, err error) {
	return h.pg.GetGatewayFirmware(ctx, model, forUpdate)
}
func (h *gs) GetGatewayFirmwareList(ctx context.Context) (list []GatewayFirmware, err error) {
	return h.pg.GetGatewayFirmwareList(ctx)
}
func (h *gs) UpdateGatewayFirmware(ctx context.Context, gwFw *GatewayFirmware) (model string, err error) {
	return h.pg.UpdateGatewayFirmware(ctx, gwFw)
}
func (h *gs) UpdateGatewayConfigByGwId(ctx context.Context, config string, mac lorawan.EUI64) error {
	return h.pg.UpdateGatewayConfigByGwId(ctx, config, mac)
}
func (h *gs) CreateGateway(ctx context.Context, gw *Gateway) error {
	return h.pg.CreateGateway(ctx, gw)
}
func (h *gs) UpdateGateway(ctx context.Context, gw *Gateway) error {
	return h.pg.UpdateGateway(ctx, gw)
}
func (h *gs) UpdateFirstHeartbeat(ctx context.Context, mac lorawan.EUI64, time int64) error {
	return h.pg.UpdateFirstHeartbeat(ctx, mac, time)
}
func (h *gs) UpdateLastHeartbeat(ctx context.Context, mac lorawan.EUI64, time int64) error {
	return h.pg.UpdateLastHeartbeat(ctx, mac, time)
}
func (h *gs) SetAutoUpdateFirmware(ctx context.Context, mac lorawan.EUI64, autoUpdateFirmware bool) error {
	return h.pg.SetAutoUpdateFirmware(ctx, mac, autoUpdateFirmware)
}
func (h *gs) DeleteGateway(ctx context.Context, mac lorawan.EUI64) error {
	return h.pg.DeleteGateway(ctx, mac)
}
func (h *gs) GetGateway(ctx context.Context, mac lorawan.EUI64, forUpdate bool) (Gateway, error) {
	return h.pg.GetGateway(ctx, mac, forUpdate)
}
func (h *gs) GetGatewayCount(ctx context.Context, filters GatewayFilters) (int, error) {
	return h.pg.GetGatewayCount(ctx, filters)
}
func (h *gs) GetGateways(ctx context.Context, filters GatewayFilters) ([]GatewayListItem, error) {
	return h.pg.GetGateways(ctx, filters)
}
func (h *gs) GetGatewayConfigByGwId(ctx context.Context, mac lorawan.EUI64) (string, error) {
	return h.pg.GetGatewayConfigByGwId(ctx, mac)
}
func (h *gs) GetFirstHeartbeat(ctx context.Context, mac lorawan.EUI64) (int64, error) {
	return h.pg.GetFirstHeartbeat(ctx, mac)
}
func (h *gs) UpdateFirstHeartbeatToZero(ctx context.Context, mac lorawan.EUI64) error {
	return h.pg.UpdateFirstHeartbeatToZero(ctx, mac)
}
func (h *gs) GetLastHeartbeat(ctx context.Context, mac lorawan.EUI64) (int64, error) {
	return h.pg.GetLastHeartbeat(ctx, mac)
}
func (h *gs) GetGatewayMiningList(ctx context.Context, time, limit int64) ([]lorawan.EUI64, error) {
	return h.pg.GetGatewayMiningList(ctx, time, limit)
}
func (h *gs) GetGatewaysLoc(ctx context.Context, limit int) ([]GatewayLocation, error) {
	return h.pg.GetGatewaysLoc(ctx, limit)
}
func (h *gs) GetGatewaysForMACs(ctx context.Context, macs []lorawan.EUI64) (map[lorawan.EUI64]Gateway, error) {
	return h.pg.GetGatewaysForMACs(ctx, macs)
}
func (h *gs) CreateGatewayPing(ctx context.Context, ping *GatewayPing) error {
	return h.pg.CreateGatewayPing(ctx, ping)
}
func (h *gs) GetGatewayPing(ctx context.Context, id int64) (GatewayPing, error) {
	return h.pg.GetGatewayPing(ctx, id)
}
func (h *gs) CreateGatewayPingRX(ctx context.Context, rx *GatewayPingRX) error {
	return h.pg.CreateGatewayPingRX(ctx, rx)
}
func (h *gs) DeleteAllGatewaysForOrganizationID(ctx context.Context, organizationID int64) error {
	return h.pg.DeleteAllGatewaysForOrganizationID(ctx, organizationID)
}
func (h *gs) GetAllGatewayMacList(ctx context.Context) ([]string, error) {
	return h.pg.GetAllGatewayMacList(ctx)
}
func (h *gs) GetGatewayPingRXForPingID(ctx context.Context, pingID int64) ([]GatewayPingRX, error) {
	return h.pg.GetGatewayPingRXForPingID(ctx, pingID)
}
func (h *gs) GetLastGatewayPingAndRX(ctx context.Context, mac lorawan.EUI64) (GatewayPing, []GatewayPingRX, error) {
	return h.pg.GetLastGatewayPingAndRX(ctx, mac)
}

// validator
func (h *gs) CheckCreateGatewayAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error) {
	return h.pg.CheckCreateGatewayAccess(ctx, username, organizationID, userID)
}
func (h *gs) CheckListGatewayAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error) {
	return h.pg.CheckListGatewayAccess(ctx, username, organizationID, userID)
}

func (h *gs) CheckReadGatewayAccess(ctx context.Context, username string, mac lorawan.EUI64, userID int64) (bool, error) {
	return h.pg.CheckReadGatewayAccess(ctx, username, mac, userID)
}
func (h *gs) CheckUpdateDeleteGatewayAccess(ctx context.Context, username string, mac lorawan.EUI64, userID int64) (bool, error) {
	return h.pg.CheckUpdateDeleteGatewayAccess(ctx, username, mac, userID)
}

func (h *gs) CheckReadOrganizationNetworkServerAccess(ctx context.Context, username string, organizationID, networkserverID, userID int64) (bool, error) {
	return h.pg.CheckReadOrganizationNetworkServerAccess(ctx, username, organizationID, networkserverID, userID)
}
