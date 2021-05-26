package mqttauth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mxc-foundation/lpwan-app-server/api/extapi"
	"github.com/mxc-foundation/lpwan-app-server/internal/auth"
	"github.com/mxc-foundation/lpwan-app-server/internal/jwt"
)

// Server defines the MosquittoAuth Service Server API structure
type Server struct {
	auth auth.Authenticator
	st   Store
	jwtv *jwt.Validator
}

// NewServer returns a new MosquittoAuth Service Server
func NewServer(st Store, auth auth.Authenticator, jwtv *jwt.Validator) *Server {
	return &Server{
		auth: auth,
		st:   st,
		jwtv: jwtv,
	}
}

// Store defines db APIs for MosquittoAuth Service
type Store interface {
}

// GetJWT returns JWT for mosquitto auth with given org id
// Only accessible for authenticated supernode user
func (s *Server) GetJWT(ctx context.Context, req *pb.GetJWTRequest) (*pb.GetJWTResponse, error) {
	cred, err := s.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrganizationId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}

	if !cred.IsOrgUser {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied: %v", err)
	}

	jwToken, err := s.jwtv.SignToken(jwt.Claims{
		UserID:   cred.UserID,
		Username: cred.Username,
		Service:  auth.EMAIL,
	}, 0, []string{"mosquitto-auth"})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't create a token: %v", err)
	}

	return &pb.GetJWTResponse{JwtMqttAuth: jwToken}, nil
}

// JWTAuthentication will be called by mosquitto auth plugin JWT backend, request and response are also defined there
func (s *Server) JWTAuthentication(ctx context.Context, req *pb.JWTAuthenticationRequest) (*pb.JWTAuthenticationResponse, error) {
	_, err := s.auth.GetCredentials(ctx, auth.NewOptions().WithAudience("mosquitto-auth"))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}
	return &pb.JWTAuthenticationResponse{}, nil
}

// CheckACL will be called by mosquitto auth plugin JWT backend, request and response are also defined there
func (s *Server) CheckACL(ctx context.Context, req *pb.CheckACLRequest) (*pb.CheckACLResponse, error) {
	_, err := s.auth.GetCredentials(ctx, auth.NewOptions().WithAudience("mosquitto-auth"))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}
	return &pb.CheckACLResponse{}, nil
}
