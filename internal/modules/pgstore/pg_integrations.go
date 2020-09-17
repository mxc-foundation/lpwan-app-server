package pgstore

import (
	"context"
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
