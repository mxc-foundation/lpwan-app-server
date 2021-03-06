package data

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"github.com/gofrs/uuid"

	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
)

var applicationNameRegexp = regexp.MustCompile(`^[\w-]+$`)

// Application represents an application.
type Application struct {
	ID                   int64     `db:"id"`
	Name                 string    `db:"name"`
	Description          string    `db:"description"`
	OrganizationID       int64     `db:"organization_id"`
	ServiceProfileID     uuid.UUID `db:"service_profile_id"`
	PayloadCodec         string    `db:"payload_codec"`
	PayloadEncoderScript string    `db:"payload_encoder_script"`
	PayloadDecoderScript string    `db:"payload_decoder_script"`
}

// ApplicationListItem devices the application as a list item.
type ApplicationListItem struct {
	Application
	ServiceProfileName string `db:"service_profile_name"`
}

// Validate validates the data of the Application.
func (a Application) Validate() error {
	if !applicationNameRegexp.MatchString(a.Name) {
		return errHandler.ErrApplicationInvalidName
	}

	return nil
}

// ApplicationFilters provides filters for filtering applications.
type ApplicationFilters struct {
	UserID         int64  `db:"user_id"`
	OrganizationID int64  `db:"organization_id"`
	Search         string `db:"search"`

	// Limit and Offset are added for convenience so that this struct can
	// be given as the arguments.
	Limit  int `db:"limit"`
	Offset int `db:"offset"`
}

// SQL returns the SQL filters.
func (f ApplicationFilters) SQL() string {
	var filters []string

	if f.UserID != 0 {
		filters = append(filters, "u.id = :user_id")
	}

	if f.OrganizationID != 0 {
		filters = append(filters, "a.organization_id = :organization_id")
	}

	if f.Search != "" {
		filters = append(filters, "(a.name ilike :search)")
	}

	if len(filters) == 0 {
		return ""
	}

	return "where " + strings.Join(filters, " and ")
}

// Integration represents an integration.
type Integration struct {
	ID            int64           `db:"id"`
	CreatedAt     time.Time       `db:"created_at"`
	UpdatedAt     time.Time       `db:"updated_at"`
	ApplicationID int64           `db:"application_id"`
	Kind          string          `db:"kind"`
	Settings      json.RawMessage `db:"settings"`
}
