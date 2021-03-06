package external

import (
	"database/sql"
	"encoding/json"
	"strconv"

	appd "github.com/mxc-foundation/lpwan-app-server/internal/modules/application/data"

	"google.golang.org/grpc/status"

	"github.com/gobuffalo/packr/v2/file/resolver/encoding/hex"
	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lib/pq/hstore"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/brocaar/chirpstack-api/go/v3/common"
	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/brocaar/lorawan"

	api "github.com/mxc-foundation/lpwan-app-server/api/extapi"
	pb "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	psPb "github.com/mxc-foundation/lpwan-app-server/api/ps-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/device"
	dps "github.com/mxc-foundation/lpwan-app-server/internal/api/external/dp"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/organization"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	"github.com/mxc-foundation/lpwan-app-server/internal/auth"
	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	"github.com/mxc-foundation/lpwan-app-server/internal/eventlog"
	"github.com/mxc-foundation/lpwan-app-server/internal/logging"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/application"
	devmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/device"
	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/device/data"
	serviceprofile "github.com/mxc-foundation/lpwan-app-server/internal/modules/service-profile"
	"github.com/mxc-foundation/lpwan-app-server/internal/mxpcli"
	"github.com/mxc-foundation/lpwan-app-server/internal/nscli"
	"github.com/mxc-foundation/lpwan-app-server/internal/pscli"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

// DeviceAPI exports the Node related functions.
type DeviceAPI struct {
	st                  *store.Handler
	auth                auth.Authenticator
	ApplicationServerID uuid.UUID
	mxpCli              pb.DSDeviceServiceClient
	psCli               psPb.DeviceProvisionClient
	nsCli               *nscli.Client
}

// NewDeviceAPI creates a new NodeAPI.
func NewDeviceAPI(applicationID uuid.UUID, h *store.Handler, mxpCli *mxpcli.Client,
	psCli *pscli.Client, nsCli *nscli.Client, auth auth.Authenticator) *DeviceAPI {
	return &DeviceAPI{
		st:                  h,
		ApplicationServerID: applicationID,
		auth:                auth,
		mxpCli:              mxpCli.GetM2MDeviceServiceClient(),
		psCli:               psCli.GetDeviceProvisionServiceClient(),
		nsCli:               nsCli,
	}
}

// Create creates the given device.
func (a *DeviceAPI) Create(ctx context.Context, req *api.CreateDeviceRequest) (*empty.Empty, error) {
	var response empty.Empty
	if req.Device == nil {
		return nil, status.Errorf(codes.InvalidArgument, "device must not be nil")
	}

	var devEUI lorawan.EUI64
	if err := devEUI.UnmarshalText([]byte(req.Device.DevEui)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	dpID, err := uuid.FromString(req.Device.DeviceProfileId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	// if Name is "", set it to the DevEUI
	if req.Device.Name == "" {
		req.Device.Name = req.Device.DevEui
	}

	// Validate that application and device-profile are under the same
	// organization ID.
	app, err := a.st.GetApplication(ctx, req.Device.ApplicationId)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	dp, err := a.st.GetDeviceProfile(ctx, dpID, false)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	if app.OrganizationID != dp.OrganizationID {
		return nil, status.Errorf(codes.InvalidArgument, "device-profile and application must be under the same organization")
	}

	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(app.OrganizationID))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}
	if !cred.IsGlobalAdmin && !cred.IsDeviceAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
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

	nsCli, err := a.nsCli.GetNetworkServerServiceClient(dp.NetworkServerID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't get network server client: %v", err)
	}

	if err := device.CreateDevice(ctx, a.st, &d, &app, a.ApplicationServerID, a.mxpCli, nsCli); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &response, nil
}

// Get returns the device matching the given DevEUI.
func (a *DeviceAPI) Get(ctx context.Context, req *api.GetDeviceRequest) (*api.GetDeviceResponse, error) {
	var eui lorawan.EUI64
	if err := eui.UnmarshalText([]byte(req.DevEui)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	d, err := a.st.GetDevice(ctx, eui, false)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	app, err := a.st.GetApplication(ctx, d.ApplicationID)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(app.OrganizationID))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}
	if !cred.IsGlobalAdmin && !cred.IsOrgUser {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	n, err := a.st.GetNetworkServerForDevEUI(ctx, d.DevEUI)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	nsClient, err := a.nsCli.GetNetworkServerServiceClient(n.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	res, err := nsClient.GetDevice(ctx, &ns.GetDeviceRequest{
		DevEui: d.DevEUI[:],
	})
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	if res.Device != nil {
		d.SkipFCntCheck = res.Device.SkipFCntCheck
		d.ReferenceAltitude = res.Device.ReferenceAltitude
	}

	response := api.GetDeviceResponse{
		Device: &api.Device{
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
		response.DeviceStatusBattery = uint32(*d.DeviceStatusBattery)
	}
	if d.DeviceStatusMargin != nil {
		response.DeviceStatusMargin = int32(*d.DeviceStatusMargin)
	}
	if d.LastSeenAt != nil {
		response.LastSeenAt = timestamppb.New(*d.LastSeenAt)
	}

	if d.Latitude != nil && d.Longitude != nil && d.Altitude != nil {
		response.Location = &common.Location{
			Latitude:  *d.Latitude,
			Longitude: *d.Longitude,
			Altitude:  *d.Altitude,
		}
	}

	for k, v := range d.Variables.Map {
		if v.Valid {
			response.Device.Variables[k] = v.String
		}
	}

	for k, v := range d.Tags.Map {
		if v.Valid {
			response.Device.Tags[k] = v.String
		}
	}

	return &response, nil
}

// List lists the available applications.
func (a *DeviceAPI) List(ctx context.Context, req *api.ListDeviceRequest) (*api.ListDeviceResponse, error) {
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
		if valid, err := application.NewValidator(a.st).ValidateApplicationAccess(ctx, authcus.Read, req.ApplicationId); !valid || err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
		}
	}

	if filters.MulticastGroupID != uuid.Nil {
		idFilter = true

		// validate that the client has access to the given multicast-group
		if valid, err := devmod.NewValidator(a.st).ValidateMulticastGroupAccess(ctx, authcus.Read, filters.MulticastGroupID); !valid || err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
		}
	}

	if filters.ServiceProfileID != uuid.Nil {
		idFilter = true

		// validate that the client has access to the given service-profile
		if valid, err := serviceprofile.NewValidator(a.st).ValidateServiceProfileAccess(ctx, authcus.Read, filters.ServiceProfileID); !valid || err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
		}
	}

	if !idFilter {
		u, err := devmod.NewValidator(a.st).GetUser(ctx)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}

		if !u.IsGlobalAdmin {
			return nil, status.Errorf(codes.Unauthenticated, "client must be global admin for unfiltered request")
		}
	}

	count, err := a.st.GetDeviceCount(ctx, filters)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	devices, err := a.st.GetDevices(ctx, filters)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return a.returnList(count, devices)
}

