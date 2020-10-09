package data

import (
	"time"

	"github.com/brocaar/lorawan"
	"github.com/gofrs/uuid"
)

type MulticastStruct struct {
	SyncInterval  time.Duration `mapstructure:"sync_interval"`
	SyncRetries   int           `mapstructure:"sync_retries"`
	SyncBatchSize int           `mapstructure:"sync_batch_size"`
}

// RemoteMulticastSetupState defines the state type.
type RemoteMulticastSetupState string

// Possible states
const (
	RemoteMulticastSetupSetup  RemoteMulticastSetupState = "SETUP"
	RemoteMulticastSetupDelete RemoteMulticastSetupState = "DELETE"
)

// RemoteMulticastSetup defines a remote multicast-setup record.
type RemoteMulticastSetup struct {
	DevEUI           lorawan.EUI64             `db:"dev_eui"`
	MulticastGroupID uuid.UUID                 `db:"multicast_group_id"`
	CreatedAt        time.Time                 `db:"created_at"`
	UpdatedAt        time.Time                 `db:"updated_at"`
	McGroupID        int                       `db:"mc_group_id"`
	McAddr           lorawan.DevAddr           `db:"mc_addr"`
	McKeyEncrypted   lorawan.AES128Key         `db:"mc_key_encrypted"`
	MinMcFCnt        uint32                    `db:"min_mc_f_cnt"`
	MaxMcFCnt        uint32                    `db:"max_mc_f_cnt"`
	State            RemoteMulticastSetupState `db:"state"`
	StateProvisioned bool                      `db:"state_provisioned"`
	RetryInterval    time.Duration             `db:"retry_interval"`
	RetryAfter       time.Time                 `db:"retry_after"`
	RetryCount       int                       `db:"retry_count"`
}

// RemoteMulticastClassCSession defines a remote multicast-setup Class-C session record.
type RemoteMulticastClassCSession struct {
	DevEUI           lorawan.EUI64 `db:"dev_eui"`
	MulticastGroupID uuid.UUID     `db:"multicast_group_id"`
	CreatedAt        time.Time     `db:"created_at"`
	UpdatedAt        time.Time     `db:"updated_at"`
	McGroupID        int           `db:"mc_group_id"`
	SessionTime      time.Time     `db:"session_time"`
	SessionTimeOut   int           `db:"session_time_out"`
	DLFrequency      int           `db:"dl_frequency"`
	DR               int           `db:"dr"`
	StateProvisioned bool          `db:"state_provisioned"`
	RetryAfter       time.Time     `db:"retry_after"`
	RetryCount       int           `db:"retry_count"`
	RetryInterval    time.Duration `db:"retry_interval"`
}
