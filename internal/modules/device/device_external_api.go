package device

import (
	"database/sql"
	"encoding/json"
	"strconv"

	"google.golang.org/grpc/status"

	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq/hstore"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/brocaar/chirpstack-api/go/v3/common"
	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/brocaar/lorawan"

	pb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	m2mServer "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/m2m_client"
	nsClient "github.com/mxc-foundation/lpwan-app-server/internal/backend/networkserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/eventlog"
	"github.com/mxc-foundation/lpwan-app-server/internal/logging"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/application"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/networkserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/organization"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/user"
)

type DeviceStore interface {
	CreateDevice(ctx context.Context, d *Device, applicationServerID uuid.UUID) error
	GetDevice(ctx context.Context, devEUI lorawan.EUI64, forUpdate, localOnly bool) (Device, error)
	GetDeviceCount(ctx context.Context, filters DeviceFilters) (int, error)
	GetAllDeviceEuis(ctx context.Context) ([]string, error)
	GetDevices(ctx context.Context, filters DeviceFilters) ([]DeviceListItem, error)
	UpdateDevice(ctx context.Context, d *Device, localOnly bool) error
	DeleteDevice(ctx context.Context, devEUI lorawan.EUI64) error
	CreateDeviceKeys(ctx context.Context, dc *DeviceKeys) error
	GetDeviceKeys(ctx context.Context, devEUI lorawan.EUI64) (DeviceKeys, error)
	UpdateDeviceKeys(ctx context.Context, dc *DeviceKeys) error
	DeleteDeviceKeys(ctx context.Context, devEUI lorawan.EUI64) error
	CreateDeviceActivation(ctx context.Context, da *DeviceActivation) error
	GetLastDeviceActivationForDevEUI(ctx context.Context, devEUI lorawan.EUI64) (DeviceActivation, error)
	DeleteAllDevicesForApplicationID(ctx context.Context, applicationID int64) error
	UpdateDeviceActivation(ctx context.Context, devEUI lorawan.EUI64, devAddr lorawan.DevAddr, appSKey lorawan.AES128Key) error
}

// DeviceAPI exports the Node related functions.
type DeviceAPI struct {
	Validator            *validator
	Store                DeviceStore
	AppplicationServerID uuid.UUID
}

// NewDeviceAPI creates a new NodeAPI.
func NewDeviceAPI(api DeviceAPI) *DeviceAPI {
	deviceAPI = DeviceAPI{
		Validator:            api.Validator,
		Store:                api.Store,
		AppplicationServerID: api.AppplicationServerID,
	}

	return &deviceAPI
}

var (
	deviceAPI DeviceAPI
)

func GetDeviceAPI() *DeviceAPI {
	return &deviceAPI
}

