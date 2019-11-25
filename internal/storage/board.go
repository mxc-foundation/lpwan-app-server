package storage

import (
	"database/sql"
	"regexp"
	"time"

	"github.com/brocaar/lorawan"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"fmt"

	log "github.com/sirupsen/logrus"
)

var boardServerRegexp = regexp.MustCompile(`^[.\w-]+$`)

// Board represents a gateway.
type Board struct {
	MAC          lorawan.EUI64 `db:"mac"`
	SN           *string       `db:"sn"`
	CreatedAt    time.Time     `db:"created_at"`
	UpdatedAt    time.Time     `db:"updated_at"`
	Model        string        `db:"model"`
	VpnAddr      string        `db:"vpn_addr"`
	QaErr        int           `db:"qa_err"`
	OsVersion    *string       `db:"os_version"`
	FPGAVersion  *string       `db:"fpga_version"`
	RootPassword *string       `db:"root_password"`
	Server       *string       `db:"server"`
}

// Validate validates the board data.
func (bd Board) Validate() error {
	if bd.Server == nil {
		return ErrBoardInvalidServer
	}

	server := *bd.Server
	if !boardServerRegexp.MatchString(server) {
		return ErrBoardInvalidServer
	}

	return nil
}

// CreateBoard creates the given board.
func CreateBoard(db sqlx.Execer, bd *Board) error {
	now := time.Now()
	fmt.Println("CreateBoard... ")
	_, err := db.Exec(`
		insert into board (
			mac,
			sn,
			created_at,
			updated_at,
			model,
			vpn_addr,
			qa_err,
			os_version,
			fpga_version,
			root_password,
			server
		) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		bd.MAC[:],
		bd.SN,
		now,
		now,
		bd.Model,
		bd.VpnAddr,
		bd.QaErr,
		bd.OsVersion,
		bd.FPGAVersion,
		bd.RootPassword,
		bd.Server,
	)
	if err != nil {
		return handlePSQLError(Insert, err, "insert error")
	}

	bd.CreatedAt = now
	bd.UpdatedAt = now

	log.WithFields(log.Fields{
		"mac": bd.MAC,
		"sn":  bd.SN,
	}).Info("Board created")
	return nil
}

// GetBoard returns the board for the given mac.
func GetBoard(db sqlx.Queryer, mac lorawan.EUI64) (Board, error) {
	var bd Board
	err := sqlx.Get(db, &bd, "select * from board where mac = $1", mac[:])
	if err != nil {
		if err == sql.ErrNoRows {
			return bd, ErrDoesNotExist
		}
		return bd, errors.Wrap(err, "Get board error")
	}
	return bd, nil
}

// GetBoardMacBySerialNumber returns the board MAC for the serial number.
func GetBoardMacBySerialNumber(db sqlx.Queryer, sn string) (Board, error) {
	var bd Board
	err := sqlx.Get(db, &bd, "select * from board where sn = $1", sn)
	if err != nil {
		if err == sql.ErrNoRows {
			return bd, ErrDoesNotExist
		}
		return bd, errors.Wrap(err, "Get board by mac error")
	}

	return bd, nil
}

// UpdateBoard updates the given board.
func UpdateBoard(db sqlx.Execer, bd *Board) error {
	now := time.Now()

	res, err := db.Exec(`
		update board
			set updated_at = $2,
			model = $3,
			os_version = $4,
			fpga_version = $5,
			root_password = $6
		where
			sn = $1`,
		bd.SN,
		now,
		bd.Model,
		bd.OsVersion,
		bd.FPGAVersion,
		bd.RootPassword,
	)
	if err != nil {
		return handlePSQLError(Update, err, "update error")
	}
	err = handlePSQLEffect(res)
	if err != nil {
		return err
	}

	bd.UpdatedAt = now
	log.WithFields(log.Fields{
		"mac": bd.MAC,
		"sn":  bd.SN,
	}).Info("Board updated")

	return nil
}

// RegisterBoardAtomic updates "server" field atomically.
func RegisterBoardAtomic(db sqlx.Execer, bd *Board) error {
	if err := bd.Validate(); err != nil {
		return errors.Wrap(err, "validate error")
	}

	now := time.Now()

	res, err := db.Exec(`
		update board
			set updated_at = $2,
			server = $3
		where
			sn = $1
			AND (server IS NULL OR server = '')`,
		bd.SN,
		now,
		bd.Server,
	)
	if err != nil {
		return handlePSQLError(Update, err, "update error")
	}
	err = handlePSQLEffect(res)
	if err != nil {
		return err
	}

	bd.UpdatedAt = now
	log.WithFields(log.Fields{
		"mac":    bd.MAC,
		"sn":     bd.SN,
		"server": bd.Server,
	}).Info("Board registered")

	return nil
}

// UnregisterBoardAtomic ...
func UnregisterBoardAtomic(db sqlx.Execer, bd *Board) error {
	if err := bd.Validate(); err != nil {
		return errors.Wrap(err, "validate error")
	}

	now := time.Now()

	res, err := db.Exec(`
		update board
			set updated_at = $2,
			server = NULL
		where
			sn = $1
			AND server = $3`,
		bd.SN,
		now,
		bd.Server,
	)
	if err != nil {
		return handlePSQLError(Update, err, "update error")
	}
	err = handlePSQLEffect(res)
	if err != nil {
		return err
	}

	bd.UpdatedAt = now
	bd.Server = nil
	log.WithFields(log.Fields{
		"mac": bd.MAC,
		"sn":  bd.SN,
	}).Info("Board unregistered")

	return nil
}
