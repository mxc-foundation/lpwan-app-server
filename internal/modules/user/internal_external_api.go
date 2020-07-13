package user

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	inpb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/email"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/organization"
)

// InternalUserAPI exports the internal User related functions.
type InternalUserAPI struct {
	Validator *Validator
	Store     UserStore
}

// NewInternalUserAPI creates a new InternalUserAPI.
func NewInternalUserAPI(api InternalUserAPI) *InternalUserAPI {
	return &InternalUserAPI{
		Validator: api.Validator,
		Store:     api.Store,
	}
}

// Login validates the login request and returns a JWT token.
func (a *InternalUserAPI) Login(ctx context.Context, req *inpb.LoginRequest) (*inpb.LoginResponse, error) {
	err := a.Store.LoginUserByPassword(ctx, req.Username, req.Password)
	if nil != err {
		return nil, helpers.ErrToRPCError(err)
	}

	user, err := a.Store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't get info about the user")
	}

	if !user.IsActive {
		return nil, status.Error(codes.Unauthenticated, "incactive user")
	}

	ttl := 60 * int64(user.SessionTTL)
	var audience []string

	is2fa, err := a.Validator.Credentials.Is2FAEnabled(ctx, user.Username)
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

	jwt, err := a.Validator.Credentials.SignJWToken(user.Username, ttl, audience)
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
	u, err := a.Validator.Credentials.GetUser(ctx, authcus.WithAudience("login-2fa"), authcus.WithValidOTP())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "not authenticated")
	}

	user, err := a.Store.GetUserByUsername(ctx, u.Username)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't get info about the user")
	}

	jwt, err := a.Validator.Credentials.SignJWToken(u.Username, 60*int64(user.SessionTTL), nil)
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
	if valid, err := a.Validator.ValidateActiveUser(ctx); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	user, err := a.Validator.Credentials.GetUser(ctx)
	if nil != err {
		return nil, helpers.ErrToRPCError(err)
	}

	prof, err := a.Store.GetProfile(ctx, user.ID)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := inpb.ProfileResponse{
		User: &inpb.User{
			Id:         prof.User.ID,
			Username:   prof.User.Username,
			Email:      prof.User.Email,
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
		Registration: config.C.ApplicationServer.Branding.Registration,
		Footer:       config.C.ApplicationServer.Branding.Footer,
		LogoPath:     os.Getenv("APPSERVER") + "/branding.png",
	}

	return &resp, nil
}

// GlobalSearch performs a global search.
func (a *InternalUserAPI) GlobalSearch(ctx context.Context, req *inpb.GlobalSearchRequest) (*inpb.GlobalSearchResponse, error) {
	if valid, err := a.Validator.ValidateActiveUser(ctx); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	user, err := a.Validator.Credentials.GetUser(ctx)
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

// RegisterUser adds new user and sends activation email
func (a *InternalUserAPI) RegisterUser(ctx context.Context, req *inpb.RegisterUserRequest) (*empty.Empty, error) {
	logInfo := "api/appserver_serves_ui/RegisterUser"

	log.WithFields(log.Fields{
		"email":     req.Email,
		"languange": req.Language,
	}).Info(logInfo)

	user := User{
		Username:   req.Email,
		Email:      req.Email,
		SessionTTL: 0,
		IsAdmin:    false,
		IsActive:   false,
	}

	token := OTPgen()

	obj, err := a.Store.GetUserByEmail(ctx, user.Email)
	if err == storage.ErrDoesNotExist {

		err := storage.Transaction(func(tx sqlx.Ext) error {
			// user has never been created yet
			err = a.Store.RegisterUser(&user, token)
			if err != nil {
				log.WithError(err).Error(logInfo)
				return helpers.ErrToRPCError(err)
			}

			// get user again
			obj, err = a.Store.GetUserByEmail(ctx, user.Email)
			if err != nil {
				log.WithError(err).Error(logInfo)
				// internal error
				return helpers.ErrToRPCError(err)
			}
			if err != nil {
				return helpers.ErrToRPCError(err)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}

	} else if err != nil && err != storage.ErrDoesNotExist {
		// internal error
		return nil, helpers.ErrToRPCError(err)
	} else if err == nil && obj.SecurityToken == nil {
		// user exists and finished registration
		return nil, helpers.ErrToRPCError(storage.ErrAlreadyExists)
	}

	err = email.SendInvite(obj.Email, *obj.SecurityToken, email.EmailLanguage(req.Language), email.RegistrationConfirmation)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// GetTOTPStatus returns info about TOTP status for the current user
func (a *InternalUserAPI) GetTOTPStatus(ctx context.Context, req *inpb.TOTPStatusRequest) (*inpb.TOTPStatusResponse, error) {
	if valid, err := a.Validator.ValidateActiveUser(ctx); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "not authenticated: %v", err)
	}
	u, err := a.Validator.Credentials.GetUser(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	enabled, err := a.Validator.Credentials.Is2FAEnabled(ctx, u.Username)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &inpb.TOTPStatusResponse{
		Enabled: enabled,
	}, nil
}

// GetTOTPConfiguration generates a new TOTP configuration for the user
func (a *InternalUserAPI) GetTOTPConfiguration(ctx context.Context, req *inpb.GetTOTPConfigurationRequest) (*inpb.GetTOTPConfigurationResponse, error) {
	if valid, err := a.Validator.ValidateActiveUser(ctx); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "not authenticated: %v", err)
	}

	u, err := a.Validator.Credentials.GetUser(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	cfg, err := a.Validator.Credentials.NewConfiguration(ctx, u.Username)
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
	if valid, err := a.Validator.ValidateActiveUser(ctx); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "not authenticated: %v", err)
	}

	u, err := a.Validator.Credentials.GetUser(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	otp, err := a.Validator.Credentials.GetOTP(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	if err := a.Validator.Credentials.EnableOTP(ctx, u.Username, otp); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "not authenticated: %v", err)
	}

	return &inpb.TOTPStatusResponse{
		Enabled: true,
	}, nil
}

