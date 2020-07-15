package serverinfo

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes/empty"

	pb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	api "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/m2m_client"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
)

// ServerInfoAPI defines the Server Info API structure
type ServerInfoAPI struct{}

// NewServerInfoAPI defines the Server Info API
func NewServerInfoAPI() *ServerInfoAPI {
	return &ServerInfoAPI{}
}

// GetAppserverVersion defines the Appserver Version response
func (s *ServerInfoAPI) GetAppserverVersion(ctx context.Context, req *empty.Empty) (*pb.GetAppserverVersionResponse, error) {
	return &pb.GetAppserverVersionResponse{Version: config.AppserverVersion}, nil
}

func (s *ServerInfoAPI) GetServerRegion(ctx context.Context, req *empty.Empty) (*pb.GetServerRegionResponse, error) {
	region := os.Getenv("SERVER_REGION")

	if region == pb.ServerRegion_name[int32(pb.ServerRegion_AVERAGE)] {
		return &pb.GetServerRegionResponse{ServerRegion: pb.ServerRegion_name[int32(pb.ServerRegion_AVERAGE)]}, nil
	}

	if region == pb.ServerRegion_name[int32(pb.ServerRegion_RESTRICTED)] {
		return &pb.GetServerRegionResponse{ServerRegion: pb.ServerRegion_name[int32(pb.ServerRegion_RESTRICTED)]}, nil
	}

	// pb.ServerRegion_NOT_DEFINED same as default region: pb.ServerRegion_AVERAGE
	return &pb.GetServerRegionResponse{ServerRegion: pb.ServerRegion_name[int32(pb.ServerRegion_NOT_DEFINED)]}, nil
}

// GetVersion defines the version of the machine to machine server API
func (s *ServerInfoAPI) GetMxprotocolServerVersion(ctx context.Context, req *empty.Empty) (*pb.GetMxprotocolServerVersionResponse, error) {
	log.WithField("", "").Info("grpc_api/GetVersion")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &pb.GetMxprotocolServerVersionResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	verClient := api.NewM2MServerInfoServiceClient(m2mClient)

	resp, err := verClient.GetVersion(ctx, req)
	if err != nil {
		return &pb.GetMxprotocolServerVersionResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &pb.GetMxprotocolServerVersionResponse{
		Version: resp.Version,
	}, nil
}
