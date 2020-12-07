package user

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/mxc-foundation/lpwan-app-server/internal/jwt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	inpb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	"github.com/mxc-foundation/lpwan-app-server/internal/auth"
	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
	"github.com/mxc-foundation/lpwan-app-server/internal/grpcauth"
)

// Login validates the login request and returns a JWT token.
func (a *Server) Login(ctx context.Context, req *inpb.LoginRequest) (*inpb.LoginResponse, error) {
	userEmail := normalizeUsername(req.Username)

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

	var audience []string
	is2fa, err := a.otpv.IsEnabled(ctx, u.Email)
	if err != nil {
		ctxlogrus.Extract(ctx).WithError(err).Error("couldn't get 2fa status")
		return nil, status.Error(codes.Internal, "couldn't get 2fa status")
	}
	if !a.config.Enable2FALogin {
		is2fa = false
	}

	var ttl int64
	if is2fa {
		// if 2fa is enabled we issue token that is only valid for 10 minutes
		// and is only good to perform second factor authentication. If second
		// factor authentication has been successful then it will return to the
		// user another token, that provides access to all the api
		ttl = 600
		audience = []string{"login-2fa"}
	}

	jwToken, err := a.jwtv.SignToken(jwt.Claims{UserID: u.ID, Username: u.Email, Service: auth.EMAIL}, ttl, audience)
	if err != nil {
		log.Errorf("SignToken returned an error: %v", err)
		return nil, status.Errorf(codes.Internal, "couldn't create a token")
	}

	_ = a.store.SetUserLastLogin(ctx, u.ID, u.Email, auth.EMAIL)

	return &inpb.LoginResponse{Jwt: jwToken, Is_2FaRequired: is2fa}, nil
}

// Login2FA performs second factor authentication. It requires u to have
// already passed password check and checks if the OTP code is valid. If it is
// it returns JWT with access to the api.
func (a *Server) Login2FA(ctx context.Context, req *inpb.Login2FARequest) (*inpb.LoginResponse, error) {
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithAudience("login-2fa").WithRequireOTP())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	jwToken, err := a.jwtv.SignToken(jwt.Claims{UserID: cred.UserID, Username: cred.Username}, 0, nil)

	if err != nil {
		log.Errorf("SignToken returned an error: %v", err)
		return nil, status.Error(codes.Internal, "couldn't create a token")
	}

	return &inpb.LoginResponse{Jwt: jwToken}, nil
}

type RecaptchaConfig struct {
	HostServer string `mapstructure:"host_server"`
	Secret     string `mapstructure:"secret"`
}

// IsPassVerifyingGoogleRecaptcha defines the response to pass the google recaptcha verification
func (a *Server) IsPassVerifyingGoogleRecaptcha(response string, remoteip string) (*inpb.GoogleRecaptchaResponse, error) {
	secret := a.config.Recaptcha.Secret
	postURL := a.config.Recaptcha.HostServer

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
func (a *Server) GetVerifyingGoogleRecaptcha(ctx context.Context, req *inpb.GoogleRecaptchaRequest) (*inpb.GoogleRecaptchaResponse, error) {
	res, err := a.IsPassVerifyingGoogleRecaptcha(req.Response, req.Remoteip)
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

func validatePass(password string) error {
	if len(password) < 8 {
		return status.Errorf(codes.InvalidArgument, "password must be at least 8 characters long")
	}
	return nil
}

// based on https://www.w3.org/TR/2016/REC-html51-20161101/sec-forms.html#email-state-typeemail
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func validateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return status.Errorf(codes.InvalidArgument, "invalid email address")
	}
	return nil
}

