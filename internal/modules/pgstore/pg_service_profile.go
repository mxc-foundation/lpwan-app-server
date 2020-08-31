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

// DeleteAllServiceProfilesForOrganizationID deletes all service-profiles
// given an organization id.
func (ps *pgstore) DeleteAllServiceProfilesForOrganizationID(ctx context.Context, organizationID int64) error {
	var sps []store.ServiceProfileMeta
	err := sqlx.SelectContext(ctx, ps.db, &sps, "select * from service_profile where organization_id = $1", organizationID)
	if err != nil {
		return errors.Wrap(err, "select error")
	}

	for _, sp := range sps {
		err = ps.DeleteServiceProfile(ctx, sp.ServiceProfileID)
		if err != nil {
			return errors.Wrap(err, "delete service-profile error")
		}
	}

	return nil
}

// DeleteServiceProfile deletes the service-profile matching the given id.
func (ps *pgstore) DeleteServiceProfile(ctx context.Context, id uuid.UUID) error {
	n, err := ps.GetNetworkServerForServiceProfileID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "get network-server error")
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

	res, err := ps.db.ExecContext(ctx, "delete from service_profile where service_profile_id = $1", id)
	if err != nil {
		return errors.Wrap(err, "select error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return errors.New("not exist")
	}

	_, err = nsClient.DeleteServiceProfile(ctx, &ns.DeleteServiceProfileRequest{
		Id: id.Bytes(),
	})
	if err != nil && status.Code(err) != codes.NotFound {
		return errors.Wrap(err, "delete service-profile error")
	}

	log.WithFields(log.Fields{
		"id":     id,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("service-profile deleted")

	return nil
}
