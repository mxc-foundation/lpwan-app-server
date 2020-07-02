package user

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jmoiron/sqlx"

	"github.com/pkg/errors"
	"golang.org/x/crypto/pbkdf2"

	inpb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/organization"
)

type UserStore interface {
	CreateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, id int64) (User, error)
	GetUserByExternalID(ctx context.Context, externalID string) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserCount(ctx context.Context) (int, error)
	GetUsers(ctx context.Context, limit, offset int) ([]User, error)
	UpdateUser(ctx context.Context, u *User) error
	DeleteUser(ctx context.Context, id int64) error
	LoginUserByPassword(ctx context.Context, email string, password string) (string, error)
	GetProfile(ctx context.Context, id int64) (UserProfile, error)
	GetUserToken(u User) (string, error)
	RegisterUser(user *User, token string) error
	GetUserByToken(token string) (User, error)
	GetTokenByUsername(ctx context.Context, username string) (string, error)
	FinishRegistration(userID int64, pwdHash string) error
}

// UserAPI exports the User related functions.
type UserAPI struct {
	Validator *Validator
	Store     UserStore
}

// NewUserAPI creates a new UserAPI.
func NewUserAPI(api UserAPI) *UserAPI {
	userAPI = UserAPI{
		Validator: api.Validator,
		Store:     api.Store,
	}

	return &userAPI
}

var (
	userAPI UserAPI
)

func GetUserAPI() *UserAPI {
	return &userAPI
}

