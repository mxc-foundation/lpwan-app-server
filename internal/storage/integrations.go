package storage

import (
	"context"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

// Integration represents an integration.
type Integration store.Integration

// CreateIntegration creates the given Integration.
func CreateIntegration(ctx context.Context, handler *store.Handler, i *Integration) error {
	return handler.CreateIntegration(ctx, (*store.Integration)(i))
}

// GetIntegration returns the Integration for the given id.
func GetIntegration(ctx context.Context, handler *store.Handler, id int64) (Integration, error) {
	res, err := handler.GetIntegration(ctx, id)
	return Integration(res), err
}

// GetIntegrationByApplicationID returns the Integration for the given
// application id and kind.
func GetIntegrationByApplicationID(ctx context.Context, handler *store.Handler, applicationID int64, kind string) (Integration, error) {
	res, err := handler.GetIntegrationByApplicationID(ctx, applicationID, kind)
	return Integration(res), err
}

// GetIntegrationsForApplicationID returns the integrations for the given
// application id.
func GetIntegrationsForApplicationID(ctx context.Context, handler *store.Handler, applicationID int64) ([]Integration, error) {
	res, err := handler.GetIntegrationsForApplicationID(ctx, applicationID)
	if err != nil {
		return nil, err
	}

	var iList []Integration
	for _, v := range res {
		iItem := Integration(v)
		iList = append(iList, iItem)
	}

	return iList, nil
}

// UpdateIntegration updates the given Integration.
func UpdateIntegration(ctx context.Context, handler *store.Handler, i *Integration) error {
	return handler.UpdateIntegration(ctx, (*store.Integration)(i))
}

// DeleteIntegration deletes the integration matching the given id.
func DeleteIntegration(ctx context.Context, handler *store.Handler, id int64) error {
	return handler.DeleteIntegration(ctx, id)
}
