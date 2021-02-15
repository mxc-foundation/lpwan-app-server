package pgstore

import (
	"context"
	"database/sql"
)

// GetUserBonusOrgID given the user's email address returns the ID of the
// organization to which the bonuses awarded to the user should be paid. If
// user does not exist returns 0, nil
func (ps *PgStore) GetUserBonusOrgID(ctx context.Context, email string) (int64, error) {
	query := `SELECT ou.organization_id
			  FROM organization_user ou
			  	JOIN "user" u ON (ou.user_id = u.id)
			  WHERE u.email = $1 AND ou.is_admin
			  ORDER BY ou.created_at ASC
			  LIMIT 1`
	var orgID int64
	err := ps.db.QueryRowContext(ctx, query, email).Scan(&orgID)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return orgID, err
}