// Create creates the given device.
func (a *DeviceAPI) Create(ctx context.Context, req *pb.CreateDeviceRequest) (*empty.Empty, error) {
	if req.Device == nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "device must not be nil")
	}

	var devEUI lorawan.EUI64
	if err := devEUI.UnmarshalText([]byte(req.Device.DevEui)); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}

	dpID, err := uuid.FromString(req.Device.DeviceProfileId)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
		validateNodesAccess(req.Device.ApplicationId, Create)); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	// if Name is "", set it to the DevEUI
	if req.Device.Name == "" {
		req.Device.Name = req.Device.DevEui
	}

	// Validate that application and device-profile are under the same
	// organization ID.
	app, err := application.GetApplicationAPI().Store.GetApplication(ctx, req.Device.ApplicationId)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	dp, err := storage.GetDeviceProfile(ctx, storage.DB(), dpID, false, true)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	if app.OrganizationID != dp.OrganizationID {
		return nil, grpc.Errorf(codes.InvalidArgument, "device-profile and application must be under the same organization")
	}

	// Set Device struct.
	d := Device{
		DevEUI:            devEUI,
		ApplicationID:     req.Device.ApplicationId,
		DeviceProfileID:   dpID,
		Name:              req.Device.Name,
		Description:       req.Device.Description,
		SkipFCntCheck:     req.Device.SkipFCntCheck,
		ReferenceAltitude: req.Device.ReferenceAltitude,
		Variables: hstore.Hstore{
			Map: make(map[string]sql.NullString),
		},
		Tags: hstore.Hstore{
			Map: make(map[string]sql.NullString),
		},
	}

	for k, v := range req.Device.Variables {
		d.Variables.Map[k] = sql.NullString{String: v, Valid: true}
	}

	for k, v := range req.Device.Tags {
		d.Tags.Map[k] = sql.NullString{String: v, Valid: true}
	}

	// A transaction is needed as:
	//  * A remote gRPC call is performed and in case of error, we want to
	//    rollback the transaction.
	//  * We want to lock the organization so that we can validate the
	//    max device count.
	err = storage.Transaction(func(tx sqlx.Ext) error {
		org, err := organization.GetOrganizationAPI().Store.GetOrganization(ctx, app.OrganizationID, true)
		if err != nil {
			return err
		}

		// Validate max. device count when != 0.
		if org.MaxDeviceCount != 0 {
			count, err := a.Store.GetDeviceCount(ctx, DeviceFilters{ApplicationID: app.OrganizationID})
			if err != nil {
				return err
			}

			if count >= org.MaxDeviceCount {
				return storage.ErrOrganizationMaxDeviceCount
			}
		}

		return a.Store.CreateDevice(ctx, &d, a.AppplicationServerID)
	})
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// Get returns the device matching the given DevEUI.
func (a *DeviceAPI) Get(ctx context.Context, req *pb.GetDeviceRequest) (*pb.GetDeviceResponse, error) {
	var eui lorawan.EUI64
	if err := eui.UnmarshalText([]byte(req.DevEui)); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
		validateNodeAccess(eui, Read)); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	d, err := a.Store.GetDevice(ctx, eui, false, false)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := pb.GetDeviceResponse{
		Device: &pb.Device{
			DevEui:            d.DevEUI.String(),
			Name:              d.Name,
			ApplicationId:     d.ApplicationID,
			Description:       d.Description,
			DeviceProfileId:   d.DeviceProfileID.String(),
			SkipFCntCheck:     d.SkipFCntCheck,
			ReferenceAltitude: d.ReferenceAltitude,
			Variables:         make(map[string]string),
			Tags:              make(map[string]string),
		},

		DeviceStatusBattery: 256,
		DeviceStatusMargin:  256,
	}

	if d.DeviceStatusBattery != nil {
		resp.DeviceStatusBattery = uint32(*d.DeviceStatusBattery)
	}
	if d.DeviceStatusMargin != nil {
		resp.DeviceStatusMargin = int32(*d.DeviceStatusMargin)
	}
	if d.LastSeenAt != nil {
		resp.LastSeenAt, err = ptypes.TimestampProto(*d.LastSeenAt)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
	}

	if d.Latitude != nil && d.Longitude != nil && d.Altitude != nil {
		resp.Location = &common.Location{
			Latitude:  *d.Latitude,
			Longitude: *d.Longitude,
			Altitude:  *d.Altitude,
		}
	}

	for k, v := range d.Variables.Map {
		if v.Valid {
			resp.Device.Variables[k] = v.String
		}
	}

	for k, v := range d.Tags.Map {
		if v.Valid {
			resp.Device.Tags[k] = v.String
		}
	}

	return &resp, nil
}

// List lists the available applications.
func (a *DeviceAPI) List(ctx context.Context, req *pb.ListDeviceRequest) (*pb.ListDeviceResponse, error) {
	var err error
	var idFilter bool

	filters := DeviceFilters{
		ApplicationID: req.ApplicationId,
		Search:        req.Search,
		/*		Tags: hstore.Hstore{
				Map: make(map[string]sql.NullString),
			},*/
		Limit:  int(req.Limit),
		Offset: int(req.Offset),
	}

	if req.MulticastGroupId != "" {
		filters.MulticastGroupID, err = uuid.FromString(req.MulticastGroupId)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
	}

	if req.ServiceProfileId != "" {
		filters.ServiceProfileID, err = uuid.FromString(req.ServiceProfileId)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
	}

	/*	for k, v := range req.Tags {
		filters.Tags.Map[k] = sql.NullString{String: v, Valid: true}
	}*/

	if filters.ApplicationID != 0 {
		idFilter = true

		// validate that the client has access to the given application
		if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
			application.ValidateApplicationAccess(req.ApplicationId, application.Read),
		); err != nil {
			return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
		}

	}

	if filters.MulticastGroupID != uuid.Nil {
		idFilter = true

		// validate that the client has access to the given multicast-group
		if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
			authcus.ValidateMulticastGroupAccess(authcus.Read, filters.MulticastGroupID),
		); err != nil {
			return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
		}
	}

	if filters.ServiceProfileID != uuid.Nil {
		idFilter = true

		// validate that the client has access to the given service-profile
		if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
			authcus.ValidateServiceProfileAccess(authcus.Read, filters.ServiceProfileID),
		); err != nil {
			return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
		}
	}

	if !idFilter {
		user, err := user.GetUserAPI().Validator.GetUser(ctx)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}

		if !user.IsAdmin {
			return nil, grpc.Errorf(codes.Unauthenticated, "client must be global admin for unfiltered request")
		}
	}

	count, err := a.Store.GetDeviceCount(ctx, filters)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	devices, err := a.Store.GetDevices(ctx, filters)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return a.returnList(count, devices)
}

