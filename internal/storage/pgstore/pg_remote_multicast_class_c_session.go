package pgstore

import (
	"context"
	"time"

	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"

	"github.com/brocaar/lorawan"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/mxc-foundation/lpwan-app-server/internal/logging"

	mg "github.com/mxc-foundation/lpwan-app-server/internal/applayer/multicastsetup/data"
)

// CreateRemoteMulticastClassCSession creates the given multicast Class-C session.
func (ps *PgStore) CreateRemoteMulticastClassCSession(ctx context.Context, sess *mg.RemoteMulticastClassCSession) error {
	now := time.Now()
	sess.CreatedAt = now
	sess.UpdatedAt = now

	_, err := ps.db.ExecContext(ctx, `
		insert into remote_multicast_class_c_session (
			dev_eui,
			multicast_group_id,
			created_at,
			updated_at,
			mc_group_id,
			session_time,
			session_time_out,
			dl_frequency,
			dr,
			state_provisioned,
			retry_after,
			retry_count,
			retry_interval
		) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		sess.DevEUI,
		sess.MulticastGroupID,
		sess.CreatedAt,
		sess.UpdatedAt,
		sess.McGroupID,
		sess.SessionTime,
		sess.SessionTimeOut,
		sess.DLFrequency,
		sess.DR,
		sess.StateProvisioned,
		sess.RetryAfter,
		sess.RetryCount,
		sess.RetryInterval,
	)
	if err != nil {
		return handlePSQLError(Insert, err, "insert error")
	}

	log.WithFields(log.Fields{
		"dev_eui":            sess.DevEUI,
		"multicast_group_id": sess.MulticastGroupID,
		"ctx_id":             ctx.Value(logging.ContextIDKey),
	}).Info("remote multicast class-c session created")

	return nil
}

// GetRemoteMulticastClassCSession returns the multicast Class-C session given
// a DevEUI and multicast-group ID.
func (ps *PgStore) GetRemoteMulticastClassCSession(ctx context.Context, devEUI lorawan.EUI64, multicastGroupID uuid.UUID, forUpdate bool) (mg.RemoteMulticastClassCSession, error) {
	var fu string
	if forUpdate {
		fu = " for update"
	}

	var sess mg.RemoteMulticastClassCSession
	if err := sqlx.GetContext(ctx, ps.db, &sess, `
		select
			*
		from
			remote_multicast_class_c_session
		where
			dev_eui = $1
			and multicast_group_id = $2`+fu,
		devEUI,
		multicastGroupID,
	); err != nil {
		return sess, handlePSQLError(Select, err, "select error")
	}

	return sess, nil
}

// GetRemoteMulticastClassCSessionByGroupID returns the multicast Class-C session given
// a DevEUI and McGroupID.
func (ps *PgStore) GetRemoteMulticastClassCSessionByGroupID(ctx context.Context, devEUI lorawan.EUI64, mcGroupID int, forUpdate bool) (mg.RemoteMulticastClassCSession, error) {
	var fu string
	if forUpdate {
		fu = " for update"
	}

	var sess mg.RemoteMulticastClassCSession
	if err := sqlx.GetContext(ctx, ps.db, &sess, `
		select
			*
		from
			remote_multicast_class_c_session
		where
			dev_eui = $1
			and mc_group_id = $2`+fu,
		devEUI,
		mcGroupID,
	); err != nil {
		return sess, handlePSQLError(Select, err, "select error")
	}

	return sess, nil
}

// GetPendingRemoteMulticastClassCSessions returns a slice of pending remote
// multicast Class-C sessions.
func (ps *PgStore) GetPendingRemoteMulticastClassCSessions(ctx context.Context, limit, maxRetryCount int) ([]mg.RemoteMulticastClassCSession, error) {
	var items []mg.RemoteMulticastClassCSession

	if err := sqlx.SelectContext(ctx, ps.db, &items, `
		select
			sess.*
		from
			remote_multicast_class_c_session sess
		inner join
			remote_multicast_setup ms
			on
				sess.dev_eui = ms.dev_eui
				and sess.multicast_group_id = ms.multicast_group_id
				and sess.mc_group_id = ms.mc_group_id
		where
			ms.state_provisioned = true
			and ms.state = $3
			and sess.state_provisioned = false
			and sess.retry_count < $1
			and sess.retry_after < $2
		limit $4
		for update
		skip locked`,
		maxRetryCount,
		time.Now(),
		mg.RemoteMulticastSetupSetup,
		limit,
	); err != nil {
		return nil, handlePSQLError(Select, err, "select error")
	}

	return items, nil
}

// UpdateRemoteMulticastClassCSession updates the given remote multicast
// Class-C session.
func (ps *PgStore) UpdateRemoteMulticastClassCSession(ctx context.Context, sess *mg.RemoteMulticastClassCSession) error {
	sess.UpdatedAt = time.Now()

	res, err := ps.db.ExecContext(ctx, `
		update
			remote_multicast_class_c_session
		set
			updated_at = $3,
			mc_group_id = $4,
			session_time = $5,
			session_time_out = $6,
			dl_frequency = $7,
			dr = $8,
			state_provisioned = $9,
			retry_after = $10,
			retry_count = $11,
			retry_interval = $12
		where
			dev_eui = $1
			and multicast_group_id = $2`,
		sess.DevEUI,
		sess.MulticastGroupID,
		sess.UpdatedAt,
		sess.McGroupID,
		sess.SessionTime,
		sess.SessionTimeOut,
		sess.DLFrequency,
		sess.DR,
		sess.StateProvisioned,
		sess.RetryAfter,
		sess.RetryCount,
		sess.RetryInterval,
	)
	if err != nil {
		return handlePSQLError(Update, err, "update error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return errHandler.ErrDoesNotExist
	}

	log.WithFields(log.Fields{
		"dev_eui":            sess.DevEUI,
		"multicast_group_id": sess.MulticastGroupID,
		"ctx_id":             ctx.Value(logging.ContextIDKey),
	}).Info("remote multicast class-c session updated")
	return nil
}

// DeleteRemoteMulticastClassCSession deletes the multicast Class-C session
// given a DevEUI and multicast-group ID.
func (ps *PgStore) DeleteRemoteMulticastClassCSession(ctx context.Context, devEUI lorawan.EUI64, multicastGroupID uuid.UUID) error {
	res, err := ps.db.ExecContext(ctx, `
		delete from remote_multicast_class_c_session
		where
			dev_eui = $1
			and multicast_group_id = $2`,
		devEUI,
		multicastGroupID,
	)
	if err != nil {
		return handlePSQLError(Delete, err, "delete error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return errHandler.ErrDoesNotExist
	}
	log.WithFields(log.Fields{
		"dev_eui":            devEUI,
		"multicast_group_id": multicastGroupID,
		"ctx_id":             ctx.Value(logging.ContextIDKey),
	}).Info("remote multicast class-c session deleted")
	return nil
}
