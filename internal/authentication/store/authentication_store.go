package store

import (
	"context"

	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"

	. "github.com/mxc-foundation/lpwan-app-server/internal/authentication/data"
)

func NewStore(pg pgstore.PgStore) *auths {
	return &auths{
		pg: pg,
	}
}

// Store provides access to information about users and their roles
type Store interface {
	// GetUser returns user's information given that there is an active user
	// with the given username
	GetUser(ctx context.Context, username string) (User, error)
	// GetOrgUser returns user's role in the listed organization
	GetOrgUser(ctx context.Context, userID int64, orgID int64) (OrgUser, error)
}

type auths struct {
	pg pgstore.AuthenticationPgStore
}

func (s *auths) GetUser(ctx context.Context, username string) (User, error) {
	return s.pg.AuthGetUser(ctx, username)
}
func (s *auths) GetOrgUser(ctx context.Context, userID int64, orgID int64) (OrgUser, error) {
	return s.pg.AuthGetOrgUser(ctx, userID, orgID)
}
