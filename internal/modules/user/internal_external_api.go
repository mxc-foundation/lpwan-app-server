package user

import (
	"context"
	"os"

	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	inpb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/oidc"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
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
	jwt, err := a.Store.LoginUserByPassword(ctx, req.Username, req.Password)
	if nil != err {
		return nil, helpers.ErrToRPCError(err)
	}

	return &inpb.LoginResponse{Jwt: jwt}, nil
}

// Profile returns the user profile.
func (a *InternalUserAPI) Profile(ctx context.Context, req *empty.Empty) (*inpb.ProfileResponse, error) {
	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
		validateActiveUser()); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	user, err := a.Validator.GetUser(ctx)
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
			Username:   prof.User.Email,
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
	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
		validateActiveUser()); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	isAdmin, err := a.Validator.GetIsAdmin(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	user, err := a.Validator.GetUser(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	results, err := storage.GlobalSearch(ctx, storage.DB(), user.ID, isAdmin, req.Search, int(req.Limit), int(req.Offset))
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
		"languange": inpb.Language_name[int32(req.Language)],
	}).Info(logInfo)

	user := User{
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

	err = email.SendInvite(obj.Email, *obj.SecurityToken, email.EmailLanguage(inpb.Language_name[int32(req.Language)]), email.RegistrationConfirmation)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// ConfirmRegistration checks provided security token and activates user
func (a *InternalUserAPI) ConfirmRegistration(ctx context.Context, req *inpb.ConfirmRegistrationRequest) (*inpb.ConfirmRegistrationResponse, error) {
	user, err := a.Store.GetUserByToken(req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}

	jwt, err := a.Store.GetUserToken(user)
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
	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx, validateUserAccess(req.UserId, FinishRegistration)); err != nil {
		log.Println("UpdatePassword", err)
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	org := organization.Organization{
		Name:            req.OrganizationName,
		DisplayName:     req.OrganizationDisplayName,
		CanHaveGateways: true,
	}

	err := storage.Transaction(func(tx sqlx.Ext) error {
		pwdHash, _ := hash(req.Password, saltSize, config.C.General.PasswordHashIterations)
		err := a.Store.FinishRegistration(req.UserId, pwdHash)
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

// GetVerifyingGoogleRecaptcha defines the request and response to verify the google recaptcha
func (a *InternalUserAPI) GetVerifyingGoogleRecaptcha(ctx context.Context, req *inpb.GoogleRecaptchaRequest) (*inpb.GoogleRecaptchaResponse, error) {
	res, err := IsPassVerifyingGoogleRecaptcha(req.Response, req.Remoteip)
	if err != nil {
		log.WithError(err).Error("Cannot verify from google recaptcha")
		return &inpb.GoogleRecaptchaResponse{}, err
	}

	return &inpb.GoogleRecaptchaResponse{Success: res.Success, ChallengeTs: res.ChallengeTs, Hostname: res.Hostname}, nil
}

// GetTOTPStatus returns info about TOTP status for the current user
func (a *InternalUserAPI) GetTOTPStatus(ctx context.Context, req *inpb.TOTPStatusRequest) (*inpb.TOTPStatusResponse, error) {
	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx, validateActiveUser()); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "not authenticated: %v", err)
	}
	username, err := a.Validator.otpValidator.JwtValidator.GetUsername(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	enabled, err := a.Validator.otpValidator.IsEnabled(ctx, username)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &inpb.TOTPStatusResponse{
		Enabled: enabled,
	}, nil
}

// GetTOTPConfiguration generates a new TOTP configuration for the user
func (a *InternalUserAPI) GetTOTPConfiguration(ctx context.Context, req *inpb.GetTOTPConfigurationRequest) (*inpb.GetTOTPConfigurationResponse, error) {
	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx, validateActiveUser()); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "not authenticated: %v", err)
	}
	username, err := a.Validator.otpValidator.JwtValidator.GetUsername(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	cfg, err := a.Validator.otpValidator.NewConfiguration(ctx, username)
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
	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx, validateActiveUser()); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "not authenticated: %v", err)
	}
	username, err := a.Validator.otpValidator.JwtValidator.GetUsername(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	otp := a.Validator.otpValidator.JwtValidator.GetOTP(ctx)
	if err := a.Validator.otpValidator.Enable(ctx, username, otp); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "not authenticated: %v", err)
	}
	return &inpb.TOTPStatusResponse{
		Enabled: true,
	}, nil
}

