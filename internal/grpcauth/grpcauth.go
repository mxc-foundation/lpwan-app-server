// Package grpcauth implements auth.Authenticator for grpc protocol
package grpcauth

import (
	"context"
	"fmt"
	"regexp"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"

	"github.com/mxc-foundation/lpwan-app-server/internal/auth"
	"github.com/mxc-foundation/lpwan-app-server/internal/jwt"
)

// OTPValidator validates one time passwords
type OTPValidator interface {
	Validate(ctx context.Context, username, otp string) error
}

type grpcAuth struct {
	store auth.Store
	jwtv  *jwt.Validator
	otpv  OTPValidator
}

// New creates and returns a new authenticator
func New(store auth.Store, jwtv *jwt.Validator, otpv OTPValidator) auth.Authenticator {
	return &grpcAuth{
		store: store,
		jwtv:  jwtv,
		otpv:  otpv,
	}
}

func (ga *grpcAuth) GetCredentials(ctx context.Context, opts *auth.Options) (*auth.Credentials, error) {
	token, err := getTokenFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("couldn't find JWT token: %v", err)
	}
	claims, err := ga.jwtv.GetClaims(token, opts.Audience)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}
	// For non-existing user we only return username, there's no point in
	// checking the OTP or anything else
	if opts.AllowNonExisting {
		return &auth.Credentials{
			Username: claims.Username,
		}, nil
	}

	if opts.RequireOTP {
		otp := GetOTPFromContext(ctx)
		if err := ga.otpv.Validate(ctx, claims.Username, otp); err != nil {
			return nil, err
		}
	}

	creds, err := auth.NewCredentials(ctx, ga.store, claims.Username, opts.OrgID)
	if err != nil {
		return nil, fmt.Errorf("user validation has failed: %v", err)
	}
	return creds, nil
}

var validAuthorizationRegexp = regexp.MustCompile(`(?i)^bearer (.*)$`)

func getTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("no metadata in context")
	}

	token, ok := md["authorization"]
	if !ok || len(token) == 0 {
		return "", fmt.Errorf("missing authorization header")
	}

	match := validAuthorizationRegexp.FindStringSubmatch(token[0])

	// authorization header should respect RFC1945
	if len(match) == 0 {
		l := len(token)
		if l > 16 {
			l = 16
		}
		logrus.Warnf("Deprecated Authorization header: %s", token[0:l])
		return token[0], nil
	}

	return match[1], nil
}

// GetOTPFromContext extracts OTP from the context
func GetOTPFromContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	if len(md["x-otp"]) == 1 {
		return md["x-otp"][0]
	}
	return ""
}
