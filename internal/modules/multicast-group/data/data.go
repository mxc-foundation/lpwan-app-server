package data

import (
	"strings"
	"time"

	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/brocaar/lorawan"
	"github.com/gofrs/uuid"

	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
)

// MulticastGroup defines the multicast-group.
type MulticastGroup struct {
	CreatedAt        time.Time         `db:"created_at"`
	UpdatedAt        time.Time         `db:"updated_at"`
	Name             string            `db:"name"`
	MCAppSKey        lorawan.AES128Key `db:"mc_app_s_key"`
	MCKey            lorawan.AES128Key `db:"mc_key"`
	ServiceProfileID uuid.UUID         `db:"service_profile_id"`
	MulticastGroup   ns.MulticastGroup `db:"-"`
}

// MulticastGroupListItem defines the multicast-group for listing.
type MulticastGroupListItem struct {
	ID                 uuid.UUID `db:"id"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
	Name               string    `db:"name"`
	ServiceProfileID   uuid.UUID `db:"service_profile_id"`
	ServiceProfileName string    `db:"service_profile_name"`
}

// Validate validates the service-profile data.
func (mg MulticastGroup) Validate() error {
	if strings.TrimSpace(mg.Name) == "" || len(mg.Name) > 100 {
		return errHandler.ErrMulticastGroupInvalidName
	}
	return nil
}

// MulticastGroupFilters provide filters that can be used to filter on
// multicast-groups. Note that empty values are not used as filters.
type MulticastGroupFilters struct {
	OrganizationID   int64         `db:"organization_id"`
	ServiceProfileID uuid.UUID     `db:"service_profile_id"`
	DevEUI           lorawan.EUI64 `db:"dev_eui"`
	Search           string        `db:"search"`

	// Limit and Offset are added for convenience so that this struct can
	// be given as the arguments.
	Limit  int `db:"limit"`
	Offset int `db:"offset"`
}

// SQL returns the SQL filter.
func (f MulticastGroupFilters) SQL() string {
	var filters []string
	var nilEUI lorawan.EUI64

	if f.OrganizationID != 0 {
		filters = append(filters, "o.id = :organization_id")
	}
	if f.ServiceProfileID != uuid.Nil {
		filters = append(filters, "mg.service_profile_id = :service_profile_id")
	}
	if f.DevEUI != nilEUI {
		filters = append(filters, "dmg.dev_eui = :dev_eui")
	}
	if f.Search != "" {
		filters = append(filters, "mg.name ilike :search")
	}

	if len(filters) == 0 {
		return ""
	}

	return "where " + strings.Join(filters, " and ")
}
