package pgstore

import (
	"context"
	"time"

	"github.com/brocaar/lorawan"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"

	pgerr "github.com/mxc-foundation/lpwan-app-server/internal/errors"
)

func (ps *PgStore) ExecuteMigrateUp(m migrate.MigrationSource) error {
	n, err := migrate.Exec(ctrl.db.DB.DB, "postgres", m, migrate.Up)
	if err != nil {
		return errors.Wrap(err, "storage: applying PostgreSQL data migrations error")
	}

	log.WithField("count", n).Info("storage: PostgreSQL data migrations applied")
	return nil
}

func (ps *PgStore) FixGorpMigrationsItemId(ctx context.Context, oldID, newID string) error {
	res, err := ps.db.ExecContext(ctx, `
			update gorp_migrations 
			set id=$1 where id=$2
		`, newID, oldID)
	if err != nil {
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

func (ps *PgStore) GetAllFromGorpMigrations(ctx context.Context) ([]string, error) {
	_, err := ps.db.ExecContext(ctx, `
		select exists ( select from gorp_migrations)
	`)
	if err != nil {
		return nil, pgerr.ErrEmptyGorpMigrations
	}

	var items []string
	err = sqlx.SelectContext(ctx, ps.db, &items, `
		select
			id
		from
			gorp_migrations
	`)
	if err != nil {
		return nil, handlePSQLError(Select, err, "select error")
	}

	return items, nil
}

func (ps *PgStore) Migrate(ctx context.Context, name string) error {
	_, err := ps.db.ExecContext(ctx, `lock table code_migration`)
	if err != nil {
		return errors.Wrap(err, "lock code migration table error")
	}

	var count int
	if err := ps.db.QueryRowxContext(ctx, `
		select count(id) from code_migration where id = $1
	`, name).Scan(&count); err == nil {
		return pgerr.ErrAlreadyExists
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

func (ps *PgStore) GetAllGatewayIDs(ctx context.Context) ([]lorawan.EUI64, error) {
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
