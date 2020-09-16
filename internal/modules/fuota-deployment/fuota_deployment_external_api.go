package fuotamod

import (
	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/application"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/device"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"

	pb "github.com/brocaar/chirpstack-api/go/v3/as/external/api"
	"github.com/brocaar/chirpstack-api/go/v3/common"
	"github.com/brocaar/lorawan"
	"github.com/brocaar/lorawan/band"

	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/networkserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
)

// FUOTADeploymentAPI exports the FUOTA deployment related functions.
type FUOTADeploymentAPI struct {
	st *store.Handler
}

// NewFUOTADeploymentAPI creates a new FUOTADeploymentAPI.
func NewFUOTADeploymentAPI() *FUOTADeploymentAPI {
	return &FUOTADeploymentAPI{
		st: Service.St,
	}
}

// CreateForDevice creates a deployment for the given DevEUI.
func (f *FUOTADeploymentAPI) CreateForDevice(ctx context.Context, req *pb.CreateFUOTADeploymentForDeviceRequest) (*pb.CreateFUOTADeploymentForDeviceResponse, error) {
	if req.FuotaDeployment == nil {
		return nil, status.Errorf(codes.InvalidArgument, "fuota_deployment must not be nil")
	}

	var devEUI lorawan.EUI64
	if err := devEUI.UnmarshalText([]byte(req.DevEui)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if valid, err := NewValidator().ValidateFUOTADeploymentsAccess(ctx, authcus.Create, 0, devEUI); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	n, err := storage.GetNetworkServerForDevEUI(ctx, storage.DB(), devEUI)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	nsClient, err := networkserver.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	versionResp, err := nsClient.GetVersion(ctx, &empty.Empty{})
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	var b band.Band

	switch versionResp.Region {
	case common.Region_EU868:
		b, err = band.GetConfig(band.EU868, false, lorawan.DwellTimeNoLimit)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
	case common.Region_US915:
		b, err = band.GetConfig(band.US915, false, lorawan.DwellTimeNoLimit)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
	case common.Region_CN779:
		b, err = band.GetConfig(band.CN779, false, lorawan.DwellTimeNoLimit)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
	case common.Region_EU433:
		b, err = band.GetConfig(band.EU433, false, lorawan.DwellTimeNoLimit)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
	case common.Region_AU915:
		b, err = band.GetConfig(band.AU915, false, lorawan.DwellTimeNoLimit)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
	case common.Region_CN470:
		b, err = band.GetConfig(band.CN470, false, lorawan.DwellTimeNoLimit)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
	case common.Region_AS923:
		b, err = band.GetConfig(band.AS923, false, lorawan.DwellTimeNoLimit)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
	case common.Region_KR920:
		b, err = band.GetConfig(band.KR920, false, lorawan.DwellTimeNoLimit)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
	case common.Region_IN865:
		b, err = band.GetConfig(band.IN865, false, lorawan.DwellTimeNoLimit)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
	case common.Region_RU864:
		b, err = band.GetConfig(band.RU864, false, lorawan.DwellTimeNoLimit)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
	default:
		return nil, status.Errorf(codes.Internal, "region %s is not implemented", versionResp.Region)
	}

	maxPLSize, err := b.GetMaxPayloadSizeForDataRateIndex("", "", int(req.FuotaDeployment.Dr))
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	fd := storage.FUOTADeployment{
		Name:             req.FuotaDeployment.Name,
		DR:               int(req.FuotaDeployment.Dr),
		Frequency:        int(req.FuotaDeployment.Frequency),
		Payload:          req.FuotaDeployment.Payload,
		FragSize:         maxPLSize.N - 3,
		Redundancy:       int(req.FuotaDeployment.Redundancy),
		MulticastTimeout: int(req.FuotaDeployment.MulticastTimeout),
	}

	switch req.FuotaDeployment.GroupType {
	case pb.MulticastGroupType_CLASS_C:
		fd.GroupType = storage.FUOTADeploymentGroupTypeC
	default:
		return nil, status.Errorf(codes.InvalidArgument, "group_type %s is not supported", req.FuotaDeployment.GroupType)
	}

	fd.UnicastTimeout, err = ptypes.Duration(req.FuotaDeployment.UnicastTimeout)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "unicast_timeout: %s", err)
	}

	err = storage.Transaction(func(ctx context.Context, handler *store.Handler) error {
		return storage.CreateFUOTADeploymentForDevice(ctx, handler, &fd, devEUI)
	})
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &pb.CreateFUOTADeploymentForDeviceResponse{
		Id: fd.ID.String(),
	}, nil
}

