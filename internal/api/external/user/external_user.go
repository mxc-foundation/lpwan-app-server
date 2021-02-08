package user

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/auth"
	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
	"github.com/mxc-foundation/lpwan-app-server/internal/jwt"
)

// ExternalUser represents struct of table external_login
type ExternalUser struct {
	UserID           int64  `db:"user_id"`
	ServiceName      string `db:"service"`
	ExternalUserID   string `db:"external_id"`
	ExternalUsername string `db:"external_username"`
	Verification     string `db:"verification"`
}

func (a *Server) authenticateWeChatUser(ctx context.Context, code, appID, secret string) (*pb.AuthenticateWeChatUserResponse, error) {
	body := auth.GetAccessTokenResponse{}
	user := auth.GetWeChatUserInfoResponse{}

	if err := auth.GetAccessTokenFromCode(ctx, auth.URLStrGetAccessTokenFromCode, code, appID, secret, &body); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
	}

	if err := auth.GetWeChatUserInfoFromAccessToken(ctx, auth.URLStrGetWeChatUserInfoFromAccessToken, body.AccessToken, body.OpenID, &user); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
	}

	body.UnionID = user.UnionID

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
			Service:      auth.WECHAT,
			ExternalCred: string(wechatCredStr),
		}

		// not bound with any existing account, sign jwt with access_token and openid
		jwtWithLimited, err := a.jwtv.SignToken(claims, 600, []string{"authenticate-external"})
		if err != nil {
			log.Errorf("SignToken returned an error: %v", err)
			return nil, status.Errorf(codes.Internal, "couldn't create a token")
		}

		log.WithFields(log.Fields{
			"jwtWithLimited": jwtWithLimited,
		}).Debug("AuthenticateWeChatUser")

		return &pb.AuthenticateWeChatUserResponse{Jwt: jwtWithLimited, BindingIsRequired: true}, nil

	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// wechat user has bound to existing account
	u, err := a.store.GetUserByID(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot get user by id: %d", userID)
	}

	if !u.IsActive {
		return nil, status.Errorf(codes.Unauthenticated, "inactive user")
	}

	jwtNormal, err := a.loginWithExternalUser(ctx, u.ID, auth.WECHAT, u.Email, user.UnionID, user.NickName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	log.WithFields(log.Fields{
		"jwt": jwtNormal,
	}).Debug("AuthenticateWeChatUser")

	return &pb.AuthenticateWeChatUserResponse{Jwt: jwtNormal, BindingIsRequired: false}, nil
}

func (a *Server) loginWithExternalUser(ctx context.Context, userID int64, service, email, externalUserID, externalUsername string) (string, error) {
	// wechat user already bound with existing account, sign jwt with username and user id
	jwtNormal, err := a.jwtv.SignToken(jwt.Claims{Username: email, UserID: userID, Service: service}, 0, nil)
	if err != nil {
		log.Errorf("SignToken returned an error: %v", err)
		return "", fmt.Errorf("couldn't create a token")
	}

	// update external username
	_ = a.store.SetExternalUsername(ctx, service, externalUserID, externalUsername)
	// from the moment user successfully login with external user account, set user display name to external user's nickname
	_ = a.store.SetUserLastLogin(ctx, userID, externalUsername, service)

	return jwtNormal, nil
}

// AuthenticateWeChatUser interacts with wechat open platform to authenticate wechat user
// then check binding status of this wechat user
func (a *Server) AuthenticateWeChatUser(ctx context.Context, req *pb.AuthenticateWeChatUserRequest) (*pb.AuthenticateWeChatUserResponse, error) {
	log.WithFields(log.Fields{
		"code":   req.Code,
		"appid":  a.config.WeChatLogin.AppID,
		"secret": a.config.WeChatLogin.Secret,
	}).Debug("AuthenticateWeChatUser")

	return a.authenticateWeChatUser(ctx, req.Code, a.config.WeChatLogin.AppID, a.config.WeChatLogin.Secret)
}

// DebugAuthenticateWeChatUser will only be called by debug mode
func (a *Server) DebugAuthenticateWeChatUser(ctx context.Context, req *pb.AuthenticateWeChatUserRequest) (*pb.AuthenticateWeChatUserResponse, error) {
	log.WithFields(log.Fields{
		"code":   req.Code,
		"appid":  a.config.DebugWeChatLogin.AppID,
		"secret": a.config.DebugWeChatLogin.Secret,
	}).Debug("DebugAuthenticateWeChatUser")

	return a.authenticateWeChatUser(ctx, req.Code, a.config.DebugWeChatLogin.AppID, a.config.DebugWeChatLogin.Secret)
}

