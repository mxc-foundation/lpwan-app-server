package pgstore

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/mxc-foundation/lpwan-app-server/internal/logging"

	appmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/application"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/device"
)

type ApplicationHandler struct {
	tx *sqlx.Tx
	db *sqlx.DB
}

func New(tx *sqlx.Tx, db *sqlx.DB) *ApplicationHandler {
	applicationHandler = ApplicationHandler{
		tx: tx,
		db: db,
	}
	return &applicationHandler
}

var applicationHandler ApplicationHandler

func Handler() *ApplicationHandler {
	return &applicationHandler
}

// CheckCreateApplicationAccess validate validates if the client has access to the applications resource.
func (h *ApplicationHandler) CheckCreateApplicationAccess(username string, userID, organizationID int64) (bool, error) {
	// global admin
	// organization admin
	// organization device admin
	var userWhere = [][]string{
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "o.id = $2", "ou.is_admin = true"},
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "o.id = $2", "ou.is_device_admin = true"},
	}

	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join organization o
			on o.id = ou.organization_id
		left join application a
			on a.organization_id = o.id
	`

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.Get(h.db, &count, userQuery, username, organizationID, userID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

// CheckListApplicationAccess :
func (h *ApplicationHandler) CheckListApplicationAccess(username string, userID, organizationID int64) (bool, error) {
	// global admin
	// organization user (when organization id is given)
	// any active user (api will filter on user)
	var userWhere = [][]string{
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "$2 > 0", "o.id = $2 or a.organization_id = $2"},
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "$2 = 0"},
	}

	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join organization o
			on o.id = ou.organization_id
		left join application a
			on a.organization_id = o.id
	`

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.Get(h.db, &count, userQuery, username, organizationID, userID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

// CheckReadApplicationAccess :
func (h *ApplicationHandler) CheckReadApplicationAccess(username string, userID, applicationID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join organization o
			on o.id = ou.organization_id
		left join application a
			on a.organization_id = o.id
	`
	// global admin
	// organization user
	var userWhere = [][]string{
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "a.id = $2"},
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.Get(h.db, &count, userQuery, username, applicationID, userID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

// CheckUpdateApplicationAccess :
func (h *ApplicationHandler) CheckUpdateApplicationAccess(username string, userID, applicationID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join organization o
			on o.id = ou.organization_id
		left join application a
			on a.organization_id = o.id
	`
	// global admin
	// organization admin
	// organization device admin
	var userWhere = [][]string{
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "ou.is_admin = true", "a.id = $2"},
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "ou.is_device_admin = true", "a.id = $2"},
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.Get(h.db, &count, userQuery, username, applicationID, userID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

// CheckDeleteApplicationAccess :
func (h *ApplicationHandler) CheckDeleteApplicationAccess(username string, userID, applicationID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join organization o
			on o.id = ou.organization_id
		left join application a
			on a.organization_id = o.id
	`
	// global admin
	// organization admin
	// organization device admin
	var userWhere = [][]string{
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "ou.is_admin = true", "a.id = $2"},
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "ou.is_device_admin = true", "a.id = $2"},
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.Get(h.db, &count, userQuery, username, applicationID, userID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

// CreateApplication creates the given Application.
func (h *ApplicationHandler) CreateApplication(ctx context.Context, item *appmod.Application) error {
	if err := item.Validate(); err != nil {
		return errors.Wrap(err, "validate error")
	}

	err := sqlx.Get(h.tx, &item.ID, `
		insert into application (
			name,
			description,
			organization_id,
			service_profile_id,
			payload_codec,
			payload_encoder_script,
			payload_decoder_script
		) values ($1, $2, $3, $4, $5, $6, $7) returning id`,
		item.Name,
		item.Description,
		item.OrganizationID,
		item.ServiceProfileID,
		item.PayloadCodec,
		item.PayloadEncoderScript,
		item.PayloadDecoderScript,
	)
	if err != nil {
		return errors.Wrap(err, "insert error")
	}

	log.WithFields(log.Fields{
		"id":     item.ID,
		"name":   item.Name,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("application created")

	return nil
}

// GetApplication returns the Application for the given id.
func (h *ApplicationHandler) GetApplication(ctx context.Context, id int64) (appmod.Application, error) {
	var app appmod.Application
	err := sqlx.Get(h.db, &app, "select * from application where id = $1", id)
	if err != nil {
		return app, errors.Wrap(err, "select error")
	}

	return app, nil
}

// GetApplicationCount returns the total number of applications.
func (h *ApplicationHandler) GetApplicationCount(ctx context.Context, filters appmod.ApplicationFilters) (int, error) {
	if filters.Search != "" {
		filters.Search = "%" + filters.Search + "%"
	}

	query, args, err := sqlx.BindNamed(sqlx.DOLLAR, `
		select
			count(distinct a.*)
		from
			application a
		left join organization_user ou
			on a.organization_id = ou.organization_id
		left join "user" u
			on ou.user_id = u.id
	`+filters.SQL(), filters)
	if err != nil {
		return 0, errors.Wrap(err, "named query error")
	}

	var count int
	err = sqlx.Get(h.db, &count, query, args...)
	if err != nil {
		return 0, errors.Wrap(err, "select error")
	}

	return count, nil
}

// GetApplications returns a slice of applications, sorted by name and
// respecting the given limit and offset.
func (h *ApplicationHandler) GetApplications(ctx context.Context, filters appmod.ApplicationFilters) ([]appmod.ApplicationListItem, error) {
	if filters.Search != "" {
		filters.Search = "%" + filters.Search + "%"
	}

	query, args, err := sqlx.BindNamed(sqlx.DOLLAR, `
		select
			a.*,
			sp.name as service_profile_name
		from
			application a
		inner join service_profile sp
			on a.service_profile_id = sp.service_profile_id
		left join organization_user ou
			on a.organization_id = ou.organization_id
		left join "user" u
			on ou.user_id = u.id
	`+filters.SQL()+`
		group by
			a.id,
			sp.name
		order by
			a.name
		limit :limit
		offset :offset
	`, filters)
	if err != nil {
		return nil, errors.Wrap(err, "named query error")
	}

	var apps []appmod.ApplicationListItem
	err = sqlx.Select(h.db, &apps, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "select error")
	}

	return apps, nil
}

// UpdateApplication updates the given Application.
func (h *ApplicationHandler) UpdateApplication(ctx context.Context, item appmod.Application) error {
	if err := item.Validate(); err != nil {
		return fmt.Errorf("validate application error: %s", err)
	}

	res, err := h.tx.Exec(`
		update application
		set
			name = $2,
			description = $3,
			organization_id = $4,
			service_profile_id = $5,
			payload_codec = $6,
			payload_encoder_script = $7,
			payload_decoder_script = $8
		where id = $1`,
		item.ID,
		item.Name,
		item.Description,
		item.OrganizationID,
		item.ServiceProfileID,
		item.PayloadCodec,
		item.PayloadEncoderScript,
		item.PayloadDecoderScript,
	)
	if err != nil {
		return errors.Wrap(err, "update error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return errors.New("ErrDoesNotExist")
	}

	log.WithFields(log.Fields{
		"id":     item.ID,
		"name":   item.Name,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("application updated")

	return nil
}

// DeleteApplication deletes the Application matching the given ID.
func (h *ApplicationHandler) DeleteApplication(ctx context.Context, id int64) error {
	err := device.GetDeviceAPI().Store.DeleteAllDevicesForApplicationID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "delete all nodes error")
	}

	res, err := h.tx.Exec("delete from application where id = $1", id)
	if err != nil {
		return errors.Wrap(err, "delete error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return errors.New("ErrDoesNotExist")
	}

	log.WithFields(log.Fields{
		"id":     id,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("application deleted")

	return nil
}

// DeleteAllApplicationsForOrganizationID deletes all applications
// given an organization id.
func (h *ApplicationHandler) DeleteAllApplicationsForOrganizationID(ctx context.Context, organizationID int64) error {
	var apps []appmod.Application
	err := sqlx.Select(h.db, &apps, "select * from application where organization_id = $1", organizationID)
	if err != nil {
		return errors.Wrap(err, "select error")
	}

	for _, app := range apps {
		err = h.DeleteApplication(ctx, app.ID)
		if err != nil {
			return errors.Wrap(err, "delete application error")
		}
	}

	return nil
}