// Update updates the device matching the given DevEUI.
func (a *DeviceAPI) Update(ctx context.Context, req *pb.UpdateDeviceRequest) (*empty.Empty, error) {
	if req.Device == nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "device must not be nil")
	}

	var devEUI lorawan.EUI64
	if err := devEUI.UnmarshalText([]byte(req.Device.DevEui)); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}

	dpID, err := uuid.FromString(req.Device.DeviceProfileId)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
		validateNodeAccess(devEUI, Update)); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	app, err := application.GetApplicationAPI().Store.GetApplication(ctx, req.Device.ApplicationId)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	dp, err := storage.GetDeviceProfile(ctx, storage.DB(), dpID, false, true)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	if app.OrganizationID != dp.OrganizationID {
		return nil, grpc.Errorf(codes.InvalidArgument, "device-profile and application must be under the same organization")
	}

	err = storage.Transaction(func(tx sqlx.Ext) error {
		d, err := a.Store.GetDevice(ctx, devEUI, true, false)
		if err != nil {
			return helpers.ErrToRPCError(err)
		}

		// If the device is moved to a different application, validate that
		// the new application is assigned to the same service-profile.
		// This to guarantee that the new application is still on the same
		// network-server and is not assigned to a different organization.
		if req.Device.ApplicationId != d.ApplicationID {
			appOld, err := application.GetApplicationAPI().Store.GetApplication(ctx, d.ApplicationID)
			if err != nil {
				return helpers.ErrToRPCError(err)
			}

			appNew, err := application.GetApplicationAPI().Store.GetApplication(ctx, req.Device.ApplicationId)
			if err != nil {
				return helpers.ErrToRPCError(err)
			}

			if appOld.ServiceProfileID != appNew.ServiceProfileID {
				return grpc.Errorf(codes.InvalidArgument, "when moving a device from application A to B, both A and B must share the same service-profile")
			}
		}

		d.ApplicationID = req.Device.ApplicationId
		d.DeviceProfileID = dpID
		d.Name = req.Device.Name
		d.Description = req.Device.Description
		d.SkipFCntCheck = req.Device.SkipFCntCheck
		d.ReferenceAltitude = req.Device.ReferenceAltitude
		d.Variables = hstore.Hstore{
			Map: make(map[string]sql.NullString),
		}
		d.Tags = hstore.Hstore{
			Map: make(map[string]sql.NullString),
		}

		for k, v := range req.Device.Variables {
			d.Variables.Map[k] = sql.NullString{String: v, Valid: true}
		}

		for k, v := range req.Device.Tags {
			d.Tags.Map[k] = sql.NullString{String: v, Valid: true}
		}

		if err := a.Store.UpdateDevice(ctx, &d, false); err != nil {
			return helpers.ErrToRPCError(err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// Delete deletes the node matching the given name.
func (a *DeviceAPI) Delete(ctx context.Context, req *pb.DeleteDeviceRequest) (*empty.Empty, error) {
	var eui lorawan.EUI64
	if err := eui.UnmarshalText([]byte(req.DevEui)); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
		validateNodeAccess(eui, Delete)); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	// as this also performs a remote call to delete the node from the
	// network-server, wrap it in a transaction
	err := storage.Transaction(func(tx sqlx.Ext) error {
		return a.Store.DeleteDevice(ctx, eui)
	})
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// CreateKeys creates the given device-keys.
func (a *DeviceAPI) CreateKeys(ctx context.Context, req *pb.CreateDeviceKeysRequest) (*empty.Empty, error) {
	if req.DeviceKeys == nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "device_keys must not be nil")
	}

	// appKey is not used for LoRaWAN 1.0
	var appKey lorawan.AES128Key
	if req.DeviceKeys.AppKey != "" {
		if err := appKey.UnmarshalText([]byte(req.DeviceKeys.AppKey)); err != nil {
			return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
		}
	}

	// genAppKey is only for LoRaWAN 1.0 devices that implement the
	// remote multicast setup specification.
	var genAppKey lorawan.AES128Key
	if req.DeviceKeys.GenAppKey != "" {
		if err := genAppKey.UnmarshalText([]byte(req.DeviceKeys.GenAppKey)); err != nil {
			return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
		}
	}

	// nwkKey
	var nwkKey lorawan.AES128Key
	if err := nwkKey.UnmarshalText([]byte(req.DeviceKeys.NwkKey)); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}

	// devEUI
	var eui lorawan.EUI64
	if err := eui.UnmarshalText([]byte(req.DeviceKeys.DevEui)); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
		validateNodeAccess(eui, Update),
	); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	err := a.Store.CreateDeviceKeys(ctx, &DeviceKeys{
		DevEUI:    eui,
		NwkKey:    nwkKey,
		AppKey:    appKey,
		GenAppKey: genAppKey,
	})
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// GetKeys returns the device-keys for the given DevEUI.
func (a *DeviceAPI) GetKeys(ctx context.Context, req *pb.GetDeviceKeysRequest) (*pb.GetDeviceKeysResponse, error) {
	var eui lorawan.EUI64
	if err := eui.UnmarshalText([]byte(req.DevEui)); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
		validateNodeAccess(eui, Update),
	); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	dk, err := a.Store.GetDeviceKeys(ctx, eui)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &pb.GetDeviceKeysResponse{
		DeviceKeys: &pb.DeviceKeys{
			DevEui:    eui.String(),
			AppKey:    dk.AppKey.String(),
			NwkKey:    dk.NwkKey.String(),
			GenAppKey: dk.GenAppKey.String(),
		},
	}, nil
}