// DisableTOTP disables TOTP for the user
func (a *InternalUserAPI) DisableTOTP(ctx context.Context, req *inpb.TOTPStatusRequest) (*inpb.TOTPStatusResponse, error) {
	if valid, err := a.Validator.ValidateActiveUser(ctx); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "not authenticated: %v", err)
	}

	u, err := a.Validator.Credentials.GetUser(ctx, authcus.WithValidOTP())
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	if err := a.Validator.Credentials.DisableOTP(ctx, u.Username); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	return &inpb.TOTPStatusResponse{
		Enabled: false,
	}, nil
}

// GetRecoveryCodes returns the list of recovery codes for the user
func (a *InternalUserAPI) GetRecoveryCodes(ctx context.Context, req *inpb.GetRecoveryCodesRequest) (*inpb.GetRecoveryCodesResponse, error) {
	if valid, err := a.Validator.ValidateActiveUser(ctx); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "not authenticated: %v", err)
	}

	u, err := a.Validator.Credentials.GetUser(ctx, authcus.WithValidOTP())
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	codes, err := a.Validator.Credentials.OTPGetRecoveryCodes(ctx, u.Username, req.Regenerate)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &inpb.GetRecoveryCodesResponse{
		RecoveryCode: codes,
	}, nil
}

func (a *InternalUserAPI) RequestPasswordReset(ctx context.Context, req *inpb.PasswordResetReq) (*inpb.PasswordResetResp, error) {

	tx, err := storage.DB().BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't begin tx: %v", err)
	}
	defer tx.Rollback()
	user, err := storage.GetUserByUsername(ctx, tx, req.Username)
	if err != nil {
		if err == storage.ErrDoesNotExist {
			ctxlogrus.Extract(ctx).Warnf("password reset request for unknown user %s", req.Username)
			if err := email.SendInvite(req.Username, "", email.EmailLanguage(req.Language), email.PasswordResetUnknown); err != nil {
				return nil, status.Errorf(codes.Internal, "couldn't send recovery email: %v", err)
			}
			return &inpb.PasswordResetResp{}, nil
		}
		return nil, status.Errorf(codes.Internal, "couldn't get user info: %v", err)
	}
	if !user.IsActive {
		ctxlogrus.Extract(ctx).Warnf("password reset request for inactive user %s", req.Username)
		return &inpb.PasswordResetResp{}, nil
	}
	pr, err := storage.GetPasswordResetRecord(tx, user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't get password reset record: %v", err)
	}
	if pr.GeneratedAt.After(time.Now().Add(-30 * 24 * time.Hour)) {
		return nil, status.Errorf(codes.PermissionDenied, "can't reset password more than once a month")
	}
	if err := pr.SetOTP(OTPgen()); err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't store reset code: %v", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't store reset code: %v", err)
	}
	if err := email.SendInvite(req.Username, pr.OTP, email.EmailLanguage(req.Language), email.PasswordReset); err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't send recovery email: %v", err)
	}
	return &inpb.PasswordResetResp{}, nil
}

