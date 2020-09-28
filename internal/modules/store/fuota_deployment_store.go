package store

import (
	"context"
	"strings"
	"time"

	"github.com/brocaar/lorawan"

	"github.com/gofrs/uuid"
)

type FUOTADeploymentStore interface {
	GetDeviceKeysFromFuotaDevelopmentDevice(ctx context.Context, id uuid.UUID) ([]DeviceKeys, error)
	CreateFUOTADeploymentForDevice(ctx context.Context, fd *FUOTADeployment, devEUI lorawan.EUI64) error
	GetFUOTADeployment(ctx context.Context, id uuid.UUID, forUpdate bool) (FUOTADeployment, error)
	GetPendingFUOTADeployments(ctx context.Context, batchSize int) ([]FUOTADeployment, error)
	UpdateFUOTADeployment(ctx context.Context, fd *FUOTADeployment) error
	GetFUOTADeploymentCount(ctx context.Context, filters FUOTADeploymentFilters) (int, error)
	GetFUOTADeployments(ctx context.Context, filters FUOTADeploymentFilters) ([]FUOTADeploymentListItem, error)
	GetFUOTADeploymentDevice(ctx context.Context, fuotaDeploymentID uuid.UUID, devEUI lorawan.EUI64) (FUOTADeploymentDevice, error)
	GetPendingFUOTADeploymentDevice(ctx context.Context, devEUI lorawan.EUI64) (FUOTADeploymentDevice, error)
	UpdateFUOTADeploymentDevice(ctx context.Context, fdd *FUOTADeploymentDevice) error
	GetFUOTADeploymentDeviceCount(ctx context.Context, fuotaDeploymentID uuid.UUID) (int, error)
	GetFUOTADeploymentDevices(ctx context.Context, fuotaDeploymentID uuid.UUID, limit, offset int) ([]FUOTADeploymentDeviceListItem, error)
	GetServiceProfileIDForFUOTADeployment(ctx context.Context, fuotaDeploymentID uuid.UUID) (uuid.UUID, error)
	SetFromRemoteMulticastSetup(ctx context.Context, fuotaDevelopmentID, multicastGroupID uuid.UUID) error
	SetFromRemoteFragmentationSession(ctx context.Context, fuotaDevelopmentID uuid.UUID, fragIdx int) error
	SetIncompleteFuotaDevelopment(ctx context.Context, fuotaDevelopmentID uuid.UUID) error

	// validator
	CheckReadFUOTADeploymentAccess(ctx context.Context, username string, id uuid.UUID, userID int64) (bool, error)
	CheckCreateFUOTADeploymentsAccess(ctx context.Context, username string, applicationID int64, devEUI lorawan.EUI64, userID int64) (bool, error)
}

func (h *Handler) SetIncompleteFuotaDevelopment(ctx context.Context, fuotaDevelopmentID uuid.UUID) error {
	return h.store.SetIncompleteFuotaDevelopment(ctx, fuotaDevelopmentID)
}

func (h *Handler) SetFromRemoteFragmentationSession(ctx context.Context, fuotaDevelopmentID uuid.UUID, fragIdx int) error {
	return h.store.SetFromRemoteFragmentationSession(ctx, fuotaDevelopmentID, fragIdx)
}
func (h *Handler) SetFromRemoteMulticastSetup(ctx context.Context, fuotaDevelopmentID, multicastGroupID uuid.UUID) error {
	return h.store.SetFromRemoteMulticastSetup(ctx, fuotaDevelopmentID, multicastGroupID)
}
func (h *Handler) CreateFUOTADeploymentForDevice(ctx context.Context, fd *FUOTADeployment, devEUI lorawan.EUI64) error {
	return h.store.CreateFUOTADeploymentForDevice(ctx, fd, devEUI)
}
func (h *Handler) GetFUOTADeployment(ctx context.Context, id uuid.UUID, forUpdate bool) (FUOTADeployment, error) {
	return h.store.GetFUOTADeployment(ctx, id, forUpdate)
}
func (h *Handler) GetPendingFUOTADeployments(ctx context.Context, batchSize int) ([]FUOTADeployment, error) {
	return h.store.GetPendingFUOTADeployments(ctx, batchSize)
}
func (h *Handler) UpdateFUOTADeployment(ctx context.Context, fd *FUOTADeployment) error {
	return h.store.UpdateFUOTADeployment(ctx, fd)
}
func (h *Handler) GetFUOTADeploymentCount(ctx context.Context, filters FUOTADeploymentFilters) (int, error) {
	return h.store.GetFUOTADeploymentCount(ctx, filters)
}
func (h *Handler) GetFUOTADeployments(ctx context.Context, filters FUOTADeploymentFilters) ([]FUOTADeploymentListItem, error) {
	return h.store.GetFUOTADeployments(ctx, filters)
}
func (h *Handler) GetFUOTADeploymentDevice(ctx context.Context, fuotaDeploymentID uuid.UUID, devEUI lorawan.EUI64) (FUOTADeploymentDevice, error) {
	return h.store.GetFUOTADeploymentDevice(ctx, fuotaDeploymentID, devEUI)
}
func (h *Handler) GetPendingFUOTADeploymentDevice(ctx context.Context, devEUI lorawan.EUI64) (FUOTADeploymentDevice, error) {
	return h.store.GetPendingFUOTADeploymentDevice(ctx, devEUI)
}
func (h *Handler) UpdateFUOTADeploymentDevice(ctx context.Context, fdd *FUOTADeploymentDevice) error {
	return h.store.UpdateFUOTADeploymentDevice(ctx, fdd)
}
func (h *Handler) GetFUOTADeploymentDeviceCount(ctx context.Context, fuotaDeploymentID uuid.UUID) (int, error) {
	return h.store.GetFUOTADeploymentDeviceCount(ctx, fuotaDeploymentID)
}
func (h *Handler) GetFUOTADeploymentDevices(ctx context.Context, fuotaDeploymentID uuid.UUID, limit, offset int) ([]FUOTADeploymentDeviceListItem, error) {
	return h.store.GetFUOTADeploymentDevices(ctx, fuotaDeploymentID, limit, offset)
}
func (h *Handler) GetServiceProfileIDForFUOTADeployment(ctx context.Context, fuotaDeploymentID uuid.UUID) (uuid.UUID, error) {
	return h.store.GetServiceProfileIDForFUOTADeployment(ctx, fuotaDeploymentID)
}
func (h *Handler) GetDeviceKeysFromFuotaDevelopmentDevice(ctx context.Context, id uuid.UUID) ([]DeviceKeys, error) {
	return h.store.GetDeviceKeysFromFuotaDevelopmentDevice(ctx, id)
}
func (h *Handler) CheckReadFUOTADeploymentAccess(ctx context.Context, username string, id uuid.UUID, userID int64) (bool, error) {
	return h.store.CheckReadFUOTADeploymentAccess(ctx, username, id, userID)
}

func (h *Handler) CheckCreateFUOTADeploymentsAccess(ctx context.Context, username string, applicationID int64, devEUI lorawan.EUI64, userID int64) (bool, error) {
	return h.store.CheckCreateFUOTADeploymentsAccess(ctx, username, applicationID, devEUI, userID)
}

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
		return ErrFUOTADeploymentInvalidName
	}
	if len(fd.Payload) <= 0 || fd.Payload == nil {
		return ErrFUOTADeploymentNullPayload
	}
	return nil
}
