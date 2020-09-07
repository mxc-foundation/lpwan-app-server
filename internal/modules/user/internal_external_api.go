package user

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/serverinfo"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	inpb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/email"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
)

// InternalUserAPI exports the internal User related functions.
type InternalUserAPI struct {
	st *store.Handler
}

// NewInternalUserAPI creates a new InternalUserAPI.
func NewInternalUserAPI() *InternalUserAPI {
	return &InternalUserAPI{
		st: Service.St,
	}
}

// Login validates the login request and returns a JWT token.
func (a *InternalUserAPI) Login(ctx context.Context, req *inpb.LoginRequest) (*inpb.LoginResponse, error) {
	userEmail := normalizeUsername(req.Username)
	err := a.st.LoginUserByPassword(ctx, userEmail, req.Password, Service.pwh)
	if nil != err {
		return nil, helpers.ErrToRPCError(err)
	}

	user, err := a.st.GetUserByUsername(ctx, userEmail)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't get info about the user: %s", err.Error())
	}

	if !user.IsActive {
		return nil, status.Error(codes.Unauthenticated, "incactive user")
	}

	ttl := 60 * int64(user.SessionTTL)
	var audience []string

	is2fa, err := NewValidator().Is2FAEnabled(ctx, user.UserEmail)
	if err != nil {
		ctxlogrus.Extract(ctx).WithError(err).Error("couldn't get 2fa status")
		return nil, status.Error(codes.Internal, "couldn't get 2fa status")
	}
	if !config.C.General.Enable2FALogin {
		is2fa = false
	}
	if is2fa {
		// if 2fa is enabled we issue token that is only valid for 10 minutes
		// and is only good to perform second factor authentication. If second
		// factor authentication has been successful then it will return to the
		// user another token, that provides access to all the api
		ttl = 600
		audience = []string{"login-2fa"}
	}

	jwt, err := NewValidator().SignJWToken(user.UserEmail, ttl, audience)
	if err != nil {
		log.Errorf("SignToken returned an error: %v", err)
		return nil, status.Errorf(codes.Internal, "couldn't create a token")
	}

	return &inpb.LoginResponse{Jwt: jwt, Is_2FaRequired: is2fa}, nil
}

// Login2FA performs second factor authentication. It requires user to have
// already passed password check and checks if the OTP code is valid. If it is
// it returns JWT with access to the api.
func (a *InternalUserAPI) Login2FA(ctx context.Context, req *inpb.Login2FARequest) (*inpb.LoginResponse, error) {
	u, err := NewValidator().GetUser(ctx, authcus.WithAudience("login-2fa"), authcus.WithValidOTP())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "not authenticated")
	}

	user, err := a.st.GetUserByUsername(ctx, u.UserEmail)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't get info about the user")
	}

	jwt, err := NewValidator().SignJWToken(u.UserEmail, 60*int64(user.SessionTTL), nil)
	if err != nil {
		log.Errorf("SignToken returned an error: %v", err)
		return nil, status.Error(codes.Internal, "couldn't create a token")
	}

	return &inpb.LoginResponse{Jwt: jwt}, nil
}

// IsPassVerifyingGoogleRecaptcha defines the response to pass the google recaptcha verification
func IsPassVerifyingGoogleRecaptcha(response string, remoteip string) (*inpb.GoogleRecaptchaResponse, error) {
	secret := config.C.Recaptcha.Secret
	postURL := config.C.Recaptcha.HostServer

	postStr := url.Values{"secret": {secret}, "response": {response}, "remoteip": {remoteip}}
	/* #nosec */
	responsePost, err := http.PostForm(postURL, postStr)

	if err != nil {
		log.Warn(err.Error())
		return &inpb.GoogleRecaptchaResponse{}, err
	}

	defer func() {
		err := responsePost.Body.Close()
		if err != nil {
			log.WithError(err).Error("cannot close the responsePost body.")
		}
	}()

	body, err := ioutil.ReadAll(responsePost.Body)

	if err != nil {
		log.Warn(err.Error())
		return &inpb.GoogleRecaptchaResponse{}, err
	}

	gresponse := &inpb.GoogleRecaptchaResponse{}
	err = json.Unmarshal(body, &gresponse)
	if err != nil {
		fmt.Println("unmarshal response", err)
	}

	return gresponse, nil
}

