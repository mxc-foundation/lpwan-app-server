package pgstore

import (
	"context"
	"database/sql"

	"github.com/gofrs/uuid"

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

// ApplicationOwnedByOrganization returns true when given orgID owns given applicationID
func (ps *PgStore) ApplicationOwnedByOrganization(ctx context.Context, orgID, applicationID int64) (bool, error) {
	var count int64
	q := `SELECT count(a.id) 
		FROM application a JOIN organization org ON a.organization_id = org.id 
		WHERE a.id=$1 AND org.id=$2`
	err := ps.db.QueryRowContext(ctx, q, applicationID, orgID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// DeviceProfileOwnedByOrganization returns true when given orgID owns given device profile
func (ps *PgStore) DeviceProfileOwnedByOrganization(ctx context.Context, orgID int64, deviceProfile uuid.UUID) (bool, error) {
	var count int64
	q := `SELECT count(dp.device_profile_id) 
		FROM device_profile dp JOIN organization org ON dp.organization_id = org.id 
		WHERE dp.device_profile_id=$1 AND org.id=$2`
	err := ps.db.QueryRowContext(ctx, q, deviceProfile, orgID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
