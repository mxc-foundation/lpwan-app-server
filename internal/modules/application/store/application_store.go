package store

import (
	"context"

	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"

	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/application/data"
)

func NewApplicationStore(pg pgstore.PgStore) *as {
	return &as{
		pg: pg,
	}
}

type as struct {
	pg pgstore.ApplicationPgStore
}

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

func (h *as) CreateApplication(ctx context.Context, item *Application) error {
	return h.pg.CreateApplication(ctx, item)
}
func (h *as) GetApplication(ctx context.Context, id int64) (Application, error) {
	return h.pg.GetApplication(ctx, id)
}
func (h *as) GetApplicationCount(ctx context.Context, filters ApplicationFilters) (int, error) {
	return h.pg.GetApplicationCount(ctx, filters)
}
func (h *as) GetApplications(ctx context.Context, filters ApplicationFilters) ([]ApplicationListItem, error) {
	return h.pg.GetApplications(ctx, filters)
}
func (h *as) UpdateApplication(ctx context.Context, item Application) error {
	return h.pg.UpdateApplication(ctx, item)
}
func (h *as) DeleteApplication(ctx context.Context, id int64) error {
	return h.pg.DeleteApplication(ctx, id)
}
func (h *as) DeleteAllApplicationsForOrganizationID(ctx context.Context, organizationID int64) error {
	return h.pg.DeleteAllApplicationsForOrganizationID(ctx, organizationID)
}
func (h *as) CheckCreateApplicationAccess(ctx context.Context, username string, userID, organizationID int64) (bool, error) {
	return h.pg.CheckCreateApplicationAccess(ctx, username, userID, organizationID)
}
func (h *as) CheckListApplicationAccess(ctx context.Context, username string, userID, organizationID int64) (bool, error) {
	return h.pg.CheckListApplicationAccess(ctx, username, userID, organizationID)
}
func (h *as) CheckReadApplicationAccess(ctx context.Context, username string, userID, applicationID int64) (bool, error) {
	return h.pg.CheckReadApplicationAccess(ctx, username, userID, applicationID)
}
func (h *as) CheckUpdateApplicationAccess(ctx context.Context, username string, userID, applicationID int64) (bool, error) {
	return h.pg.CheckUpdateApplicationAccess(ctx, username, userID, applicationID)
}
func (h *as) CheckDeleteApplicationAccess(ctx context.Context, username string, userID, applicationID int64) (bool, error) {
	return h.pg.CheckDeleteApplicationAccess(ctx, username, userID, applicationID)
}