// GetVerifyingGoogleRecaptcha defines the request and response to verify the google recaptcha
func (a *InternalUserAPI) GetVerifyingGoogleRecaptcha(ctx context.Context, req *inpb.GoogleRecaptchaRequest) (*inpb.GoogleRecaptchaResponse, error) {
	res, err := IsPassVerifyingGoogleRecaptcha(req.Response, req.Remoteip)
	if err != nil {
		log.WithError(err).Error("Cannot verify from google recaptcha")
		return &inpb.GoogleRecaptchaResponse{}, err
	}

	return &inpb.GoogleRecaptchaResponse{Success: res.Success, ChallengeTs: res.ChallengeTs, Hostname: res.Hostname}, nil
}

func OTPgen() string {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	otp := make([]byte, 6)
	n, err := io.ReadAtLeast(rand.Reader, otp, 6)
	if n != 6 {
		panic(err)
	}
	for i := 0; i < len(otp); i++ {
		otp[i] = table[int(otp[i])%len(table)]
	}
	return string(otp)
}

// Profile returns the user profile.
func (a *InternalUserAPI) Profile(ctx context.Context, req *empty.Empty) (*inpb.ProfileResponse, error) {
	if valid, err := NewValidator().ValidateActiveUser(ctx); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	user, err := NewValidator().GetUser(ctx)
	if nil != err {
		return nil, helpers.ErrToRPCError(err)
	}

	prof, err := a.st.GetProfile(ctx, user.ID)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := inpb.ProfileResponse{
		User: &inpb.User{
			Id:         prof.User.ID,
			Email:      prof.User.UserEmail,
			SessionTtl: prof.User.SessionTTL,
			IsAdmin:    prof.User.IsAdmin,
			IsActive:   prof.User.IsActive,
		},
	}

	for _, org := range prof.Organizations {
		row := inpb.OrganizationLink{
			OrganizationId:   org.ID,
			OrganizationName: org.Name,
			IsAdmin:          org.IsAdmin,
			IsDeviceAdmin:    org.IsDeviceAdmin,
			IsGatewayAdmin:   org.IsGatewayAdmin,
		}

		row.CreatedAt, err = ptypes.TimestampProto(org.CreatedAt)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
		row.UpdatedAt, err = ptypes.TimestampProto(org.UpdatedAt)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}

		resp.Organizations = append(resp.Organizations, &row)
	}

	return &resp, nil
}

// Branding returns UI branding.
func (a *InternalUserAPI) Branding(ctx context.Context, req *empty.Empty) (*inpb.BrandingResponse, error) {
	resp := inpb.BrandingResponse{
		LogoPath: "https://" + serverinfo.Service.SupernodeAddr + "/branding.png",
	}

	return &resp, nil
}