// UpdateKeys updates the device-keys.
func (a *DeviceAPI) UpdateKeys(ctx context.Context, req *pb.UpdateDeviceKeysRequest) (*empty.Empty, error) {
	if req.DeviceKeys == nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "device_keys must not be nil")
	}

	var appKey lorawan.AES128Key
	// appKey is not used for LoRaWAN 1.0
	if req.DeviceKeys.AppKey != "" {
		if err := appKey.UnmarshalText([]byte(req.DeviceKeys.AppKey)); err != nil {
			return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
		}
	}

	// genAppKey is only for LoRaWAN 1.0 devices that implement the
	// remote multicast setup specification.
	var genAppKey lorawan.AES128Key
	if req.DeviceKeys.GenAppKey != "" {
		if err := genAppKey.UnmarshalText([]byte(req.DeviceKeys.GenAppKey)); err != nil {
			return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
		}
	}

	var nwkKey lorawan.AES128Key
	if err := nwkKey.UnmarshalText([]byte(req.DeviceKeys.NwkKey)); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}

	var eui lorawan.EUI64
	if err := eui.UnmarshalText([]byte(req.DeviceKeys.DevEui)); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
		validateNodeAccess(eui, Update),
	); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	dk, err := a.Store.GetDeviceKeys(ctx, eui)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	dk.NwkKey = nwkKey
	dk.AppKey = appKey
	dk.GenAppKey = genAppKey

	err = a.Store.UpdateDeviceKeys(ctx, &dk)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// DeleteKeys deletes the device-keys for the given DevEUI.
func (a *DeviceAPI) DeleteKeys(ctx context.Context, req *pb.DeleteDeviceKeysRequest) (*empty.Empty, error) {
	var eui lorawan.EUI64
	if err := eui.UnmarshalText([]byte(req.DevEui)); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
		validateNodeAccess(eui, Delete),
	); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	if err := a.Store.DeleteDeviceKeys(ctx, eui); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// Deactivate de-activates the device.
func (a *DeviceAPI) Deactivate(ctx context.Context, req *pb.DeactivateDeviceRequest) (*empty.Empty, error) {
	var devEUI lorawan.EUI64
	if err := devEUI.UnmarshalText([]byte(req.DevEui)); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
		validateNodeAccess(devEUI, Update),
	); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	d, err := a.Store.GetDevice(ctx, devEUI, false, true)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	n, err := networkserver.GetNetworkServerAPI().Store.GetNetworkServerForDevEUI(ctx, d.DevEUI)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	client, err := nsClient.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	_, _ = client.DeactivateDevice(ctx, &ns.DeactivateDeviceRequest{
		DevEui: d.DevEUI[:],
	})

	return &empty.Empty{}, nil
}