func (a *DeviceAPI) verifyDeviceProfileChange(ctx context.Context, dpIDNew, oldDpID uuid.UUID) error {
	var dpNew, dpOld dps.DeviceProfile
	var err error
	if dpIDNew != oldDpID {
		dpOld, err = a.st.GetDeviceProfile(ctx, oldDpID, false)
		if err != nil {
			return status.Errorf(codes.Internal, err.Error())
		}
		dpNew, err = a.st.GetDeviceProfile(ctx, dpIDNew, false)
		if err != nil {
			return status.Errorf(codes.Internal, err.Error())
		}
		if dpOld.NetworkServerID != dpNew.NetworkServerID || dpOld.OrganizationID != dpNew.OrganizationID {
			return status.Errorf(codes.InvalidArgument,
				"when updating device profile for a device, "+
					"new device profile and old device profile must share same network server and organization")
		}
	}
	return nil
}

func (a *DeviceAPI) verifyApplicationChange(ctx context.Context, appNew appd.Application, applicationIDOld int64) error {
	if appNew.ID != applicationIDOld {
		appOld, err := a.st.GetApplication(ctx, applicationIDOld)
		if err != nil {
			return status.Errorf(codes.Internal, err.Error())
		}

		if appOld.ServiceProfileID != appNew.ServiceProfileID {
			return status.Errorf(codes.InvalidArgument,
				"when moving a device from application A to B, both A and B must share the same service-profile")
		}
	}
	return nil
}

