package pgstore

import (
	"context"
	"time"

	"github.com/brocaar/lorawan"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mxc-foundation/lpwan-app-server/internal/logging"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// CreateRemoteMulticastSetup creates the given multicast-setup.
func (ps *pgstore) CreateRemoteMulticastSetup(ctx context.Context, dms *store.RemoteMulticastSetup) error {
	now := time.Now()
	dms.CreatedAt = now
	dms.UpdatedAt = now

	_, err := ps.db.ExecContext(ctx, `
		insert into remote_multicast_setup (
			dev_eui,
			multicast_group_id,
			created_at,
			updated_at,
			mc_group_id,
			mc_addr,
			mc_key_encrypted,
			min_mc_f_cnt,
			max_mc_f_cnt,
			state,
			state_provisioned,
			retry_after,
			retry_count,
			retry_interval
		) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
		dms.DevEUI[:],
		dms.MulticastGroupID,
		dms.CreatedAt,
		dms.UpdatedAt,
		dms.McGroupID,
		dms.McAddr[:],
		dms.McKeyEncrypted[:],
		dms.MinMcFCnt,
		dms.MaxMcFCnt,
		dms.State,
		dms.StateProvisioned,
		dms.RetryAfter,
		dms.RetryCount,
		dms.RetryInterval,
	)
	if err != nil {
		return handlePSQLError(Insert, err, "insert error")
	}

	log.WithFields(log.Fields{
		"dev_eui":            dms.DevEUI,
		"multicast_group_id": dms.MulticastGroupID,
		"ctx_id":             ctx.Value(logging.ContextIDKey),
	}).Info("remote multicast-setup created")
	return nil
}

// GetRemoteMulticastSetup returns the multicast-setup given a multicast-group ID and DevEUI.
func (ps *pgstore) GetRemoteMulticastSetup(ctx context.Context, devEUI lorawan.EUI64, multicastGroupID uuid.UUID, forUpdate bool) (store.RemoteMulticastSetup, error) {
	var fu string
	if forUpdate {
		fu = " for update"
	}

	var dmg store.RemoteMulticastSetup
	if err := sqlx.GetContext(ctx, ps.db, &dmg, `
		select
			*
		from
			remote_multicast_setup
		where
			dev_eui = $1
			and multicast_group_id = $2`+fu,
		devEUI,
		multicastGroupID,
	); err != nil {
		return dmg, handlePSQLError(Select, err, "select error")
	}

	return dmg, nil
}

// GetRemoteMulticastSetupByGroupID returns the multicast-setup given a DevEUI and McGroupID.
func (ps *pgstore) GetRemoteMulticastSetupByGroupID(ctx context.Context, devEUI lorawan.EUI64, mcGroupID int, forUpdate bool) (store.RemoteMulticastSetup, error) {
	var fu string
	if forUpdate {
		fu = " for update"
	}

	var dmg store.RemoteMulticastSetup
	if err := sqlx.GetContext(ctx, ps.db, &dmg, `
		select
			*
		from
			remote_multicast_setup
		where
			dev_eui = $1
			and mc_group_id = $2`+fu,
		devEUI,
		mcGroupID,
	); err != nil {
		return dmg, handlePSQLError(Select, err, "select error")
	}

	return dmg, nil
}

// GetPendingRemoteMulticastSetupItems returns a slice of pending remote multicast-setup items.
// The selected items will be locked.
func (ps *pgstore) GetPendingRemoteMulticastSetupItems(ctx context.Context, limit, maxRetryCount int) ([]store.RemoteMulticastSetup, error) {
	var items []store.RemoteMulticastSetup

	if err := sqlx.SelectContext(ctx, ps.db, &items, `
		select
			*
		from
			remote_multicast_setup
		where
			state_provisioned = false
			and retry_count < $1
			and retry_after < $2
		limit $3
		for update
		skip locked`,
		maxRetryCount,
		time.Now(),
		limit,
	); err != nil {
		return nil, handlePSQLError(Select, err, "select error")
	}

	return items, nil
}

// UpdateRemoteMulticastSetup updates the given update multicast-group setup.
func (ps *pgstore) UpdateRemoteMulticastSetup(ctx context.Context, dmg *store.RemoteMulticastSetup) error {
	dmg.UpdatedAt = time.Now()

	res, err := ps.db.ExecContext(ctx, `
		update
			remote_multicast_setup
		set
			updated_at = $3,
			mc_group_id = $4,
			mc_addr = $5,
			mc_key_encrypted = $6,
			min_mc_f_cnt = $7,
			max_mc_f_cnt = $8,
			state = $9,
			state_provisioned = $10,
			retry_after = $11,
			retry_count = $12,
			retry_interval = $13
		where
			dev_eui = $1
			and multicast_group_id = $2`,
		dmg.DevEUI,
		dmg.MulticastGroupID,
		dmg.UpdatedAt,
		dmg.McGroupID,
		dmg.McAddr[:],
		dmg.McKeyEncrypted[:],
		dmg.MinMcFCnt,
		dmg.MaxMcFCnt,
		dmg.State,
		dmg.StateProvisioned,
		dmg.RetryAfter,
		dmg.RetryCount,
		dmg.RetryInterval,
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
		"dev_eui":            dmg.DevEUI,
		"multicast_group_id": dmg.MulticastGroupID,
		"ctx_id":             ctx.Value(logging.ContextIDKey),
	}).Info("remote multicast-setup updated")
	return nil
}

// DeleteRemoteMulticastSetup deletes the multicast-setup given a multicast-group ID and DevEUI.
func (ps *pgstore) DeleteRemoteMulticastSetup(ctx context.Context, devEUI lorawan.EUI64, multicastGroupID uuid.UUID) error {
	res, err := ps.db.ExecContext(ctx, `
		delete from remote_multicast_setup
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
		return store.ErrDoesNotExist
	}

	log.WithFields(log.Fields{
		"dev_eui":            devEUI,
		"multicast_group_id": multicastGroupID,
		"ctx_id":             ctx.Value(logging.ContextIDKey),
	}).Info("remote multicast-setup deleted")
	return nil
}

// GetDevEUIsWithMulticastSetup query all devices with complete multicast setup
func (ps *pgstore) GetDevEUIsWithMulticastSetup(ctx context.Context, id *uuid.UUID) ([]lorawan.EUI64, error) {
	var devEUIs []lorawan.EUI64
	err := sqlx.SelectContext(ctx, ps.db, &devEUIs, `
		select
			dev_eui
		from
			remote_multicast_setup
		where
			multicast_group_id = $1
			and state = $2
			and state_provisioned = $3`,
		id,
		store.RemoteMulticastSetupSetup,
		true,
	)
	if err != nil {
		return nil, handlePSQLError(Select, err, "select error")
	}

	return devEUIs, nil
}

// GetDevEUIsWithFragmentationSessionSetup query all devices with complete fragmentation session setup
func (ps *pgstore) GetDevEUIsWithFragmentationSessionSetup(ctx context.Context, id *uuid.UUID, fragIdx int) ([]lorawan.EUI64, error) {
	var devEUIs []lorawan.EUI64
	err := sqlx.SelectContext(ctx, ps.db, &devEUIs, `
		select
			rms.dev_eui
		from
			remote_multicast_setup rms
		inner join
			remote_fragmentation_session rfs
		on
			rfs.dev_eui = rms.dev_eui
			and rfs.frag_index = $1
		where
			rms.multicast_group_id = $2
			and rms.state = $3
			and rms.state_provisioned = $4
			and rfs.state = $3
			and rms.state_provisioned = $4`,
		fragIdx,
		id,
		store.RemoteMulticastSetupSetup,
		true,
	)
	if err != nil {
		return nil, handlePSQLError(Select, err, "select error")
	}

	return devEUIs, nil
}
