package jwt

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"

	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
)

// defaultSessionTTL defines the default session TTL
const defaultSessionTTL = 86400

var validAuthorizationRegexp = regexp.MustCompile(`(?i)^bearer (.*)$`)

// Claims defines the struct containing the token claims.
type Claims struct {
	// Username defines the identity of the user.
	Username string `json:"username"`
	// UserID defines the ID of th user.
	UserID int64 `json:"user_id"`

	// APIKeyID defines the API key ID.
	APIKeyID uuid.UUID `json:"api_key_id"`
}

// JWTValidator validates JWT tokens.
type JWTValidator struct {
	secret    interface{}
	algorithm jwa.SignatureAlgorithm
}

// NewJWTValidator creates a new JWTValidator.
func NewJWTValidator(algorithm jwa.SignatureAlgorithm, secret interface{}) *JWTValidator {
	return &JWTValidator{
		secret:    secret,
		algorithm: algorithm,
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

// GetUsername returns the username of the authenticated user.
func (v JWTValidator) GetUsername(ctx context.Context) (string, error) {
	claims, err := v.GetClaims(ctx, "")
	if err != nil {
		return "", err
	}

	return claims.Username, nil
}

func (v JWTValidator) GetClaims(ctx context.Context, audience string) (*Claims, error) {
	tokenStr, err := getTokenFromContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get token from context error")
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

// GetSubject returns the claim subject.
func (v JWTValidator) GetSubject(ctx context.Context) (string, error) {
	return "user", nil
}

// GetUser returns the user object.
func (v JWTValidator) GetUser(ctx context.Context) (storage.User, error) {
	claims, err := v.GetClaims(ctx, "")
	if err != nil {
		return storage.User{}, err
	}

	if claims.UserID != 0 {
		return storage.GetUser(ctx, storage.DB().DB, claims.UserID)
	}

	if claims.Username != "" {
		return storage.GetUserByEmail(ctx, storage.DB().DB, claims.Username)
	}

	return storage.User{}, errors.New("no username or user_id in claims")
}

// GetAPIKeyID returns the API key ID.
func (v JWTValidator) GetAPIKeyID(ctx context.Context) (uuid.UUID, error) {
	claims, err := v.GetClaims(ctx, "")
	if err != nil {
		return uuid.Nil, err
	}

	return claims.APIKeyID, nil
}