// Update updates the device matching the given DevEUI.
func (a *DeviceAPI) Update(ctx context.Context, req *api.UpdateDeviceRequest) (*empty.Empty, error) {
	log.WithFields(log.Fields{
		"application_id":    req.Device.ApplicationId,
		"device_profile_id": req.Device.DeviceProfileId,
	}).Debug("update device")
	var response empty.Empty
	if req.Device == nil {
		return nil, status.Errorf(codes.InvalidArgument, "device must not be nil")
	}

	var devEUI lorawan.EUI64
	if err := devEUI.UnmarshalText([]byte(req.Device.DevEui)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	// get device
	d, err := a.st.GetDevice(ctx, devEUI, true)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	// If the device is moved to a different application, validate that
	// the new application is assigned to the same service-profile.
	// This to guarantee that the new application is still on the same
	// network-server and is not assigned to a different organization.
	appNew, err := a.st.GetApplication(ctx, req.Device.ApplicationId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	err = a.verifyApplicationChange(ctx, appNew, d.ApplicationID)
	if err != nil {
		return nil, err
	}
	sp, err := a.st.GetServiceProfile(ctx, appNew.ServiceProfileID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "service profile %s not found: %v", appNew.ServiceProfileID.String(), err)
	}

	// If the device is assigned with a different device profile, validate that
	// the new device profile is still on the same network-server and organization.
	dpIDNew, err := uuid.FromString(req.Device.DeviceProfileId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid device profile : %v", err)
	}
	err = a.verifyDeviceProfileChange(ctx, dpIDNew, d.DeviceProfileID)
	if err != nil {
		return nil, err
	}

	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(sp.OrganizationID))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}
	if !cred.IsGlobalAdmin && !cred.IsDeviceAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	// get network server client
	nsClient, err := a.nsCli.GetNetworkServerServiceClient(sp.NetworkServerID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// update device attributes
	d.ApplicationID = req.Device.ApplicationId
	d.DeviceProfileID = dpIDNew
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

	if err := a.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		if err := handler.UpdateDevice(ctx, &d); err != nil {
			return status.Errorf(codes.Internal, "%v", err)
		}
		// update external server at last, if internal process succeeds, but external service call fails,
		// then internal process will roll back, no changes on both sides.
		// If internal process fails, external service won't be called, still no changes on both sides.
		_, err = nsClient.UpdateDevice(ctx, &ns.UpdateDeviceRequest{
			Device: &ns.Device{
				DevEui:            d.DevEUI[:],
				DeviceProfileId:   d.DeviceProfileID.Bytes(),
				ServiceProfileId:  appNew.ServiceProfileID.Bytes(),
				RoutingProfileId:  a.ApplicationServerID.Bytes(),
				SkipFCntCheck:     d.SkipFCntCheck,
				ReferenceAltitude: d.ReferenceAltitude,
			},
		})
		if err != nil {
			return status.Errorf(codes.Unknown, "update device error: %v", err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &response, nil
}

// Delete deletes the node matching the given name.
func (a *DeviceAPI) Delete(ctx context.Context, req *api.DeleteDeviceRequest) (*empty.Empty, error) {
	var response empty.Empty
	var eui lorawan.EUI64
	if err := eui.UnmarshalText([]byte(req.DevEui)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	dev, err := a.st.GetDevice(ctx, eui, false)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	app, err := a.st.GetApplication(ctx, dev.ApplicationID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(app.OrganizationID))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}
	if !cred.IsGlobalAdmin && !cred.IsDeviceAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	// as this also performs a remote call to delete the node from the
	// network-server, wrap it in a transaction
	err = device.DeleteDevice(ctx, a.st, eui, a.mxpCli, a.psCli, a.nsCli)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &response, nil
}

// CreateKeys creates the given device-keys.
func (a *DeviceAPI) CreateKeys(ctx context.Context, req *api.CreateDeviceKeysRequest) (*empty.Empty, error) {
	var response empty.Empty
	if req.DeviceKeys == nil {
		return nil, status.Errorf(codes.InvalidArgument, "device_keys must not be nil")
	}

	// appKey is not used for LoRaWAN 1.0
	var appKey lorawan.AES128Key
	if req.DeviceKeys.AppKey != "" {
		if err := appKey.UnmarshalText([]byte(req.DeviceKeys.AppKey)); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
	}

	// genAppKey is only for LoRaWAN 1.0 devices that implement the
	// remote multicast setup specification.
	var genAppKey lorawan.AES128Key
	if req.DeviceKeys.GenAppKey != "" {
		if err := genAppKey.UnmarshalText([]byte(req.DeviceKeys.GenAppKey)); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
	}

	// nwkKey
	var nwkKey lorawan.AES128Key
	if err := nwkKey.UnmarshalText([]byte(req.DeviceKeys.NwkKey)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	// devEUI
	var eui lorawan.EUI64
	if err := eui.UnmarshalText([]byte(req.DeviceKeys.DevEui)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if valid, err := devmod.NewValidator(a.st).ValidateNodeAccess(ctx, authcus.Update, eui); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	err := a.st.CreateDeviceKeys(ctx, &DeviceKeys{
		DevEUI:    eui,
		NwkKey:    nwkKey,
		AppKey:    appKey,
		GenAppKey: genAppKey,
	})
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetKeys returns the device-keys for the given DevEUI.
func (a *DeviceAPI) GetKeys(ctx context.Context, req *api.GetDeviceKeysRequest) (*api.GetDeviceKeysResponse, error) {
	var eui lorawan.EUI64
	if err := eui.UnmarshalText([]byte(req.DevEui)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if valid, err := devmod.NewValidator(a.st).ValidateNodeAccess(ctx, authcus.Update, eui); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	dk, err := a.st.GetDeviceKeys(ctx, eui)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &api.GetDeviceKeysResponse{
		DeviceKeys: &api.DeviceKeys{
			DevEui:    eui.String(),
			AppKey:    dk.AppKey.String(),
			NwkKey:    dk.NwkKey.String(),
			GenAppKey: dk.GenAppKey.String(),
		},
	}, nil
}

// UpdateKeys updates the device-keys.
func (a *DeviceAPI) UpdateKeys(ctx context.Context, req *api.UpdateDeviceKeysRequest) (*empty.Empty, error) {
	var response empty.Empty
	if req.DeviceKeys == nil {
		return nil, status.Errorf(codes.InvalidArgument, "device_keys must not be nil")
	}

	var appKey lorawan.AES128Key
	// appKey is not used for LoRaWAN 1.0
	if req.DeviceKeys.AppKey != "" {
		if err := appKey.UnmarshalText([]byte(req.DeviceKeys.AppKey)); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
	}

	// genAppKey is only for LoRaWAN 1.0 devices that implement the
	// remote multicast setup specification.
	var genAppKey lorawan.AES128Key
	if req.DeviceKeys.GenAppKey != "" {
		if err := genAppKey.UnmarshalText([]byte(req.DeviceKeys.GenAppKey)); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
	}

	var nwkKey lorawan.AES128Key
	if err := nwkKey.UnmarshalText([]byte(req.DeviceKeys.NwkKey)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	var eui lorawan.EUI64
	if err := eui.UnmarshalText([]byte(req.DeviceKeys.DevEui)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if valid, err := devmod.NewValidator(a.st).ValidateNodeAccess(ctx, authcus.Update, eui); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	dk, err := a.st.GetDeviceKeys(ctx, eui)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	dk.NwkKey = nwkKey
	dk.AppKey = appKey
	dk.GenAppKey = genAppKey

	err = a.st.UpdateDeviceKeys(ctx, &dk)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteKeys deletes the device-keys for the given DevEUI.
func (a *DeviceAPI) DeleteKeys(ctx context.Context, req *api.DeleteDeviceKeysRequest) (*empty.Empty, error) {
	var response empty.Empty
	var eui lorawan.EUI64
	if err := eui.UnmarshalText([]byte(req.DevEui)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if valid, err := devmod.NewValidator(a.st).ValidateNodeAccess(ctx, authcus.Delete, eui); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	err := a.st.DeleteDeviceKeys(ctx, eui)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Deactivate de-activates the device.
func (a *DeviceAPI) Deactivate(ctx context.Context, req *api.DeactivateDeviceRequest) (*empty.Empty, error) {
	var response empty.Empty
	var devEUI lorawan.EUI64
	if err := devEUI.UnmarshalText([]byte(req.DevEui)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if valid, err := devmod.NewValidator(a.st).ValidateNodeAccess(ctx, authcus.Update, devEUI); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	d, err := a.st.GetDevice(ctx, devEUI, false)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	n, err := a.st.GetNetworkServerForDevEUI(ctx, d.DevEUI)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	nsClient, err := a.nsCli.GetNetworkServerServiceClient(n.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	_, _ = nsClient.DeactivateDevice(ctx, &ns.DeactivateDeviceRequest{
		DevEui: d.DevEUI[:],
	})

	return &response, nil
}

// Activate activates the node (ABP only).
func (a *DeviceAPI) Activate(ctx context.Context, req *api.ActivateDeviceRequest) (*empty.Empty, error) {
	var response empty.Empty
	if req.DeviceActivation == nil {
		return nil, status.Errorf(codes.InvalidArgument, "device_activation must not be nil")
	}

	var devAddr lorawan.DevAddr
	var devEUI lorawan.EUI64
	var appSKey lorawan.AES128Key
	var nwkSEncKey lorawan.AES128Key
	var sNwkSIntKey lorawan.AES128Key
	var fNwkSIntKey lorawan.AES128Key

	if err := devAddr.UnmarshalText([]byte(req.DeviceActivation.DevAddr)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "devAddr: %s", err)
	}
	if err := devEUI.UnmarshalText([]byte(req.DeviceActivation.DevEui)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "devEUI: %s", err)
	}
	if err := appSKey.UnmarshalText([]byte(req.DeviceActivation.AppSKey)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "appSKey: %s", err)
	}
	if err := nwkSEncKey.UnmarshalText([]byte(req.DeviceActivation.NwkSEncKey)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "nwkSEncKey: %s", err)
	}
	if err := sNwkSIntKey.UnmarshalText([]byte(req.DeviceActivation.SNwkSIntKey)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "sNwkSIntKey: %s", err)
	}
	if err := fNwkSIntKey.UnmarshalText([]byte(req.DeviceActivation.FNwkSIntKey)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "fNwkSIntKey: %s", err)
	}

	if valid, err := devmod.NewValidator(a.st).ValidateNodeAccess(ctx, authcus.Update, devEUI); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	d, err := a.st.GetDevice(ctx, devEUI, false)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	n, err := a.st.GetNetworkServerForDevEUI(ctx, devEUI)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	nsClient, err := a.nsCli.GetNetworkServerServiceClient(n.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	_, _ = nsClient.DeactivateDevice(ctx, &ns.DeactivateDeviceRequest{
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

	if err := a.st.UpdateDeviceActivation(ctx, d.DevEUI, devAddr, appSKey); err != nil {
		return nil, status.Errorf(codes.Unknown, "%v", err)
	}

	_, err = nsClient.ActivateDevice(ctx, &actReq)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%v", err)
	}

	log.WithFields(log.Fields{
		"dev_addr": devAddr,
		"dev_eui":  d.DevEUI,
		"ctx_id":   ctx.Value(logging.ContextIDKey),
	}).Info("device activated")

	return &response, nil
}

// GetActivation returns the device activation for the given DevEUI.
func (a *DeviceAPI) GetActivation(ctx context.Context, req *api.GetDeviceActivationRequest) (*api.GetDeviceActivationResponse, error) {
	var devAddr lorawan.DevAddr
	var devEUI lorawan.EUI64
	var sNwkSIntKey lorawan.AES128Key
	var fNwkSIntKey lorawan.AES128Key
	var nwkSEncKey lorawan.AES128Key

	if err := devEUI.UnmarshalText([]byte(req.DevEui)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "devEUI: %s", err)
	}

	if valid, err := devmod.NewValidator(a.st).ValidateNodeAccess(ctx, authcus.Read, devEUI); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	d, err := a.st.GetDevice(ctx, devEUI, false)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	n, err := a.st.GetNetworkServerForDevEUI(ctx, devEUI)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	nsClient, err := a.nsCli.GetNetworkServerServiceClient(n.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	devAct, err := nsClient.GetDeviceActivation(ctx, &ns.GetDeviceActivationRequest{
		DevEui: d.DevEUI[:],
	})
	if err != nil {
		return nil, err
	}

	copy(devAddr[:], devAct.DeviceActivation.DevAddr)
	copy(nwkSEncKey[:], devAct.DeviceActivation.NwkSEncKey)
	copy(sNwkSIntKey[:], devAct.DeviceActivation.SNwkSIntKey)
	copy(fNwkSIntKey[:], devAct.DeviceActivation.FNwkSIntKey)

	return &api.GetDeviceActivationResponse{
		DeviceActivation: &api.DeviceActivation{
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
func (a *DeviceAPI) StreamFrameLogs(req *api.StreamDeviceFrameLogsRequest, srv api.DeviceService_StreamFrameLogsServer) error {
	var devEUI lorawan.EUI64

	if err := devEUI.UnmarshalText([]byte(req.DevEui)); err != nil {
		return status.Errorf(codes.InvalidArgument, "devEUI: %s", err)
	}

	if valid, err := devmod.NewValidator(a.st).ValidateNodeAccess(srv.Context(), authcus.Read, devEUI); !valid || err != nil {
		return status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	n, err := a.st.GetNetworkServerForDevEUI(srv.Context(), devEUI)
	if err != nil {
		return helpers.ErrToRPCError(err)
	}

	nsClient, err := a.nsCli.GetNetworkServerServiceClient(n.ID)
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	streamClient, err := nsClient.StreamFrameLogsForDevice(srv.Context(), &ns.StreamFrameLogsForDeviceRequest{
		DevEui: devEUI[:],
	})
	if err != nil {
		return err
	}

	for {
		resp, err := streamClient.Recv()
		if err != nil {
			if status.Code(err) == codes.Canceled {
				return nil
			}

			return err
		}

		up, down, err := ConvertUplinkAndDownlinkFrames(resp.GetUplinkFrameSet(), resp.GetDownlinkFrame(), true)
		if err != nil {
			return helpers.ErrToRPCError(err)
		}

		var frameResp api.StreamDeviceFrameLogsResponse
		if up != nil {
			frameResp.Frame = &api.StreamDeviceFrameLogsResponse_UplinkFrame{
				UplinkFrame: up,
			}
		}

		if down != nil {
			frameResp.Frame = &api.StreamDeviceFrameLogsResponse_DownlinkFrame{
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
func (a *DeviceAPI) StreamEventLogs(req *api.StreamDeviceEventLogsRequest, srv api.DeviceService_StreamEventLogsServer) error {
	var devEUI lorawan.EUI64

	if err := devEUI.UnmarshalText([]byte(req.DevEui)); err != nil {
		return status.Errorf(codes.InvalidArgument, "devEUI: %s", err)
	}

	if valid, err := devmod.NewValidator(a.st).ValidateNodeAccess(srv.Context(), authcus.Read, devEUI); !valid || err != nil {
		return status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
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
			return status.Errorf(codes.Internal, "marshal json error: %s", err)
		}

		resp := api.StreamDeviceEventLogsResponse{
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
func (a *DeviceAPI) GetRandomDevAddr(ctx context.Context, req *api.GetRandomDevAddrRequest) (*api.GetRandomDevAddrResponse, error) {
	var devEUI lorawan.EUI64

	if err := devEUI.UnmarshalText([]byte(req.DevEui)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "devEUI: %s", err)
	}

	n, err := a.st.GetNetworkServerForDevEUI(ctx, devEUI)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	nsClient, err := a.nsCli.GetNetworkServerServiceClient(n.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	var emptyReqeust empty.Empty
	resp, err := nsClient.GetRandomDevAddr(ctx, &emptyReqeust)
	if err != nil {
		return nil, err
	}

	var devAddr lorawan.DevAddr
	copy(devAddr[:], resp.DevAddr)

	return &api.GetRandomDevAddrResponse{
		DevAddr: devAddr.String(),
	}, nil
}

func (a *DeviceAPI) returnList(count int, devices []DeviceListItem) (*api.ListDeviceResponse, error) {
	resp := api.ListDeviceResponse{
		TotalCount: int64(count),
	}
	for _, device := range devices {
		item := api.DeviceListItem{
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
			item.LastSeenAt = timestamppb.New(*device.LastSeenAt)
		}

		resp.Result = append(resp.Result, &item)
	}
	return &resp, nil
}

func ConvertUplinkAndDownlinkFrames(up *ns.UplinkFrameLog, down *ns.DownlinkFrameLog, decodeMACCommands bool) (*api.UplinkFrameLog, *api.DownlinkFrameLog, error) {
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
		uplinkFrameLog := api.UplinkFrameLog{
			TxInfo:         up.TxInfo,
			RxInfo:         up.RxInfo,
			PhyPayloadJson: string(phyJSON),
		}

		return &uplinkFrameLog, nil, nil
	}

	if down != nil {
		var gatewayID lorawan.EUI64
		copy(gatewayID[:], down.GatewayId)

		downlinkFrameLog := api.DownlinkFrameLog{
			TxInfo:         down.TxInfo,
			PhyPayloadJson: string(phyJSON),
		}

		return nil, &downlinkFrameLog, nil
	}

	return nil, nil, nil
}

// GetDeviceList defines the get device list request and response
func (a *DeviceAPI) GetDeviceList(ctx context.Context, req *api.GetDeviceListRequest) (*api.GetDeviceListResponse, error) {
	logInfo := "api/appserver_serves_ui/GetDeviceList org=" + strconv.FormatInt(req.OrgId, 10)

	// verify if user is global admin
	u, err := devmod.NewValidator(a.st).GetUser(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetDeviceListResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if !u.IsGlobalAdmin {
		if valid, err := organization.NewValidator(a.st).ValidateOrganizationAccess(ctx, authcus.Read, req.OrgId); !valid || err != nil {
			return &api.GetDeviceListResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
		}
	}

	devClient := mxpcli.Global.GetM2MDeviceServiceClient()

	resp, err := devClient.GetDeviceList(ctx, &pb.GetDeviceListRequest{
		OrgId:  req.OrgId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetDeviceListResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	var deviceProfileList []*api.DSDeviceProfile
	for _, item := range resp.DevProfile {
		deviceProfile := &api.DSDeviceProfile{
			Id:            item.Id,
			DevEui:        item.DevEui,
			FkWallet:      item.Id,
			Mode:          api.DeviceMode(item.Mode),
			CreatedAt:     item.CreatedAt,
			LastSeenAt:    item.LastSeenAt,
			ApplicationId: item.ApplicationId,
			Name:          item.Name,
		}

		// Get Last Seen from DB
		euiarray, err := hex.DecodeString(item.DevEui)
		if err == nil {
			var eui lorawan.EUI64
			copy(eui[:], euiarray)
			d, err := a.st.GetDevice(ctx, eui, false)
			if err == nil && d.LastSeenAt != nil {
				deviceProfile.LastSeenAt = d.LastSeenAt.String()
			}
		}

		deviceProfileList = append(deviceProfileList, deviceProfile)
	}

	return &api.GetDeviceListResponse{
		DevProfile: deviceProfileList,
		Count:      resp.Count,
	}, status.Error(codes.OK, "")
}

// GetDeviceProfile defines the function to get the device profile
func (a *DeviceAPI) GetDeviceProfile(ctx context.Context, req *api.GetDSDeviceProfileRequest) (*api.GetDSDeviceProfileResponse, error) {
	logInfo := "api/appserver_serves_ui/GetDeviceProfile org=" + strconv.FormatInt(req.OrgId, 10)

	// verify if user is global admin
	u, err := devmod.NewValidator(a.st).GetUser(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetDSDeviceProfileResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if !u.IsGlobalAdmin {
		if valid, err := organization.NewValidator(a.st).ValidateOrganizationAccess(ctx, authcus.Read, req.OrgId); !valid || err != nil {
			return &api.GetDSDeviceProfileResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
		}
	}

	devClient := mxpcli.Global.GetM2MDeviceServiceClient()

	resp, err := devClient.GetDeviceProfile(ctx, &pb.GetDSDeviceProfileRequest{
		OrgId: req.OrgId,
		DevId: req.DevId,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetDSDeviceProfileResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetDSDeviceProfileResponse{
		DevProfile: &api.DSDeviceProfile{
			Id:            resp.DevProfile.Id,
			DevEui:        resp.DevProfile.DevEui,
			FkWallet:      resp.DevProfile.FkWallet,
			Mode:          api.DeviceMode(resp.DevProfile.Mode),
			CreatedAt:     resp.DevProfile.CreatedAt,
			LastSeenAt:    resp.DevProfile.LastSeenAt,
			ApplicationId: resp.DevProfile.ApplicationId,
			Name:          resp.DevProfile.Name,
		},
	}, status.Error(codes.OK, "")
}

// GetDeviceHistory defines the get device history request and response
func (a *DeviceAPI) GetDeviceHistory(ctx context.Context, req *api.GetDeviceHistoryRequest) (*api.GetDeviceHistoryResponse, error) {
	logInfo := "api/appserver_serves_ui/GetDeviceHistory org=" + strconv.FormatInt(req.OrgId, 10)

	// verify if user is global admin
	u, err := devmod.NewValidator(a.st).GetUser(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetDeviceHistoryResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if !u.IsGlobalAdmin {
		if valid, err := organization.NewValidator(a.st).ValidateOrganizationAccess(ctx, authcus.Read, req.OrgId); !valid || err != nil {
			return &api.GetDeviceHistoryResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
		}
	}

	devClient := mxpcli.Global.GetM2MDeviceServiceClient()

	resp, err := devClient.GetDeviceHistory(ctx, &pb.GetDeviceHistoryRequest{
		OrgId:  req.OrgId,
		DevId:  req.DevId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetDeviceHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetDeviceHistoryResponse{
		DevHistory: resp.DevHistory,
	}, status.Error(codes.OK, "")
}

// SetDeviceMode defines the set device mode request and response
func (a *DeviceAPI) SetDeviceMode(ctx context.Context, req *api.SetDeviceModeRequest) (*api.SetDeviceModeResponse, error) {
	logInfo := "api/appserver_serves_ui/SetDeviceMode org=" + strconv.FormatInt(req.OrgId, 10)

	// verify if user is global admin
	u, err := devmod.NewValidator(a.st).GetUser(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.SetDeviceModeResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if !u.IsGlobalAdmin {
		if valid, err := organization.NewValidator(a.st).ValidateOrganizationAccess(ctx, authcus.Read, req.OrgId); !valid || err != nil {
			return &api.SetDeviceModeResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
		}
	}

	devClient := mxpcli.Global.GetM2MDeviceServiceClient()

	resp, err := devClient.SetDeviceMode(ctx, &pb.SetDeviceModeRequest{
		OrgId:   req.OrgId,
		DevId:   req.DevId,
		DevMode: pb.DeviceMode(req.DevMode),
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.SetDeviceModeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.SetDeviceModeResponse{
		Status: resp.Status,
	}, status.Error(codes.OK, "")
}