// Profile returns the u profile.
func (a *Server) Profile(ctx context.Context, req *empty.Empty) (*inpb.ProfileResponse, error) {
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	user, err := a.store.GetUserByID(ctx, cred.UserID)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := inpb.ProfileResponse{
		User: &inpb.User{
			Id:               user.ID,
			Username:         user.DisplayName,
			Email:            user.Email,
			IsAdmin:          user.IsAdmin,
			IsActive:         user.IsActive,
			LastLoginService: cred.Service,
		},
	}

	orgs, err := a.store.GetUserOrganizations(ctx, cred.UserID)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	for _, org := range orgs {
		row := inpb.OrganizationLink{
			OrganizationId:   org.OrganizationID,
			OrganizationName: org.OrganizationName,
			IsAdmin:          org.IsOrgAdmin,
			IsDeviceAdmin:    org.IsDeviceAdmin,
			IsGatewayAdmin:   org.IsGatewayAdmin,
			CreatedAt:        &timestamp.Timestamp{Seconds: org.CreatedAt.Unix()},
			UpdatedAt:        &timestamp.Timestamp{Seconds: org.UpdatedAt.Unix()},
		}

		resp.Organizations = append(resp.Organizations, &row)
	}

	externalUsers, err := a.store.GetExternalUsersByUserID(ctx, user.ID)
	if err != nil {
		if err != errHandler.ErrDoesNotExist {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
	}

	for _, eu := range externalUsers {
		item := inpb.ExternalUserAccount{
			ExternalUserId:   eu.ExternalUserID,
			ExternalUsername: eu.ExternalUsername,
			Service:          eu.ServiceName,
		}

		resp.ExternalUserAccounts = append(resp.ExternalUserAccounts, &item)
	}

	return &resp, nil
}

// Branding returns UI branding.
func (a *Server) Branding(ctx context.Context, req *empty.Empty) (*inpb.BrandingResponse, error) {
	resp := inpb.BrandingResponse{
		LogoPath: a.config.OperatorLogoPath,
	}

	return &resp, nil
}

// GlobalSearch performs a global search.
func (a *Server) GlobalSearch(ctx context.Context, req *inpb.GlobalSearchRequest) (*inpb.GlobalSearchResponse, error) {
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	results, err := a.store.GlobalSearch(ctx, cred.UserID, cred.IsGlobalAdmin, req.Search, int(req.Limit), int(req.Offset))
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

// RegisterUser adds new u and sends activation userEmail
func (a *Server) RegisterUser(ctx context.Context, req *inpb.RegisterUserRequest) (*empty.Empty, error) { // nolint: gocyclo
	logInfo := "api/appserver_serves_ui/RegisterUser"

	userEmail := normalizeUsername(req.Email)

	log.WithFields(log.Fields{
		"userEmail": userEmail,
		"languange": req.Language,
	}).Info(logInfo)
	if err := validateEmail(userEmail); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	user, err := a.store.GetUserByEmail(ctx, userEmail)
	// internal error
	if err != nil && err != errHandler.ErrDoesNotExist {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err == nil {
		// user exists but haven't finished registration
		if !user.IsActive && user.SecurityToken != "" {
			err := a.mailer.SendRegistrationConfirmation(user.Email, req.Language, user.SecurityToken)
			if err != nil {
				log.WithError(err).Error(logInfo)
				return nil, helpers.ErrToRPCError(err)
			}

			return &empty.Empty{}, nil
		}
		// user exists and finished registration
		return nil, status.Errorf(codes.AlreadyExists, "")
	}

	// user doesn't exist
	token := OTPgen()
	u := User{
		Email:         userEmail,
		SecurityToken: token,
	}

	if _, err := a.store.CreateUser(ctx, u, nil); err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't register user: %v", err)
	}

	if err := a.mailer.SendRegistrationConfirmation(userEmail, req.Language, token); err != nil {
		log.WithError(err).Error(logInfo)
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// GetTOTPStatus returns info about TOTP status for the current u
func (a *Server) GetTOTPStatus(ctx context.Context, req *inpb.TOTPStatusRequest) (*inpb.TOTPStatusResponse, error) {
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	enabled, err := a.otpv.IsEnabled(ctx, cred.Username)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &inpb.TOTPStatusResponse{
		Enabled: enabled,
	}, nil
}

// GetTOTPConfiguration generates a new TOTP configuration for the u
func (a *Server) GetTOTPConfiguration(ctx context.Context, req *inpb.GetTOTPConfigurationRequest) (*inpb.GetTOTPConfigurationResponse, error) {
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	cfg, err := a.otpv.NewConfiguration(ctx, cred.Username)
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

// EnableTOTP enables TOTP for the u
func (a *Server) EnableTOTP(ctx context.Context, req *inpb.TOTPStatusRequest) (*inpb.TOTPStatusResponse, error) {
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	otp := grpcauth.GetOTPFromContext(ctx)
	if err := a.otpv.Enable(ctx, cred.Username, otp); err != nil {
		return nil, status.Errorf(codes.Unknown, "%v", err)
	}

	return &inpb.TOTPStatusResponse{
		Enabled: true,
	}, nil
}

// DisableTOTP disables TOTP for the u
func (a *Server) DisableTOTP(ctx context.Context, req *inpb.TOTPStatusRequest) (*inpb.TOTPStatusResponse, error) {
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithRequireOTP())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	if err := a.otpv.Disable(ctx, cred.Username); err != nil {
		return nil, status.Errorf(codes.Unknown, " %v", err)
	}

	return &inpb.TOTPStatusResponse{
		Enabled: false,
	}, nil
}

// GetRecoveryCodes returns the list of recovery codes for the u
func (a *Server) GetRecoveryCodes(ctx context.Context, req *inpb.GetRecoveryCodesRequest) (*inpb.GetRecoveryCodesResponse, error) {
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithRequireOTP())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	codes, err := a.otpv.GetRecoveryCodes(ctx, cred.Username, req.Regenerate)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &inpb.GetRecoveryCodesResponse{
		RecoveryCode: codes,
	}, nil
}

func (a *Server) RequestPasswordReset(ctx context.Context, req *inpb.PasswordResetReq) (*inpb.PasswordResetResp, error) {
	userEmail := normalizeUsername(req.Username)
	user, err := a.store.GetUserByEmail(ctx, userEmail)
	if err != nil {
		if err == errHandler.ErrDoesNotExist {
			ctxlogrus.Extract(ctx).Warnf("password reset request for unknown user %s", userEmail)
			if err := a.mailer.SendPasswordResetUnknown(userEmail, req.Language); err != nil {
				return nil, status.Errorf(codes.Internal, "couldn't send recovery email: %v", err)
			}
			return &inpb.PasswordResetResp{}, nil
		}
		return nil, status.Errorf(codes.Internal, "couldn't get user info: %v", err)
	}

	if !user.IsActive {
		ctxlogrus.Extract(ctx).Warnf("password reset request for inactive user %s", userEmail)
		return nil, status.Errorf(codes.PermissionDenied, "permission denied: inactive")
	}

	otp, err := a.store.GetOrSetPasswordResetOTP(ctx, user.ID, OTPgen())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	if err := a.mailer.SendPasswordReset(userEmail, req.Language, otp); err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't send recovery email: %v", err)
	}

	return &inpb.PasswordResetResp{}, nil
}

func (a *Server) ConfirmPasswordReset(ctx context.Context, req *inpb.ConfirmPasswordResetReq) (*inpb.PasswordResetResp, error) {
	userEmail := normalizeUsername(req.Username)
	user, err := a.store.GetUserByEmail(ctx, userEmail)
	if err != nil {
		if err == errHandler.ErrDoesNotExist {
			ctxlogrus.Extract(ctx).Warnf("password reset request for unknown user %s", userEmail)
			return nil, status.Errorf(codes.PermissionDenied, "no match found")
		}
		return nil, status.Errorf(codes.Internal, "couldn't get user info: %v", err)
	}

	if err := validatePass(req.NewPassword); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	ph, err := a.pwhasher.HashPassword(req.NewPassword)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't hash password: %v", err)
	}
	if err := a.store.SetUserPasswordIfOTPMatch(ctx, user.ID, req.Otp, ph); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &inpb.PasswordResetResp{}, nil
}

// ConfirmRegistration checks provided security token and activates u
func (a *Server) ConfirmRegistration(ctx context.Context, req *inpb.ConfirmRegistrationRequest) (*inpb.ConfirmRegistrationResponse, error) {
	u, err := a.store.GetUserByToken(ctx, req.Token)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// issue a token that is valid only to finish the registration process
	jwToken, err := a.jwtv.SignToken(jwt.Claims{UserID: u.ID, Username: u.Email, Service: auth.EMAIL}, 86400, []string{"registration", "lora-app-server"})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &inpb.ConfirmRegistrationResponse{
		Id:       u.ID,
		Username: u.Email,
		IsAdmin:  u.IsAdmin,
		IsActive: u.IsActive,
		Jwt:      jwToken,
	}, nil
}

// FinishRegistration sets new u password and creates a new organization
func (a *Server) FinishRegistration(ctx context.Context, req *inpb.FinishRegistrationRequest) (*empty.Empty, error) {
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithAllowNonExisting().WithAudience("registration"))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	user, err := a.store.GetUserByEmail(ctx, cred.Username)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	if req.UserId != user.ID {
		return nil, status.Errorf(codes.PermissionDenied, "user id mismatch")
	}
	if user.IsActive {
		return nil, status.Errorf(codes.AlreadyExists, "user has been activated already")
	}
	if err := validatePass(req.Password); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	ph, err := a.pwhasher.HashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't hash the password: %v", err)
	}
	if err := a.store.ActivateUser(ctx, user.ID, ph, req.OrganizationName, req.OrganizationDisplayName); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}
