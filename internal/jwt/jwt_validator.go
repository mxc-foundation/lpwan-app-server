package jwt

import (
	"fmt"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

// defaultSessionTTL defines the default token ttl if not specified in config
const defaultSessionTTL = 86400

// Claims defines the struct containing the token claims.
type Claims struct {
	// Username defines the identity of the user.
	Username string `json:"username"`
}

// Validator validates JWT tokens.
type Validator struct {
	secret     interface{}
	defaultTTL int64
	algorithm  jwa.SignatureAlgorithm
}

// NewValidator creates a new jwt.Validator.
func NewValidator(algorithm jwa.SignatureAlgorithm, secret interface{}, defaultTTL int64) *Validator {
	if defaultTTL == 0 {
		defaultTTL = defaultSessionTTL
	}
	return &Validator{
		secret:     secret,
		algorithm:  algorithm,
		defaultTTL: defaultTTL,
	}
}

// SignToken creates and signs a new JWT token for user
func (v Validator) SignToken(username string, ttl int64, audience []string) (string, error) {
	t := jwt.New()
	if ttl == 0 {
		ttl = v.defaultTTL
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

func (v Validator) GetClaims(tokenEncoded, audience string) (*Claims, error) {
	token, err := jwt.ParseVerify(strings.NewReader(tokenEncoded), v.algorithm, v.secret)
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
