package data

import (
	"strings"
	"time"

	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/gofrs/uuid"

	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
)

// ServiceProfile defines the service-profile.
type ServiceProfile struct {
	NetworkServerID int64             `db:"network_server_id"`
	OrganizationID  int64             `db:"organization_id"`
	CreatedAt       time.Time         `db:"created_at"`
	UpdatedAt       time.Time         `db:"updated_at"`
	Name            string            `db:"name"`
	ServiceProfile  ns.ServiceProfile `db:"-"`
}

// ServiceProfileMeta defines the service-profile meta record.
type ServiceProfileMeta struct {
	ServiceProfileID  uuid.UUID `db:"service_profile_id"`
	NetworkServerID   int64     `db:"network_server_id"`
	OrganizationID    int64     `db:"organization_id"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
	Name              string    `db:"name"`
	NetworkServerName string    `db:"network_server_name"`
}

// Validate validates the service-profile data.
func (sp ServiceProfile) Validate() error {
	if strings.TrimSpace(sp.Name) == "" || len(sp.Name) > 100 {
		return errHandler.ErrServiceProfileInvalidName
	}
	return nil
}
