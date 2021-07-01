package mxpapisrv

import (
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/brocaar/lorawan"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-m2m"
	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"
)

// DeviceM2MAPI exports the API to mxprotocol client
type DeviceM2MAPI struct {
	st *pgstore.PgStore
}

// NewDeviceM2MAPI creates new DeviceM2MAPI
func NewDeviceM2MAPI(h *pgstore.PgStore) *DeviceM2MAPI {
	return &DeviceM2MAPI{
		st: h,
	}
}

// GetDeviceDevEuiList defines the response of the Device DevEui list
func (a *DeviceM2MAPI) GetDeviceDevEuiList(ctx context.Context, req *empty.Empty) (*pb.GetDeviceDevEuiListResponse, error) {
	devEuiList, err := a.st.GetAllDeviceEuis(ctx)
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

	device, err := a.st.GetDevice(ctx, devEui, false)
	if err == errHandler.ErrDoesNotExist {
		return &resp, nil
	} else if err != nil {
		return &resp, status.Errorf(codes.Unknown, err.Error())
	}

	app, err := a.st.GetApplication(ctx, device.ApplicationID)
	if err != nil {
		return &resp, status.Errorf(codes.Unknown, err.Error())
	}

	resp.OrgId = app.OrganizationID
	resp.DevProfile.DevEui = req.DevEui
	resp.DevProfile.Name = device.Name
	resp.DevProfile.ApplicationId = device.ApplicationID
	resp.DevProfile.CreatedAt = timestamppb.New(device.CreatedAt)

	return &resp, nil
}
