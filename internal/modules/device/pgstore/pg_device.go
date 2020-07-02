package storage

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/brocaar/lorawan"

	m2m_api "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"

	"github.com/mxc-foundation/lpwan-app-server/internal/backend/m2m_client"
	nsClient "github.com/mxc-foundation/lpwan-app-server/internal/backend/networkserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/logging"
	devmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/device"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/application"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/networkserver"
)

type DeviceHandler struct {
	tx *sqlx.Tx
	db *sqlx.DB
}

func New(tx *sqlx.Tx, db *sqlx.DB) *DeviceHandler {
	deviceHandler = DeviceHandler{
		tx: tx,
		db: db,
	}
	return &deviceHandler
}

var deviceHandler DeviceHandler

func Handler() *DeviceHandler {
	return &deviceHandler
}

// UpdateDeviceActivation updates the device address and the AppSKey.
func (h *DeviceHandler) UpdateDeviceActivation(ctx context.Context, devEUI lorawan.EUI64, devAddr lorawan.DevAddr, appSKey lorawan.AES128Key) error {
	res, err := h.tx.Exec(`
		update device
		set
			dev_addr = $2,
			app_s_key = $3
		where
			dev_eui = $1`,
		devEUI[:],
		devAddr[:],
		appSKey[:],
	)
	if err != nil {
		return errors.Wrap(err, "update last-seen and dr error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return errors.New("ErrDoesNotExist")
	}

	log.WithFields(log.Fields{
		"dev_eui":  devEUI,
		"dev_addr": devAddr,
		"ctx_id":   ctx.Value(logging.ContextIDKey),
	}).Info("device activation updated")

	return nil
}

// CreateDevice creates the given device.
func (h *DeviceHandler) CreateDevice(ctx context.Context, d *devmod.Device, applicationServerID uuid.UUID) error {
	if err := d.Validate(); err != nil {
		return errors.Wrap(err, "validate error")
	}

	now := time.Now()
	d.CreatedAt = now
	timestampCreatedAt, _ := ptypes.TimestampProto(d.CreatedAt)

	d.UpdatedAt = now

	_, err := h.tx.Exec(`
        insert into device (
            dev_eui,
            created_at,
            updated_at,
            application_id,
            device_profile_id,
            name,
			description,
			device_status_battery,
			device_status_margin,
			device_status_external_power_source,
			last_seen_at,
			latitude,
			longitude,
			altitude,
			dr,
			variables,
			tags
        ) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`,
		d.DevEUI[:],
		d.CreatedAt,
		d.UpdatedAt,
		d.ApplicationID,
		d.DeviceProfileID,
		d.Name,
		d.Description,
		d.DeviceStatusBattery,
		d.DeviceStatusMargin,
		d.DeviceStatusExternalPower,
		d.LastSeenAt,
		d.Latitude,
		d.Longitude,
		d.Altitude,
		d.DR,
		d.Variables,
		d.Tags,
	)
	if err != nil {
		return errors.Wrap(err, "insert error")
	}

	app, err := application.GetApplicationAPI().Store.GetApplication(ctx, d.ApplicationID)
	if err != nil {
		return errors.Wrap(err, "get application error")
	}

	n, err := networkserver.GetNetworkServerAPI().Store.GetNetworkServerForDevEUI(ctx, d.DevEUI)
	if err != nil {
		return errors.Wrap(err, "get network-server error")
	}

	// add this device to network server
	client, err := nsClient.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err != nil {
		return errors.Wrap(err, "get network-server client error")
	}

	_, err = client.CreateDevice(ctx, &ns.CreateDeviceRequest{
		Device: &ns.Device{
			DevEui:            d.DevEUI[:],
			DeviceProfileId:   d.DeviceProfileID.Bytes(),
			ServiceProfileId:  app.ServiceProfileID.Bytes(),
			RoutingProfileId:  applicationServerID.Bytes(),
			SkipFCntCheck:     d.SkipFCntCheck,
			ReferenceAltitude: d.ReferenceAltitude,
		},
	})
	if err != nil {
		return errors.Wrap(err, "create device error")
	}

	// add this device to m2m server, this procedure should not block insert device into appserver once it's added to
	// network server successfully
	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	dvClient := m2m_api.NewM2MServerServiceClient(m2mClient)
	if err == nil {
		_, err = dvClient.AddDeviceInM2MServer(context.Background(), &m2m_api.AddDeviceInM2MServerRequest{
			OrgId: app.OrganizationID,
			DevProfile: &m2m_api.AppServerDeviceProfile{
				DevEui:        d.DevEUI.String(),
				ApplicationId: d.ApplicationID,
				Name:          d.Name,
				CreatedAt:     timestampCreatedAt,
			},
		})
		if err != nil {
			log.WithError(err).Error("m2m server create device api error")
		}
	} else {
		log.WithError(err).Error("get m2m-server client error")
	}

	log.WithFields(log.Fields{
		"dev_eui": d.DevEUI,
		"ctx_id":  ctx.Value(logging.ContextIDKey),
	}).Info("device created")

	return nil
}

// GetDevice returns the device matching the given DevEUI.
// When forUpdate is set to true, then tx must be a tx transaction.
// When localOnly is set to true, no call to the network-server is made to
// retrieve additional device data.
func (h *DeviceHandler) GetDevice(ctx context.Context, devEUI lorawan.EUI64, forUpdate, localOnly bool) (devmod.Device, error) {
	var fu string
	if forUpdate {
		fu = " for update"
	}

	var d devmod.Device
	err := sqlx.Get(h.db, &d, "select * from device where dev_eui = $1"+fu, devEUI[:])
	if err != nil {
		return d, errors.Wrap(err, "select error")
	}

	if localOnly {
		return d, nil
	}

	n, err := networkserver.GetNetworkServerAPI().Store.GetNetworkServerForDevEUI(ctx, d.DevEUI)
	if err != nil {
		return d, errors.Wrap(err, "get network-server error")
	}

	client, err := nsClient.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err != nil {
		return d, errors.Wrap(err, "get network-server client error")
	}

	resp, err := client.GetDevice(ctx, &ns.GetDeviceRequest{
		DevEui: d.DevEUI[:],
	})
	if err != nil {
		return d, err
	}

	if resp.Device != nil {
		d.SkipFCntCheck = resp.Device.SkipFCntCheck
		d.ReferenceAltitude = resp.Device.ReferenceAltitude
	}

	return d, nil
}

// GetDeviceCount returns the number of devices.
func (h *DeviceHandler) GetDeviceCount(ctx context.Context, filters devmod.DeviceFilters) (int, error) {
	if filters.Search != "" {
		filters.Search = "%" + filters.Search + "%"
	}

	query, args, err := sqlx.BindNamed(sqlx.DOLLAR, `
		select
			count(distinct d.*)
		from device d
		inner join application a
			on d.application_id = a.id
		left join device_multicast_group dmg
			on d.dev_eui = dmg.dev_eui
	`+filters.SQL(), filters)
	if err != nil {
		return 0, errors.Wrap(err, "named query error")
	}

	var count int
	err = sqlx.Get(h.db, &count, query, args...)
	if err != nil {
		return 0, errors.Wrap(err, "select query error")
	}

	return count, nil
}

// GetAllDeviceEuis returns a slice of devices.
func (h *DeviceHandler) GetAllDeviceEuis(ctx context.Context) ([]string, error) {
	var devEuiList []string
	var list []lorawan.EUI64
	err := sqlx.Select(h.db, &list, "select dev_eui from device ORDER BY created_at DESC")
	if err != nil {
		return nil, errors.Wrap(err, "select error")
	}

	for _, devEui := range list {
		devEuiList = append(devEuiList, devEui.String())
	}

	return devEuiList, nil
}

// GetDevices returns a slice of devices.
func (h *DeviceHandler) GetDevices(ctx context.Context, filters devmod.DeviceFilters) ([]devmod.DeviceListItem, error) {
	if filters.Search != "" {
		filters.Search = "%" + filters.Search + "%"
	}

	query, args, err := sqlx.BindNamed(sqlx.DOLLAR, `
		select
			distinct d.*,
			dp.name as device_profile_name
		from
			device d
		inner join device_profile dp
			on dp.device_profile_id = d.device_profile_id
		inner join application a
			on d.application_id = a.id
		left join device_multicast_group dmg
			on d.dev_eui = dmg.dev_eui
		`+filters.SQL()+`
		order by
			d.name
		limit :limit
		offset :offset
	`, filters)
	if err != nil {
		return nil, errors.Wrap(err, "named query error")
	}

	var devices []devmod.DeviceListItem
	err = sqlx.Select(h.db, &devices, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "select error")
	}

	return devices, nil
}