// Create creates the given user.
func (a *UserAPI) Create(ctx context.Context, req *inpb.CreateUserRequest) (*inpb.CreateUserResponse, error) {
	if req.User == nil {
		return nil, status.Errorf(codes.InvalidArgument, "user must not be nil")
	}

	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
		validateUsersAccess(Create)); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	user := User{
		Username:   req.User.Username,
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
		err := a.Store.CreateUser(ctx, &user)
		if err != nil {
			return err
		}

		for _, org := range req.Organizations {
			if err := organization.GetOrganizationAPI().Store.CreateOrganizationUser(ctx, org.OrganizationId, user.Username, org.IsAdmin, org.IsDeviceAdmin, org.IsGatewayAdmin); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &inpb.CreateUserResponse{Id: user.ID}, nil
}

// Get returns the user matching the given ID.
func (a *UserAPI) Get(ctx context.Context, req *inpb.GetUserRequest) (*inpb.GetUserResponse, error) {
	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
		validateUserAccess(req.Id, Read)); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	user, err := a.Store.GetUser(ctx, req.Id)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := inpb.GetUserResponse{
		User: &inpb.User{
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
	u, err := a.Store.GetUserByEmail(ctx, req.UserEmail)
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
func (a *UserAPI) List(ctx context.Context, req *inpb.ListUserRequest) (*inpb.ListUserResponse, error) {
	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
		validateUsersAccess(List)); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	users, err := a.Store.GetUsers(ctx, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	totalUserCount, err := a.Store.GetUserCount(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := inpb.ListUserResponse{
		TotalCount: int64(totalUserCount),
	}

	for _, u := range users {
		row := inpb.UserListItem{
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
func (a *UserAPI) Update(ctx context.Context, req *inpb.UpdateUserRequest) (*empty.Empty, error) {
	if req.User == nil {
		return nil, status.Errorf(codes.InvalidArgument, "user must not be nil")
	}

	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
		validateUserAccess(req.User.Id, Update)); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	err := storage.Transaction(func(tx sqlx.Ext) error {
		user, err := a.Store.GetUser(ctx, req.User.Id)
		if err != nil {
			return helpers.ErrToRPCError(err)
		}

		user.IsAdmin = req.User.IsAdmin
		user.IsActive = req.User.IsActive
		user.SessionTTL = req.User.SessionTtl
		user.Email = req.User.Email
		user.Note = req.User.Note

		err = a.Store.UpdateUser(ctx, &user)
		if err != nil {
			return helpers.ErrToRPCError(err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// Delete deletes the user matching the given ID.
func (a *UserAPI) Delete(ctx context.Context, req *inpb.DeleteUserRequest) (*empty.Empty, error) {
	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
		validateUserAccess(req.Id, Delete)); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	err := storage.Transaction(func(tx sqlx.Ext) error {
		err := a.Store.DeleteUser(ctx, req.Id)
		if err != nil {
			return helpers.ErrToRPCError(err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// UpdatePassword updates the password for the user matching the given ID.
func (a *UserAPI) UpdatePassword(ctx context.Context, req *inpb.UpdateUserPasswordRequest) (*empty.Empty, error) {
	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
		validateUserAccess(req.UserId, UpdateProfile)); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	user, err := a.Store.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	if err := user.SetPasswordHash(req.Password); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
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

func (a *UserAPI) GetOTPCode(ctx context.Context, req *inpb.GetOTPCodeRequest) (*inpb.GetOTPCodeResponse, error) {
	otp, err := a.Store.GetTokenByUsername(ctx, req.UserEmail)
	if err != nil {
		return nil, err
	}

	return &inpb.GetOTPCodeResponse{OtpCode: otp}, nil
}

const (
	// saltSize defines the salt size
	saltSize = 16
	//  defines the default session TTL
	DefaultSessionTTL = time.Hour * 24
)

var (
	// Any printable characters, at least 6 characters.
	passwordValidator = regexp.MustCompile(`^.{6,}$`)

	// Must contain @ (this is far from perfect)
	emailValidator = regexp.MustCompile(`.+@.+`)
)

// Validate validates the user data.
func (u *User) Validate() error {
	if !emailValidator.MatchString(u.Email) {
		return errors.New("invalid email")
	}

	return nil
}

// SetPasswordHash hashes the given password and sets it.
func (u *User) SetPasswordHash(pw string) error {
	if !passwordValidator.MatchString(pw) {
		return errors.New("invalid user password length")
	}

	pwHash, err := hash(pw, saltSize, config.C.General.PasswordHashIterations)
	if err != nil {
		return err
	}

	u.PasswordHash = pwHash

	return nil
}

// hashCompare verifies that passed password hashes to the same value as the
// passed passwordHash.
func (u *User) HashCompare(password string, passwordHash string) bool {
	// SPlit the hash string into its parts.
	hashSplit := strings.Split(passwordHash, "$")

	// Get the iterations and the salt and use them to encode the password
	// being compared.cre
	iterations, _ := strconv.Atoi(hashSplit[2])
	salt, _ := base64.StdEncoding.DecodeString(hashSplit[3])
	newHash := hashWithSalt(password, salt, iterations)
	return newHash == passwordHash
}

// Generate the hash of a password for storage in the database.
// NOTE: We Store the details of the hashing algorithm with the hash itself,
// making it easy to recreate the hash for password checking, even if we change
// the default criteria here.
func hash(password string, saltSize int, iterations int) (string, error) {
	// Generate a random salt value, 128 bits.
	salt := make([]byte, saltSize)
	_, err := rand.Read(salt)
	if err != nil {
		return "", errors.Wrap(err, "read random bytes error")
	}

	return hashWithSalt(password, salt, iterations), nil
}

func hashWithSalt(password string, salt []byte, iterations int) string {
	// Generate the hash.  This should be a little painful, adjust ITERATIONS
	// if it needs performance tweeking.  Greatly depends on the hardware.
	// NOTE: We Store these details with the returned hash, so changes will not
	// affect our ability to do password compares.
	hash := pbkdf2.Key([]byte(password), salt, iterations, sha512.Size, sha512.New)

	// Build up the parameters and hash into a single string so we can compare
	// other string to the same hash.  Note that the hash algorithm is hard-
	// coded here, as it is above.  Introducing alternate encodings must support
	// old encodings as well, and build this string appropriately.
	var buffer bytes.Buffer

	buffer.WriteString("PBKDF2$")
	buffer.WriteString("sha512$")
	buffer.WriteString(strconv.Itoa(iterations))
	buffer.WriteString("$")
	buffer.WriteString(base64.StdEncoding.EncodeToString(salt))
	buffer.WriteString("$")
	buffer.WriteString(base64.StdEncoding.EncodeToString(hash))

	return buffer.String()
}
