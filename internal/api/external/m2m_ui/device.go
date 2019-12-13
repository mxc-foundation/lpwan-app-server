package m2m_ui

import (
	"context"
	m2m_api "github.com/mxc-foundation/lpwan-app-server/api/m2m_server"
	api "github.com/mxc-foundation/lpwan-app-server/api/m2m_ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external"
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

	getUserProfile, err := external.InternalUserAPI{}.Profile(ctx, nil)
	if err != nil {
		log.WithError(err).Error("Cannot get userprofile")
		return &api.GetDeviceListResponse{}, err
	}

	userProfile := api.GetDeviceListResponse.GetUserProfile(getUserProfile)

	return &api.GetDeviceListResponse{
		DevProfile:  devProfile,
		Count:       resp.Count,
		UserProfile: userProfile,
	}, nil
}

func (s *DeviceServerAPI) GetDeviceProfile(ctx context.Context, req *api.GetDeviceProfileRequest) (*api.GetDeviceProfileResponse, error) {
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

	getUserProfile, err := external.InternalUserAPI{}.Profile(ctx, nil)
	if err != nil {
		log.WithError(err).Error("Cannot get userprofile")
		return &api.GetDeviceProfileResponse{}, err
	}

	userProfile := api.GetDeviceProfileResponse.GetUserProfile(getUserProfile)

	return &api.GetDeviceProfileResponse{
		DevProfile:  devProfile,
		UserProfile: userProfile,
	}, nil
}

func (s *DeviceServerAPI) GetDeviceHistory(ctx context.Context, req *api.GetDeviceHistoryRequest) (*api.GetDeviceHistoryResponse, error) {
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

	getUserProfile, err := external.InternalUserAPI{}.Profile(ctx, nil)
	if err != nil {
		log.WithError(err).Error("Cannot get userprofile")
		return &api.GetDeviceHistoryResponse{}, err
	}

	userProfile := api.GetDeviceHistoryResponse.GetUserProfile(getUserProfile)

	return &api.GetDeviceHistoryResponse{
		DevHistory:  resp.DevHistory,
		UserProfile: userProfile,
	}, nil
}

func (s *DeviceServerAPI) SetDeviceMode(ctx context.Context, req *api.SetDeviceModeRequest) (*api.SetDeviceModeResponse, error) {
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

	getUserProfile, err := external.InternalUserAPI{}.Profile(ctx, nil)
	if err != nil {
		log.WithError(err).Error("Cannot get userprofile")
		return &api.SetDeviceModeResponse{}, err
	}

	userProfile := api.SetDeviceModeResponse.GetUserProfile(getUserProfile)

	return &api.SetDeviceModeResponse{
		Status:      resp.Status,
		UserProfile: userProfile,
	}, nil
}
