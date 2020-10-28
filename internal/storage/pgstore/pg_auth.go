package pgstore

import (
	"context"

	"github.com/mxc-foundation/lpwan-app-server/internal/authentication"
)

func (ps *PgStore) AuthGetUser(ctx context.Context, username string) (authentication.User, error) {
	q := `SELECT id, is_admin FROM "user" WHERE email=$1 AND is_active=true`
	row := ps.db.QueryRowContext(ctx, q, username)
	var res authentication.User
	if err := row.Scan(&res.ID, &res.IsGlobalAdmin); err != nil {
		return res, err
	}
	return res, nil
}

func (ps *PgStore) AuthGetOrgUser(ctx context.Context, userID int64, orgID int64) (authentication.OrgUser, error) {
	q := `SELECT is_admin, is_device_admin, is_gateway_admin
		FROM organization_user WHERE user_id=$1 AND organization_id=$2`
	row := ps.db.QueryRowContext(ctx, q, userID, orgID)
	var ou authentication.OrgUser
	err := row.Scan(&ou.IsOrgAdmin, &ou.IsDeviceAdmin, &ou.IsGatewayAdmin)
	return ou, err
}
