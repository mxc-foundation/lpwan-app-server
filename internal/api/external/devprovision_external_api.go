package external

import (
	"context"
	"database/sql"
	"encoding/hex"

	"github.com/brocaar/lorawan"
	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lib/pq/hstore"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	psPb "github.com/mxc-foundation/lpwan-app-server/api/ps-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	"github.com/mxc-foundation/lpwan-app-server/internal/auth"
	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/application"
	appmoddata "github.com/mxc-foundation/lpwan-app-server/internal/modules/application/data"
	devmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/device"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/device/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

// ProvisionedDeviceAPI exports the Gateway related functions.
type ProvisionedDeviceAPI struct {
	st                  *store.Handler
	ApplicationServerID uuid.UUID
	ServerAddr          string
	auth                auth.Authenticator
	pscli               psPb.DeviceProvisionClient
}

// NewDeviceProvisionAPI - creates a new DeviceProvisionAPI.
func NewDeviceProvisionAPI(h *store.Handler, auth auth.Authenticator,
	applicationID uuid.UUID, serverAddr string, pscli psPb.DeviceProvisionClient) *ProvisionedDeviceAPI {
	return &ProvisionedDeviceAPI{
		st:                  h,
		auth:                auth,
		ApplicationServerID: applicationID,
		ServerAddr:          serverAddr,
		pscli:               pscli,
	}
}

//
func (a *ProvisionedDeviceAPI) createDevice(ctx context.Context, app appmoddata.Application, devEUI lorawan.EUI64,
	req *api.CreateProvisionedDeviceRequest, deviceProfileID uuid.UUID, provisioned *psPb.GetDeviceResponse) error {
	// Set Device struct.
	d := data.Device{
		DevEUI:            devEUI,
		ApplicationID:     req.Device.ApplicationId,
		DeviceProfileID:   deviceProfileID,
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
	log.Debugf("data.Device: %v", d)

	for k, v := range req.Device.Variables {
		d.Variables.Map[k] = sql.NullString{String: v, Valid: true}
	}

	for k, v := range req.Device.Tags {
		d.Tags.Map[k] = sql.NullString{String: v, Valid: true}
	}

	if err := devmod.CreateDevice(ctx, &d, &app, a.ApplicationServerID); err != nil {
		return status.Errorf(codes.Unknown, err.Error())
	}

	// Set Keys
	var appKey lorawan.AES128Key
	var nwkKey lorawan.AES128Key
	var genAppKey lorawan.AES128Key

	copy(appKey[:], provisioned.AppKey)
	copy(nwkKey[:], provisioned.NwkKey)
	copy(genAppKey[:], provisioned.NwkKey)

	if err := a.st.CreateDeviceKeys(ctx, &data.DeviceKeys{
		DevEUI:    devEUI,
		NwkKey:    nwkKey,
		AppKey:    appKey,
		GenAppKey: genAppKey,
	}); err != nil {
		return err
	}

	return nil
}

// Create - creates the given device.
func (a *ProvisionedDeviceAPI) Create(ctx context.Context, req *api.CreateProvisionedDeviceRequest) (*empty.Empty, error) {
	log.Debugf("ProvisionedDeviceServiceServer.Create() called")
	// first check whether user is an authorized user
	_, err := a.auth.GetCredentials(ctx, auth.NewOptions())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	// Get Device provisioned data from PS
	provisioneddata, err := a.pscli.GetDeviceByID(ctx, &psPb.GetDeviceByIdRequest{ProvisionId: req.Device.ProvisionId})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get device from PS: %v, %v", req.Device.ProvisionId, err)
	}
	if provisioneddata.Status != "PROVISIONED" {
		return nil, status.Errorf(codes.NotFound, "Device not provisioned: %v", req.Device.ProvisionId)
	}
	if provisioneddata.Server != "" {
		return nil, status.Errorf(codes.Internal, "Device already registered at %v: %v", provisioneddata.Server, req.Device.ProvisionId)
	}

	//
	var devEUI lorawan.EUI64
	copy(devEUI[:], provisioneddata.DevEUI)
	if len(provisioneddata.DevEUI) != 8 {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid DevEUI %v", hex.EncodeToString(provisioneddata.DevEUI))
	}

	dpID, err := uuid.FromString(req.Device.DeviceProfileId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if valid, err := devmod.NewValidator(a.st).ValidateGlobalNodesAccess(ctx, authcus.Create, req.Device.ApplicationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	// if Name is "", set it to the DevEUI
	if req.Device.Name == "" {
		req.Device.Name = provisioneddata.Model + "_" + provisioneddata.SerialNumber
	}

	// Validate that application and device-profile are under the same
	// organization ID.
	app, err := application.GetApplication(ctx, req.Device.ApplicationId)
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

	//
	err = a.createDevice(ctx, app, devEUI, req, dpID, provisioneddata)
	if err != nil {
		return nil, err
	}

	//
	_, err = a.pscli.SetDeviceServer(ctx, &psPb.SetDeviceServerRequest{ProvisionId: req.Device.ProvisionId, Server: a.ServerAddr})
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// Get - returns the device matching the given Provision ID.
func (a *ProvisionedDeviceAPI) Get(ctx context.Context, req *api.GetProvisionedDeviceRequest) (*api.GetProvisionedDeviceResponse, error) {
	log.Debugf("ProvisionedDeviceServiceServer.Get() called")
	// first check whether user is an authorized user
	_, err := a.auth.GetCredentials(ctx, auth.NewOptions())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	//
	respcheck, err := a.pscli.IsDeviceExist(ctx, &psPb.IsDeviceExistRequest{ProvisionId: req.ProvisionId})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to check device existence: %v, %v", req.ProvisionId, err)
	}
	if !respcheck.Exist {
		return nil, status.Errorf(codes.NotFound, "Device not found: %v", req.ProvisionId)
	}

	//
	respdev, err := a.pscli.GetDeviceByID(ctx, &psPb.GetDeviceByIdRequest{ProvisionId: req.ProvisionId})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get device from PS: %v, %v", req.ProvisionId, err)
	}

	respmfg, err := a.pscli.GetManufacturerByID(ctx, &psPb.GetMfgByIdRequest{ManufacturerId: respdev.ManufacturerId})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get mfg info from PS: id=%v, %v", respdev.ManufacturerId, err)
	}
	if len(respmfg.Info) == 0 {
		return nil, status.Errorf(codes.Internal, "mfg info not found from PS: id=%v", respdev.ManufacturerId)
	}

	ret := api.GetProvisionedDeviceResponse{
		ProvisionId:      respdev.ProvisionId,
		ManufacturerId:   respdev.ManufacturerId,
		ManufacturerName: respmfg.Info[0].Name,
		Model:            respdev.Model,
		SerialNumber:     respdev.SerialNumber,
		FixedDevEUI:      respdev.FixedDevEUI,
		DevEUI:           hex.EncodeToString(respdev.DevEUI),
		Status:           respdev.Status,
		Server:           respdev.Server,
		TimeCreated:      respdev.TimeCreated,
		TimeProvisioned:  respdev.TimeProvisioned,
		TimeAddToServer:  respdev.TimeAddToServer,
	}

	return &ret, nil
}
