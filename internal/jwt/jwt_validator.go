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
	// UserID
	UserID int64 `json:"userId"`
	// Username defines the identity of the user.
	Username string `json:"username"`
	// Service
	Service string `json:"service"`
	// ExternalCred defines key credentials to verify a wechat user
	ExternalCred string `json:"externalCred"`
	// OrganizationID is used when organization id is required for signing JWT and with audience "mosquitto-auth"
	OrganizationID int64 `json:"organizationId"`
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
func (v Validator) SignToken(claims Claims, ttl int64, audience []string) (string, error) {
	t := jwt.New()
	if ttl == 0 {
		ttl = v.defaultTTL
	}
	_ = t.Set(jwt.IssuerKey, "lora-app-server")
	if len(audience) == 0 {
		_ = t.Set(jwt.AudienceKey, "lora-app-server")
	} else {
		_ = t.Set(jwt.AudienceKey, audience)
	}
	_ = t.Set(jwt.IssuedAtKey, time.Now())
	_ = t.Set(jwt.ExpirationKey, time.Now().Add(time.Duration(ttl)*time.Second))
	_ = t.Set("username", claims.Username)
	_ = t.Set("userId", claims.UserID)
	_ = t.Set("service", claims.Service)
	_ = t.Set("externalCred", claims.ExternalCred)
	_ = t.Set("organizationId", claims.OrganizationID)

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

	claims := &Claims{}

	// use of a claim name is optional depending on context
	username, ok := token.Get("username")
	if ok {
		usernameStr, ok := username.(string)
		if !ok {
			return nil, fmt.Errorf("username is not a string")
		}
		claims.Username = usernameStr
	}

	userID, ok := token.Get("userId")
	if ok {
		userIDFloat, ok := userID.(float64)
		if !ok {
			return nil, fmt.Errorf("userId is not a number")
		}
		claims.UserID = int64(userIDFloat)
	}

	serviceName, ok := token.Get("service")
	if ok {
		serviceNameStr, ok := serviceName.(string)
		if !ok {
			return nil, fmt.Errorf("serviceName is not a string")
		}
		claims.Service = serviceNameStr
	}

	externalCred, ok := token.Get("externalCred")
	if ok {
		externalCredStr, ok := externalCred.(string)
		if !ok {
			return nil, fmt.Errorf("externalCred is not a string")
		}
		claims.ExternalCred = externalCredStr
	}

	organizationID, ok := token.Get("organizationId")
	if ok {
		organizationIDFloat, ok := organizationID.(float64)
		if !ok {
			return nil, fmt.Errorf("organizationId is not a number")
		}
		claims.OrganizationID = int64(organizationIDFloat)
	}

	return claims, nil
}
