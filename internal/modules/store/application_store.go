package store

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/gofrs/uuid"

	"github.com/mxc-foundation/lpwan-app-server/internal/codec"
)

type ApplicationStore interface {
	CreateApplication(ctx context.Context, item *Application) error
	GetApplication(ctx context.Context, id int64) (Application, error)
	GetApplicationCount(ctx context.Context, filters ApplicationFilters) (int, error)
	GetApplications(ctx context.Context, filters ApplicationFilters) ([]ApplicationListItem, error)
	UpdateApplication(ctx context.Context, item Application) error
	DeleteApplication(ctx context.Context, id int64) error
	DeleteAllApplicationsForOrganizationID(ctx context.Context, organizationID int64) error

	// validator
	CheckCreateApplicationAccess(ctx context.Context, username string, userID, organizationID int64) (bool, error)
	CheckListApplicationAccess(ctx context.Context, username string, userID, organizationID int64) (bool, error)

	CheckReadApplicationAccess(ctx context.Context, username string, userID, applicationID int64) (bool, error)
	CheckUpdateApplicationAccess(ctx context.Context, username string, userID, applicationID int64) (bool, error)
	CheckDeleteApplicationAccess(ctx context.Context, username string, userID, applicationID int64) (bool, error)
}

func (h *Handler) CreateApplication(ctx context.Context, item *Application) error {
	return h.store.CreateApplication(ctx, item)
}
func (h *Handler) GetApplication(ctx context.Context, id int64) (Application, error) {
	return h.store.GetApplication(ctx, id)
}
func (h *Handler) GetApplicationCount(ctx context.Context, filters ApplicationFilters) (int, error) {
	return h.store.GetApplicationCount(ctx, filters)
}
func (h *Handler) GetApplications(ctx context.Context, filters ApplicationFilters) ([]ApplicationListItem, error) {
	return h.store.GetApplications(ctx, filters)
}
func (h *Handler) UpdateApplication(ctx context.Context, item Application) error {
	return h.store.UpdateApplication(ctx, item)
}
func (h *Handler) DeleteApplication(ctx context.Context, id int64) error {
	return h.store.DeleteApplication(ctx, id)
}
func (h *Handler) DeleteAllApplicationsForOrganizationID(ctx context.Context, organizationID int64) error {
	return h.store.DeleteAllApplicationsForOrganizationID(ctx, organizationID)
}
func (h *Handler) CheckCreateApplicationAccess(ctx context.Context, username string, userID, organizationID int64) (bool, error) {
	return h.store.CheckCreateApplicationAccess(ctx, username, userID, organizationID)
}
func (h *Handler) CheckListApplicationAccess(ctx context.Context, username string, userID, organizationID int64) (bool, error) {
	return h.store.CheckListApplicationAccess(ctx, username, userID, organizationID)
}
func (h *Handler) CheckReadApplicationAccess(ctx context.Context, username string, userID, applicationID int64) (bool, error) {
	return h.store.CheckReadApplicationAccess(ctx, username, userID, applicationID)
}
func (h *Handler) CheckUpdateApplicationAccess(ctx context.Context, username string, userID, applicationID int64) (bool, error) {
	return h.store.CheckUpdateApplicationAccess(ctx, username, userID, applicationID)
}
func (h *Handler) CheckDeleteApplicationAccess(ctx context.Context, username string, userID, applicationID int64) (bool, error) {
	return h.store.CheckDeleteApplicationAccess(ctx, username, userID, applicationID)
}

var applicationNameRegexp = regexp.MustCompile(`^[\w-]+$`)

// Application represents an application.
type Application struct {
	ID                   int64      `db:"id"`
	Name                 string     `db:"name"`
	Description          string     `db:"description"`
	OrganizationID       int64      `db:"organization_id"`
	ServiceProfileID     uuid.UUID  `db:"service_profile_id"`
	PayloadCodec         codec.Type `db:"payload_codec"`
	PayloadEncoderScript string     `db:"payload_encoder_script"`
	PayloadDecoderScript string     `db:"payload_decoder_script"`
}

// ApplicationListItem devices the application as a list item.
type ApplicationListItem struct {
	Application
	ServiceProfileName string `db:"service_profile_name"`
}

// Validate validates the data of the Application.
func (a Application) Validate() error {
	if !applicationNameRegexp.MatchString(a.Name) {
		return errors.New("ErrApplicationInvalidName")
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
