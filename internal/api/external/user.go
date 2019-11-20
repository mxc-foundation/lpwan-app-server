package external

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jmoiron/sqlx"
	"github.com/mxc-foundation/lpwan-server/api/ns"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	pb "github.com/mxc-foundation/lpwan-app-server/api"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/auth"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"

	"github.com/gofrs/uuid"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/email"
	log "github.com/sirupsen/logrus"
)

// UserAPI exports the User related functions.
type UserAPI struct {
	validator auth.Validator
}

// InternalUserAPI exports the internal User related functions.
type InternalUserAPI struct {
	validator auth.Validator
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
		return nil, grpc.Errorf(codes.InvalidArgument, "user must not be nil")
	}

	if err := a.validator.Validate(ctx,
		auth.ValidateUsersAccess(auth.Create)); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	// validate if the client has admin rights for the given organizations
	// to which the user must be linked
	for _, org := range req.Organizations {
		if err := a.validator.Validate(ctx,
			auth.ValidateIsOrganizationAdmin(org.OrganizationId)); err != nil {
			return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
		}
	}

	user := storage.User{
		Username:   req.User.Username,
		SessionTTL: req.User.SessionTtl,
		IsAdmin:    req.User.IsAdmin,
		IsActive:   req.User.IsActive,
		Email:      req.User.Email,
		Note:       req.User.Note,
	}

	isAdmin, err := a.validator.GetIsAdmin(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	if !isAdmin {
		// non-admin users are not able to modify the fields below
		user.IsAdmin = false
		user.IsActive = true
		user.SessionTTL = 0
	}

	var userID int64

	err = storage.Transaction(func(tx sqlx.Ext) error {
		userID, err = storage.CreateUser(ctx, tx, &user, req.Password)
		if err != nil {
			return err
		}

		for _, org := range req.Organizations {
			if err := storage.CreateOrganizationUser(ctx, tx, org.OrganizationId, userID, org.IsAdmin, org.IsDeviceAdmin, org.IsGatewayAdmin); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &pb.CreateUserResponse{Id: userID}, nil
}

// Get returns the user matching the given ID.
func (a *UserAPI) Get(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if err := a.validator.Validate(ctx,
		auth.ValidateUserAccess(req.Id, auth.Read)); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	user, err := storage.GetUser(ctx, storage.DB(), req.Id)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := pb.GetUserResponse{
		User: &pb.User{
			Id:         user.ID,
			Username:   user.Username,
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

// List lists the users.
func (a *UserAPI) List(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserResponse, error) {
	if err := a.validator.Validate(ctx,
		auth.ValidateUsersAccess(auth.List)); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	users, err := storage.GetUsers(ctx, storage.DB(), int(req.Limit), int(req.Offset), req.Search)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	totalUserCount, err := storage.GetUserCount(ctx, storage.DB(), req.Search)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := pb.ListUserResponse{
		TotalCount: int64(totalUserCount),
	}

	for _, u := range users {
		row := pb.UserListItem{
			Id:         u.ID,
			Username:   u.Username,
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
		return nil, grpc.Errorf(codes.InvalidArgument, "user must not be nil")
	}

	if err := a.validator.Validate(ctx,
		auth.ValidateUserAccess(req.User.Id, auth.Update)); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	userUpdate := storage.UserUpdate{
		ID:         req.User.Id,
		Username:   req.User.Username,
		IsAdmin:    req.User.IsAdmin,
		IsActive:   req.User.IsActive,
		SessionTTL: req.User.SessionTtl,
		Email:      req.User.Email,
		Note:       req.User.Note,
	}

	err := storage.UpdateUser(ctx, storage.DB(), userUpdate)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// Delete deletes the user matching the given ID.
func (a *UserAPI) Delete(ctx context.Context, req *pb.DeleteUserRequest) (*empty.Empty, error) {
	if err := a.validator.Validate(ctx,
		auth.ValidateUserAccess(req.Id, auth.Delete)); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
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
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	user, err := storage.GetUser(ctx, storage.DB(), req.UserId)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	if user.Username == storage.DemoUser {
		return nil, helpers.ErrToRPCError(errors.New(fmt.Sprintf("User %s can not change password.", storage.DemoUser)))
	}

	err = storage.UpdatePassword(ctx, storage.DB(), req.UserId, req.Password)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// NewInternalUserAPI creates a new InternalUserAPI.
func NewInternalUserAPI(validator auth.Validator) *InternalUserAPI {
	return &InternalUserAPI{
		validator: validator,
	}
}

// Login validates the login request and returns a JWT token.
func (a *InternalUserAPI) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	jwt, err := storage.LoginUser(ctx, storage.DB(), req.Username, req.Password)
	if nil != err {
		return nil, helpers.ErrToRPCError(err)
	}

	return &pb.LoginResponse{Jwt: jwt}, nil
}

func IsPassVerifyingGoogleRecaptcha(response string, remoteip string) (*pb.GoogleRecaptchaResponse, error) {
	secret := config.C.Recaptcha.Secret
	postURL := config.C.Recaptcha.HostServer

	postStr := url.Values{"secret": {secret}, "response": {response}, "remoteip": {remoteip}}
	responsePost, err := http.PostForm(postURL, postStr)

	if err != nil {
		log.Warn(err)
		return &pb.GoogleRecaptchaResponse{}, err
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
		return &pb.GoogleRecaptchaResponse{}, err
	}

	g_response := &pb.GoogleRecaptchaResponse{}
	err = json.Unmarshal(body, &g_response)
	if err != nil {
		fmt.Println("unmarshal response", err)
	}

	return g_response, nil
}

func (a *InternalUserAPI) GetVerifyingGoogleRecaptcha(ctx context.Context, req *pb.GoogleRecaptchaRequest) (*pb.GoogleRecaptchaResponse, error) {
	res, err := IsPassVerifyingGoogleRecaptcha(req.Response, req.Remoteip)
	if err != nil {
		log.WithError(err).Error("Cannot verify from google recaptcha")
		return &pb.GoogleRecaptchaResponse{}, err
	}

	return &pb.GoogleRecaptchaResponse{Success: res.Success, ChallengeTs: res.ChallengeTs, Hostname: res.Hostname}, nil
}

type claims struct {
	Username string `json:"username"`
}

// Profile returns the user profile.
func (a *InternalUserAPI) Profile(ctx context.Context, req *empty.Empty) (*pb.ProfileResponse, error) {
	if err := a.validator.Validate(ctx,
		auth.ValidateActiveUser()); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	username, err := a.validator.GetUsername(ctx)
	if nil != err {
		return nil, helpers.ErrToRPCError(err)
	}

	// Get the user id based on the username.
	user, err := storage.GetUserByUsername(ctx, storage.DB(), username)
	if nil != err {
		return nil, helpers.ErrToRPCError(err)
	}

	prof, err := storage.GetProfile(ctx, storage.DB(), user.ID)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := pb.ProfileResponse{
		User: &pb.User{
			Id:         prof.User.ID,
			Username:   prof.User.Username,
			SessionTtl: prof.User.SessionTTL,
			IsAdmin:    prof.User.IsAdmin,
			IsActive:   prof.User.IsActive,
		},
		Settings: &pb.ProfileSettings{
			DisableAssignExistingUsers: auth.DisableAssignExistingUsers,
		},
	}

	for _, org := range prof.Organizations {
		row := pb.OrganizationLink{
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
func (a *InternalUserAPI) Branding(ctx context.Context, req *empty.Empty) (*pb.BrandingResponse, error) {
	resp := pb.BrandingResponse{
		Logo:         brandingHeader,
		Registration: brandingRegistration,
		Footer:       brandingFooter,
	}

	return &resp, nil
}

// GlobalSearch performs a global search.
func (a *InternalUserAPI) GlobalSearch(ctx context.Context, req *pb.GlobalSearchRequest) (*pb.GlobalSearchResponse, error) {
	if err := a.validator.Validate(ctx,
		auth.ValidateActiveUser()); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	isAdmin, err := a.validator.GetIsAdmin(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	username, err := a.validator.GetUsername(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	results, err := storage.GlobalSearch(ctx, storage.DB(), username, isAdmin, req.Search, int(req.Limit), int(req.Offset))
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

// RegisterUser adds new user and sends activation email
func (a *InternalUserAPI) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*empty.Empty, error) {
	user := storage.User{
		Username:   req.Email,
		SessionTTL: 0,
		IsAdmin:    false,
		IsActive:   false,
	}

	u, err := uuid.NewV4()
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	token := u.String()

	obj, err := storage.GetUserByUsername(ctx, storage.DB(), user.Username)
	if err == storage.ErrDoesNotExist {
		// user has never been created yet
		err = storage.RegisterUser(storage.DB(), &user, token)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}

		// get user again
		obj, err = storage.GetUserByUsername(ctx, storage.DB(), user.Username)
		if err != nil {
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

	err = email.SendInvite(obj.Username, *obj.SecurityToken, int32(req.Language))
	if err != nil {
		log.WithFields(log.Fields{
			"username": user.Username,
		}).Info("Send email error!")
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// ConfirmRegistration checks provided security token and activates user
func (a *InternalUserAPI) ConfirmRegistration(ctx context.Context, req *pb.ConfirmRegistrationRequest) (*pb.ConfirmRegistrationResponse, error) {
	user, err := storage.GetUserByToken(storage.DB(), req.Token)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	log.Println("Confirming GetJwt", user.Username)
	jwt, err := storage.MakeJWT(user.Username, user.SessionTTL)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &pb.ConfirmRegistrationResponse{
		Id:       user.ID,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
		IsActive: user.IsActive,
		Jwt:      jwt,
	}, nil
}

// FinishRegistration sets new user password and creates a new organization
func (a *InternalUserAPI) FinishRegistration(ctx context.Context, req *pb.FinishRegistrationRequest) (*empty.Empty, error) {
	if err := a.validator.Validate(ctx, auth.ValidateUserAccess(req.UserId, auth.FinishRegistration)); err != nil {
		log.Println("UpdatePassword", err)
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
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

		// add admin user into this organization
		adminUser, err := storage.GetUserByUsername(ctx, tx, "admin")
		if err == nil {
			err = storage.CreateOrganizationUser(ctx, tx, org.ID, adminUser.ID, false, false, false)
			if err != nil {
				log.WithError(err).Error("Insert admin into organization ", org.ID, " failed")
			}
		} else {
			log.WithError(err).Error("Get user by username 'admin' failed")
		}

		err = storage.CreateOrganizationUser(ctx, tx, org.ID, req.UserId, true, false, false)
		if err != nil {
			return helpers.ErrToRPCError(err)
		}

		// add service profile for this organization
		networkServerList, err := storage.GetNetworkServers(ctx, tx, 10, 0)
		if err == nil && len(networkServerList) >= 1 {
			sp := storage.ServiceProfile{
				OrganizationID:  org.ID,
				NetworkServerID: networkServerList[0].ID,
				Name:            "service_profile_" + org.Name,
				ServiceProfile: ns.ServiceProfile{
					UlRate:                 0,
					UlBucketSize:           0,
					DlRate:                 0,
					DlBucketSize:           0,
					AddGwMetadata:          true,
					DevStatusReqFreq:       0,
					ReportDevStatusBattery: true,
					ReportDevStatusMargin:  true,
					DrMin:                  0,
					DrMax:                  0,
					ChannelMask:            []byte(""),
					PrAllowed:              true,
					HrAllowed:              true,
					RaAllowed:              true,
					NwkGeoLoc:              true,
					TargetPer:              0,
					MinGwDiversity:         0,
					UlRatePolicy:           ns.RatePolicy_DROP,
					DlRatePolicy:           ns.RatePolicy_DROP,
				},
			}

			err := storage.CreateServiceProfile(ctx, tx, &sp)
			if err != nil {
				log.WithError(err).Error("Add service profile for organization_id = ", org.ID, " failed")
			}
		} else {
			log.WithError(err).Error("Get network server for organization_id = 0 failed")
		}

		return nil
	})

	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}
