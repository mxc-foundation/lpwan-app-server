package m2m_ui

import (
	"context"
	m2m_api "github.com/mxc-foundation/lpwan-app-server/api/m2m_server"
	api "github.com/mxc-foundation/lpwan-app-server/api/m2m_ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/auth"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/m2m_client"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
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
	userProfile, res := auth.VerifyRequestViaAuthServer(ctx, s.serviceName, req.OrgId)

	if err := s.validator.Validate(ctx,
		auth.ValidateActiveUser()); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	log.WithField("orgId", req.OrgId).Info("grpc_api/GetDeviceList")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetDeviceListResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := m2mClient.GetDeviceList(ctx, &m2m_api.GetDeviceListRequest{
		OrgId:  req.OrgId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		return &api.GetDeviceListResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	devProfile := api.GetDeviceListResponse.GetDevProfile(&resp.DevProfile)
	
	return &api.GetDeviceListResponse{
		DevProfile:  devProfile,
		Count:       resp.Count,
		UserProfile: resp.UserProfile,
	}, nil
}

func (s *DeviceServerAPI) GetDeviceProfile(ctx context.Context, req *api.GetDeviceProfileRequest) (*api.GetDeviceProfileResponse, error) {
	if err := s.validator.Validate(ctx,
		auth.ValidateActiveUser()); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	log.WithField("orgId", req.OrgId).Info("grpc_api/GetDeviceProfile")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetDeviceProfileResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := m2mClient.GetDeviceProfile(ctx, &m2m_api.GetDeviceProfileRequest{
		OrgId: req.OrgId,
		DevId: req.DevId,
	})
	if err != nil {
		return &api.GetDeviceProfileResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	devProfile := api.GetDeviceProfileResponse.GetDevProfile(&resp.DevProfile)

	return &api.GetDeviceProfileResponse{
		DevProfile:  devProfile,
		UserProfile: resp.UserProfile,
	}, nil
}

func (s *DeviceServerAPI) GetDeviceHistory(ctx context.Context, req *api.GetDeviceHistoryRequest) (*api.GetDeviceHistoryResponse, error) {
	if err := s.validator.Validate(ctx,
		auth.ValidateActiveUser()); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	log.WithField("orgId", req.OrgId).Info("grpc_api/GetDeviceHistory")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetDeviceHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := m2mClient.GetDeviceHistory(ctx, &m2m_api.GetDeviceHistoryRequest{
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
		UserProfile: resp.UserProfile,
	}, nil
}

func (s *DeviceServerAPI) SetDeviceMode(ctx context.Context, req *api.SetDeviceModeRequest) (*api.SetDeviceModeResponse, error) {
	if err := s.validator.Validate(ctx,
		auth.ValidateActiveUser()); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	log.WithField("orgId", req.OrgId).Info("grpc_api/SetDeviceMode")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.SetDeviceModeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	devMode := m2m_api.DeviceMode(req.DevMode)

	resp, err := m2mClient.SetDeviceMode(ctx, &m2m_api.SetDeviceModeRequest{
		OrgId:   req.OrgId,
		DevId:   req.DevId,
		DevMode: devMode,
	})
	if err != nil {
		return &api.SetDeviceModeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.SetDeviceModeResponse{
		Status:      resp.Status,
		UserProfile: resp.UserProfile,
	}, nil
}
