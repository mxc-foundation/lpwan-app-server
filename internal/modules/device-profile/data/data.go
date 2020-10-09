package data

import (
	"strings"
	"time"

	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/gofrs/uuid"
	"github.com/lib/pq/hstore"

	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
)

// DeviceProfileMeta defines the device-profile meta record.
type DeviceProfileMeta struct {
	DeviceProfileID   uuid.UUID `db:"device_profile_id"`
	NetworkServerID   int64     `db:"network_server_id"`
	OrganizationID    int64     `db:"organization_id"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
	Name              string    `db:"name"`
	NetworkServerName string    `db:"network_server_name"`
}

// DeviceProfile defines the device-profile.
type DeviceProfile struct {
	NetworkServerID      int64            `db:"network_server_id"`
	OrganizationID       int64            `db:"organization_id"`
	CreatedAt            time.Time        `db:"created_at"`
	UpdatedAt            time.Time        `db:"updated_at"`
	Name                 string           `db:"name"`
	PayloadCodec         string           `db:"payload_codec"`
	PayloadEncoderScript string           `db:"payload_encoder_script"`
	PayloadDecoderScript string           `db:"payload_decoder_script"`
	Tags                 hstore.Hstore    `db:"tags"`
	UplinkInterval       time.Duration    `db:"uplink_interval"`
	DeviceProfile        ns.DeviceProfile `db:"-"`
}

// Validate validates the device-profile data.
func (dp DeviceProfile) Validate() error {
	if strings.TrimSpace(dp.Name) == "" || len(dp.Name) > 100 {
		return errHandler.ErrDeviceProfileInvalidName
	}
	return nil
}

// DeviceProfileFilters provide filders for filtering device-profiles.
type DeviceProfileFilters struct {
	ApplicationID  int64 `db:"application_id"`
	OrganizationID int64 `db:"organization_id"`
	UserID         int64 `db:"user_id"`

	// Limit and Offset are added for convenience so that this struct can
	// be given as the arguments.
	Limit  int `db:"limit"`
	Offset int `db:"offset"`
}

// SQL returns the SQL filters.
func (f DeviceProfileFilters) SQL() string {
	var filters []string

	if f.ApplicationID != 0 {
		// Filter on organization_id too since dp > network-server > service-profile > application
		// join.
		filters = append(filters, "a.id = :application_id and dp.organization_id = a.organization_id")
	}

	if f.OrganizationID != 0 {
		filters = append(filters, "o.id = :organization_id")
	}

	if f.UserID != 0 {
		filters = append(filters, "u.id = :user_id")
	}

	if len(filters) == 0 {
		return ""
	}

	return "where " + strings.Join(filters, " and ")
}
