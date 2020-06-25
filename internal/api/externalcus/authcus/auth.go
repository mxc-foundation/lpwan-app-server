package authcus

import (
	"fmt"
	"regexp"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"

	"github.com/mxc-foundation/lpwan-app-server/internal/otp"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
)

var validAuthorizationRegexp = regexp.MustCompile(`(?i)^bearer (.*)$`)

// Claims defines the struct containing the token claims.
type Claims struct {
	jwt.StandardClaims

	// Username defines the identity of the user.
	Username string `json:"username"`
	// UserID defines the ID of th user.
	UserID int64 `json:"user_id"`

	// APIKeyID defines the API key ID.
	APIKeyID uuid.UUID `json:"api_key_id"`

	// OTP code if it is present, not a part of JWT
	OTP string `json:"-"`
}

// Validator defines the interface a validator needs to implement.
type Validator interface {
	// Validate validates the given set of validators against the given context.
	// Must return after the first validator function either returns true or
	// and error. The way how the validation must be seens is:
	//   if validatorFunc1 || validatorFunc2 || validatorFunc3 ...
	// In case multiple validators must validate to true, then a validator
	// func needs to be implemented which validates a given set of funcs as:
	//   if validatorFunc1 && validatorFunc2 && ValidatorFunc3 ...
	Validate(context.Context, ...ValidatorFunc) error

	// GetSubject returns the claim subject.
	GetSubject(context.Context) (string, error)

	// GetUser returns the user object.
	GetUser(context.Context) (storage.User, error)

	// GetAPIKey returns the API key ID.
	GetAPIKeyID(context.Context) (uuid.UUID, error)

	// GetUsername returns the name of the authenticated user.
	GetUsername(context.Context) (string, error)

	// GetOTP returns OTP code
	GetOTP(context.Context) string

	// ValidateOTP validates OTP and returns the error if it is not valid
	ValidateOTP(context.Context) error

	// GetIsAdmin returns if the authenticated user is a global admin.
	GetIsAdmin(context.Context) (bool, error)
}

// ValidatorFunc defines the signature of a claim validator function.
// It returns a bool indicating if the validation passed or failed and an
// error in case an error occurred (e.g. db connectivity).
type ValidatorFunc func(sqlx.Queryer, *Claims) (bool, error)

// JWTValidator validates JWT tokens.
type JWTValidator struct {
	db           sqlx.Ext
	secret       string
	algorithm    string
	otpValidator *otp.Validator
}

// NewJWTValidator creates a new JWTValidator.
func NewJWTValidator(db sqlx.Ext, algorithm, secret string, otpValidator *otp.Validator) *JWTValidator {
	return &JWTValidator{
		db:           db,
		secret:       secret,
		algorithm:    algorithm,
		otpValidator: otpValidator,
	}
}

// Validate validates the token from the given context against the given
// validator funcs.
func (v JWTValidator) Validate(ctx context.Context, funcs ...ValidatorFunc) error {
	claims, err := v.getClaims(ctx)
	if err != nil {
		return err
	}

	for _, f := range funcs {
		ok, err := f(v.db, claims)
		if err != nil {
			return errors.Wrap(err, "validator func error")
		}
		if ok {
			return nil
		}
	}

	return ErrNotAuthorized
}

// GetSubject returns the subject of the claim.
func (v JWTValidator) GetSubject(ctx context.Context) (string, error) {
	claims, err := v.getClaims(ctx)
	if err != nil {
		return "", err
	}

	return claims.Subject, nil
}

// GetAPIKeyID returns the API key of the token.
func (v JWTValidator) GetAPIKeyID(ctx context.Context) (uuid.UUID, error) {
	claims, err := v.getClaims(ctx)
	if err != nil {
		return uuid.Nil, err
	}

	return claims.APIKeyID, nil
}

// GetUser returns the user object.
func (v JWTValidator) GetUser(ctx context.Context) (storage.User, error) {
	claims, err := v.getClaims(ctx)
	if err != nil {
		return storage.User{}, err
	}

	if claims.Subject != SubjectUser {
		return storage.User{}, errors.New("subject must be user")
	}

	if claims.UserID != 0 {
		return storage.GetUser(ctx, v.db, claims.UserID)
	}

	if claims.Username != "" {
		return storage.GetUserByEmail(ctx, v.db, claims.Username)
	}

	return storage.User{}, errors.New("no username or user_id in claims")
}

// GetUsername returns the username of the authenticated user.
func (v JWTValidator) GetUsername(ctx context.Context) (string, error) {
	claims, err := v.getClaims(ctx)
	if err != nil {
		return "", err
	}

	return claims.Username, nil
}

// GetIsAdmin returns if the authenticated user is a global amin.
func (v JWTValidator) GetIsAdmin(ctx context.Context) (bool, error) {
	claims, err := v.getClaims(ctx)
	if err != nil {
		return false, err
	}

	user, err := storage.GetUserByUsername(ctx, v.db, claims.Username)
	if err != nil {
		return false, errors.Wrap(err, "get user by username error")
	}

	return user.IsAdmin, nil
}

// GetOTP returns OTP from the context
func (v JWTValidator) GetOTP(ctx context.Context) string {
	return getOTPFromContext(ctx)
}

// ValidateOTP validates OTP and returns the error if it is not valid
func (v JWTValidator) ValidateOTP(ctx context.Context) error {
	claims, err := v.getClaims(ctx)
	if err != nil {
		return err
	}
	enabled, err := v.otpValidator.IsEnabled(ctx, claims.Username)
	if err != nil {
		return err
	}
	if !enabled {
		return errors.New("OTP is not enabled")
	}
	return v.otpValidator.Validate(ctx, claims.Username, claims.OTP)
}

func (v JWTValidator) getClaims(ctx context.Context) (*Claims, error) {
	tokenStr, err := getTokenFromContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get token from context error")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Header["alg"] != v.algorithm {
			return nil, ErrInvalidAlgorithm
		}
		return []byte(v.secret), nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "jwt parse error")
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		// no need to use a static error, this should never happen
		return nil, fmt.Errorf("api/auth: expected *Claims, got %T", token.Claims)
	}

	claims.OTP = getOTPFromContext(ctx)

	return claims, nil
}

func getTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", ErrNoMetadataInContext
	}

	token, ok := md["authorization"]
	if !ok || len(token) == 0 {
		return "", ErrNoAuthorizationInMetadata
	}

	match := validAuthorizationRegexp.FindStringSubmatch(token[0])

	// authorization header should respect RFC1945
	if len(match) == 0 {
		log.Warning("Deprecated Authorization header : RFC1945 format expected : Authorization: <type> <credentials>")
		return token[0], nil
	}

	return match[1], nil
}

func getOTPFromContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	if len(md["x-otp"]) == 1 {
		return md["x-otp"][0]
	}
	return ""
}
