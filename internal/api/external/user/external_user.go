package user

import (
	"context"
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/auth"
	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
	"github.com/mxc-foundation/lpwan-app-server/internal/jwt"
)

// AuthenticateWeChatUser interacts with wechat open platform to authenticate wechat user
// then check binding status of this wechat user
func (a *Server) AuthenticateWeChatUser(ctx context.Context, req *pb.AuthenticateWeChatUserRequest) (*pb.AuthenticateWeChatUserResponse, error) {
	body := auth.GetAccessTokenResponse{}
	if err := auth.GetAccessTokenFromCode(ctx, req.Code, a.config.WeChatLogin.AppID, a.config.WeChatLogin.Secret, &body); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
	}

	if body.UnionID == "" {
		user := auth.GetWeChatUserInfoResponse{}
		if err := auth.GetWeChatUserInfoFromAccessToken(ctx, body.AccessToken, body.OpenID, &user); err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}

		body.UnionID = user.UnionID
	}

	// check whether wechat user has already bound to existing account
	userID, err := a.store.GetUserIDByExternalUserID(ctx, auth.WECHAT, body.UnionID)
	if err == errHandler.ErrDoesNotExist {
		// authorized wechat user is not bound to any existing account
		wechatCredStr, err := json.Marshal(auth.WeChatAuth{
			AccessToken: body.AccessToken,
			OpenID:      body.OpenID,
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}

		claims := jwt.Claims{
			ExternalService: auth.WECHAT,
			ExternalCred:    string(wechatCredStr),
		}

		// not bound with any existing account, sign jwt with access_token and openid
		jwtWithLimited, err := a.jwtv.SignToken(claims, 600, []string{"authenticate-external"})
		if err != nil {
			log.Errorf("SignToken returned an error: %v", err)
			return nil, status.Errorf(codes.Internal, "couldn't create a token")
		}

		return &pb.AuthenticateWeChatUserResponse{Jwt: jwtWithLimited, BindingIsRequired: true}, nil

	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// wechat user has bound to existing account
	u, err := a.store.GetUserByID(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot get user by id: %d", userID)
	}

	// wechat user already bound with existing account, sign jwt with username and user id
	jwtNormal, err := a.jwtv.SignToken(jwt.Claims{Username: u.Email, UserID: u.ID}, 0, nil)
	if err != nil {
		log.Errorf("SignToken returned an error: %v", err)
		return nil, status.Errorf(codes.Internal, "couldn't create a token")
	}

	return &pb.AuthenticateWeChatUserResponse{Jwt: jwtNormal, BindingIsRequired: false}, nil
}

// BindExternalUser binds external user id to supernode user
func (a *Server) BindExternalUser(ctx context.Context, req *pb.BindExternalUserRequest) (*pb.BindExternalUserResponse, error) {
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithExternalLimited().WithAudience("authenticate-external"))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	if cred.ExternalUserService == auth.WECHAT {
		wechatAuthCred := auth.WeChatAuth{}
		if err := json.Unmarshal([]byte(cred.ExternalUserID), &wechatAuthCred); err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}

		// verify access_token and openid
		user := auth.GetWeChatUserInfoResponse{}
		if err := auth.GetWeChatUserInfoFromAccessToken(ctx, wechatAuthCred.AccessToken, wechatAuthCred.OpenID, &user); err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}

		// verify user credentials: req.Email, req.Password
		userEmail := normalizeUsername(req.Email)
		u, err := a.store.GetUserByEmail(ctx, userEmail)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "couldn't get info about the user: %s", err.Error())
		}
		if !u.IsActive {
			return nil, status.Error(codes.Unauthenticated, "inactive user")
		}
		if err := a.pwhasher.Validate(req.Password, u.PasswordHash); err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid email or password")
		}

		// check whether user has already bound with wechat account
		_, err = a.store.GetExternalUserIDByUserID(ctx, cred.ExternalUserService, u.ID)
		if err == nil {
			return &pb.BindExternalUserResponse{TryDifferentCredentials: true, Jwt: ""}, nil
		} else if err != errHandler.ErrDoesNotExist {
			return nil, status.Errorf(codes.Internal, err.Error())
		}

		// Bind wechat account with supernode account
		if err := a.store.AddExternalUserLogin(ctx, cred.ExternalUserService, u.ID, user.UnionID); err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}

		jwToken, err := a.jwtv.SignToken(jwt.Claims{UserID: u.ID, Username: u.Email}, 0, nil)
		if err != nil {
			log.Errorf("SignToken returned an error: %v", err)
			return nil, status.Errorf(codes.Internal, "couldn't create a token")
		}

		return &pb.BindExternalUserResponse{TryDifferentCredentials: false, Jwt: jwToken}, nil

	}

	return nil, status.Errorf(codes.Unavailable, "external service authentication not supported: %s", cred.ExternalUserService)
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
