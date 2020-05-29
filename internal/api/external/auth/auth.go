package auth

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

// defaultSessionTTL defines the default session TTL
const defaultSessionTTL = 86400

var validAuthorizationRegexp = regexp.MustCompile(`(?i)^bearer (.*)$`)

// Claims defines the struct containing the token claims.
type Claims struct {
	// Username defines the identity of the user.
	Username string `json:"username"`

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

	// GetUsername returns the name of the authenticated user.
	GetUsername(context.Context) (string, error)

	// GetOTP returns OTP code
	GetOTP(context.Context) string

	// ValidateOTP validates OTP and returns the error if it is not valid
	ValidateOTP(context.Context) error

	// GetIsAdmin returns if the authenticated user is a global admin.
	GetIsAdmin(context.Context) (bool, error)

	// GetCredentials returns users credentials
	GetCredentials(context.Context, ...Option) (Credentials, error)

	// SignToken returns a signed token for the user
	SignToken(username string, ttl int64, audience []string) (string, error)
}

type options struct {
	audience    string
	requireOTP  bool
	limitedCred bool
}

// Option is used to configure validator checks
type Option func(opts *options)

// WithAudience requires that credentials presented included all the listed audiences
func WithAudience(audience string) Option {
	return func(opts *options) {
		opts.audience = audience
	}
}

// WithValidOTP requires that the request included valid OTP code
func WithValidOTP() Option {
	return func(opts *options) {
		opts.requireOTP = true
	}
}

// WithLimitedCredentials creates limited credentials
//
// Deprecated: do not use, this is only for the purposes of the registration
// process
func WithLimitedCredentials() Option {
	return func(opts *options) {
		opts.limitedCred = true
	}
}

// ValidatorFunc defines the signature of a claim validator function.
// It returns a bool indicating if the validation passed or failed and an
// error in case an error occured (e.g. db connectivity).
type ValidatorFunc func(sqlx.Queryer, *Claims) (bool, error)

// OTPValidator provides methods to check if 2FA is enabled and if OTP is valid
type OTPValidator interface {
	// IsEnabled returns true if 2FA for the given user is enabled
	IsEnabled(ctx context.Context, username string) (bool, error)
	// Validate checks that the OTP for the given user is valid, if not it
	// returns an error
	Validate(ctx context.Context, username, otp string) error
}

// JWTValidator validates JWT tokens.
type JWTValidator struct {
	db           sqlx.Ext
	userStore    Store
	secret       interface{}
	algorithm    jwa.SignatureAlgorithm
	otpValidator OTPValidator
}

// NewJWTValidator creates a new JWTValidator.
func NewJWTValidator(db sqlx.Ext, algorithm jwa.SignatureAlgorithm, secret interface{}, otpValidator OTPValidator, userStore Store) *JWTValidator {
	return &JWTValidator{
		db:           db,
		secret:       secret,
		algorithm:    algorithm,
		otpValidator: otpValidator,
		userStore:    userStore,
	}
}

// SignToken creates and signs a new JWT token for user
func (v JWTValidator) SignToken(username string, ttl int64, audience []string) (string, error) {
	t := jwt.New()
	if ttl == 0 {
		ttl = defaultSessionTTL
	}
	t.Set(jwt.IssuerKey, "lora-app-server")
	if len(audience) == 0 {
		t.Set(jwt.AudienceKey, "lora-app-server")
	} else {
		t.Set(jwt.AudienceKey, audience)
	}
	t.Set(jwt.IssuedAtKey, time.Now())
	t.Set(jwt.ExpirationKey, time.Now().Add(time.Duration(ttl)*time.Second))
	t.Set("username", username)
	token, err := jwt.Sign(t, v.algorithm, v.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %v", err)
	}
	return string(token), nil
}

// Validate validates the token from the given context against the given
// validator funcs.
func (v JWTValidator) Validate(ctx context.Context, funcs ...ValidatorFunc) error {
	claims, err := v.getClaims(ctx, "")
	if err != nil {
		return err
	}

	for _, f := range funcs {
		ok, err := f(v.db, claims)
		if err != nil {
			return fmt.Errorf("validation has failed: %v", err)
		}
		if ok {
			return nil
		}
	}

	return ErrNotAuthorized
}

// GetUsername returns the username of the authenticated user.
func (v JWTValidator) GetUsername(ctx context.Context) (string, error) {
	claims, err := v.getClaims(ctx, "")
	if err != nil {
		return "", err
	}

	return claims.Username, nil
}

// GetIsAdmin returns if the authenticated user is a global amin.
func (v JWTValidator) GetIsAdmin(ctx context.Context) (bool, error) {
	cred, err := v.GetCredentials(ctx)
	if err != nil {
		return false, err
	}
	if err = cred.IsGlobalAdmin(ctx); err != nil {
		return false, nil
	}
	return true, nil
}

// GetOTP returns OTP from the context
func (v JWTValidator) GetOTP(ctx context.Context) string {
	return getOTPFromContext(ctx)
}

// ValidateOTP validates OTP and returns the error if it is not valid
func (v JWTValidator) ValidateOTP(ctx context.Context) error {
	claims, err := v.getClaims(ctx, "")
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

func (v JWTValidator) GetCredentials(ctx context.Context, opts ...Option) (Credentials, error) {
	cfg := options{audience: "lora-app-server"}
	for _, o := range opts {
		o(&cfg)
	}
	claims, err := v.getClaims(ctx, cfg.audience)
	if err != nil {
		return nil, err
	}

	var cred Credentials
	if cfg.limitedCred {
		cred, err = GetLimitedCredentials(ctx, nil, claims.Username)
	} else {
		cred, err = GetCredentials(ctx, v.userStore, claims.Username)
	}
	if err != nil {
		return nil, err
	}

	if cfg.requireOTP {
		if claims.OTP == "" {
			return nil, fmt.Errorf("OTP is required")
		}
		if enabled, err := v.otpValidator.IsEnabled(ctx, claims.Username); !enabled || err != nil {
			return nil, fmt.Errorf("two-factor authentication is not enabled")
		}
		if err := v.otpValidator.Validate(ctx, claims.Username, claims.OTP); err != nil {
			return nil, fmt.Errorf("OTP is not valid")
		}
	}
	return cred, nil
}

func (v JWTValidator) getClaims(ctx context.Context, audience string) (*Claims, error) {
	tokenStr, err := getTokenFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("get token from context error: %v", err)
	}

	token, err := jwt.ParseVerify(strings.NewReader(tokenStr), v.algorithm, v.secret)
	if err != nil {
		return nil, err
	}
	if audience == "" {
		audience = "lora-app-server"
	}
	if err := jwt.Verify(token, jwt.WithAudience(audience)); err != nil {
		return nil, err
	}

	username, ok := token.Get("username")
	if !ok {
		return nil, fmt.Errorf("username is missing from the token")
	}
	usernameStr, ok := username.(string)
	if !ok {
		return nil, fmt.Errorf("username is not a string")
	}

	claims := &Claims{
		Username: usernameStr,
		OTP:      getOTPFromContext(ctx),
	}

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
