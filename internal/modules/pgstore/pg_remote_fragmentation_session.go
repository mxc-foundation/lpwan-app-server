package pgstore

import (
	"context"
	"fmt"
	"time"

	"github.com/brocaar/lorawan"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/mxc-foundation/lpwan-app-server/internal/logging"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

// CreateRemoteFragmentationSession creates the given fragmentation session.
func (ps *pgstore) CreateRemoteFragmentationSession(ctx context.Context, sess *store.RemoteFragmentationSession) error {
	now := time.Now()
	sess.CreatedAt = now
	sess.UpdatedAt = now

	_, err := ps.db.ExecContext(ctx, `
		insert into remote_fragmentation_session (
			dev_eui,
			frag_index,
			created_at,
			updated_at,
			mc_group_ids,
			nb_frag,
			frag_size,
			fragmentation_matrix,
			block_ack_delay,
			padding,
			descriptor,
			state,
			state_provisioned,
			retry_after,
			retry_count,
			retry_interval
		) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`,
		sess.DevEUI,
		sess.FragIndex,
		sess.CreatedAt,
		sess.UpdatedAt,
		pq.Array(sess.MCGroupIDs),
		sess.NbFrag,
		sess.FragSize,
		[]byte{sess.FragmentationMatrix},
		sess.BlockAckDelay,
		sess.Padding,
		sess.Descriptor[:],
		sess.State,
		sess.StateProvisioned,
		sess.RetryAfter,
		sess.RetryCount,
		sess.RetryInterval,
	)
	if err != nil {
		return handlePSQLError(Insert, err, "insert error")
	}

	log.WithFields(log.Fields{
		"dev_eui":    sess.DevEUI,
		"frag_index": sess.FragIndex,
		"ctx_id":     ctx.Value(logging.ContextIDKey),
	}).Info("remote fragmentation session created")
	return nil
}

// GetRemoteFragmentationSession returns the fragmentation session given a
// DevEUI and fragmentation index.
func (ps *pgstore) GetRemoteFragmentationSession(ctx context.Context, devEUI lorawan.EUI64, fragIndex int, forUpdate bool) (store.RemoteFragmentationSession, error) {
	var fu string
	if forUpdate {
		fu = " for update"
	}

	row := ps.db.QueryRowxContext(ctx, `
		select
			dev_eui,
			frag_index,
			created_at,
			updated_at,
			mc_group_ids,
			nb_frag,
			frag_size,
			fragmentation_matrix,
			block_ack_delay,
			padding,
			descriptor,
			state,
			state_provisioned,
			retry_after,
			retry_count,
			retry_interval
		from
			remote_fragmentation_session
		where
			dev_eui = $1
			and frag_index = $2`+fu,
		devEUI,
		fragIndex,
	)

	return ps.scanRemoteFragmentationSession(row)
}

// GetPendingRemoteFragmentationSessions returns a slice of pending remote
// fragmentation sessions.
func (ps *pgstore) GetPendingRemoteFragmentationSessions(ctx context.Context, limit, maxRetryCount int) ([]store.RemoteFragmentationSession, error) {
	var items []store.RemoteFragmentationSession

	rows, err := ps.db.QueryxContext(ctx, `
		select
			fs.dev_eui,
			fs.frag_index,
			fs.created_at,
			fs.updated_at,
			fs.mc_group_ids,
			fs.nb_frag,
			fs.frag_size,
			fs.fragmentation_matrix,
			fs.block_ack_delay,
			fs.padding,
			fs.descriptor,
			fs.state,
			fs.state_provisioned,
			fs.retry_after,
			fs.retry_count,
			fs.retry_interval
		from
			remote_fragmentation_session fs
		where
			fs.state_provisioned = false
			and fs.retry_count < $1
			and fs.retry_after < $2
			and (
				-- in case of unicast
				array_length(fs.mc_group_ids, 1) is null

				-- in case of multicast
				or exists (
					select
						1
					from
						remote_multicast_setup ms
					where
						ms.dev_eui = fs.dev_eui
						and ms.state_provisioned = true
						and ms.mc_group_id = any(fs.mc_group_ids)
				)
			)
		limit $3
		for update of fs
		skip locked`,
		maxRetryCount,
		time.Now(),
		limit,
	)
	if err != nil {
		return nil, handlePSQLError(Select, err, "select error")
	}
	defer rows.Close()

	for rows.Next() {
		item, err := ps.scanRemoteFragmentationSession(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// UpdateRemoteFragmentationSession updates the given fragmentation session.
func (ps *pgstore) UpdateRemoteFragmentationSession(ctx context.Context, sess *store.RemoteFragmentationSession) error {
	sess.UpdatedAt = time.Now()

	res, err := ps.db.ExecContext(ctx, `
		update
			remote_fragmentation_session
		set
			updated_at = $3,
			mc_group_ids = $4,
			nb_frag = $5,
			frag_size = $6,
			fragmentation_matrix = $7,
			block_ack_delay = $8,
			padding = $9,
			descriptor = $10,
			state = $11,
			state_provisioned = $12,
			retry_after = $13,
			retry_count = $14,
			retry_interval = $15
		where
			dev_eui = $1
			and frag_index = $2`,
		sess.DevEUI,
		sess.FragIndex,
		sess.UpdatedAt,
		pq.Array(sess.MCGroupIDs),
		sess.NbFrag,
		sess.FragSize,
		[]byte{sess.FragmentationMatrix},
		sess.BlockAckDelay,
		sess.Padding,
		sess.Descriptor[:],
		sess.State,
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
		return store.ErrDoesNotExist
	}

	log.WithFields(log.Fields{
		"dev_eui":    sess.DevEUI,
		"frag_index": sess.FragIndex,
		"ctx_id":     ctx.Value(logging.ContextIDKey),
	}).Info("remote fragmentation session updated")
	return nil
}

// DeleteRemoteFragmentationSession removes the fragmentation session for the
// given DevEUI / fragmentation index combination.
func (ps *pgstore) DeleteRemoteFragmentationSession(ctx context.Context, devEUI lorawan.EUI64, fragIndex int) error {
	res, err := ps.db.ExecContext(ctx, `
		delete from remote_fragmentation_session
		where
			dev_eui = $1
			and frag_index = $2`,
		devEUI,
		fragIndex,
	)
	if err != nil {
		return handlePSQLError(Delete, err, "delete error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affacted error")
	}
	if ra == 0 {
		return store.ErrDoesNotExist
	}

	log.WithFields(log.Fields{
		"dev_eui":    devEUI,
		"frag_index": fragIndex,
		"ctx_id":     ctx.Value(logging.ContextIDKey),
	}).Info("remote fragmentation session deleted")
	return nil
}

func (ps *pgstore) scanRemoteFragmentationSession(row sqlx.ColScanner) (store.RemoteFragmentationSession, error) {
	var sess store.RemoteFragmentationSession

	var mcGroupIDs []int64
	var fragmentationMatrix []byte
	var descriptor []byte

	err := row.Scan(
		&sess.DevEUI,
		&sess.FragIndex,
		&sess.CreatedAt,
		&sess.UpdatedAt,
		pq.Array(&mcGroupIDs),
		&sess.NbFrag,
		&sess.FragSize,
		&fragmentationMatrix,
		&sess.BlockAckDelay,
		&sess.Padding,
		&descriptor,
		&sess.State,
		&sess.StateProvisioned,
		&sess.RetryAfter,
		&sess.RetryCount,
		&sess.RetryInterval,
	)
	if err != nil {
		return sess, handlePSQLError(Select, err, "select error")
	}

	for _, v := range mcGroupIDs {
		sess.MCGroupIDs = append(sess.MCGroupIDs, int(v))
	}

	if len(fragmentationMatrix) != 1 {
		return sess, fmt.Errorf("FragmentationMatrix must have length 1, got %d", len(fragmentationMatrix))
	}
	sess.FragmentationMatrix = fragmentationMatrix[0]

	if len(descriptor) != len(sess.Descriptor) {
		return sess, fmt.Errorf("Descriptor must have length %d, got %d", len(sess.Descriptor), len(descriptor))
	}
	copy(sess.Descriptor[:], descriptor)

	return sess, nil
}
