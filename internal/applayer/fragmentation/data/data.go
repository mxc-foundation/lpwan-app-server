package data

import (
	"time"

	"github.com/brocaar/lorawan"

	mss "github.com/mxc-foundation/lpwan-app-server/internal/applayer/multicastsetup/data"
)

type FragmentationStruct struct {
	SyncInterval  time.Duration `mapstructure:"sync_interval"`
	SyncRetries   int           `mapstructure:"sync_retries"`
	SyncBatchSize int           `mapstructure:"sync_batch_size"`
}

// RemoteFragmentationSession defines a remote fragmentation session record.
type RemoteFragmentationSession struct {
	DevEUI              lorawan.EUI64                 `db:"dev_eui"`
	FragIndex           int                           `db:"frag_index"`
	CreatedAt           time.Time                     `db:"created_at"`
	UpdatedAt           time.Time                     `db:"updated_at"`
	MCGroupIDs          []int                         `db:"mc_group_ids"`
	NbFrag              int                           `db:"nb_frag"`
	FragSize            int                           `db:"frag_size"`
	FragmentationMatrix uint8                         `db:"fragmentation_matrix"`
	BlockAckDelay       int                           `db:"block_ack_delay"`
	Padding             int                           `db:"padding"`
	Descriptor          [4]byte                       `db:"descriptor"`
	State               mss.RemoteMulticastSetupState `db:"state"`
	StateProvisioned    bool                          `db:"state_provisioned"`
	RetryAfter          time.Time                     `db:"retry_after"`
	RetryCount          int                           `db:"retry_count"`
	RetryInterval       time.Duration                 `db:"retry_interval"`
}