// GlobalSearch performs a global search.
func (a *InternalUserAPI) GlobalSearch(ctx context.Context, req *inpb.GlobalSearchRequest) (*inpb.GlobalSearchResponse, error) {
	if valid, err := NewValidator().ValidateActiveUser(ctx); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	user, err := NewValidator().GetUser(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	results, err := storage.GlobalSearch(ctx, storage.DB(), user.ID, user.IsGlobalAdmin, req.Search, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	var out inpb.GlobalSearchResponse

	for _, r := range results {
		res := inpb.GlobalSearchResult{
			Kind:  r.Kind,
			Score: float32(r.Score),
		}

		if r.OrganizationID != nil {
			res.OrganizationId = *r.OrganizationID
		}
		if r.OrganizationName != nil {
			res.OrganizationName = *r.OrganizationName
		}

		if r.ApplicationID != nil {
			res.ApplicationId = *r.ApplicationID
		}
		if r.ApplicationName != nil {
			res.ApplicationName = *r.ApplicationName
		}

		if r.DeviceDevEUI != nil {
			res.DeviceDevEui = r.DeviceDevEUI.String()
		}
		if r.DeviceName != nil {
			res.DeviceName = *r.DeviceName
		}

		if r.GatewayMAC != nil {
			res.GatewayMac = r.GatewayMAC.String()
		}
		if r.GatewayName != nil {
			res.GatewayName = *r.GatewayName
		}

		out.Result = append(out.Result, &res)
	}

	return &out, nil
}

func normalizeUsername(userEmail string) string {
	return strings.ToLower(userEmail)
}

// RegisterUser adds new user and sends activation userEmail
func (a *InternalUserAPI) RegisterUser(ctx context.Context, req *inpb.RegisterUserRequest) (*empty.Empty, error) {
	logInfo := "api/appserver_serves_ui/RegisterUser"

	userEmail := normalizeUsername(req.Email)

	log.WithFields(log.Fields{
		"userEmail": userEmail,
		"languange": req.Language,
	}).Info(logInfo)

	user := store.User{
		UserEmail:  userEmail,
		SessionTTL: 0,
		IsAdmin:    false,
		IsActive:   false,
	}

	token := OTPgen()

	obj, err := a.st.GetUserByEmail(ctx, user.UserEmail)
	if err == storage.ErrDoesNotExist {
		if err := a.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
			// user has never been created yet
			err = handler.RegisterUser(ctx, &user, token)
			if err != nil {
				log.WithError(err).Error(logInfo)
				return status.Errorf(codes.Unknown, "%v", err)
			}

			// get user again
			obj, err = handler.GetUserByEmail(ctx, user.UserEmail)
			if err != nil {
				log.WithError(err).Error(logInfo)
				// internal error
				return status.Errorf(codes.Unknown, "%v", err)
			}

			err = email.SendInvite(obj.UserEmail, email.Param{Token: *obj.SecurityToken}, email.EmailLanguage(req.Language), email.RegistrationConfirmation)
			if err != nil {
				log.WithError(err).Error(logInfo)
				return helpers.ErrToRPCError(err)
			}

			return nil
		}); err != nil {
			return nil, status.Errorf(codes.Unknown, err.Error())
		}

	} else if err != nil && err != storage.ErrDoesNotExist {
		// internal error
		return nil, helpers.ErrToRPCError(err)
	} else if err == nil && obj.SecurityToken == nil {
		// user exists and finished registration
		return nil, helpers.ErrToRPCError(storage.ErrAlreadyExists)
	}

	return &empty.Empty{}, nil
}

// GetTOTPStatus returns info about TOTP status for the current user
func (a *InternalUserAPI) GetTOTPStatus(ctx context.Context, req *inpb.TOTPStatusRequest) (*inpb.TOTPStatusResponse, error) {
	if valid, err := NewValidator().ValidateActiveUser(ctx); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "not authenticated: %v", err)
	}
	u, err := NewValidator().GetUser(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	enabled, err := NewValidator().Is2FAEnabled(ctx, u.UserEmail)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &inpb.TOTPStatusResponse{
		Enabled: enabled,
	}, nil
}

// GetTOTPConfiguration generates a new TOTP configuration for the user
func (a *InternalUserAPI) GetTOTPConfiguration(ctx context.Context, req *inpb.GetTOTPConfigurationRequest) (*inpb.GetTOTPConfigurationResponse, error) {
	if valid, err := NewValidator().ValidateActiveUser(ctx); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "not authenticated: %v", err)
	}

	u, err := NewValidator().GetUser(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	cfg, err := NewValidator().NewConfiguration(ctx, u.UserEmail)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &inpb.GetTOTPConfigurationResponse{
		Url:          cfg.URL,
		Secret:       cfg.Secret,
		QrCode:       cfg.QRCode,
		RecoveryCode: cfg.RecoveryCodes,
	}, nil
}

