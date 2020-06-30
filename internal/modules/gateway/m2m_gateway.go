package gateway

import (
	"context"
	"strconv"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/user"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	m2mServer "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/m2m_client"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
)

// GetGatewayList defines the get Gateway list request and response
func (s *GatewayAPI) GetGatewayList(ctx context.Context, req *api.GetGatewayListRequest) (*api.GetGatewayListResponse, error) {
	logInfo := "api/appserver_serves_ui/GetGatewayList org=" + strconv.FormatInt(req.OrgId, 10)

	// verify if user is global admin
	userIsAdmin, err := user.GetUserAPI().Validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetGatewayListResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if userIsAdmin == false {
		if err := s.Validator.otpValidator.JwtValidator.Validate(ctx, authcus.ValidateOrganizationAccess(authcus.Read, req.OrgId)); err != nil {
			log.WithError(err).Error(logInfo)
			return &api.GetGatewayListResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetGatewayListResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	gwClient := m2mServer.NewGSGatewayServiceClient(m2mClient)

	resp, err := gwClient.GetGatewayList(ctx, &m2mServer.GetGatewayListRequest{
		OrgId:  req.OrgId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetGatewayListResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	var gatewayProfileList []*api.GSGatewayProfile
	for _, item := range resp.GwProfile {
		gatewayProfile := &api.GSGatewayProfile{
			Id:          item.Id,
			Mac:         item.Mac,
			FkGwNs:      item.FkGwNs,
			FkWallet:    item.FkWallet,
			Mode:        api.GatewayMode(item.Mode),
			CreateAt:    item.CreateAt,
			LastSeenAt:  item.LastSeenAt,
			OrgId:       item.OrgId,
			Description: item.Description,
			Name:        item.Name,
		}

		gatewayProfileList = append(gatewayProfileList, gatewayProfile)
	}

	return &api.GetGatewayListResponse{
		GwProfile: gatewayProfileList,
		Count:     resp.Count,
	}, status.Error(codes.OK, "")
}

// GetGatewayProfile defines the get Gateway Profile request and response
func (s *GatewayAPI) GetGatewayProfile(ctx context.Context, req *api.GetGSGatewayProfileRequest) (*api.GetGSGatewayProfileResponse, error) {
	logInfo := "api/appserver_serves_ui/GetGatewayProfile org=" + strconv.FormatInt(req.OrgId, 10)

	// verify if user is global admin
	userIsAdmin, err := user.GetUserAPI().Validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetGSGatewayProfileResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if userIsAdmin == false {
		if err := s.Validator.otpValidator.JwtValidator.Validate(ctx, authcus.ValidateOrganizationAccess(authcus.Read, req.OrgId)); err != nil {
			log.WithError(err).Error(logInfo)
			return &api.GetGSGatewayProfileResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetGSGatewayProfileResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	gwClient := m2mServer.NewGSGatewayServiceClient(m2mClient)

	resp, err := gwClient.GetGatewayProfile(ctx, &m2mServer.GetGSGatewayProfileRequest{
		OrgId:  req.OrgId,
		GwId:   req.GwId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetGSGatewayProfileResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetGSGatewayProfileResponse{
		GwProfile: &api.GSGatewayProfile{
			Id:          resp.GwProfile.Id,
			Mac:         resp.GwProfile.Mac,
			FkGwNs:      resp.GwProfile.FkGwNs,
			FkWallet:    resp.GwProfile.FkWallet,
			Mode:        api.GatewayMode(resp.GwProfile.Mode),
			CreateAt:    resp.GwProfile.CreateAt,
			LastSeenAt:  resp.GwProfile.LastSeenAt,
			OrgId:       resp.GwProfile.OrgId,
			Description: resp.GwProfile.Description,
			Name:        resp.GwProfile.Name,
		},
	}, status.Error(codes.OK, "")
}

// GetGatewayHistory defines the get Gateway History request and response
func (s *GatewayAPI) GetGatewayHistory(ctx context.Context, req *api.GetGatewayHistoryRequest) (*api.GetGatewayHistoryResponse, error) {
	logInfo := "api/appserver_serves_ui/GetGatewayHistory org=" + strconv.FormatInt(req.OrgId, 10)

	// verify if user is global admin
	userIsAdmin, err := user.GetUserAPI().Validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetGatewayHistoryResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if userIsAdmin == false {
		if err := s.Validator.otpValidator.JwtValidator.Validate(ctx, authcus.ValidateOrganizationAccess(authcus.Read, req.OrgId)); err != nil {
			log.WithError(err).Error(logInfo)
			return &api.GetGatewayHistoryResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetGatewayHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	gwClient := m2mServer.NewGSGatewayServiceClient(m2mClient)

	resp, err := gwClient.GetGatewayHistory(ctx, &m2mServer.GetGatewayHistoryRequest{
		OrgId:  req.OrgId,
		GwId:   req.GwId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetGatewayHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetGatewayHistoryResponse{
		GwHistory: resp.GwHistory,
	}, status.Error(codes.OK, "")
}

// SetGatewayMode defines the set Gateway mode request and response
func (s *GatewayAPI) SetGatewayMode(ctx context.Context, req *api.SetGatewayModeRequest) (*api.SetGatewayModeResponse, error) {
	logInfo := "api/appserver_serves_ui/SetGatewayMode org=" + strconv.FormatInt(req.OrgId, 10)

	// verify if user is global admin
	userIsAdmin, err := user.GetUserAPI().Validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.SetGatewayModeResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if userIsAdmin == false {
		if err := s.Validator.otpValidator.JwtValidator.Validate(ctx, authcus.ValidateOrganizationAccess(authcus.Read, req.OrgId)); err != nil {
			log.WithError(err).Error(logInfo)
			return &api.SetGatewayModeResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.SetGatewayModeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	gwClient := m2mServer.NewGSGatewayServiceClient(m2mClient)

	resp, err := gwClient.SetGatewayMode(ctx, &m2mServer.SetGatewayModeRequest{
		OrgId:  req.OrgId,
		GwId:   req.GwId,
		GwMode: m2mServer.GatewayMode(req.GwMode),
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.SetGatewayModeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.SetGatewayModeResponse{
		Status: resp.Status,
	}, status.Error(codes.OK, "")
}
