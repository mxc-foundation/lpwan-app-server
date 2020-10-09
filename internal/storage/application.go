package storage

import (
	"context"

	apps "github.com/mxc-foundation/lpwan-app-server/internal/modules/application/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

// Application represents an application.
type Application apps.Application

// ApplicationListItem devices the application as a list item.
type ApplicationListItem apps.ApplicationListItem

// Validate validates the data of the Application.
func (a Application) Validate() error {
	return apps.Application(a).Validate()
}

// CreateApplication creates the given Application.
func CreateApplication(ctx context.Context, handler *store.Handler, item *Application) error {
	return handler.CreateApplication(ctx, (*apps.Application)(item))
}

// GetApplication returns the Application for the given id.
func GetApplication(ctx context.Context, handler *store.Handler, id int64) (Application, error) {
	app, err := handler.GetApplication(ctx, id)
	return Application(app), err
}

// ApplicationFilters provides filters for filtering applications.
type ApplicationFilters apps.ApplicationFilters

// SQL returns the SQL filters.
func (f ApplicationFilters) SQL() string {
	return apps.ApplicationFilters(f).SQL()
}

// GetApplicationCount returns the total number of applications.
func GetApplicationCount(ctx context.Context, handler *store.Handler, filters ApplicationFilters) (int, error) {
	return handler.GetApplicationCount(ctx, apps.ApplicationFilters(filters))
}

// GetApplications returns a slice of applications, sorted by name and
// respecting the given limit and offset.
func GetApplications(ctx context.Context, handler *store.Handler, filters ApplicationFilters) ([]ApplicationListItem, error) {
	res, err := handler.GetApplications(ctx, apps.ApplicationFilters(filters))
	if err != nil {
		return nil, err
	}

	var apps []ApplicationListItem
	for _, v := range res {
		appItem := ApplicationListItem(v)
		apps = append(apps, appItem)
	}
	return apps, nil
}

// UpdateApplication updates the given Application.
func UpdateApplication(ctx context.Context, handler *store.Handler, item Application) error {
	return handler.UpdateApplication(ctx, apps.Application(item))
}

// DeleteApplication deletes the Application matching the given ID.
func DeleteApplication(ctx context.Context, handler *store.Handler, id int64) error {
	return handler.DeleteApplication(ctx, id)
}

// DeleteAllApplicationsForOrganizationID deletes all applications
// given an organization id.
func DeleteAllApplicationsForOrganizationID(ctx context.Context, handler *store.Handler, organizationID int64) error {
	return handler.DeleteAllApplicationsForOrganizationID(ctx, organizationID)
}
