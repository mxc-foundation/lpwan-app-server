package m2m_ui

import (
	"context"
	api "github.com/mxc-foundation/lpwan-app-server/api/m2m_ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/auth"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/m2m_client"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeviceServerAPI struct {
	validator auth.Validator
}

func NewDeviceServerAPI(validator auth.Validator) *DeviceServerAPI {
	return &DeviceServerAPI{
		validator: validator,
	}
}

func (s *DeviceServerAPI) GetDeviceList(ctx context.Context, req *api.GetDeviceListRequest) (*api.GetDeviceListResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetDeviceList")

	prof, err := getUserProfileByJwt(s.validator, ctx, req.OrgId)
	if err != nil{
		return &api.GetDeviceListResponse{}, status.Errorf(codes.Unauthenticated, err.Error())
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetDeviceListResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	devClient := api.NewDeviceServiceClient(m2mClient)

	resp, err := devClient.GetDeviceList(ctx, &api.GetDeviceListRequest{
		OrgId:  req.OrgId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		return &api.GetDeviceListResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetDeviceListResponse{
		DevProfile:  resp.DevProfile,
		Count:       resp.Count,
		UserProfile: &prof,
	}, nil
}

func (s *DeviceServerAPI) GetDeviceProfile(ctx context.Context, req *api.GetDeviceProfileRequest) (*api.GetDeviceProfileResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetDeviceProfile")

	prof, err := getUserProfileByJwt(s.validator, ctx, req.OrgId)
	if err != nil{
		return &api.GetDeviceProfileResponse{}, status.Errorf(codes.Unauthenticated, err.Error())
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetDeviceProfileResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	devClient := api.NewDeviceServiceClient(m2mClient)

	resp, err := devClient.GetDeviceProfile(ctx, &api.GetDeviceProfileRequest{
		OrgId: req.OrgId,
		DevId: req.DevId,
	})
	if err != nil {
		return &api.GetDeviceProfileResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetDeviceProfileResponse{
		DevProfile:  resp.DevProfile,
		UserProfile: &prof,
	}, nil
}

func (s *DeviceServerAPI) GetDeviceHistory(ctx context.Context, req *api.GetDeviceHistoryRequest) (*api.GetDeviceHistoryResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetDeviceHistory")

	prof, err := getUserProfileByJwt(s.validator, ctx, req.OrgId)
	if err != nil{
		return &api.GetDeviceHistoryResponse{}, status.Errorf(codes.Unauthenticated, err.Error())
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetDeviceHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	devClient := api.NewDeviceServiceClient(m2mClient)

	resp, err := devClient.GetDeviceHistory(ctx, &api.GetDeviceHistoryRequest{
		OrgId:  req.OrgId,
		DevId:  req.DevId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		return &api.GetDeviceHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetDeviceHistoryResponse{
		DevHistory:  resp.DevHistory,
		UserProfile: &prof,
	}, nil
}

func (s *DeviceServerAPI) SetDeviceMode(ctx context.Context, req *api.SetDeviceModeRequest) (*api.SetDeviceModeResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/SetDeviceMode")

	prof, err := getUserProfileByJwt(s.validator, ctx, req.OrgId)
	if err != nil{
		return &api.SetDeviceModeResponse{}, status.Errorf(codes.Unauthenticated, err.Error())
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.SetDeviceModeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	devClient := api.NewDeviceServiceClient(m2mClient)

	resp, err := devClient.SetDeviceMode(ctx, &api.SetDeviceModeRequest{
		OrgId:   req.OrgId,
		DevId:   req.DevId,
		DevMode: req.DevMode,
	})
	if err != nil {
		return &api.SetDeviceModeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.SetDeviceModeResponse{
		Status:      resp.Status,
		UserProfile: &prof,
	}, nil
}