// EnableTOTP enables TOTP for the user
func (a *InternalUserAPI) EnableTOTP(ctx context.Context, req *inpb.TOTPStatusRequest) (*inpb.TOTPStatusResponse, error) {
	if valid, err := NewValidator().ValidateActiveUser(ctx); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "not authenticated: %v", err)
	}

	err := NewValidator().Enable2FA(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%v", err)
	}

	return &inpb.TOTPStatusResponse{
		Enabled: true,
	}, nil
}

// DisableTOTP disables TOTP for the user
func (a *InternalUserAPI) DisableTOTP(ctx context.Context, req *inpb.TOTPStatusRequest) (*inpb.TOTPStatusResponse, error) {
	if valid, err := NewValidator().ValidateActiveUser(ctx); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "not authenticated: %v", err)
	}

	if err := NewValidator().Disable2FA(ctx); err != nil {
		return nil, status.Errorf(codes.Unknown, " %v", err)
	}

	return &inpb.TOTPStatusResponse{
		Enabled: false,
	}, nil
}

// GetRecoveryCodes returns the list of recovery codes for the user
func (a *InternalUserAPI) GetRecoveryCodes(ctx context.Context, req *inpb.GetRecoveryCodesRequest) (*inpb.GetRecoveryCodesResponse, error) {
	if valid, err := NewValidator().ValidateActiveUser(ctx); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "not authenticated: %v", err)
	}

	u, err := NewValidator().GetUser(ctx, authcus.WithValidOTP())
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	codes, err := NewValidator().OTPGetRecoveryCodes(ctx, u.UserEmail, req.Regenerate)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &inpb.GetRecoveryCodesResponse{
		RecoveryCode: codes,
	}, nil
}

func (a *InternalUserAPI) RequestPasswordReset(ctx context.Context, req *inpb.PasswordResetReq) (*inpb.PasswordResetResp, error) {
	if err := a.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		userEmail := normalizeUsername(req.Username)
		user, err := handler.GetUserByUsername(ctx, userEmail)
		if err != nil {
			if err == storage.ErrDoesNotExist {
				ctxlogrus.Extract(ctx).Warnf("password reset request for unknown user %s", userEmail)
				if err := email.SendInvite(userEmail, email.Param{}, email.EmailLanguage(req.Language), email.PasswordResetUnknown); err != nil {
					return status.Errorf(codes.Internal, "couldn't send recovery email: %v", err)
				}
				return nil
			}
			return status.Errorf(codes.Internal, "couldn't get user info: %v", err)
		}
		if !user.IsActive {
			ctxlogrus.Extract(ctx).Warnf("password reset request for inactive user %s", userEmail)
			return nil
		}
		pr, err := handler.GetPasswordResetRecord(ctx, user.ID)
		if err != nil {
			return status.Errorf(codes.Internal, "couldn't get password reset record: %v", err)
		}
		if pr.GeneratedAt.After(time.Now().Add(-30 * 24 * time.Hour)) {
			return status.Errorf(codes.PermissionDenied, "can't reset password more than once a month")
		}
		if err := pr.SetOTP(ctx, OTPgen()); err != nil {
			return status.Errorf(codes.Internal, "couldn't store reset code: %v", err)
		}

		if err := email.SendInvite(userEmail, email.Param{Token: pr.OTP}, email.EmailLanguage(req.Language), email.PasswordReset); err != nil {
			return status.Errorf(codes.Internal, "couldn't send recovery email: %v", err)
		}

		return nil
	}); err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}

	return &inpb.PasswordResetResp{}, nil
}

