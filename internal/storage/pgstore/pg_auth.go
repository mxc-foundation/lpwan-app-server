package pgstore

import (
	"context"
	"database/sql"

	"github.com/mxc-foundation/lpwan-app-server/internal/auth"
)

func (ps *PgStore) AuthGetUser(ctx context.Context, username string) (auth.User, error) {
	q := `SELECT id, email, is_admin FROM "user" WHERE email=$1 AND is_active=true`
	row := ps.db.QueryRowContext(ctx, q, username)
	var res auth.User
	if err := row.Scan(&res.ID, &res.Email, &res.IsGlobalAdmin); err != nil {
		return res, err
	}
	return res, nil
}

func (ps *PgStore) AuthGetOrgUser(ctx context.Context, userID int64, orgID int64) (auth.OrgUser, error) {
	q := `SELECT is_admin, is_device_admin, is_gateway_admin
		FROM organization_user WHERE user_id=$1 AND organization_id=$2`
	row := ps.db.QueryRowContext(ctx, q, userID, orgID)
	var ou auth.OrgUser
	err := row.Scan(&ou.IsOrgAdmin, &ou.IsDeviceAdmin, &ou.IsGatewayAdmin)
	if err == nil {
		ou.IsOrgUser = true
	} else if err == sql.ErrNoRows {
		// if user is not an org member, then we just return an empty OrgUser,
		// it's not an error
		err = nil
	}
	return ou, err
}