// Activate activates the node (ABP only).
func (a *DeviceAPI) Activate(ctx context.Context, req *pb.ActivateDeviceRequest) (*empty.Empty, error) {
	if req.DeviceActivation == nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "device_activation must not be nil")
	}

	var devAddr lorawan.DevAddr
	var devEUI lorawan.EUI64
	var appSKey lorawan.AES128Key
	var nwkSEncKey lorawan.AES128Key
	var sNwkSIntKey lorawan.AES128Key
	var fNwkSIntKey lorawan.AES128Key

	if err := devAddr.UnmarshalText([]byte(req.DeviceActivation.DevAddr)); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "devAddr: %s", err)
	}
	if err := devEUI.UnmarshalText([]byte(req.DeviceActivation.DevEui)); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "devEUI: %s", err)
	}
	if err := appSKey.UnmarshalText([]byte(req.DeviceActivation.AppSKey)); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "appSKey: %s", err)
	}
	if err := nwkSEncKey.UnmarshalText([]byte(req.DeviceActivation.NwkSEncKey)); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "nwkSEncKey: %s", err)
	}
	if err := sNwkSIntKey.UnmarshalText([]byte(req.DeviceActivation.SNwkSIntKey)); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "sNwkSIntKey: %s", err)
	}
	if err := fNwkSIntKey.UnmarshalText([]byte(req.DeviceActivation.FNwkSIntKey)); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "fNwkSIntKey: %s", err)
	}

	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
		validateNodeAccess(devEUI, Update)); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	d, err := a.Store.GetDevice(ctx, devEUI, false, true)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	n, err := networkserver.GetNetworkServerAPI().Store.GetNetworkServerForDevEUI(ctx, devEUI)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	client, err := nsClient.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	_, _ = client.DeactivateDevice(ctx, &ns.DeactivateDeviceRequest{
		DevEui: d.DevEUI[:],
	})

	actReq := ns.ActivateDeviceRequest{
		DeviceActivation: &ns.DeviceActivation{
			DevEui:      d.DevEUI[:],
			DevAddr:     devAddr[:],
			NwkSEncKey:  nwkSEncKey[:],
			SNwkSIntKey: sNwkSIntKey[:],
			FNwkSIntKey: fNwkSIntKey[:],
			FCntUp:      req.DeviceActivation.FCntUp,
			NFCntDown:   req.DeviceActivation.NFCntDown,
			AFCntDown:   req.DeviceActivation.AFCntDown,
		},
	}

	err = storage.Transaction(func(db sqlx.Ext) error {
		if err := a.Store.UpdateDeviceActivation(ctx, d.DevEUI, devAddr, appSKey); err != nil {
			return helpers.ErrToRPCError(err)
		}

		_, err := client.ActivateDevice(ctx, &actReq)
		if err != nil {
			return helpers.ErrToRPCError(err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"dev_addr": devAddr,
		"dev_eui":  d.DevEUI,
		"ctx_id":   ctx.Value(logging.ContextIDKey),
	}).Info("device activated")

	return &empty.Empty{}, nil
}

// GetActivation returns the device activation for the given DevEUI.
func (a *DeviceAPI) GetActivation(ctx context.Context, req *pb.GetDeviceActivationRequest) (*pb.GetDeviceActivationResponse, error) {
	var devAddr lorawan.DevAddr
	var devEUI lorawan.EUI64
	var sNwkSIntKey lorawan.AES128Key
	var fNwkSIntKey lorawan.AES128Key
	var nwkSEncKey lorawan.AES128Key

	if err := devEUI.UnmarshalText([]byte(req.DevEui)); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "devEUI: %s", err)
	}

	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
		validateNodeAccess(devEUI, Read)); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	d, err := a.Store.GetDevice(ctx, devEUI, false, true)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	n, err := networkserver.GetNetworkServerAPI().Store.GetNetworkServerForDevEUI(ctx, devEUI)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	client, err := nsClient.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	devAct, err := client.GetDeviceActivation(ctx, &ns.GetDeviceActivationRequest{
		DevEui: d.DevEUI[:],
	})
	if err != nil {
		return nil, err
	}

	copy(devAddr[:], devAct.DeviceActivation.DevAddr)
	copy(nwkSEncKey[:], devAct.DeviceActivation.NwkSEncKey)
	copy(sNwkSIntKey[:], devAct.DeviceActivation.SNwkSIntKey)
	copy(fNwkSIntKey[:], devAct.DeviceActivation.FNwkSIntKey)

	return &pb.GetDeviceActivationResponse{
		DeviceActivation: &pb.DeviceActivation{
			DevEui:  d.DevEUI.String(),
			DevAddr: devAddr.String(),
			//AppSKey:     d.AppSKey.String(),
			NwkSEncKey:  nwkSEncKey.String(),
			SNwkSIntKey: sNwkSIntKey.String(),
			FNwkSIntKey: fNwkSIntKey.String(),
			FCntUp:      devAct.DeviceActivation.FCntUp,
			NFCntDown:   devAct.DeviceActivation.NFCntDown,
			AFCntDown:   devAct.DeviceActivation.AFCntDown,
		},
	}, nil
}

