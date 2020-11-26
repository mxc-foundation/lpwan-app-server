package user

import (
	"context"

	pb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
)

// AuthenticateWeChatUser interacts with wechat open platform to authenticate wechat user
// then check binding status of this wechat user
func (a *Server) AuthenticateWeChatUser(ctx context.Context, req *pb.AuthenticateWeChatUserRequest) (*pb.AuthenticateWeChatUserResponse, error) {

	return &pb.AuthenticateWeChatUserResponse{}, nil
}

// BindExternalUser binds external user id to supernode user
func (a *Server) BindExternalUser(ctx context.Context, req *pb.BindExternalUserRequest) (*pb.BindExternalUserResponse, error) {

	return &pb.BindExternalUserResponse{}, nil
}

// RegisterExternalUser creates new supernode account then bind it with external user id
func (a *Server) RegisterExternalUser(ctx context.Context, req *pb.RegisterExternalUserRequest) (*pb.RegisterExternalUserResponse, error) {

	return &pb.RegisterExternalUserResponse{}, nil
}

// GetExternalUserIDFromUserID gets external user id by supernode user id
func (a *Server) GetExternalUserIDFromUserID(ctx context.Context, req *pb.GetExternalUserIDFromUserIDRequest) (*pb.GetExternalUserIDFromUserIDResponse, error) {

	return &pb.GetExternalUserIDFromUserIDResponse{}, nil
}

// UnbindExternalUser unbinds external user and supernode user account
func (a *Server) UnbindExternalUser(ctx context.Context, req *pb.UnbindExternalUserRequest) (*pb.UnbindExternalUserResponse, error) {

	return &pb.UnbindExternalUserResponse{}, nil
}