func (a *InternalUserAPI) ConfirmPasswordReset(ctx context.Context, req *inpb.ConfirmPasswordResetReq) (*inpb.PasswordResetResp, error) {
	tx, err := storage.DB().BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't begin tx: %v", err)
	}
	defer tx.Rollback()
	user, err := storage.GetUserByUsername(ctx, tx, req.Username)
	if err != nil {
		if err == storage.ErrDoesNotExist {
			ctxlogrus.Extract(ctx).Warnf("password reset request for unknown user %s", req.Username)
			return nil, status.Errorf(codes.PermissionDenied, "no match found")
		}
		return nil, status.Errorf(codes.Internal, "couldn't get user info: %v", err)
	}
	pr, err := storage.GetPasswordResetRecord(tx, user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't get password reset record: %v", err)
	}
	if pr.AttemptsLeft < 1 {
		return nil, status.Errorf(codes.PermissionDenied, "no match found")
	}
	if err := pr.ReduceAttempts(); err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't update db: %v", err)
	}
	if subtle.ConstantTimeCompare([]byte(pr.OTP), []byte(req.Otp)) == 1 {
		if err := pr.SetOTP(""); err != nil {
			return nil, status.Errorf(codes.Internal, "couldn't update db: %v", err)
		}
		if err := storage.UpdatePassword(ctx, tx, pr.UserID, req.NewPassword); err != nil {
			return nil, status.Errorf(codes.Internal, "couldn't update db: %v", err)
		}
		if err := tx.Commit(); err != nil {
			return nil, status.Errorf(codes.Internal, "couldn't update db: %v", err)
		}
		return &inpb.PasswordResetResp{}, nil
	}
	if err := tx.Commit(); err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't update db: %v", err)
	}
	return nil, status.Errorf(codes.PermissionDenied, "no match found")
}

// ConfirmRegistration checks provided security token and activates user
func (a *InternalUserAPI) ConfirmRegistration(ctx context.Context, req *inpb.ConfirmRegistrationRequest) (*inpb.ConfirmRegistrationResponse, error) {
	user, err := storage.GetUserByToken(storage.DB(), req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}

	log.Println("Confirming GetJwt", user.Username)
	// give user a token that is valid only to finish the registration process
	jwt, err := a.validator.SignToken(user.Username, 86400, []string{"registration", "lora-app-server"})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &inpb.ConfirmRegistrationResponse{
		Id:       user.ID,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
		IsActive: user.IsActive,
		Jwt:      jwt,
	}, status.Errorf(codes.OK, "")
}

// FinishRegistration sets new user password and creates a new organization
func (a *InternalUserAPI) FinishRegistration(ctx context.Context, req *inpb.FinishRegistrationRequest) (*empty.Empty, error) {
	cred, err := a.validator.GetCredentials(ctx,
		auth.WithLimitedCredentials(), // nolint: staticcheck
		auth.WithAudience("registration"),
	)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	// Get the user id based on the username and check that it matches the one
	// in the request and that user is not active
	user, err := storage.GetUserByUsername(ctx, storage.DB(), cred.Username())
	if nil != err {
		return nil, helpers.ErrToRPCError(err)
	}
	if user.ID != req.UserId {
		return nil, status.Errorf(codes.PermissionDenied, "user id mismatch")
	}
	if user.IsActive {
		return nil, status.Error(codes.PermissionDenied, "user has been registered already")
	}

	org := organization.Organization{
		Name:            req.OrganizationName,
		DisplayName:     req.OrganizationDisplayName,
		CanHaveGateways: true,
	}

	err := storage.Transaction(func(tx sqlx.Ext) error {
		err := a.Store.FinishRegistration(req.UserId, req.Password)
		if err != nil {
			return helpers.ErrToRPCError(err)
		}

		err = organization.GetOrganizationAPI().Store.CreateOrganization(ctx, &org)
		if err != nil {
			return helpers.ErrToRPCError(err)
		}

		u, err := a.Store.GetUser(ctx, req.UserId)
		if err != nil {
			return helpers.ErrToRPCError(err)
		}

		err = organization.GetOrganizationAPI().Store.CreateOrganizationUser(ctx, org.ID, u.Username, true, false, false)
		if err != nil {
			return helpers.ErrToRPCError(err)
		}

		return nil
	})

	if err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}

	return &empty.Empty{}, status.Errorf(codes.OK, "")
}
