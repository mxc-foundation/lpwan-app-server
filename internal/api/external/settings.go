package external

import (
	"context"

	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	m2mServer "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/auth"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/m2m_client"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
)

// SettingsServerAPI defines the settings of the Server API structure
type SettingsServerAPI struct {
	validator auth.Validator
}

// NewSettingsServerAPI defines the SettingsServerAPI validator
func NewSettingsServerAPI(validator auth.Validator) *SettingsServerAPI {
	return &SettingsServerAPI{
		validator: validator,
	}
}

// GetSettings defines the settings of the Server API request and response
func (s *SettingsServerAPI) GetSettings(ctx context.Context, req *api.GetSettingsRequest) (*api.GetSettingsResponse, error) {
	logInfo := "api/appserver_serves_ui/GetSettings"

	// verify if user is global admin
	userIsAdmin, err := s.validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetSettingsResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if userIsAdmin == false {
		return &api.GetSettingsResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed")
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetSettingsResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	settingClient := m2mServer.NewSettingsServiceClient(m2mClient)

	resp, err := settingClient.GetSettings(ctx, &m2mServer.GetSettingsRequest{})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetSettingsResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	compensation, err := decimal.NewFromString(resp.Compensation)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't parse compensation: %v", err)
	}
	compFloat, _ := compensation.Float64()
	return &api.GetSettingsResponse{
		LowBalanceWarning:                resp.LowBalanceWarning,
		DownlinkPrice:                    resp.DownlinkPrice,
		SupernodeIncomeRatio:             resp.SupernodeIncomeRatio,
		StakingPercentage:                resp.StakingPercentage,
		StakingExpectedRevenuePercentage: resp.StakingExpectedRevenuePercentage,
		Compensation:                     compFloat,
	}, status.Error(codes.OK, "")
}

// ModifySettings defines the modification of the Server API settings
func (s *SettingsServerAPI) ModifySettings(ctx context.Context, req *api.ModifySettingsRequest) (*api.ModifySettingsResponse, error) {
	logInfo := "api/appserver_serves_ui/ModifySettings"

	// verify if user is global admin
	userIsAdmin, err := s.validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.ModifySettingsResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if userIsAdmin == false {
		return &api.ModifySettingsResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed")
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.ModifySettingsResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	settingClient := m2mServer.NewSettingsServiceClient(m2mClient)

	resp, err := settingClient.ModifySettings(ctx, &m2mServer.ModifySettingsRequest{
		LowBalanceWarning:          req.LowBalanceWarning,
		DownlinkFee:                req.DownlinkFee,
		TransactionPercentageShare: req.TransactionPercentageShare,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.ModifySettingsResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.ModifySettingsResponse{
		Status: resp.Status,
	}, status.Error(codes.OK, "")
}