// DisableTOTP disables TOTP for the user
func (a *InternalUserAPI) DisableTOTP(ctx context.Context, req *inpb.TOTPStatusRequest) (*inpb.TOTPStatusResponse, error) {
	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx, validateActiveUser()); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "not authenticated: %v", err)
	}
	if err := a.Validator.otpValidator.ValidateOTP(ctx); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "OTP is not present or not valid")
	}
	username, err := a.Validator.otpValidator.JwtValidator.GetUsername(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	if err := a.Validator.otpValidator.Disable(ctx, username); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	return &inpb.TOTPStatusResponse{
		Enabled: false,
	}, nil
}

// GetRecoveryCodes returns the list of recovery codes for the user
func (a *InternalUserAPI) GetRecoveryCodes(ctx context.Context, req *inpb.GetRecoveryCodesRequest) (*inpb.GetRecoveryCodesResponse, error) {
	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx, validateActiveUser()); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "not authenticated: %v", err)
	}
	if err := a.Validator.otpValidator.ValidateOTP(ctx); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "OTP is not present or not valid")
	}
	username, err := a.Validator.otpValidator.JwtValidator.GetUsername(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	codes, err := a.Validator.otpValidator.GetRecoveryCodes(ctx, username, req.Regenerate)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &inpb.GetRecoveryCodesResponse{
		RecoveryCode: codes,
	}, nil
}

// CreateAPIKey creates the given API key.
func (a *InternalUserAPI) CreateAPIKey(ctx context.Context, req *inpb.CreateAPIKeyRequest) (*inpb.CreateAPIKeyResponse, error) {
	apiKey := req.GetApiKey()

	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
		validateAPIKeysAccess(Create, apiKey.GetOrganizationId(), apiKey.GetApplicationId())); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	if apiKey.GetIsAdmin() && (apiKey.GetOrganizationId() != 0 || apiKey.GetApplicationId() != 0) {
		return nil, status.Errorf(codes.InvalidArgument, "when is_admin is true, organization_id and application_id must be left blank")
	}

	var organizationID *int64
	var applicationID *int64

	if id := apiKey.GetOrganizationId(); id != 0 {
		organizationID = &id
	}

	if id := apiKey.GetApplicationId(); id != 0 {
		applicationID = &id
	}

	ak := storage.APIKey{
		Name:           apiKey.GetName(),
		IsAdmin:        apiKey.GetIsAdmin(),
		OrganizationID: organizationID,
		ApplicationID:  applicationID,
	}

	jwtToken, err := storage.CreateAPIKey(ctx, storage.DB(), &ak)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &inpb.CreateAPIKeyResponse{
		Id:       ak.ID.String(),
		JwtToken: jwtToken,
	}, nil
}

// ListAPIKeys lists the API keys.
func (a *InternalUserAPI) ListAPIKeys(ctx context.Context, req *inpb.ListAPIKeysRequest) (*inpb.ListAPIKeysResponse, error) {
	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
		validateAPIKeysAccess(List, req.GetOrganizationId(), req.GetApplicationId())); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	if req.GetIsAdmin() && (req.GetOrganizationId() != 0 || req.GetApplicationId() != 0) {
		return nil, status.Errorf(codes.InvalidArgument, "when is_admin is true, organization_id and application_id must be left blank")
	}

	filters := storage.APIKeyFilters{
		IsAdmin: req.GetIsAdmin(),
		Limit:   int(req.GetLimit()),
		Offset:  int(req.GetOffset()),
	}

	if id := req.GetOrganizationId(); id != 0 {
		filters.OrganizationID = &id
	}

	if id := req.GetApplicationId(); id != 0 {
		filters.ApplicationID = &id
	}

	count, err := storage.GetAPIKeyCount(ctx, storage.DB(), filters)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	apiKeys, err := storage.GetAPIKeys(ctx, storage.DB(), filters)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := inpb.ListAPIKeysResponse{
		TotalCount: int64(count),
	}

	for _, apiKey := range apiKeys {
		ak := inpb.APIKey{
			Id:      apiKey.ID.String(),
			Name:    apiKey.Name,
			IsAdmin: apiKey.IsAdmin,
		}

		if apiKey.OrganizationID != nil {
			ak.OrganizationId = *apiKey.OrganizationID
		}

		if apiKey.ApplicationID != nil {
			ak.ApplicationId = *apiKey.ApplicationID
		}

		resp.Result = append(resp.Result, &ak)
	}

	return &resp, nil
}

