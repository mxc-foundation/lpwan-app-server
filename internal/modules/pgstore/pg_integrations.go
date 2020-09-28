package pgstore

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

func (ps *pgstore) CreateIntegration(ctx context.Context, i *store.Integration) error {
	now := time.Now()
	err := sqlx.GetContext(ctx, ps.db, &i.ID, `
		insert into integration (
			created_at,
			updated_at,
			application_id,
			kind,
			settings
		) values ($1, $2, $3, $4, $5) returning id`,
		now,
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
				return store.ErrAlreadyExists
			default:
				return errors.Wrap(err, "insert error")
			}
		default:
			return errors.Wrap(err, "insert error")
		}
	}

	i.CreatedAt = now
	i.UpdatedAt = now
	log.WithFields(log.Fields{
		"id":             i.ID,
		"kind":           i.Kind,
		"application_id": i.ApplicationID,
		"ctx_id":         ctx.Value(logging.ContextIDKey),
	}).Info("integration created")
	return nil
}

// GetIntegration returns the Integration for the given id.
func (ps *pgstore) GetIntegration(ctx context.Context, id int64) (store.Integration, error) {
	var i store.Integration
	err := sqlx.GetContext(ctx, ps.db, &i, "select * from integration where id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return i, store.ErrDoesNotExist
		}
		return i, errors.Wrap(err, "select error")
	}
	return i, nil
}

// GetIntegrationByApplicationID returns the Integration for the given
// application id and kind.
func (ps *pgstore) GetIntegrationByApplicationID(ctx context.Context, applicationID int64, kind string) (store.Integration, error) {
	var i store.Integration
	err := sqlx.GetContext(ctx, ps.db, &i, "select * from integration where application_id = $1 and kind = $2", applicationID, kind)
	if err != nil {
		if err == sql.ErrNoRows {
			return i, store.ErrDoesNotExist
		}
		return i, errors.Wrap(err, "select error")
	}
	return i, nil
}

// GetIntegrationsForApplicationID returns the integrations for the given
// application id.
func (ps *pgstore) GetIntegrationsForApplicationID(ctx context.Context, applicationID int64) ([]store.Integration, error) {
	var is []store.Integration
	err := sqlx.SelectContext(ctx, ps.db, &is, `
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
func (ps *pgstore) UpdateIntegration(ctx context.Context, i *store.Integration) error {
	now := time.Now()
	res, err := ps.db.ExecContext(ctx, `
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
				return store.ErrAlreadyExists
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
		return store.ErrDoesNotExist
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
func (ps *pgstore) DeleteIntegration(ctx context.Context, id int64) error {
	res, err := ps.db.ExecContext(ctx, "delete from integration where id = $1", id)
	if err != nil {
		return errors.Wrap(err, "delete error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return store.ErrDoesNotExist
	}

	log.WithFields(log.Fields{
		"id":     id,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("integration deleted")
	return nil
}
