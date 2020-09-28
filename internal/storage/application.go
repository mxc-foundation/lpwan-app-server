package storage

import (
	"context"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

// Application represents an application.
type Application store.Application

// ApplicationListItem devices the application as a list item.
type ApplicationListItem store.ApplicationListItem

// Validate validates the data of the Application.
func (a Application) Validate() error {
	return store.Application(a).Validate()
}

// CreateApplication creates the given Application.
func CreateApplication(ctx context.Context, handler *store.Handler, item *Application) error {
	return handler.CreateApplication(ctx, (*store.Application)(item))
}

// GetApplication returns the Application for the given id.
func GetApplication(ctx context.Context, handler *store.Handler, id int64) (Application, error) {
	app, err := handler.GetApplication(ctx, id)
	return Application(app), err
}

// ApplicationFilters provides filters for filtering applications.
type ApplicationFilters store.ApplicationFilters

// SQL returns the SQL filters.
func (f ApplicationFilters) SQL() string {
	return store.ApplicationFilters(f).SQL()
}

// GetApplicationCount returns the total number of applications.
func GetApplicationCount(ctx context.Context, handler *store.Handler, filters ApplicationFilters) (int, error) {
	return handler.GetApplicationCount(ctx, store.ApplicationFilters(filters))
}

// GetApplications returns a slice of applications, sorted by name and
// respecting the given limit and offset.
func GetApplications(ctx context.Context, handler *store.Handler, filters ApplicationFilters) ([]ApplicationListItem, error) {
	res, err := handler.GetApplications(ctx, store.ApplicationFilters(filters))
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
	return handler.UpdateApplication(ctx, store.Application(item))
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
