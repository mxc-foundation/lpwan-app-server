package external

import (
	"context"
	"strconv"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	pb "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	m2mcli "github.com/mxc-foundation/lpwan-app-server/internal/mxp_portal"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/topup"
)

// TopUpServerAPI defines the topup server api structure
type TopUpServerAPI struct{}

// NewTopUpServerAPI validates the topup server api
func NewTopUpServerAPI() *TopUpServerAPI {
	return &TopUpServerAPI{}
}

// GetTopUpHistory defines the topup history request and response
func (s *TopUpServerAPI) GetTopUpHistory(ctx context.Context, req *api.GetTopUpHistoryRequest) (*api.GetTopUpHistoryResponse, error) {
	logInfo := "api/appserver_serves_ui/GetTopUpHistory org=" + strconv.FormatInt(req.OrgId, 10)

	if err := topup.NewValidator().IsOrgAdmin(ctx, req.OrgId); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	topupClient, err := m2mcli.GetTopupServiceClient()
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetTopUpHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := topupClient.GetTopUpHistory(ctx, &pb.GetTopUpHistoryRequest{
		OrgId:    req.OrgId,
		Currency: req.Currency,
		From:     req.From,
		Till:     req.Till,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetTopUpHistoryResponse{}, status.Errorf(codes.Unavailable, "call API in m2m failed: %v", err)
	}

	var topUpHistoryList []*api.TopUpHistory
	for _, item := range resp.TopupHistory {
		topUpHistory := &api.TopUpHistory{
			Amount:    item.Amount,
			Timestamp: item.Timestamp,
			TxHash:    item.TxHash,
		}

		topUpHistoryList = append(topUpHistoryList, topUpHistory)
	}

	return &api.GetTopUpHistoryResponse{
		TopupHistory: topUpHistoryList,
	}, status.Error(codes.OK, "")
}

// GetTopUpDestination defines the topup destination request and response
func (s *TopUpServerAPI) GetTopUpDestination(ctx context.Context, req *api.GetTopUpDestinationRequest) (*api.GetTopUpDestinationResponse, error) {
	logInfo := "api/appserver_serves_ui/GetTopUpDestination org=" + strconv.FormatInt(req.OrgId, 10)

	if err := topup.NewValidator().IsOrgAdmin(ctx, req.OrgId); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	topupClient, err := m2mcli.GetTopupServiceClient()
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetTopUpDestinationResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := topupClient.GetTopUpDestination(ctx, &pb.GetTopUpDestinationRequest{
		OrgId: req.OrgId,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetTopUpDestinationResponse{}, status.Errorf(codes.Unavailable, "call API in m2m failed: %v", err)
	}

	return &api.GetTopUpDestinationResponse{
		ActiveAccount: resp.ActiveAccount,
	}, status.Error(codes.OK, "")
}
