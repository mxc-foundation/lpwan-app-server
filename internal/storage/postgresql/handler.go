package postgresql

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// PgDB :
var PgDB *PostgresHandler

// PostgresHandler :
type PostgresHandler struct {
	*sqlx.DB
}

func logQuery(query string, duration time.Duration, args ...interface{}) {
	log.WithFields(log.Fields{
		"query":    query,
		"args":     args,
		"duration": duration,
	}).Debug("sql query executed")
}

// AddDB :
func (h *PostgresHandler) AddDB(d *sqlx.DB) {
	PgDB.DB = d
}

// OpenDB:
func (h *PostgresHandler) OpenDB() error {
	d, err := sqlx.Open("postgres", config.C.PostgreSQL.DSN)
	if err != nil {
		return errors.Wrap(err, "storage: PostgreSQL connection error")
	}

	d.SetMaxOpenConns(config.C.PostgreSQL.MaxOpenConnections)
	d.SetMaxIdleConns(config.C.PostgreSQL.MaxIdleConnections)
	for {
		if err := d.Ping(); err != nil {
			log.WithError(err).Warning("storage: ping PostgreSQL database error, will retry in 2s")
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}
	PgDB.DB = d

	return nil
}

// GetDB :
func (h *PostgresHandler) GetDB() *sqlx.DB {
	return PgDB.DB
}

// Beginx returns a transaction with logging.
func (h *PostgresHandler) Beginx() (*TxLogger, error) {
	tx, err := h.DB.Beginx()
	return &TxLogger{tx}, err
}

// Query logs the queries executed by the Query method.
func (h *PostgresHandler) Query(query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	rows, err := h.DB.Query(query, args...)
	logQuery(query, time.Since(start), args...)
	return rows, err
}

// Queryx logs the queries executed by the Queryx method.
func (h *PostgresHandler) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	start := time.Now()
	rows, err := h.DB.Queryx(query, args...)
	logQuery(query, time.Since(start), args...)
	return rows, err
}

// QueryRowx logs the queries executed by the QueryRowx method.
func (h *PostgresHandler) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	start := time.Now()
	row := h.DB.QueryRowx(query, args...)
	logQuery(query, time.Since(start), args...)
	return row
}

// Exec logs the queries executed by the Exec method.
func (h *PostgresHandler) Exec(query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	res, err := h.DB.Exec(query, args...)
	logQuery(query, time.Since(start), args...)
	return res, err
}

// TxLogger logs the executed sql queries and their duration.
type TxLogger struct {
	*sqlx.Tx
}

// Query logs the queries executed by the Query method.
func (q *TxLogger) Query(query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	rows, err := q.Tx.Query(query, args...)
	logQuery(query, time.Since(start), args...)
	return rows, err
}

// Queryx logs the queries executed by the Queryx method.
func (q *TxLogger) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	start := time.Now()
	rows, err := q.Tx.Queryx(query, args...)
	logQuery(query, time.Since(start), args...)
	return rows, err
}

// QueryRowx logs the queries executed by the QueryRowx method.
func (q *TxLogger) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	start := time.Now()
	row := q.Tx.QueryRowx(query, args...)
	logQuery(query, time.Since(start), args...)
	return row
}

// Exec logs the queries executed by the Exec method.
func (q *TxLogger) Exec(query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	res, err := q.Tx.Exec(query, args...)
	logQuery(query, time.Since(start), args...)
	return res, err
}
