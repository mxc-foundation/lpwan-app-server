package application

import (
	"golang.org/x/net/context"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
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

type Controller struct {
	St        ApplicationStore
	Validator Validator
}

var Service *Controller

func Setup() error {
	Service.St = store.New(storage.DB().DB)
	return nil
}

func (c *Controller) GetApplication(ctx context.Context, applicationID int64) (Application, error) {
	return c.St.GetApplication(ctx, applicationID)
}
