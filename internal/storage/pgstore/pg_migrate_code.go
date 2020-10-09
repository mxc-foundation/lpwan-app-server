package pgstore

import (
	"context"
	"time"

	"github.com/brocaar/lorawan"
	"github.com/jmoiron/sqlx"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type MigrateCodePgStore interface {
	Migrate(ctx context.Context, name string) error
	GetAllGatewayIDs(ctx context.Context) ([]lorawan.EUI64, error)
}

func (ps *pgstore) Migrate(ctx context.Context, name string) error {
	_, err := ps.db.ExecContext(ctx, `lock table code_migration`)
	if err != nil {
		return errors.Wrap(err, "lock code migration table error")
	}

	res, err := ps.db.ExecContext(ctx, `
			insert into code_migration (
				id,
				applied_at
			) values ($1, $2)
			on conflict
				do nothing
		`, name, time.Now())
	if err != nil {
		switch err := err.(type) {
		case *pq.Error:
			switch err.Code.Name() {
			case "unique_violation":
				return nil
			}
		}

		return err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if ra == 0 {
		return nil
	}

	return nil
}

func (ps *pgstore) GetAllGatewayIDs(ctx context.Context) ([]lorawan.EUI64, error) {
	var ids []lorawan.EUI64
	err := sqlx.SelectContext(ctx, ps.db, &ids, `
		select
			mac
		from
			gateway
	`)
	if err != nil {
		return nil, handlePSQLError(Select, err, "select error")
	}

	return ids, nil
}
