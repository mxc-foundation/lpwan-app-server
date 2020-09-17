package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/mxc-foundation/lpwan-app-server/internal/logging"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

// Integration represents an integration.
type Integration store.Integration

// CreateIntegration creates the given Integration.
func CreateIntegration(ctx context.Context, handler *store.Handler, i *Integration) error {
	return handler.CreateIntegration(ctx, i)
}

// GetIntegration returns the Integration for the given id.
func GetIntegration(ctx context.Context, handler *store.Handler, id int64) (Integration, error) {
	var i Integration
	err := sqlx.Get(db, &i, "select * from integration where id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return i, ErrDoesNotExist
		}
		return i, errors.Wrap(err, "select error")
	}
	return i, nil
}

// GetIntegrationByApplicationID returns the Integration for the given
// application id and kind.
func GetIntegrationByApplicationID(ctx context.Context, handler *store.Handler, applicationID int64, kind string) (Integration, error) {
	var i Integration
	err := sqlx.Get(db, &i, "select * from integration where application_id = $1 and kind = $2", applicationID, kind)
	if err != nil {
		if err == sql.ErrNoRows {
			return i, ErrDoesNotExist
		}
		return i, errors.Wrap(err, "select error")
	}
	return i, nil
}

// GetIntegrationsForApplicationID returns the integrations for the given
// application id.
func GetIntegrationsForApplicationID(ctx context.Context, handler *store.Handler, applicationID int64) ([]Integration, error) {
	var is []Integration
	err := sqlx.Select(db, &is, `
		select *
		from integration
		where application_id = $1
		order by kind`,
		applicationID,
	)
	if err != nil {
		return nil, errors.Wrap(err, "select error")
	}
	return is, nil
}

// UpdateIntegration updates the given Integration.
func UpdateIntegration(ctx context.Context, handler *store.Handler, i *Integration) error {
	now := time.Now()
	res, err := db.Exec(`
		update integration
		set
			updated_at = $2,
			application_id = $3,
			kind = $4,
			settings = $5
		where
			id = $1`,
		i.ID,
		now,
		i.ApplicationID,
		i.Kind,
		i.Settings,
	)

	if err != nil {
		switch err := err.(type) {
		case *pq.Error:
			switch err.Code.Name() {
			case "unique_violation":
				return ErrAlreadyExists
			default:
				return errors.Wrap(err, "update error")
			}
		default:
			return errors.Wrap(err, "update error")
		}
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return ErrDoesNotExist
	}

	i.UpdatedAt = now
	log.WithFields(log.Fields{
		"id":             i.ID,
		"kind":           i.Kind,
		"application_id": i.ApplicationID,
		"ctx_id":         ctx.Value(logging.ContextIDKey),
	}).Info("integration updated")
	return nil
}

// DeleteIntegration deletes the integration matching the given id.
func DeleteIntegration(ctx context.Context, handler *store.Handler, id int64) error {
	res, err := db.Exec("delete from integration where id = $1", id)
	if err != nil {
		return errors.Wrap(err, "delete error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return ErrDoesNotExist
	}

	log.WithFields(log.Fields{
		"id":     id,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("integration deleted")
	return nil
}