// StreamFrameLogs streams the uplink and downlink frame-logs for the given DevEUI.
// Note: these are the raw LoRaWAN frames and this endpoint is intended for debugging.
func (a *DeviceAPI) StreamFrameLogs(req *pb.StreamDeviceFrameLogsRequest, srv pb.DeviceService_StreamFrameLogsServer) error {
	var devEUI lorawan.EUI64

	if err := devEUI.UnmarshalText([]byte(req.DevEui)); err != nil {
		return grpc.Errorf(codes.InvalidArgument, "devEUI: %s", err)
	}

	if err := a.Validator.otpValidator.JwtValidator.Validate(srv.Context(),
		validateNodeAccess(devEUI, Read)); err != nil {
		return grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	n, err := networkserver.GetNetworkServerAPI().Store.GetNetworkServerForDevEUI(srv.Context(), devEUI)
	if err != nil {
		return helpers.ErrToRPCError(err)
	}

	client, err := nsClient.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err != nil {
		return helpers.ErrToRPCError(err)
	}

	streamClient, err := client.StreamFrameLogsForDevice(srv.Context(), &ns.StreamFrameLogsForDeviceRequest{
		DevEui: devEUI[:],
	})
	if err != nil {
		return err
	}

	for {
		resp, err := streamClient.Recv()
		if err != nil {
			if grpc.Code(err) == codes.Canceled {
				return nil
			}

			return err
		}

		up, down, err := ConvertUplinkAndDownlinkFrames(resp.GetUplinkFrameSet(), resp.GetDownlinkFrame(), true)
		if err != nil {
			return helpers.ErrToRPCError(err)
		}

		var frameResp pb.StreamDeviceFrameLogsResponse
		if up != nil {
			frameResp.Frame = &pb.StreamDeviceFrameLogsResponse_UplinkFrame{
				UplinkFrame: up,
			}
		}

		if down != nil {
			frameResp.Frame = &pb.StreamDeviceFrameLogsResponse_DownlinkFrame{
				DownlinkFrame: down,
			}
		}

		err = srv.Send(&frameResp)
		if err != nil {
			return err
		}
	}
}

// StreamEventLogs stream the device events (uplink payloads, ACKs, joins, errors).
// Note: this endpoint is intended for debugging and should not be used for building
// integrations.
func (a *DeviceAPI) StreamEventLogs(req *pb.StreamDeviceEventLogsRequest, srv pb.DeviceService_StreamEventLogsServer) error {
	var devEUI lorawan.EUI64

	if err := devEUI.UnmarshalText([]byte(req.DevEui)); err != nil {
		return grpc.Errorf(codes.InvalidArgument, "devEUI: %s", err)
	}

	if err := a.Validator.otpValidator.JwtValidator.Validate(srv.Context(),
		validateNodeAccess(devEUI, Read)); err != nil {
		return grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	eventLogChan := make(chan eventlog.EventLog)
	go func() {
		err := eventlog.GetEventLogForDevice(srv.Context(), devEUI, eventLogChan)
		if err != nil {
			log.WithError(err).Error("get event-log for device error")
		}
		close(eventLogChan)
	}()

	for el := range eventLogChan {
		b, err := json.Marshal(el.Payload)
		if err != nil {
			return grpc.Errorf(codes.Internal, "marshal json error: %s", err)
		}

		resp := pb.StreamDeviceEventLogsResponse{
			Type:        el.Type,
			PayloadJson: string(b),
		}

		err = srv.Send(&resp)
		if err != nil {
			log.WithError(err).Error("error sending event-log response")
		}
	}

	return nil
}

// GetRandomDevAddr returns a random DevAddr taking the NwkID prefix into account.
func (a *DeviceAPI) GetRandomDevAddr(ctx context.Context, req *pb.GetRandomDevAddrRequest) (*pb.GetRandomDevAddrResponse, error) {
	var devEUI lorawan.EUI64

	if err := devEUI.UnmarshalText([]byte(req.DevEui)); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "devEUI: %s", err)
	}

	n, err := networkserver.GetNetworkServerAPI().Store.GetNetworkServerForDevEUI(ctx, devEUI)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	client, err := nsClient.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp, err := client.GetRandomDevAddr(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}

	var devAddr lorawan.DevAddr
	copy(devAddr[:], resp.DevAddr)

	return &pb.GetRandomDevAddrResponse{
		DevAddr: devAddr.String(),
	}, nil
}