// DeleteAPIKey deletes the given API key.
func (a *InternalUserAPI) DeleteAPIKey(ctx context.Context, req *inpb.DeleteAPIKeyRequest) (*empty.Empty, error) {
	apiKeyID, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "api_key: %s", err)
	}

	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
		validateAPIKeyAccess(Delete, apiKeyID)); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	if err := storage.DeleteAPIKey(ctx, storage.DB(), apiKeyID); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// Settings returns the global settings.
func (a *InternalUserAPI) Settings(ctx context.Context, _ *empty.Empty) (*inpb.SettingsResponse, error) {
	return &inpb.SettingsResponse{
		/*		Branding: &inpb.Branding{
					Registration: brandingRegistration,
					Footer:       brandingFooter,
				},
				OpenidConnect: &inpb.OpenIDConnect{
					Enabled:    openIDConnectEnabled,
					LoginLabel: openIDLoginLabel,
					LoginUrl:   "/auth/oidc/login",
				},*/
	}, nil
}

// OpenIDConnectLogin performs an OpenID Connect login.
func (a *InternalUserAPI) OpenIDConnectLogin(ctx context.Context, req *inpb.OpenIDConnectLoginRequest) (*inpb.OpenIDConnectLoginResponse, error) {
	/*oidcUser, err := oidc.GetUser(ctx, req.Code, req.State)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	if !oidcUser.EmailVerified {
		return nil, status.Errorf(codes.FailedPrecondition, "email address must be verified before you can login")
	}

	var user User

	// try to get the user by external ID.
	user, err = a.Store.GetUserByExternalID(ctx, oidcUser.ExternalID)
	if err != nil {
		if err == storage.ErrDoesNotExist {
			// try to get the user by email and set the external id.
			user, err = a.Store.GetUserByEmail(ctx, oidcUser.Email)
			if err != nil {
				// we did not find the user by external_id or email and registration is enabled.
				if err == storage.ErrDoesNotExist && registrationEnabled {
					user, err = a.createAndProvisionUser(ctx, oidcUser)
					if err != nil {
						return nil, helpers.ErrToRPCError(err)
					}
				} else {
					return nil, helpers.ErrToRPCError(err)
				}
			}
			user.ExternalID = &oidcUser.ExternalID
		} else {
			return nil, helpers.ErrToRPCError(err)
		}
	}

	// update the user
	user.Email = oidcUser.Email
	user.EmailVerified = oidcUser.EmailVerified
	if err := a.Store.UpdateUser(ctx, &user); err != nil {
		fmt.Println("SDFSDFSDFSDF")
		return nil, helpers.ErrToRPCError(err)
	}

	// get the jwt token
	token, err := a.Store.GetUserToken(user)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &inpb.OpenIDConnectLoginResponse{
		JwtToken: token,
	}, nil*/

	return &inpb.OpenIDConnectLoginResponse{}, nil
}

func (a *InternalUserAPI) createAndProvisionUser(ctx context.Context, user oidc.User) (User, error) {
	/*u := User{
		IsActive:      true,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		ExternalID:    &user.ExternalID,
	}

	if err := a.Store.CreateUser(ctx, &u); err != nil {
		return User{}, errors.Wrap(err, "create user error")
	}

	if registrationCallbackURL == "" {
		return u, nil
	}

	if err := a.provisionUser(ctx, u); err != nil {
		if err := a.Store.DeleteUser(ctx, u.ID); err != nil {
			return User{}, errors.Wrap(err, "delete user error after failed user provisioning")
		}

		log.WithError(err).Error("api/external: provision user error")

		return User{}, errors.New("error provisioning user")
	}

	return u, nil*/
	return User{}, nil
}

func (a *InternalUserAPI) provisionUser(ctx context.Context, u User) error {
	/*	req, err := http.NewRequestWithContext(ctx, "POST", registrationCallbackURL, nil)
		if err != nil {
			return errors.Wrap(err, "new request error")
		}
		q := req.URL.Query()
		q.Add("user_id", fmt.Sprintf("%d", u.ID))
		req.URL.RawQuery = q.Encode()

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return errors.Wrap(err, "make registration callback request error")
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("registration callback must return 200, received: %d (%s)", resp.StatusCode, resp.Status)
		}*/

	return nil
}