func (a *InternalUserAPI) ConfirmPasswordReset(ctx context.Context, req *inpb.ConfirmPasswordResetReq) (*inpb.PasswordResetResp, error) {
	if err := a.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		userEmail := normalizeUsername(req.Username)
		user, err := handler.GetUserByUsername(ctx, userEmail)
		if err != nil {
			if err == storage.ErrDoesNotExist {
				ctxlogrus.Extract(ctx).Warnf("password reset request for unknown user %s", userEmail)
				return status.Errorf(codes.PermissionDenied, "no match found")
			}
			return status.Errorf(codes.Internal, "couldn't get user info: %v", err)
		}
		pr, err := handler.GetPasswordResetRecord(ctx, user.ID)
		if err != nil {
			return status.Errorf(codes.Internal, "couldn't get password reset record: %v", err)
		}
		if pr.AttemptsLeft < 1 {
			return status.Errorf(codes.PermissionDenied, "no match found")
		}
		if err := pr.ReduceAttempts(ctx); err != nil {
			return status.Errorf(codes.Internal, "couldn't update db: %v", err)
		}
		if subtle.ConstantTimeCompare([]byte(pr.OTP), []byte(req.Otp)) == 1 {
			if err := pr.SetOTP(ctx, ""); err != nil {
				return status.Errorf(codes.Internal, "couldn't update db: %v", err)
			}
			if err := handler.UpdatePassword(ctx, pr.UserID, req.NewPassword, Service.pwh); err != nil {
				return status.Errorf(codes.Internal, "couldn't update db: %v", err)
			}
			return nil
		}

		return nil
	}); err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}

	return nil, status.Errorf(codes.PermissionDenied, "no match found")
}

// ConfirmRegistration checks provided security token and activates user
func (a *InternalUserAPI) ConfirmRegistration(ctx context.Context, req *inpb.ConfirmRegistrationRequest) (*inpb.ConfirmRegistrationResponse, error) {
	user, err := a.st.GetUserByToken(ctx, req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}

	log.Println("Confirming GetJwt", user.UserEmail)
	// give user a token that is valid only to finish the registration process
	jwt, err := NewValidator().SignJWToken(user.UserEmail, 86400, []string{"registration", "lora-app-server"})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &inpb.ConfirmRegistrationResponse{
		Id:       user.ID,
		Username: user.UserEmail,
		IsAdmin:  user.IsAdmin,
		IsActive: user.IsActive,
		Jwt:      jwt,
	}, status.Errorf(codes.OK, "")
}

// FinishRegistration sets new user password and creates a new organization
func (a *InternalUserAPI) FinishRegistration(ctx context.Context, req *inpb.FinishRegistrationRequest) (*empty.Empty, error) {
	u, err := NewValidator().GetUser(ctx,
		authcus.WithLimitedCredentials(), // nolint: staticcheck
		authcus.WithAudience("registration"),
	)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	if err := a.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		// Get the user id based on the userEmail and check that it matches the one
		// in the request and that user is not active
		user, err := handler.GetUserByUsername(ctx, u.UserEmail)
		if nil != err {
			return helpers.ErrToRPCError(err)
		}
		if user.ID != req.UserId {
			return status.Errorf(codes.PermissionDenied, "user id mismatch")
		}
		if user.IsActive {
			return status.Error(codes.PermissionDenied, "user has been registered already")
		}

		org := store.Organization{
			Name:            req.OrganizationName,
			DisplayName:     req.OrganizationDisplayName,
			CanHaveGateways: true,
		}

		err = handler.FinishRegistration(ctx, req.UserId, req.Password, Service.pwh)
		if err != nil {
			return status.Errorf(codes.Unknown, "%v", err)
		}

		err = handler.CreateOrganization(ctx, &org)
		if err != nil {
			return status.Errorf(codes.Unknown, "%v", err)
		}

		err = handler.CreateOrganizationUser(ctx, org.ID, user.UserEmail, true, false, false)
		if err != nil {
			return status.Errorf(codes.Unknown, "%v", err)
		}

		return nil
	}); err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}

	return &empty.Empty{}, status.Errorf(codes.OK, "")
}
