package external

import (
	"context"

	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	pb "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/serverinfo"
	"github.com/mxc-foundation/lpwan-app-server/internal/mxpcli"
)

// SettingsServerAPI defines the settings of the Server API structure
type SettingsServerAPI struct{}

// NewSettingsServerAPI defines the SettingsServerAPI Validator
func NewSettingsServerAPI() *SettingsServerAPI {
	return &SettingsServerAPI{}
}

// GetSettings defines the settings of the Server API request and response
func (s *SettingsServerAPI) GetSettings(ctx context.Context, req *api.GetSettingsRequest) (*api.GetSettingsResponse, error) {
	logInfo := "api/appserver_serves_ui/GetSettings"

	if err := serverinfo.NewValidator().IsGlobalAdmin(ctx); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	settingClient := mxpcli.Global.GetSettingsServiceClient()

	resp, err := settingClient.GetSettings(ctx, &pb.GetSettingsRequest{})
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
		LowBalanceWarning:    resp.LowBalanceWarning,
		DownlinkPrice:        resp.DownlinkPrice,
		SupernodeIncomeRatio: resp.SupernodeIncomeRatio,
		StakingInterest:      resp.StakingInterest,
		Compensation:         compFloat,
	}, status.Error(codes.OK, "")
}

// ModifySettings defines the modification of the Server API settings
func (s *SettingsServerAPI) ModifySettings(ctx context.Context, req *api.ModifySettingsRequest) (*api.ModifySettingsResponse, error) {
	logInfo := "api/appserver_serves_ui/ModifySettings"

	if err := serverinfo.NewValidator().IsGlobalAdmin(ctx); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	settingClient := mxpcli.Global.GetSettingsServiceClient()

	resp, err := settingClient.ModifySettings(ctx, &pb.ModifySettingsRequest{
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
