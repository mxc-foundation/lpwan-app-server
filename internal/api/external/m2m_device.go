package external

import (
	"context"
	"strconv"

	api "github.com/mxc-foundation/lpwan-app-server/api/appserver_serves_ui"
	m2mServer "github.com/mxc-foundation/lpwan-app-server/api/m2m_serves_appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/auth"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/m2m_client"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeviceServerAPI defines the device server api structure
type DeviceServerAPI struct {
	validator auth.Validator
}

// NewDeviceServerAPI validates the new devices server api
func NewDeviceServerAPI(validator auth.Validator) *DeviceServerAPI {
	return &DeviceServerAPI{
		validator: validator,
	}
}

// GetDeviceList defines the get device list request and response
func (s *DeviceServerAPI) GetDeviceList(ctx context.Context, req *api.GetDeviceListRequest) (*api.GetDeviceListResponse, error) {
	logInfo := "api/appserver_serves_ui/GetDeviceList org=" + strconv.FormatInt(req.OrgId, 10)

	// verify if user is global admin
	userIsAdmin, err := s.validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetDeviceListResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if userIsAdmin == false {
		if err := s.validator.Validate(ctx, auth.ValidateOrganizationAccess(auth.Read, req.OrgId)); err != nil {
			log.WithError(err).Error(logInfo)
			return &api.GetDeviceListResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetDeviceListResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	devClient := m2mServer.NewDSDeviceServiceClient(m2mClient)

	resp, err := devClient.GetDeviceList(ctx, &m2mServer.GetDeviceListRequest{
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

		deviceProfileList = append(deviceProfileList, deviceProfile)
	}

	return &api.GetDeviceListResponse{
		DevProfile: deviceProfileList,
		Count:      resp.Count,
	}, status.Error(codes.OK, "")
}

// GetDeviceProfile defines the function to get the device profile
func (s *DeviceServerAPI) GetDeviceProfile(ctx context.Context, req *api.GetDSDeviceProfileRequest) (*api.GetDSDeviceProfileResponse, error) {
	logInfo := "api/appserver_serves_ui/GetDeviceProfile org=" + strconv.FormatInt(req.OrgId, 10)

	// verify if user is global admin
	userIsAdmin, err := s.validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetDSDeviceProfileResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if userIsAdmin == false {
		if err := s.validator.Validate(ctx, auth.ValidateOrganizationAccess(auth.Read, req.OrgId)); err != nil {
			log.WithError(err).Error(logInfo)
			return &api.GetDSDeviceProfileResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetDSDeviceProfileResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	devClient := m2mServer.NewDSDeviceServiceClient(m2mClient)

	resp, err := devClient.GetDeviceProfile(ctx, &m2mServer.GetDSDeviceProfileRequest{
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
func (s *DeviceServerAPI) GetDeviceHistory(ctx context.Context, req *api.GetDeviceHistoryRequest) (*api.GetDeviceHistoryResponse, error) {
	logInfo := "api/appserver_serves_ui/GetDeviceHistory org=" + strconv.FormatInt(req.OrgId, 10)

	// verify if user is global admin
	userIsAdmin, err := s.validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetDeviceHistoryResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if userIsAdmin == false {
		if err := s.validator.Validate(ctx, auth.ValidateOrganizationAccess(auth.Read, req.OrgId)); err != nil {
			log.WithError(err).Error(logInfo)
			return &api.GetDeviceHistoryResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetDeviceHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
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
		return &api.GetDeviceHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetDeviceHistoryResponse{
		DevHistory:  resp.DevHistory,
	}, status.Error(codes.OK, "")
}

// SetDeviceMode defines the set device mode request and response
func (s *DeviceServerAPI) SetDeviceMode(ctx context.Context, req *api.SetDeviceModeRequest) (*api.SetDeviceModeResponse, error) {
	logInfo := "api/appserver_serves_ui/SetDeviceMode org=" + strconv.FormatInt(req.OrgId, 10)

	// verify if user is global admin
	userIsAdmin, err := s.validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.SetDeviceModeResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if userIsAdmin == false {
		if err := s.validator.Validate(ctx, auth.ValidateOrganizationAccess(auth.Read, req.OrgId)); err != nil {
			log.WithError(err).Error(logInfo)
			return &api.SetDeviceModeResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
	} else {
		// global admin should not chagne mode of a device for organization
		return &api.SetDeviceModeResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed")
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.SetDeviceModeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	devClient := m2mServer.NewDSDeviceServiceClient(m2mClient)

	resp, err := devClient.SetDeviceMode(ctx, &m2mServer.SetDeviceModeRequest{
		OrgId:   req.OrgId,
		DevId:   req.DevId,
		DevMode: m2mServer.DeviceMode(req.DevMode),
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.SetDeviceModeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.SetDeviceModeResponse{
		Status:      resp.Status,
	}, status.Error(codes.OK, "")
}