// Get returns the fuota deployment for the given id.
func (f *FUOTADeploymentAPI) Get(ctx context.Context, req *pb.GetFUOTADeploymentRequest) (*pb.GetFUOTADeploymentResponse, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "id: %s", err)
	}

	if valid, err := NewValidator().ValidateFUOTADeploymentAccess(ctx, authcus.Read, id); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	fd, err := storage.GetFUOTADeployment(ctx, storage.DB(), id, false)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := pb.GetFUOTADeploymentResponse{
		FuotaDeployment: &pb.FUOTADeployment{
			Id:               fd.ID.String(),
			Name:             fd.Name,
			Dr:               uint32(fd.DR),
			Frequency:        uint32(fd.Frequency),
			Payload:          fd.Payload,
			Redundancy:       uint32(fd.Redundancy),
			MulticastTimeout: uint32(fd.MulticastTimeout),
			UnicastTimeout:   ptypes.DurationProto(fd.UnicastTimeout),
			State:            string(fd.State),
		},
	}

	resp.CreatedAt, err = ptypes.TimestampProto(fd.CreatedAt)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	resp.UpdatedAt, err = ptypes.TimestampProto(fd.UpdatedAt)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	resp.FuotaDeployment.NextStepAfter, err = ptypes.TimestampProto(fd.NextStepAfter)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	switch fd.GroupType {
	case storage.FUOTADeploymentGroupTypeB:
		resp.FuotaDeployment.GroupType = pb.MulticastGroupType_CLASS_B
	case storage.FUOTADeploymentGroupTypeC:
		resp.FuotaDeployment.GroupType = pb.MulticastGroupType_CLASS_C
	default:
		return nil, status.Errorf(codes.Internal, "unexpected group-type: %s", fd.GroupType)
	}

	return &resp, nil
}

// List lists the fuota deployments.
func (f *FUOTADeploymentAPI) List(ctx context.Context, req *pb.ListFUOTADeploymentRequest) (*pb.ListFUOTADeploymentResponse, error) {
	var err error
	var idFilter bool

	filters := storage.FUOTADeploymentFilters{
		Limit:  int(req.Limit),
		Offset: int(req.Offset),
	}

	if req.ApplicationId != 0 {
		idFilter = true
		filters.ApplicationID = req.ApplicationId

		// validate that the client has access to the given application
		if valid, err := application.NewValidator().ValidateApplicationAccess(ctx, authcus.Read, req.ApplicationId); !valid || err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
		}
	}

	if req.DevEui != "" {
		idFilter = true
		if err := filters.DevEUI.UnmarshalText([]byte(req.DevEui)); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "dev_eui: %s", err)
		}

		// validate that the client has access to the given devEUI
		if valid, err := device.NewValidator().ValidateNodeAccess(ctx, authcus.Read, filters.DevEUI); !valid || err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
		}
	}

	if !idFilter {
		user, err := NewValidator().GetUser(ctx)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}

		if !user.IsGlobalAdmin {
			return nil, status.Errorf(codes.Unauthenticated, "client must be global admin for unfiltered request")
		}
	}

	count, err := storage.GetFUOTADeploymentCount(ctx, storage.DB(), filters)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	deployments, err := storage.GetFUOTADeployments(ctx, storage.DB(), filters)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return f.returnList(count, deployments)
}

