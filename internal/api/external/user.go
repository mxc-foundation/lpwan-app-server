package external

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	authcus "github.com/mxc-foundation/lpwan-app-server/internal/api/external/auth_customized"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/brocaar/chirpstack-api/go/v3/as/external/api"

	inpb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/auth"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/email"
	"github.com/mxc-foundation/lpwan-app-server/internal/otp"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
)

// UserAPI exports the User related functions.
type UserAPI struct {
	validator auth.Validator
}

// InternalUserAPI exports the internal User related functions.
type InternalUserAPI struct {
	validator    auth.Validator
	validatorCus authcus.Validator
	otpValidator *otp.Validator
}

// NewUserAPI creates a new UserAPI.
func NewUserAPI(validator auth.Validator) *UserAPI {
	return &UserAPI{
		validator: validator,
	}
}

// Create creates the given user.
func (a *UserAPI) Create(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if req.User == nil {
		return nil, status.Errorf(codes.InvalidArgument, "user must not be nil")
	}

	if err := a.validator.Validate(ctx,
		auth.ValidateUsersAccess(auth.Create)); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	user := storage.User{
		SessionTTL: req.User.SessionTtl,
		IsAdmin:    req.User.IsAdmin,
		IsActive:   req.User.IsActive,
		Email:      req.User.Email,
		Note:       req.User.Note,
	}

	if err := user.SetPasswordHash(req.Password); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	err := storage.Transaction(func(tx sqlx.Ext) error {
		err := storage.CreateUser(ctx, tx, &user)
		if err != nil {
			return err
		}

		for _, org := range req.Organizations {
			if err := storage.CreateOrganizationUser(ctx, tx, org.OrganizationId, user.ID, org.IsAdmin, org.IsDeviceAdmin, org.IsGatewayAdmin); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &pb.CreateUserResponse{Id: user.ID}, nil
}

// Get returns the user matching the given ID.
func (a *UserAPI) Get(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if err := a.validator.Validate(ctx,
		auth.ValidateUserAccess(req.Id, auth.Read)); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	user, err := storage.GetUser(ctx, storage.DB(), req.Id)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := pb.GetUserResponse{
		User: &pb.User{
			Id:         user.ID,
			SessionTtl: user.SessionTTL,
			IsAdmin:    user.IsAdmin,
			IsActive:   user.IsActive,
			Email:      user.Email,
			Note:       user.Note,
		},
	}

	resp.CreatedAt, err = ptypes.TimestampProto(user.CreatedAt)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	resp.UpdatedAt, err = ptypes.TimestampProto(user.UpdatedAt)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &resp, nil
}

// GetUserEmail returns true if user does not exist
func (a *UserAPI) GetUserEmail(ctx context.Context, req *inpb.GetUserEmailRequest) (*inpb.GetUserEmailResponse, error) {
	u, err := storage.GetUserByEmail(ctx, storage.DB(), req.UserEmail)
	if err != nil {
		if err == storage.ErrDoesNotExist {
			return &inpb.GetUserEmailResponse{Status: true}, nil
		}
		return nil, helpers.ErrToRPCError(err)
	}
	if u.SecurityToken != nil {
		// user exists but has not finished registration
		return &inpb.GetUserEmailResponse{Status: true}, nil
	}

	return &inpb.GetUserEmailResponse{Status: false}, nil
}

// List lists the users.
func (a *UserAPI) List(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserResponse, error) {
	if err := a.validator.Validate(ctx,
		auth.ValidateUsersAccess(auth.List)); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	users, err := storage.GetUsers(ctx, storage.DB(), int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	totalUserCount, err := storage.GetUserCount(ctx, storage.DB())
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := pb.ListUserResponse{
		TotalCount: int64(totalUserCount),
	}

	for _, u := range users {
		row := pb.UserListItem{
			Id:         u.ID,
			SessionTtl: u.SessionTTL,
			IsAdmin:    u.IsAdmin,
			IsActive:   u.IsActive,
		}

		row.CreatedAt, err = ptypes.TimestampProto(u.CreatedAt)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
		row.UpdatedAt, err = ptypes.TimestampProto(u.UpdatedAt)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}

		resp.Result = append(resp.Result, &row)
	}

	return &resp, nil
}

// Update updates the given user.
func (a *UserAPI) Update(ctx context.Context, req *pb.UpdateUserRequest) (*empty.Empty, error) {
	if req.User == nil {
		return nil, status.Errorf(codes.InvalidArgument, "user must not be nil")
	}

	if err := a.validator.Validate(ctx,
		auth.ValidateUserAccess(req.User.Id, auth.Update)); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	user, err := storage.GetUser(ctx, storage.DB(), req.User.Id)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	user.IsAdmin = req.User.IsAdmin
	user.IsActive = req.User.IsActive
	user.SessionTTL = req.User.SessionTtl
	user.Email = req.User.Email
	user.Note = req.User.Note

	if err := storage.UpdateUser(ctx, storage.DB(), &user); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// Delete deletes the user matching the given ID.
func (a *UserAPI) Delete(ctx context.Context, req *pb.DeleteUserRequest) (*empty.Empty, error) {
	if err := a.validator.Validate(ctx,
		auth.ValidateUserAccess(req.Id, auth.Delete)); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	err := storage.DeleteUser(ctx, storage.DB(), req.Id)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// UpdatePassword updates the password for the user matching the given ID.
func (a *UserAPI) UpdatePassword(ctx context.Context, req *pb.UpdateUserPasswordRequest) (*empty.Empty, error) {
	if err := a.validator.Validate(ctx,
		auth.ValidateUserAccess(req.UserId, auth.UpdateProfile)); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	user, err := storage.GetUser(ctx, storage.DB(), req.UserId)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	if err := user.SetPasswordHash(req.Password); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// NewInternalUserAPI creates a new InternalUserAPI.
func NewInternalUserAPI(validator auth.Validator, otpValidator *otp.Validator) *InternalUserAPI {
	return &InternalUserAPI{
		validator:    validator,
		otpValidator: otpValidator,
	}
}

// Login validates the login request and returns a JWT token.
func (a *InternalUserAPI) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	jwt, err := storage.LoginUserByPassword(ctx, storage.DB(), req.Email, req.Password)
	if nil != err {
		return nil, helpers.ErrToRPCError(err)
	}

	return &pb.LoginResponse{Jwt: jwt}, nil
}

// IsPassVerifyingGoogleRecaptcha defines the response to pass the google recaptcha verification
func IsPassVerifyingGoogleRecaptcha(response string, remoteip string) (*inpb.GoogleRecaptchaResponse, error) {
	secret := config.C.Recaptcha.Secret
	postURL := config.C.Recaptcha.HostServer

	postStr := url.Values{"secret": {secret}, "response": {response}, "remoteip": {remoteip}}
	responsePost, err := http.PostForm(postURL, postStr)

	if err != nil {
		log.Warn(err)
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
		log.Warn(err)
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

type claims struct {
	Username string `json:"username"`
}

// Profile returns the user profile.
func (a *InternalUserAPI) Profile(ctx context.Context, req *empty.Empty) (*inpb.ProfileResponse, error) {
	if err := a.validator.Validate(ctx,
		auth.ValidateActiveUser()); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	user, err := a.validator.GetUser(ctx)
	if nil != err {
		return nil, helpers.ErrToRPCError(err)
	}

	prof, err := storage.GetProfile(ctx, storage.DB(), user.ID)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := inpb.ProfileResponse{
		User: &inpb.User{
			Id:         prof.User.ID,
			Username:   prof.User.Email,
			SessionTtl: prof.User.SessionTTL,
			IsAdmin:    prof.User.IsAdmin,
			IsActive:   prof.User.IsActive,
		},
		Settings: &inpb.ProfileSettings{
			DisableAssignExistingUsers: authcus.DisableAssignExistingUsers,
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
		Registration: brandingRegistration,
		Footer:       brandingFooter,
		LogoPath:     os.Getenv("APPSERVER") + "/branding.png",
	}

	return &resp, nil
}

// GlobalSearch performs a global search.
func (a *InternalUserAPI) GlobalSearch(ctx context.Context, req *pb.GlobalSearchRequest) (*pb.GlobalSearchResponse, error) {
	if err := a.validator.Validate(ctx,
		auth.ValidateActiveUser()); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	isAdmin, err := a.validatorCus.GetIsAdmin(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	user, err := a.validator.GetUser(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	results, err := storage.GlobalSearch(ctx, storage.DB(), user.ID, isAdmin, req.Search, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	var out pb.GlobalSearchResponse

	for _, r := range results {
		res := pb.GlobalSearchResult{
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

// RegisterUser adds new user and sends activation email
func (a *InternalUserAPI) RegisterUser(ctx context.Context, req *inpb.RegisterUserRequest) (*empty.Empty, error) {
	logInfo := "api/appserver_serves_ui/RegisterUser"

	log.WithFields(log.Fields{
		"email":     req.Email,
		"languange": inpb.Language_name[int32(req.Language)],
	}).Info(logInfo)

	user := storage.User{
		Email:      req.Email,
		SessionTTL: 0,
		IsAdmin:    false,
		IsActive:   false,
	}

	u := OTPgen()
	// if err != nil {
	// 	log.WithError(err).Error(logInfo)
	// 	return nil, helpers.ErrToRPCError(err)
	// }
	token := u

	obj, err := storage.GetUserByEmail(ctx, storage.DB(), user.Email)
	if err == storage.ErrDoesNotExist {
		// user has never been created yet
		err = storage.RegisterUser(storage.DB(), &user, token)
		if err != nil {
			log.WithError(err).Error(logInfo)
			return nil, helpers.ErrToRPCError(err)
		}

		// get user again
		obj, err = storage.GetUserByEmail(ctx, storage.DB(), user.Email)
		if err != nil {
			log.WithError(err).Error(logInfo)
			// internal error
			return nil, helpers.ErrToRPCError(err)
		}

	} else if err != nil && err != storage.ErrDoesNotExist {
		// internal error
		return nil, helpers.ErrToRPCError(err)
	} else if err == nil && obj.SecurityToken == nil {
		// user exists and finished registration
		return nil, helpers.ErrToRPCError(storage.ErrAlreadyExists)
	}

	err = email.SendInvite(obj.Email, *obj.SecurityToken, email.EmailLanguage(inpb.Language_name[int32(req.Language)]), email.RegistrationConfirmation)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

func (a *UserAPI) GetOTPCode(ctx context.Context, req *inpb.GetOTPCodeRequest) (*inpb.GetOTPCodeResponse, error) {
	otp, err := storage.GetTokenByUsername(ctx, storage.DB(), req.UserEmail)
	if err != nil {
		return nil, err
	}

	return &inpb.GetOTPCodeResponse{OtpCode: otp}, nil
}

// GetTOTPStatus returns info about TOTP status for the current user
func (a *InternalUserAPI) GetTOTPStatus(ctx context.Context, req *inpb.TOTPStatusRequest) (*inpb.TOTPStatusResponse, error) {
	if err := a.validator.Validate(ctx, auth.ValidateActiveUser()); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "not authenticated: %v", err)
	}
	username, err := a.validatorCus.GetUsername(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	enabled, err := a.otpValidator.IsEnabled(ctx, username)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &inpb.TOTPStatusResponse{
		Enabled: enabled,
	}, nil
}

// GetTOTPConfiguration generates a new TOTP configuration for the user
func (a *InternalUserAPI) GetTOTPConfiguration(ctx context.Context, req *inpb.GetTOTPConfigurationRequest) (*inpb.GetTOTPConfigurationResponse, error) {
	if err := a.validator.Validate(ctx, auth.ValidateActiveUser()); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "not authenticated: %v", err)
	}
	username, err := a.validatorCus.GetUsername(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	cfg, err := a.otpValidator.NewConfiguration(ctx, username)
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
	if err := a.validator.Validate(ctx, auth.ValidateActiveUser()); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "not authenticated: %v", err)
	}
	username, err := a.validatorCus.GetUsername(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	otp := a.validatorCus.GetOTP(ctx)
	if err := a.otpValidator.Enable(ctx, username, otp); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "not authenticated: %v", err)
	}
	return &inpb.TOTPStatusResponse{
		Enabled: true,
	}, nil
}

// DisableTOTP disables TOTP for the user
func (a *InternalUserAPI) DisableTOTP(ctx context.Context, req *inpb.TOTPStatusRequest) (*inpb.TOTPStatusResponse, error) {
	if err := a.validator.Validate(ctx, auth.ValidateActiveUser()); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "not authenticated: %v", err)
	}
	if err := a.validatorCus.ValidateOTP(ctx); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "OTP is not present or not valid")
	}
	username, err := a.validatorCus.GetUsername(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	if err := a.otpValidator.Disable(ctx, username); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	return &inpb.TOTPStatusResponse{
		Enabled: false,
	}, nil
}

// GetRecoveryCodes returns the list of recovery codes for the user
func (a *InternalUserAPI) GetRecoveryCodes(ctx context.Context, req *inpb.GetRecoveryCodesRequest) (*inpb.GetRecoveryCodesResponse, error) {
	if err := a.validator.Validate(ctx, auth.ValidateActiveUser()); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "not authenticated: %v", err)
	}
	if err := a.validatorCus.ValidateOTP(ctx); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "OTP is not present or not valid")
	}
	username, err := a.validatorCus.GetUsername(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	codes, err := a.otpValidator.GetRecoveryCodes(ctx, username, req.Regenerate)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &inpb.GetRecoveryCodesResponse{
		RecoveryCode: codes,
	}, nil
}

// ConfirmRegistration checks provided security token and activates user
func (a *InternalUserAPI) ConfirmRegistration(ctx context.Context, req *inpb.ConfirmRegistrationRequest) (*inpb.ConfirmRegistrationResponse, error) {
	user, err := storage.GetUserByToken(storage.DB(), req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}

	jwt, err := storage.GetUserToken(user)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}

	return &inpb.ConfirmRegistrationResponse{
		Id:       user.ID,
		IsAdmin:  user.IsAdmin,
		IsActive: user.IsActive,
		Jwt:      jwt,
	}, status.Errorf(codes.OK, "")
}

// FinishRegistration sets new user password and creates a new organization
func (a *InternalUserAPI) FinishRegistration(ctx context.Context, req *inpb.FinishRegistrationRequest) (*empty.Empty, error) {
	if err := a.validator.Validate(ctx, auth.ValidateUserAccess(req.UserId, auth.FinishRegistration)); err != nil {
		log.Println("UpdatePassword", err)
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	org := storage.Organization{
		Name:            req.OrganizationName,
		DisplayName:     req.OrganizationDisplayName,
		CanHaveGateways: true,
	}

	err := storage.Transaction(func(tx sqlx.Ext) error {
		err := storage.FinishRegistration(tx, req.UserId, req.Password)
		if err != nil {
			return helpers.ErrToRPCError(err)
		}

		err = storage.CreateOrganization(ctx, tx, &org)
		if err != nil {
			return helpers.ErrToRPCError(err)
		}

		err = storage.CreateOrganizationUser(ctx, tx, org.ID, req.UserId, true, false, false)
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
