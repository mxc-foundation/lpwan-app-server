package gatewayprofile

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
)

type GatewayProfileStore interface {
	CreateGatewayProfile(ctx context.Context, gp *GatewayProfile) error
	GetGatewayProfile(ctx context.Context, id uuid.UUID) (GatewayProfile, error)
	UpdateGatewayProfile(ctx context.Context, gp *GatewayProfile) error
	DeleteGatewayProfile(ctx context.Context, id uuid.UUID) error
	GetGatewayProfileCount(ctx context.Context) (int, error)
	GetGatewayProfileCountForNetworkServerID(ctx context.Context, networkServerID int64) (int, error)
	GetGatewayProfiles(ctx context.Context, limit, offset int) ([]GatewayProfileMeta, error)
	GetGatewayProfilesForNetworkServerID(ctx context.Context, networkServerID int64, limit, offset int) ([]GatewayProfileMeta, error)

	// validator
	CheckCreateUpdateDeleteGatewayProfileAccess(ctx context.Context, username string, userID int64) (bool, error)
	CheckReadListGatewayProfileAccess(ctx context.Context, username string, userID int64) (bool, error)
}

type Controller struct {
	St        GatewayProfileStore
	Validator Validator
}

var Service *Controller

func Setup() error {
	Service.St = store.New(storage.DB().DB)
	return nil
}
