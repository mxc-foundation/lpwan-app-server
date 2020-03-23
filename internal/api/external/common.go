package external

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/timestamp"
	api "github.com/mxc-foundation/lpwan-app-server/api/appserver_serves_ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/auth"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
)

func getUserProfileByJwt(ctx context.Context, v auth.Validator, organizationID int64) (api.ProfileResponse, error) {
	username, err := v.GetUsername(ctx)
	if nil != err {
		return api.ProfileResponse{}, err
	}

	// Get the user id based on the username.
	user, err := storage.GetUserByUsername(ctx, storage.DB(), username)
	if nil != err {
		return api.ProfileResponse{}, err
	}

	prof, err := storage.GetProfile(ctx, storage.DB(), user.ID)
	if err != nil {
		return api.ProfileResponse{}, err
	}

	userProfile := api.ProfileResponse{}

	userProfile.User = &api.User{
		Id:         prof.User.ID,
		Username:   prof.User.Username,
		SessionTtl: prof.User.SessionTTL,
		IsAdmin:    prof.User.IsAdmin,
		IsActive:   prof.User.IsActive,
	}
	userProfile.Settings = &api.ProfileSettings{
		DisableAssignExistingUsers: auth.DisableAssignExistingUsers,
	}

	orgDeleted := true
	for _, v := range prof.Organizations {
		id := v.ID
		org := api.OrganizationLink{}
		org.OrganizationId = id
		org.IsAdmin = v.IsAdmin
		org.OrganizationName = v.Name
		org.CreatedAt = &timestamp.Timestamp{Seconds: int64(v.CreatedAt.Second()), Nanos: int32(v.CreatedAt.Nanosecond())}
		org.UpdatedAt = &timestamp.Timestamp{Seconds: int64(v.UpdatedAt.Second()), Nanos: int32(v.UpdatedAt.Nanosecond())}
		userProfile.Organizations = append(userProfile.Organizations, &org)

		if id == organizationID {
			orgDeleted = false
		}
	}

	if organizationID == 0 && userProfile.User.IsAdmin {
		orgDeleted = false
	}

	if orgDeleted {
		return userProfile, fmt.Errorf("User does not have persmission to modify this organization: %d", organizationID)
	}

	return userProfile, nil
}
