package storage

import (
	"context"

	sps "github.com/mxc-foundation/lpwan-app-server/internal/modules/service-profile/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

// ServiceProfile defines the service-profile.
type ServiceProfile sps.ServiceProfile

// ServiceProfileMeta defines the service-profile meta record.
type ServiceProfileMeta sps.ServiceProfileMeta

// Validate validates the service-profile data.
func (sp ServiceProfile) Validate() error {
	return sps.ServiceProfile(sp).Validate()
}

// GetServiceProfileCount returns the total number of service-profiles.
func GetServiceProfileCount(ctx context.Context, handler *store.Handler) (int, error) {
	return handler.GetServiceProfileCount(ctx)
}

// GetServiceProfileCountForOrganizationID returns the total number of
// service-profiles for the given organization id.
func GetServiceProfileCountForOrganizationID(ctx context.Context, handler *store.Handler, organizationID int64) (int, error) {
	return handler.GetServiceProfileCountForOrganizationID(ctx, organizationID)
}

// GetServiceProfileCountForUser returns the total number of service-profiles
// for the given user ID.
func GetServiceProfileCountForUser(ctx context.Context, handler *store.Handler, userID int64) (int, error) {
	return handler.GetServiceProfileCountForUser(ctx, userID)
}

// GetServiceProfiles returns a slice of service-profiles.
func GetServiceProfiles(ctx context.Context, handler *store.Handler, limit, offset int) ([]ServiceProfileMeta, error) {
	res, err := handler.GetServiceProfiles(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	var spMetaList []ServiceProfileMeta
	for _, v := range res {
		spMetaItem := ServiceProfileMeta(v)
		spMetaList = append(spMetaList, spMetaItem)
	}
	return spMetaList, nil
}

// GetServiceProfilesForOrganizationID returns a slice of service-profiles
// for the given organization id.
func GetServiceProfilesForOrganizationID(ctx context.Context, handler *store.Handler, organizationID int64, limit, offset int) ([]ServiceProfileMeta, error) {
	res, err := handler.GetServiceProfilesForOrganizationID(ctx, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}

	var spMetaList []ServiceProfileMeta
	for _, v := range res {
		spMetaItem := ServiceProfileMeta(v)
		spMetaList = append(spMetaList, spMetaItem)
	}
	return spMetaList, nil
}

// GetServiceProfilesForUser returns a slice of service-profile for the given
// user ID.
func GetServiceProfilesForUser(ctx context.Context, handler *store.Handler, userID int64, limit, offset int) ([]ServiceProfileMeta, error) {
	res, err := handler.GetServiceProfilesForUser(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	var spMetaList []ServiceProfileMeta
	for _, v := range res {
		spMetaItem := ServiceProfileMeta(v)
		spMetaList = append(spMetaList, spMetaItem)
	}
	return spMetaList, nil
}

// DeleteAllServiceProfilesForOrganizationID deletes all service-profiles
// given an organization id.
func DeleteAllServiceProfilesForOrganizationID(ctx context.Context, handler *store.Handler, organizationID int64) error {
	return handler.DeleteAllServiceProfilesForOrganizationID(ctx, organizationID)
}