// BindExternalUser binds external user id to supernode user
func (a *Server) BindExternalUser(ctx context.Context, req *pb.BindExternalUserRequest) (*pb.BindExternalUserResponse, error) {
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithExternalLimited().WithAudience("authenticate-external"))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	if cred.Service == auth.WECHAT {
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
		externalUser, err := a.store.GetExternalUserByUserIDAndService(ctx, cred.Service, u.ID)
		if err == nil {
			if externalUser.ExternalUserID != cred.ExternalUserID {
				return &pb.BindExternalUserResponse{Jwt: ""}, nil
			}

			// when this API is called, if wechat user has been verified and bound to existing user, return jwt
			// so that caller logic can proceed
			jwToken, err := a.loginWithExternalUser(ctx, u.ID, cred.Service, u.Email, cred.ExternalUserID, cred.ExternalUsername)
			if err != nil {
				return nil, status.Errorf(codes.Internal, err.Error())
			}

			return &pb.BindExternalUserResponse{Jwt: jwToken}, nil
		} else if err != errHandler.ErrDoesNotExist {
			return nil, status.Errorf(codes.Internal, err.Error())
		}

		// Bind wechat account with supernode account
		if err := a.store.AddExternalUserLogin(ctx, ExternalUser{
			UserID:           u.ID,
			ServiceName:      cred.Service,
			ExternalUserID:   cred.ExternalUserID,
			ExternalUsername: cred.ExternalUsername,
		}); err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}

		jwToken, err := a.loginWithExternalUser(ctx, u.ID, cred.Service, u.Email, cred.ExternalUserID, cred.ExternalUsername)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}

		return &pb.BindExternalUserResponse{Jwt: jwToken}, nil

	}

	return nil, status.Errorf(codes.Unavailable, "external service authentication not supported: %s", cred.Service)
}

// RegisterExternalUser creates new supernode account then bind it with external user id
func (a *Server) RegisterExternalUser(ctx context.Context, req *pb.RegisterExternalUserRequest) (*pb.RegisterExternalUserResponse, error) {
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithExternalLimited().WithAudience("authenticate-external"))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	if cred.Service == auth.WECHAT {
		// create new user
		u, err := a.store.CreateUser(ctx, User{
			Email: req.Email,
		}, nil)
		if err != nil {
			if err == errHandler.ErrAlreadyExists {
				// if user already created but not yet activated, get user infomation and proceed with activation
				u, err = a.store.GetUserByEmail(ctx, req.Email)
				if err != nil {
					return nil, status.Error(codes.Internal, err.Error())
				}

				if u.IsActive {
					return nil, status.Errorf(codes.AlreadyExists,
						"user(%s) already exist, please bind the user instead of creating new", req.Email)
				}
			} else {
				return nil, status.Errorf(codes.Internal, err.Error())
			}
		}

		err = a.store.ActivateUser(ctx, u.ID, "", req.OrganizationName, req.OrganizationName)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}

		// bind new user with wechat account
		if err := a.store.AddExternalUserLogin(ctx, ExternalUser{
			UserID:           u.ID,
			ServiceName:      cred.Service,
			ExternalUserID:   cred.ExternalUserID,
			ExternalUsername: cred.ExternalUsername,
		}); err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}

		jwToken, err := a.loginWithExternalUser(ctx, u.ID, cred.Service, u.Email, cred.ExternalUserID, cred.ExternalUsername)
		if err != nil {
			log.Errorf("SignToken returned an error: %v", err)
			return nil, status.Errorf(codes.Internal, "couldn't create a token")
		}

		return &pb.RegisterExternalUserResponse{Jwt: jwToken}, nil
	}

	return nil, status.Errorf(codes.Unavailable, "external service authentication not supported: %s", cred.Service)
}

