package data

import (
	"strings"
	"time"

	"github.com/brocaar/lorawan"
	"github.com/gofrs/uuid"

	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
)

// FUOTADeploymentState defines the fuota deployment state.
type FUOTADeploymentState string

// FUOTA deployment states.
const (
	FUOTADeploymentMulticastCreate        FUOTADeploymentState = "MC_CREATE"
	FUOTADeploymentMulticastSetup         FUOTADeploymentState = "MC_SETUP"
	FUOTADeploymentFragmentationSessSetup FUOTADeploymentState = "FRAG_SESS_SETUP"
	FUOTADeploymentMulticastSessCSetup    FUOTADeploymentState = "MC_SESS_C_SETUP"
	FUOTADeploymentEnqueue                FUOTADeploymentState = "ENQUEUE"
	FUOTADeploymentStatusRequest          FUOTADeploymentState = "STATUS_REQUEST"
	FUOTADeploymentSetDeviceStatus        FUOTADeploymentState = "SET_DEVICE_STATUS"
	FUOTADeploymentCleanup                FUOTADeploymentState = "CLEANUP"
	FUOTADeploymentDone                   FUOTADeploymentState = "DONE"
)

// FUOTADeploymentDeviceState defines the fuota deployment device state.
type FUOTADeploymentDeviceState string

// FUOTA deployment device states.
const (
	FUOTADeploymentDevicePending FUOTADeploymentDeviceState = "PENDING"
	FUOTADeploymentDeviceSuccess FUOTADeploymentDeviceState = "SUCCESS"
	FUOTADeploymentDeviceError   FUOTADeploymentDeviceState = "ERROR"
)

// FUOTADeploymentGroupType defines the group-type.
type FUOTADeploymentGroupType string

// FUOTA deployment group types.
const (
	FUOTADeploymentGroupTypeB FUOTADeploymentGroupType = "B"
	FUOTADeploymentGroupTypeC FUOTADeploymentGroupType = "C"
)

// FUOTADeployment defiles a firmware update over the air deployment.
type FUOTADeployment struct {
	ID                  uuid.UUID                `db:"id"`
	CreatedAt           time.Time                `db:"created_at"`
	UpdatedAt           time.Time                `db:"updated_at"`
	Name                string                   `db:"name"`
	MulticastGroupID    *uuid.UUID               `db:"multicast_group_id"`
	GroupType           FUOTADeploymentGroupType `db:"group_type"`
	DR                  int                      `db:"dr"`
	Frequency           int                      `db:"frequency"`
	PingSlotPeriod      int                      `db:"ping_slot_period"`
	FragmentationMatrix uint8                    `db:"fragmentation_matrix"`
	Descriptor          [4]byte                  `db:"descriptor"`
	Payload             []byte                   `db:"payload"`
	FragSize            int                      `db:"frag_size"`
	Redundancy          int                      `db:"redundancy"`
	BlockAckDelay       int                      `db:"block_ack_delay"`
	MulticastTimeout    int                      `db:"multicast_timeout"`
	State               FUOTADeploymentState     `db:"state"`
	UnicastTimeout      time.Duration            `db:"unicast_timeout"`
	NextStepAfter       time.Time                `db:"next_step_after"`
}

// FUOTADeploymentListItem defines a FUOTA deployment item for listing.
type FUOTADeploymentListItem struct {
	ID            uuid.UUID            `db:"id"`
	CreatedAt     time.Time            `db:"created_at"`
	UpdatedAt     time.Time            `db:"updated_at"`
	Name          string               `db:"name"`
	State         FUOTADeploymentState `db:"state"`
	NextStepAfter time.Time            `db:"next_step_after"`
}

// FUOTADeploymentDevice defines the device record of a FUOTA deployment.
type FUOTADeploymentDevice struct {
	FUOTADeploymentID uuid.UUID                  `db:"fuota_deployment_id"`
	DevEUI            lorawan.EUI64              `db:"dev_eui"`
	CreatedAt         time.Time                  `db:"created_at"`
	UpdatedAt         time.Time                  `db:"updated_at"`
	State             FUOTADeploymentDeviceState `db:"state"`
	ErrorMessage      string                     `db:"error_message"`
}

// FUOTADeploymentDeviceListItem defines the Device as FUOTA deployment list item.
type FUOTADeploymentDeviceListItem struct {
	CreatedAt         time.Time                  `db:"created_at"`
	UpdatedAt         time.Time                  `db:"updated_at"`
	FUOTADeploymentID uuid.UUID                  `db:"fuota_deployment_id"`
	DevEUI            lorawan.EUI64              `db:"dev_eui"`
	DeviceName        string                     `db:"device_name"`
	State             FUOTADeploymentDeviceState `db:"state"`
	ErrorMessage      string                     `db:"error_message"`
}

// FUOTADeploymentFilters provides filters that can be used to filter on
// FUOTA deployments. Note that empty values are not used as filters.
type FUOTADeploymentFilters struct {
	DevEUI        lorawan.EUI64 `db:"dev_eui"`
	ApplicationID int64         `db:"application_id"`

	// Limit and Offset are added for convenience so that this struct can
	// be given as the arguments.
	Limit  int `db:"limit"`
	Offset int `db:"offset"`
}

// SQL returns the SQL filter.
func (f FUOTADeploymentFilters) SQL() string {
	var filters []string
	var nullDevEUI lorawan.EUI64

	if f.DevEUI != nullDevEUI {
		filters = append(filters, "fdd.dev_eui = :dev_eui")
	}

	if f.ApplicationID != 0 {
		filters = append(filters, "d.application_id = :application_id")
	}

	if len(filters) == 0 {
		return ""
	}

	return "where " + strings.Join(filters, " and ")
}

// Validate validates the FUOTADeployment data.
func (fd FUOTADeployment) Validate() error {
	if strings.TrimSpace(fd.Name) == "" {
		return errHandler.ErrFUOTADeploymentInvalidName
	}
	if len(fd.Payload) <= 0 || fd.Payload == nil {
		return errHandler.ErrFUOTADeploymentNullPayload
	}
	return nil
}
