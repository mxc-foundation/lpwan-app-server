package pgstore

import (
	"context"
	"strings"

	"github.com/brocaar/lorawan"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

func (ps *pgstore) CheckReadFUOTADeploymentAccess(ctx context.Context, username string, id uuid.UUID, userID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join application a
			on a.organization_id = ou.organization_id
		left join device d
			on a.id = d.application_id
		left join fuota_deployment_device fdd
			on d.dev_eui = fdd.dev_eui
		left join fuota_deployment fd
			on fdd.fuota_deployment_id = fd.id
	`
	// global admin
	// organization user
	userWhere := [][]string{
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "fd.id = $2"},
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.GetContext(ctx, ps.db, &count, userQuery, username, id, userID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil

}

func (ps *pgstore) CheckCreateFUOTADeploymentsAccess(ctx context.Context, username string, applicationID int64, devEUI lorawan.EUI64, userID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join application a
			on a.organization_id = ou.organization_id
		left join device d
			on a.id = d.application_id
	`
	// global admin
	// organization admin
	userWhere := [][]string{
		{"(u.email = $1 or u.id = $4)", "u.is_active = true", "u.is_admin = true"},
		{"(u.email = $1 or u.id = $4)", "u.is_active = true", "ou.is_admin = true", "$2 > 0", "a.id = $2"},
		{"(u.email = $1 or u.id = $4)", "u.is_active = true", "ou.is_admin = true", "$2 = 0", "d.dev_eui = $3"},
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.GetContext(ctx, ps.db, &count, userQuery, username, applicationID, devEUI, userID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

func (ps *pgstore) GetDeviceKeysFromFuotaDevelopmentDevice(ctx context.Context, id uuid.UUID) ([]store.DeviceKeys, error) {
	// query all device-keys that relate to this FUOTA deployment
	var deviceKeys []store.DeviceKeys
	err := sqlx.SelectContext(ctx, ps.db, &deviceKeys, `
		select
			dk.*
		from
			fuota_deployment_device dd
		inner join
			device_keys dk
			on dd.dev_eui = dk.dev_eui
		where
			dd.fuota_deployment_id = $1`,
		id,
	)
	if err != nil {
		return nil, handlePSQLError(Select, err, "select error")
	}

	return deviceKeys, nil
}