// UnbindExternalUser unbinds external user and supernode user account
func (a *Server) UnbindExternalUser(ctx context.Context, req *pb.UnbindExternalUserRequest) (*pb.UnbindExternalUserResponse, error) {

	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrganizationId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed : %s", err.Error())
	}

	if !cred.IsOrgAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	if err := a.store.DeleteExternalUserLogin(ctx, cred.UserID, req.Service); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	if cred.Service == req.Service {
		if err := a.store.SetUserLastLogin(ctx, cred.UserID, cred.Username, auth.EMAIL); err == nil {
			// re-sign jwt if user's last login service is updated successfully
			jwToken, err := a.jwtv.SignToken(jwt.Claims{
				UserID:   cred.UserID,
				Username: cred.Username,
				Service:  auth.EMAIL,
			}, 0, nil)
			if err == nil {
				return &pb.UnbindExternalUserResponse{Status: jwToken}, nil
			}
		}
	}

	return &pb.UnbindExternalUserResponse{}, nil
}

// VerifyEmail sends confirmation message to given email address
func (a *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrganizationId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed : %s", err.Error())
	}

	if !cred.IsOrgAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	// check whether email registered in shopify store
	_, err = a.getShopifyCustomerIDFromEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	// check whether external email already exist
	extUser, err := a.store.GetExternalUserByUsername(ctx, auth.SHOPIFY, req.Email)
	if err == nil {
		if extUser.Verification == "" {
			return nil, status.Errorf(codes.AlreadyExists, "email already exists")
		}
		// email exists, but haven't confirmed, just need to resend security token
		// send token to given email address
		if err = a.mailer.SendVerifyEmailConfirmation(req.Email, req.Language, extUser.Verification); err != nil {
			return nil, status.Errorf(codes.Unknown, err.Error())
		}
		return &pb.VerifyEmailResponse{Status: fmt.Sprintf("verification token has been sent to %s", req.Email)}, nil
	} else if err != errHandler.ErrDoesNotExist {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// check whether user already bound to one external email
	_, err = a.store.GetExternalUserByUserIDAndService(ctx, auth.SHOPIFY, cred.UserID)
	if err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "only one email can be bound")
	} else if err != errHandler.ErrDoesNotExist {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	token := OTPgen()
	if err = a.store.AddExternalUserLogin(ctx, ExternalUser{
		UserID:           cred.UserID,
		ServiceName:      auth.SHOPIFY,
		ExternalUserID:   "",
		ExternalUsername: req.Email,
		Verification:     token,
	}); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	// send token to given email address
	if err = a.mailer.SendVerifyEmailConfirmation(req.Email, req.Language, token); err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}

	return &pb.VerifyEmailResponse{Status: fmt.Sprintf("verification token has been sent to %s", req.Email)}, nil
}

func (a *Server) getShopifyCustomerIDFromEmail(ctx context.Context, email string) (string, error) {
	url := fmt.Sprintf("https://%s:%s@%s/admin/api/%s/customers/search.json?query=email:%s",
		a.config.ShopifyConfig.AdminAPI.APIKey, a.config.ShopifyConfig.AdminAPI.Secret, a.config.ShopifyConfig.AdminAPI.Hostname,
		a.config.ShopifyConfig.AdminAPI.APIVersion, email)

	var customers ShopifyCustomerList
	if err := auth.GetHTTPResponse(url, &customers, false); err != nil {
		return "", status.Errorf(codes.Internal, err.Error())
	}

	if len(customers.Customers) == 0 {
		return "", status.Errorf(codes.NotFound, "customer %s not found on store %s", email, a.config.ShopifyConfig.AdminAPI.StoreName)
	}

	return fmt.Sprintf("%d", customers.Customers[0].ID), nil
}

// ConfirmBindingEmail checks given token and bind
func (a *Server) ConfirmBindingEmail(ctx context.Context, req *pb.ConfirmBindingEmailRequest) (*pb.ConfirmBindingEmailResponse, error) {
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrganizationId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
	}

	if !cred.IsOrgAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	extUser, err := a.store.GetExternalUserByToken(ctx, auth.SHOPIFY, req.Token)
	if err != nil {
		if err == errHandler.ErrDoesNotExist {
			return nil, status.Errorf(codes.InvalidArgument, "token not found")
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	if cred.UserID != extUser.UserID {
		return nil, status.Errorf(codes.InvalidArgument, "incorrect token")
	}

	// get shopify customer id from email
	customerID, err := a.getShopifyCustomerIDFromEmail(ctx, extUser.ExternalUsername)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	extUser.ExternalUserID = customerID
	if err = a.store.ConfirmExternalUserID(ctx, extUser); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	go CheckNewOrders(ctx, cred.OrgID, extUser.UserID, a.config.ShopifyConfig, a.store)

	return &pb.ConfirmBindingEmailResponse{}, nil
}