func (a *DeviceAPI) returnList(count int, devices []DeviceListItem) (*pb.ListDeviceResponse, error) {
	resp := pb.ListDeviceResponse{
		TotalCount: int64(count),
	}
	for _, device := range devices {
		item := pb.DeviceListItem{
			DevEui:                          device.DevEUI.String(),
			Name:                            device.Name,
			Description:                     device.Description,
			ApplicationId:                   device.ApplicationID,
			DeviceProfileId:                 device.DeviceProfileID.String(),
			DeviceProfileName:               device.DeviceProfileName,
			DeviceStatusBattery:             256,
			DeviceStatusMargin:              256,
			DeviceStatusExternalPowerSource: device.DeviceStatusExternalPower,
		}

		if !device.DeviceStatusExternalPower && device.DeviceStatusBattery == nil {
			item.DeviceStatusBattery = 255
			item.DeviceStatusBatteryLevelUnavailable = true
		}

		if device.DeviceStatusExternalPower {
			item.DeviceStatusBattery = 0
		}

		if device.DeviceStatusBattery != nil {
			item.DeviceStatusBattery = uint32(254 / *device.DeviceStatusBattery * 100)
			item.DeviceStatusBatteryLevel = *device.DeviceStatusBattery
		}

		if device.DeviceStatusBattery != nil {
			item.DeviceStatusBattery = uint32(*device.DeviceStatusBattery)
		}
		if device.DeviceStatusMargin != nil {
			item.DeviceStatusMargin = int32(*device.DeviceStatusMargin)
		}
		if device.LastSeenAt != nil {
			var err error
			item.LastSeenAt, err = ptypes.TimestampProto(*device.LastSeenAt)
			if err != nil {
				return nil, helpers.ErrToRPCError(err)
			}
		}

		resp.Result = append(resp.Result, &item)
	}
	return &resp, nil
}

func ConvertUplinkAndDownlinkFrames(up *ns.UplinkFrameLog, down *ns.DownlinkFrameLog, decodeMACCommands bool) (*pb.UplinkFrameLog, *pb.DownlinkFrameLog, error) {
	var phy lorawan.PHYPayload

	if up != nil {
		if err := phy.UnmarshalBinary(up.PhyPayload); err != nil {
			return nil, nil, errors.Wrap(err, "unmarshal phypayload error")
		}
	}

	if down != nil {
		if err := phy.UnmarshalBinary(down.PhyPayload); err != nil {
			return nil, nil, errors.Wrap(err, "unmarshal phypayload error")
		}
	}

	if decodeMACCommands {
		switch v := phy.MACPayload.(type) {
		case *lorawan.MACPayload:
			if err := phy.DecodeFOptsToMACCommands(); err != nil {
				return nil, nil, errors.Wrap(err, "decode fopts to mac-commands error")
			}

			if v.FPort != nil && *v.FPort == 0 {
				if err := phy.DecodeFRMPayloadToMACCommands(); err != nil {
					return nil, nil, errors.Wrap(err, "decode frmpayload to mac-commands error")
				}
			}
		}
	}

	phyJSON, err := json.Marshal(phy)
	if err != nil {
		return nil, nil, errors.Wrap(err, "marshal phypayload error")
	}

	if up != nil {
		uplinkFrameLog := pb.UplinkFrameLog{
			TxInfo:         up.TxInfo,
			RxInfo:         up.RxInfo,
			PhyPayloadJson: string(phyJSON),
		}

		return &uplinkFrameLog, nil, nil
	}

	if down != nil {
		var gatewayID lorawan.EUI64
		copy(gatewayID[:], down.GatewayId)

		downlinkFrameLog := pb.DownlinkFrameLog{
			TxInfo:         down.TxInfo,
			PhyPayloadJson: string(phyJSON),
		}

		return nil, &downlinkFrameLog, nil
	}

	return nil, nil, nil
}

