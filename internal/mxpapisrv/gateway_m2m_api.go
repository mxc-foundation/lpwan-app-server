package mxpapisrv

import (
	"context"

	"github.com/brocaar/lorawan"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-m2m"
	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"
)

// GatewayM2MAPI exports the API for mxprotocol-server
type GatewayM2MAPI struct {
	st *pgstore.PgStore
}

// NewGatewayM2MAPI creates new GatewayM2MAPI
func NewGatewayM2MAPI(h *pgstore.PgStore) *GatewayM2MAPI {
	return &GatewayM2MAPI{
		st: h,
	}
}

// GetGatewayMacList defines the response of the Gateway MAC list
func (a *GatewayM2MAPI) GetGatewayMacList(ctx context.Context, req *empty.Empty) (*pb.GetGatewayMacListResponse, error) {
	gwMacList, err := a.st.GetAllGatewayMacList(ctx)
	if err != nil {
		return &pb.GetGatewayMacListResponse{}, status.Errorf(codes.DataLoss, err.Error())
	}

	return &pb.GetGatewayMacListResponse{GatewayMac: gwMacList}, nil
}

// GetGatewayByMac defines the request and response to the the gateway by MAC
func (a *GatewayM2MAPI) GetGatewayByMac(ctx context.Context, req *pb.GetGatewayByMacRequest) (*pb.GetGatewayByMacResponse, error) {
	var mac lorawan.EUI64
	resp := pb.GetGatewayByMacResponse{GwProfile: &pb.AppServerGatewayProfile{}}

	if err := mac.UnmarshalText([]byte(req.Mac)); err != nil {
		return &resp, status.Errorf(codes.InvalidArgument, err.Error())
	}

	gateway, err := a.st.GetGateway(ctx, mac, false)
	if err == errHandler.ErrDoesNotExist {
		return &resp, nil
	} else if err != nil {
		return &resp, status.Errorf(codes.InvalidArgument, err.Error())
	}

	resp.OrgId = gateway.OrganizationID
	resp.GwProfile.OrgId = gateway.OrganizationID
	resp.GwProfile.Mac = req.Mac
	resp.GwProfile.Name = gateway.Name
	resp.GwProfile.Description = gateway.Description
	resp.GwProfile.CreatedAt, _ = ptypes.TimestampProto(gateway.CreatedAt)

	return &resp, nil
}