// GetDeploymentDevice returns the deployment device.
func (f *FUOTADeploymentAPI) GetDeploymentDevice(ctx context.Context, req *pb.GetFUOTADeploymentDeviceRequest) (*pb.GetFUOTADeploymentDeviceResponse, error) {
	fuotaDeploymentID, err := uuid.FromString(req.FuotaDeploymentId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "fuota_deployment_id: %s", err)
	}

	var devEUI lorawan.EUI64
	err = devEUI.UnmarshalText([]byte(req.DevEui))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "dev_eui: %s", err)
	}

	if valid, err := NewValidator().ValidateFUOTADeploymentAccess(ctx, authcus.Read, fuotaDeploymentID); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	d, err := storage.GetDevice(ctx, storage.DB(), devEUI, false, true)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	fdd, err := storage.GetFUOTADeploymentDevice(ctx, storage.DB(), fuotaDeploymentID, devEUI)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := pb.GetFUOTADeploymentDeviceResponse{
		DeploymentDevice: &pb.FUOTADeploymentDeviceListItem{
			DevEui:       d.DevEUI.String(),
			DeviceName:   d.Name,
			ErrorMessage: fdd.ErrorMessage,
		},
	}

	switch fdd.State {
	case storage.FUOTADeploymentDevicePending:
		resp.DeploymentDevice.State = pb.FUOTADeploymentDeviceState_PENDING
	case storage.FUOTADeploymentDeviceSuccess:
		resp.DeploymentDevice.State = pb.FUOTADeploymentDeviceState_SUCCESS
	case storage.FUOTADeploymentDeviceError:
		resp.DeploymentDevice.State = pb.FUOTADeploymentDeviceState_ERROR
	default:
		return nil, status.Errorf(codes.Internal, "unexpected state: %s", fdd.State)
	}

	resp.DeploymentDevice.CreatedAt, err = ptypes.TimestampProto(fdd.CreatedAt)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	resp.DeploymentDevice.UpdatedAt, err = ptypes.TimestampProto(fdd.UpdatedAt)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &resp, nil
}

// ListDeploymentDevices lists the devices (and status) for the given fuota deployment ID.
func (f *FUOTADeploymentAPI) ListDeploymentDevices(ctx context.Context, req *pb.ListFUOTADeploymentDevicesRequest) (*pb.ListFUOTADeploymentDevicesResponse, error) {
	fuotaDeploymentID, err := uuid.FromString(req.FuotaDeploymentId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "fuota_deployment_id %s", err)
	}

	if valid, err := NewValidator().ValidateFUOTADeploymentAccess(ctx, authcus.Read, fuotaDeploymentID); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	count, err := storage.GetFUOTADeploymentDeviceCount(ctx, storage.DB(), fuotaDeploymentID)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	devices, err := storage.GetFUOTADeploymentDevices(ctx, storage.DB(), fuotaDeploymentID, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	out := pb.ListFUOTADeploymentDevicesResponse{
		TotalCount: int64(count),
		Result:     make([]*pb.FUOTADeploymentDeviceListItem, len(devices)),
	}

	for i := range devices {
		var err error

		dd := pb.FUOTADeploymentDeviceListItem{
			DevEui:       devices[i].DevEUI.String(),
			DeviceName:   devices[i].DeviceName,
			ErrorMessage: devices[i].ErrorMessage,
		}

		switch devices[i].State {
		case storage.FUOTADeploymentDevicePending:
			dd.State = pb.FUOTADeploymentDeviceState_PENDING
		case storage.FUOTADeploymentDeviceSuccess:
			dd.State = pb.FUOTADeploymentDeviceState_SUCCESS
		case storage.FUOTADeploymentDeviceError:
			dd.State = pb.FUOTADeploymentDeviceState_ERROR
		default:
			return nil, status.Errorf(codes.Internal, "unexpected state: %s", devices[i].State)
		}

		dd.CreatedAt, err = ptypes.TimestampProto(devices[i].CreatedAt)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}

		dd.UpdatedAt, err = ptypes.TimestampProto(devices[i].UpdatedAt)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}

		out.Result[i] = &dd
	}

	return &out, nil
}

func (f *FUOTADeploymentAPI) returnList(count int, deployments []storage.FUOTADeploymentListItem) (*pb.ListFUOTADeploymentResponse, error) {
	var err error

	resp := pb.ListFUOTADeploymentResponse{
		TotalCount: int64(count),
	}

	for _, fd := range deployments {
		item := pb.FUOTADeploymentListItem{
			Id:    fd.ID.String(),
			Name:  fd.Name,
			State: string(fd.State),
		}

		item.CreatedAt, err = ptypes.TimestampProto(fd.CreatedAt)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
		item.UpdatedAt, err = ptypes.TimestampProto(fd.UpdatedAt)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
		item.NextStepAfter, err = ptypes.TimestampProto(fd.NextStepAfter)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}

		resp.Result = append(resp.Result, &item)
	}

	return &resp, nil
}