// GetDeviceList defines the get device list request and response
func (s *DeviceAPI) GetDeviceList(ctx context.Context, req *pb.GetDeviceListRequest) (*pb.GetDeviceListResponse, error) {
	logInfo := "api/appserver_serves_ui/GetDeviceList org=" + strconv.FormatInt(req.OrgId, 10)

	// verify if user is global admin
	userIsAdmin, err := user.GetUserAPI().Validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &pb.GetDeviceListResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if userIsAdmin == false {
		if err := s.Validator.otpValidator.JwtValidator.Validate(ctx, organization.ValidateOrganizationAccess(organization.Read, req.OrgId)); err != nil {
			log.WithError(err).Error(logInfo)
			return &pb.GetDeviceListResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &pb.GetDeviceListResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	devClient := m2mServer.NewDSDeviceServiceClient(m2mClient)

	resp, err := devClient.GetDeviceList(ctx, &m2mServer.GetDeviceListRequest{
		OrgId:  req.OrgId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &pb.GetDeviceListResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	var deviceProfileList []*pb.DSDeviceProfile
	for _, item := range resp.DevProfile {
		deviceProfile := &pb.DSDeviceProfile{
			Id:            item.Id,
			DevEui:        item.DevEui,
			FkWallet:      item.Id,
			Mode:          pb.DeviceMode(item.Mode),
			CreatedAt:     item.CreatedAt,
			LastSeenAt:    item.LastSeenAt,
			ApplicationId: item.ApplicationId,
			Name:          item.Name,
		}

		deviceProfileList = append(deviceProfileList, deviceProfile)
	}

	return &pb.GetDeviceListResponse{
		DevProfile: deviceProfileList,
		Count:      resp.Count,
	}, status.Error(codes.OK, "")
}

// GetDeviceProfile defines the function to get the device profile
func (s *DeviceAPI) GetDeviceProfile(ctx context.Context, req *pb.GetDSDeviceProfileRequest) (*pb.GetDSDeviceProfileResponse, error) {
	logInfo := "api/appserver_serves_ui/GetDeviceProfile org=" + strconv.FormatInt(req.OrgId, 10)

	// verify if user is global admin
	userIsAdmin, err := user.GetUserAPI().Validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &pb.GetDSDeviceProfileResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if userIsAdmin == false {
		if err := s.Validator.otpValidator.JwtValidator.Validate(ctx, organization.ValidateOrganizationAccess(organization.Read, req.OrgId)); err != nil {
			log.WithError(err).Error(logInfo)
			return &pb.GetDSDeviceProfileResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &pb.GetDSDeviceProfileResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	devClient := m2mServer.NewDSDeviceServiceClient(m2mClient)

	resp, err := devClient.GetDeviceProfile(ctx, &m2mServer.GetDSDeviceProfileRequest{
		OrgId: req.OrgId,
		DevId: req.DevId,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &pb.GetDSDeviceProfileResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &pb.GetDSDeviceProfileResponse{
		DevProfile: &pb.DSDeviceProfile{
			Id:            resp.DevProfile.Id,
			DevEui:        resp.DevProfile.DevEui,
			FkWallet:      resp.DevProfile.FkWallet,
			Mode:          pb.DeviceMode(resp.DevProfile.Mode),
			CreatedAt:     resp.DevProfile.CreatedAt,
			LastSeenAt:    resp.DevProfile.LastSeenAt,
			ApplicationId: resp.DevProfile.ApplicationId,
			Name:          resp.DevProfile.Name,
		},
	}, status.Error(codes.OK, "")
}

// GetDeviceHistory defines the get device history request and response
func (s *DeviceAPI) GetDeviceHistory(ctx context.Context, req *pb.GetDeviceHistoryRequest) (*pb.GetDeviceHistoryResponse, error) {
	logInfo := "api/appserver_serves_ui/GetDeviceHistory org=" + strconv.FormatInt(req.OrgId, 10)

	// verify if user is global admin
	userIsAdmin, err := user.GetUserAPI().Validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &pb.GetDeviceHistoryResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if userIsAdmin == false {
		if err := s.Validator.otpValidator.JwtValidator.Validate(ctx, organization.ValidateOrganizationAccess(organization.Read, req.OrgId)); err != nil {
			log.WithError(err).Error(logInfo)
			return &pb.GetDeviceHistoryResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &pb.GetDeviceHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	devClient := m2mServer.NewDSDeviceServiceClient(m2mClient)

	resp, err := devClient.GetDeviceHistory(ctx, &m2mServer.GetDeviceHistoryRequest{
		OrgId:  req.OrgId,
		DevId:  req.DevId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &pb.GetDeviceHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &pb.GetDeviceHistoryResponse{
		DevHistory: resp.DevHistory,
	}, status.Error(codes.OK, "")
}

// SetDeviceMode defines the set device mode request and response
func (s *DeviceAPI) SetDeviceMode(ctx context.Context, req *pb.SetDeviceModeRequest) (*pb.SetDeviceModeResponse, error) {
	logInfo := "api/appserver_serves_ui/SetDeviceMode org=" + strconv.FormatInt(req.OrgId, 10)

	// verify if user is global admin
	userIsAdmin, err := user.GetUserAPI().Validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &pb.SetDeviceModeResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if userIsAdmin == false {
		if err := s.Validator.otpValidator.JwtValidator.Validate(ctx, organization.ValidateOrganizationAccess(organization.Read, req.OrgId)); err != nil {
			log.WithError(err).Error(logInfo)
			return &pb.SetDeviceModeResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &pb.SetDeviceModeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	devClient := m2mServer.NewDSDeviceServiceClient(m2mClient)

	resp, err := devClient.SetDeviceMode(ctx, &m2mServer.SetDeviceModeRequest{
		OrgId:   req.OrgId,
		DevId:   req.DevId,
		DevMode: m2mServer.DeviceMode(req.DevMode),
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &pb.SetDeviceModeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &pb.SetDeviceModeResponse{
		Status: resp.Status,
	}, status.Error(codes.OK, "")
}
