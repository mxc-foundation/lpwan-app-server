package device

import (
	"context"

	"github.com/brocaar/lorawan"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-m2m"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/application"
)

// DeviceM2MAPI exports the API for mxprotocol-server
type DeviceM2MAPI struct {
	Store DeviceStore
}

// NewDeviceM2MAPI creates new DeviceM2MAPI
func NewDeviceM2MAPI(api DeviceM2MAPI) *DeviceM2MAPI {
	return &DeviceM2MAPI{
		Store: api.Store,
	}
}

// GetDeviceDevEuiList defines the response of the Device DevEui list
func (a *DeviceM2MAPI) GetDeviceDevEuiList(ctx context.Context, req *empty.Empty) (*pb.GetDeviceDevEuiListResponse, error) {
	devEuiList, err := a.Store.GetAllDeviceEuis(ctx)
	if err != nil {
		return &pb.GetDeviceDevEuiListResponse{}, status.Errorf(codes.DataLoss, err.Error())
	}

	return &pb.GetDeviceDevEuiListResponse{DevEui: devEuiList}, nil
}

// GetDeviceByDevEui defines the request and response of the Device DevEui
func (a *DeviceM2MAPI) GetDeviceByDevEui(ctx context.Context, req *pb.GetDeviceByDevEuiRequest) (*pb.GetDeviceByDevEuiResponse, error) {
	var devEui lorawan.EUI64
	resp := pb.GetDeviceByDevEuiResponse{DevProfile: &pb.AppServerDeviceProfile{}}

	if err := devEui.UnmarshalText([]byte(req.DevEui)); err != nil {
		return &resp, status.Errorf(codes.InvalidArgument, err.Error())
	}

	device, err := a.Store.GetDevice(ctx, devEui, false, true)
	if err == storage.ErrDoesNotExist {
		return &resp, nil
	} else if err != nil {
		return &resp, status.Errorf(codes.Unknown, err.Error())
	}

	app, err := application.GetApplicationAPI().Store.GetApplication(ctx, device.ApplicationID)
	if err != nil {
		return &resp, status.Errorf(codes.Unknown, err.Error())
	}

	resp.OrgId = app.OrganizationID
	resp.DevProfile.DevEui = req.DevEui
	resp.DevProfile.Name = device.Name
	resp.DevProfile.ApplicationId = device.ApplicationID
	resp.DevProfile.CreatedAt, _ = ptypes.TimestampProto(device.CreatedAt)

	return &resp, nil
}
