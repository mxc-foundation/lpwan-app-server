package external

import (
	"context"
	"os"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
)

// ServerInfoAPI defines the Server Info API structure
type ServerInfoAPI struct {
}

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