package auth

import (
	"context"
	"database/sql"

	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/auth"
)

type store struct {
	db *sql.DB
}

func New(db *sql.DB) auth.Store {
	return &store{db: db}
}

func (s *store) GetUser(ctx context.Context, username string) (*auth.User, error) {
	q := `SELECT id, is_admin FROM "user" WHERE username=$1 AND is_active=true`
	row := s.db.QueryRowContext(ctx, q, username)
	var res auth.User
	if err := row.Scan(&res.ID, &res.IsAdmin); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *store) GetOrgUser(ctx context.Context, userID int64, orgID int64) (auth.OrgUser, error) {
	q := `SELECT is_admin, is_device_admin, is_gateway_admin
		FROM organization_user WHERE user_id=$1 AND organization_id=$2`
	row := s.db.QueryRowContext(ctx, q, userID, orgID)
	var ou auth.OrgUser
	err := row.Scan(&ou.IsOrgAdmin, &ou.IsDeviceAdmin, &ou.IsGatewayAdmin)
	return ou, err
}
