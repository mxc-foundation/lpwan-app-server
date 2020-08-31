package pgstore

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/brocaar/chirpstack-api/go/v3/ns"

	nscli "github.com/mxc-foundation/lpwan-app-server/internal/clients/networkserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/logging"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

// DeleteAllDeviceProfilesForOrganizationID deletes all device-profiles
// given an organization id.
func (pg *pgstore) DeleteAllDeviceProfilesForOrganizationID(ctx context.Context, organizationID int64) error {
	var dps []store.DeviceProfileMeta
	err := sqlx.SelectContext(ctx, pg.db, &dps, `
		select
			device_profile_id,
			network_server_id,
			organization_id,
			created_at,
			updated_at,
			name
		from
			device_profile
		where
			organization_id = $1`,
		organizationID)
	if err != nil {
		return errors.Wrap(err, "select error")
	}

	for _, dp := range dps {
		err = pg.DeleteDeviceProfile(ctx, dp.DeviceProfileID)
		if err != nil {
			return errors.Wrap(err, "delete device-profile error")
		}
	}

	return nil
}

// DeleteDeviceProfile deletes the device-profile matching the given id.
func (pg *pgstore) DeleteDeviceProfile(ctx context.Context, id uuid.UUID) error {
	n, err := pg.GetNetworkServerForDeviceProfileID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "get network-server error")
	}

	res, err := pg.db.ExecContext(ctx, "delete from device_profile where device_profile_id = $1", id)
	if err != nil {
		return errors.Wrap(err, "delete error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return errors.New("not exist")
	}

	nsStruct := nscli.NSStruct{
		Server:  n.Server,
		CACert:  n.CACert,
		TLSCert: n.TLSCert,
		TLSKey:  n.TLSKey,
	}
	nsClient, err := nsStruct.GetNetworkServiceClient()
	if err != nil {
		return errors.Wrap(err, "get network-server client error")
	}

	_, err = nsClient.DeleteDeviceProfile(ctx, &ns.DeleteDeviceProfileRequest{
		Id: id.Bytes(),
	})
	if err != nil && status.Code(err) != codes.NotFound {
		return errors.Wrap(err, "delete device-profile error")
	}

	log.WithFields(log.Fields{
		"id":     id,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("device-profile deleted")

	return nil
}