// UpdateDevice updates the given device.
// When localOnly is set, it will not update the device on the network-server.
func (h *DeviceHandler) UpdateDevice(ctx context.Context, d *devmod.Device, localOnly bool) error {
	if err := d.Validate(); err != nil {
		return errors.Wrap(err, "validate error")
	}

	d.UpdatedAt = time.Now()

	res, err := h.tx.Exec(`
        update device
        set
            updated_at = $2,
            application_id = $3,
            device_profile_id = $4,
            name = $5,
			description = $6,
			device_status_battery = $7,
			device_status_margin = $8,
			last_seen_at = $9,
			latitude = $10,
			longitude = $11,
			altitude = $12,
			device_status_external_power_source = $13,
			dr = $14,
			variables = $15,
			tags = $16
        where
            dev_eui = $1`,
		d.DevEUI[:],
		d.UpdatedAt,
		d.ApplicationID,
		d.DeviceProfileID,
		d.Name,
		d.Description,
		d.DeviceStatusBattery,
		d.DeviceStatusMargin,
		d.LastSeenAt,
		d.Latitude,
		d.Longitude,
		d.Altitude,
		d.DeviceStatusExternalPower,
		d.DR,
		d.Variables,
		d.Tags,
	)
	if err != nil {
		return errors.Wrap(err, "update error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return errors.New("not exist")
	}

	// update the device on the network-server
	if !localOnly {
		app, err := application.GetApplicationAPI().Store.GetApplication(ctx, d.ApplicationID)
		if err != nil {
			return errors.Wrap(err, "get application error")
		}

		n, err := networkserver.GetNetworkServerAPI().Store.GetNetworkServerForDevEUI(ctx, d.DevEUI)
		if err != nil {
			return errors.Wrap(err, "get network-server error")
		}

		client, err := nsClient.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
		if err != nil {
			return errors.Wrap(err, "get network-server client error")
		}

		rpID, err := uuid.FromString(config.C.ApplicationServer.ID)
		if err != nil {
			return errors.Wrap(err, "uuid from string error")
		}

		_, err = client.UpdateDevice(ctx, &ns.UpdateDeviceRequest{
			Device: &ns.Device{
				DevEui:            d.DevEUI[:],
				DeviceProfileId:   d.DeviceProfileID.Bytes(),
				ServiceProfileId:  app.ServiceProfileID.Bytes(),
				RoutingProfileId:  rpID.Bytes(),
				SkipFCntCheck:     d.SkipFCntCheck,
				ReferenceAltitude: d.ReferenceAltitude,
			},
		})
		if err != nil {
			return errors.Wrap(err, "update device error")
		}
	}

	log.WithFields(log.Fields{
		"dev_eui": d.DevEUI,
		"ctx_id":  ctx.Value(logging.ContextIDKey),
	}).Info("device updated")

	return nil
}

// DeleteDevice deletes the device matching the given DevEUI.
func (h *DeviceHandler) DeleteDevice(ctx context.Context, devEUI lorawan.EUI64) error {
	n, err := networkserver.GetNetworkServerAPI().Store.GetNetworkServerForDevEUI(ctx, devEUI)
	if err != nil {
		return errors.Wrap(err, "get network-server error")
	}

	res, err := h.tx.Exec("delete from device where dev_eui = $1", devEUI[:])
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

	// delete device from networkserver
	client, err := nsClient.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err != nil {
		return errors.Wrap(err, "get network-server client error")
	}

	_, err = client.DeleteDevice(ctx, &ns.DeleteDeviceRequest{
		DevEui: devEUI[:],
	})
	if err != nil && grpc.Code(err) != codes.NotFound {
		return errors.Wrap(err, "delete device error")
	}

	// delete device from m2m server, this procedure should not block delete device from appserver once it's deleted from
	// network server successfully
	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	dvClient := m2m_api.NewM2MServerServiceClient(m2mClient)
	if err == nil {
		_, err = dvClient.DeleteDeviceInM2MServer(context.Background(), &m2m_api.DeleteDeviceInM2MServerRequest{
			DevEui: devEUI.String(),
		})
		if err != nil && grpc.Code(err) != codes.NotFound {
			log.WithError(err).Error("m2m-server delete device api error")
		}
	} else {
		log.WithError(err).Error("get m2m-server client error")
	}

	log.WithFields(log.Fields{
		"dev_eui": devEUI,
		"ctx_id":  ctx.Value(logging.ContextIDKey),
	}).Info("device deleted")

	return nil
}

// CreateDeviceKeys creates the keys for the given device.
func (h *DeviceHandler) CreateDeviceKeys(ctx context.Context, dc *devmod.DeviceKeys) error {
	now := time.Now()
	dc.CreatedAt = now
	dc.UpdatedAt = now

	_, err := h.tx.Exec(`
        insert into device_keys (
            created_at,
            updated_at,
            dev_eui,
			nwk_key,
			app_key,
			join_nonce,
			gen_app_key
        ) values ($1, $2, $3, $4, $5, $6, $7)`,
		dc.CreatedAt,
		dc.UpdatedAt,
		dc.DevEUI[:],
		dc.NwkKey[:],
		dc.AppKey[:],
		dc.JoinNonce,
		dc.GenAppKey[:],
	)
	if err != nil {
		return errors.Wrap(err, "insert error")
	}

	log.WithFields(log.Fields{
		"dev_eui": dc.DevEUI,
		"ctx_id":  ctx.Value(logging.ContextIDKey),
	}).Info("device-keys created")

	return nil
}

// GetDeviceKeys returns the device-keys for the given DevEUI.
func (h *DeviceHandler) GetDeviceKeys(ctx context.Context, devEUI lorawan.EUI64) (devmod.DeviceKeys, error) {
	var dc devmod.DeviceKeys

	err := sqlx.Get(h.db, &dc, "select * from device_keys where dev_eui = $1", devEUI[:])
	if err != nil {
		return dc, errors.Wrap(err, "select error")
	}

	return dc, nil
}

// UpdateDeviceKeys updates the given device-keys.
func (h *DeviceHandler) UpdateDeviceKeys(ctx context.Context, dc *devmod.DeviceKeys) error {
	dc.UpdatedAt = time.Now()

	res, err := h.tx.Exec(`
        update device_keys
        set
            updated_at = $2,
			nwk_key = $3,
			app_key = $4,
			join_nonce = $5,
			gen_app_key = $6
        where
            dev_eui = $1`,
		dc.DevEUI[:],
		dc.UpdatedAt,
		dc.NwkKey[:],
		dc.AppKey[:],
		dc.JoinNonce,
		dc.GenAppKey[:],
	)
	if err != nil {
		return errors.Wrap(err, "update error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return errors.New("not exist")
	}

	log.WithFields(log.Fields{
		"dev_eui": dc.DevEUI,
		"ctx_id":  ctx.Value(logging.ContextIDKey),
	}).Info("device-keys updated")

	return nil
}

// DeleteDeviceKeys deletes the device-keys for the given DevEUI.
func (h *DeviceHandler) DeleteDeviceKeys(ctx context.Context, devEUI lorawan.EUI64) error {
	res, err := h.tx.Exec("delete from device_keys where dev_eui = $1", devEUI[:])
	if err != nil {
		return errors.Wrap(err, "delete error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected errro")
	}
	if ra == 0 {
		return errors.New("not exist")
	}

	log.WithFields(log.Fields{
		"dev_eui": devEUI,
		"ctx_id":  ctx.Value(logging.ContextIDKey),
	}).Info("device-keys deleted")

	return nil
}

// CreateDeviceActivation creates the given device-activation.
func (h *DeviceHandler) CreateDeviceActivation(ctx context.Context, da *devmod.DeviceActivation) error {
	da.CreatedAt = time.Now()

	err := sqlx.Get(h.tx, &da.ID, `
        insert into device_activation (
            created_at,
            dev_eui,
            dev_addr,
			app_s_key
        ) values ($1, $2, $3, $4)
        returning id`,
		da.CreatedAt,
		da.DevEUI[:],
		da.DevAddr[:],
		da.AppSKey[:],
	)
	if err != nil {
		return errors.Wrap(err, "insert error")
	}

	log.WithFields(log.Fields{
		"id":      da.ID,
		"dev_eui": da.DevEUI,
		"ctx_id":  ctx.Value(logging.ContextIDKey),
	}).Info("device-activation created")

	return nil
}

// GetLastDeviceActivationForDevEUI returns the most recent device-activation for the given DevEUI.
func (h *DeviceHandler) GetLastDeviceActivationForDevEUI(ctx context.Context, devEUI lorawan.EUI64) (devmod.DeviceActivation, error) {
	var da devmod.DeviceActivation

	err := sqlx.Get(h.db, &da, `
        select *
        from device_activation
        where
            dev_eui = $1
        order by
            created_at desc
        limit 1`,
		devEUI[:],
	)
	if err != nil {
		return da, errors.Wrap(err, "select error")
	}

	return da, nil
}

// DeleteAllDevicesForApplicationID deletes all devices given an application id.
func (h *DeviceHandler) DeleteAllDevicesForApplicationID(ctx context.Context, applicationID int64) error {
	var devs []devmod.Device
	err := sqlx.Select(h.db, &devs, "select * from device where application_id = $1", applicationID)
	if err != nil {
		return errors.Wrap(err, "select error")
	}

	for _, dev := range devs {
		err = h.DeleteDevice(ctx, dev.DevEUI)
		if err != nil {
			return errors.Wrap(err, "delete device error")
		}
	}

	return nil
}
