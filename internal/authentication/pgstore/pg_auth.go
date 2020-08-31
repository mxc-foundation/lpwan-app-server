package pgstore

import (
	"context"

	"github.com/jmoiron/sqlx"

	authmod "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
)

type Handler struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) GetUser(ctx context.Context, username string) (authmod.User, error) {
	q := `SELECT id, is_admin FROM "user" WHERE email=$1 AND is_active=true`
	row := h.db.QueryRowContext(ctx, q, username)
	var res authmod.User
	if err := row.Scan(&res.ID, &res.IsGlobalAdmin); err != nil {
		return res, err
	}
	return res, nil
}

func (h *Handler) GetOrgUser(ctx context.Context, userID int64, orgID int64) (authmod.OrgUser, error) {
	q := `SELECT is_admin, is_device_admin, is_gateway_admin
		FROM organization_user WHERE user_id=$1 AND organization_id=$2`
	row := h.db.QueryRowContext(ctx, q, userID, orgID)
	var ou authmod.OrgUser
	err := row.Scan(&ou.IsOrgAdmin, &ou.IsDeviceAdmin, &ou.IsGatewayAdmin)
	return ou, err
}
