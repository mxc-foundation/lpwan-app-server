package data

import "github.com/brocaar/lorawan"

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

type GatewayMining struct {
	GatewayMac lorawan.EUI64 `db:"mac"`
	OwnerOrgID int64         `db:"organization_id"`
	StcOrgID   *int64        `db:"stc_org_id"`
}
