package device

import (
	"context"
	"encoding/hex"

	"github.com/mxc-foundation/lpwan-app-server/internal/nscli"

	"github.com/brocaar/lorawan"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/mxc-foundation/lpwan-app-server/api/extapi"
	pb "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	psPb "github.com/mxc-foundation/lpwan-app-server/api/ps-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	"github.com/mxc-foundation/lpwan-app-server/internal/auth"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/device/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/mxpcli"
	"github.com/mxc-foundation/lpwan-app-server/internal/pscli"
)

// ProvisionedDeviceAPI exports the Gateway related functions.
type ProvisionedDeviceAPI struct {
	st                  Store
	ApplicationServerID uuid.UUID
	ServerAddr          string
	auth                auth.Authenticator
	psCli               psPb.DeviceProvisionClient
	mxpCli              pb.DSDeviceServiceClient
	nsCli               *nscli.Client
}

// NewDeviceProvisionAPI - creates a new DeviceProvisionAPI.
func NewDeviceProvisionAPI(st Store, auth auth.Authenticator, applicationID uuid.UUID,
	serverAddr string, psCli *pscli.Client, mxpCli *mxpcli.Client, nsCli *nscli.Client) *ProvisionedDeviceAPI {
	return &ProvisionedDeviceAPI{
		st:                  st,
		auth:                auth,
		ApplicationServerID: applicationID,
		ServerAddr:          serverAddr,
		psCli:               psCli.GetDeviceProvisionServiceClient(),
		mxpCli:              mxpCli.GetM2MDeviceServiceClient(),
		nsCli:               nsCli,
	}
}

// Create - creates the given device.
func (a *ProvisionedDeviceAPI) Create(ctx context.Context, req *api.CreateRequest) (*api.CreateResponse, error) {
	log.WithFields(log.Fields{
		"provision_id":   req.ProvisionId,
		"application_id": req.ApplicationId,
	}).Debugf("ProvisionedDeviceServiceServer.Create() called")

	// first check whether user is an authorized user
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrganizationId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	// only organizaiton admin or device admin can create device
	if !cred.IsDeviceAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied: neither org admin nor device admin")
	}

	application, err := a.st.GetApplicationWithIDAndOrganizationID(ctx, req.ApplicationId, req.OrganizationId)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	dpID, err := uuid.FromString(req.DeviceProfileId)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	deviceProfile, err := a.st.GetDeviceProfileWithIDAndOrganizationID(ctx, dpID, req.OrganizationId, false)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	// get device
	d, dKeys, err := a.getDeviceAttributes(ctx, req.ProvisionId, req.DeviceProfileId, req.ApplicationId)
	if err != nil {
		return nil, err
	}

	nsClient, err := a.nsCli.GetNetworkServerServiceClient(deviceProfile.NetworkServerID)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	if err := CreateDevice(ctx, a.st, d, &application, a.ApplicationServerID, a.mxpCli, nsClient); err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	// add additional attributes for device provisioning
	if err := a.st.UpdateDeviceWithDevProvisioingAttr(ctx, d); err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	if err := a.st.CreateDeviceKeys(ctx, dKeys); err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	_, err = a.psCli.SetDeviceServer(ctx, &psPb.SetDeviceServerRequest{ProvisionId: req.ProvisionId, Server: a.ServerAddr})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	return &api.CreateResponse{}, nil
}

func (a *ProvisionedDeviceAPI) getDeviceAttributes(ctx context.Context, provisionID,
	deviceProfileID string, applicationID int64) (*data.Device, *data.DeviceKeys, error) {
	respcheck, err := a.psCli.IsDeviceExist(ctx, &psPb.IsDeviceExistRequest{ProvisionId: provisionID})
	if err != nil {
		return nil, nil, status.Errorf(codes.Internal, "Failed to check device existence: %v, %v", provisionID, err)
	}
	if !respcheck.Exist {
		return nil, nil, status.Errorf(codes.NotFound, "Device not found: %v", provisionID)
	}
	respdev, err := a.psCli.GetDeviceByID(ctx, &psPb.GetDeviceByIdRequest{ProvisionId: provisionID})
	if err != nil {
		return nil, nil, status.Errorf(codes.Internal, "Failed to get device from PS: %v, %v", provisionID, err)
	}
	respmfg, err := a.psCli.GetManufacturerByID(ctx, &psPb.GetMfgByIdRequest{ManufacturerId: respdev.ManufacturerId})
	if err != nil {
		return nil, nil, status.Errorf(codes.Internal, "Failed to get mfg info from PS: id=%v, %v", respdev.ManufacturerId, err)
	}
	if len(respmfg.Info) == 0 {
		return nil, nil, status.Errorf(codes.Internal, "mfg info not found from PS: id=%v", respdev.ManufacturerId)
	}
	if respdev.Status != "PROVISIONED" {
		return nil, nil, status.Errorf(codes.FailedPrecondition, "Device not provisioned: %v", provisionID)
	}
	if respdev.Server != "" {
		return nil, nil, status.Errorf(codes.FailedPrecondition, "Device already registered at %v: %v", respdev.Server, provisionID)
	}

	var devEUI lorawan.EUI64
	copy(devEUI[:], respdev.DevEUI)
	if len(respdev.DevEUI) != 8 {
		return nil, nil, status.Errorf(codes.InvalidArgument, "Invalid DevEUI %v", hex.EncodeToString(respdev.DevEUI))
	}
	dpID, err := uuid.FromString(deviceProfileID)
	if err != nil {
		return nil, nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	// Set Device struct.
	d := data.Device{
		DevEUI:          devEUI,
		ApplicationID:   applicationID,
		DeviceProfileID: dpID,
		Name:            respdev.Model + "_" + respdev.SerialNumber,
		Description:     "",
		// attributes for device provisioning only
		ProvisionID:  provisionID,
		Model:        respdev.Model,
		SerialNumber: respdev.SerialNumber,
		Manufacturer: respmfg.Info[0].Name,
	}

	// Set Keys
	var appKey lorawan.AES128Key
	var nwkKey lorawan.AES128Key
	var genAppKey lorawan.AES128Key

	copy(appKey[:], respdev.AppKey)
	copy(nwkKey[:], respdev.NwkKey)
	copy(genAppKey[:], respdev.NwkKey)

	dKeys := data.DeviceKeys{
		DevEUI:    devEUI,
		NwkKey:    nwkKey,
		AppKey:    appKey,
		GenAppKey: genAppKey,
	}
	return &d, &dKeys, nil
}
