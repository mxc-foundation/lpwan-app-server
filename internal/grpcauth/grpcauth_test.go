package grpcauth

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	ljwt "github.com/lestrrat-go/jwx/jwt"
	"google.golang.org/grpc/metadata"

	"github.com/mxc-foundation/lpwan-app-server/internal/auth"
	"github.com/mxc-foundation/lpwan-app-server/internal/jwt"
)

var (
	testJWTKeyEnc = []byte("BlV5At5TU+LWXSEkiXZVvjuhWy6zBHJzA1jBvDbses4=")
)

type testOTPV struct{}

func (to testOTPV) Validate(ctx context.Context, username, otp string) error {
	if otp != "123456" {
		return fmt.Errorf("invalid OTP")
	}
	return nil
}

type testStore struct{}

func (ts testStore) AuthGetUser(ctx context.Context, username string) (auth.User, error) {
	if username != "alice@example.com" {
		return auth.User{}, fmt.Errorf("not found")
	}
	return auth.User{ID: 17, Email: "alice@example.com"}, nil
}

func (ts testStore) AuthGetOrgUser(ctx context.Context, userID, orgID int64) (auth.OrgUser, error) {
	var ou auth.OrgUser
	if userID == 17 && orgID == 3 {
		ou.IsOrgUser = true
		ou.IsOrgAdmin = true
	}
	return ou, nil
}

func TestAuthenticator(t *testing.T) {
	jwtv := jwt.NewValidator(jwa.HS256, testJWTKeyEnc, 86400)
	aliceTok, err := jwtv.SignToken(jwt.Claims{UserID: 17, Username: "alice@example.com", Service: auth.EMAIL}, 0, []string{"lora-app-server"})
	if err != nil {
		t.Fatal(err)
	}
	tt, err := ljwt.ParseVerify(strings.NewReader(aliceTok), jwa.HS256, testJWTKeyEnc)
	if err != nil {
		t.Fatal(err)
	}
	tt.Set(ljwt.IssuedAtKey, time.Now().Add(-72*time.Hour))
	tt.Set(ljwt.ExpirationKey, time.Now().Add(-48*time.Hour))
	aliceExp, err := ljwt.Sign(tt, jwa.HS256, testJWTKeyEnc)
	if err != nil {
		t.Fatal(err)
	}
	bobTok, err := jwtv.SignToken(jwt.Claims{UserID: 19, Username: "bob@example.com", Service: auth.EMAIL}, 0, []string{"registration"})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		token  string
		opts   *auth.Options
		otp    string
		errExp string
		creds  auth.Credentials
	}{
		{
			name:   "expired token should be rejected",
			token:  string(aliceExp),
			errExp: "exp",
		},
		{
			name:  "normal token, default options",
			token: aliceTok,
			creds: auth.Credentials{
				UserID:     17,
				Username:   "alice@example.com",
				IsExisting: true,
				Service:    auth.EMAIL,
			},
		},
		{
			name:  "normal token, with orgID option",
			token: aliceTok,
			opts:  auth.NewOptions().WithOrgID(3),
			creds: auth.Credentials{
				UserID:         17,
				Username:       "alice@example.com",
				IsExisting:     true,
				OrgID:          3,
				IsOrgUser:      true,
				IsOrgAdmin:     true,
				IsGatewayAdmin: true,
				IsDeviceAdmin:  true,
				Service:        auth.EMAIL,
			},
		},
		{
			name:   "normal token, require OTP, with no OTP",
			token:  aliceTok,
			opts:   auth.NewOptions().WithRequireOTP(),
			errExp: "OTP",
		},
		{
			name:   "normal token, require OTP, with invalid OTP",
			token:  aliceTok,
			opts:   auth.NewOptions().WithRequireOTP(),
			otp:    "111111",
			errExp: "OTP",
		},
		{
			name:  "normal token, require OTP, with valid OTP",
			token: aliceTok,
			opts:  auth.NewOptions().WithRequireOTP(),
			otp:   "123456",
			creds: auth.Credentials{
				UserID:     17,
				Username:   "alice@example.com",
				IsExisting: true,
				Service:    auth.EMAIL,
			},
		},
		{
			name:   "non-existing user, default options",
			token:  bobTok,
			errExp: "invalid token",
			creds: auth.Credentials{
				Service: auth.EMAIL,
			},
		},
		{
			name:   "non-existing user, with correct audience option",
			token:  bobTok,
			opts:   auth.NewOptions().WithAudience("registration"),
			errExp: "user validation",
			creds: auth.Credentials{
				Service: auth.EMAIL,
			},
		},
		{
			name:  "non-existing user, with audience and allow non-existent options",
			token: bobTok,
			opts:  auth.NewOptions().WithAudience("registration").WithAllowNonExisting(),
			creds: auth.Credentials{
				Username: "bob@example.com",
			},
		},
	}

	ga := New(testStore{}, jwtv, testOTPV{})
	for _, tc := range tests {
		t.Logf("test: %s", tc.name)
		opts := tc.opts
		if opts == nil {
			opts = auth.NewOptions()
		}
		md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", tc.token))
		if tc.otp != "" {
			md.Append("x-otp", tc.otp)
		}
		ctx := metadata.NewIncomingContext(context.Background(), md)
		cred, err := ga.GetCredentials(ctx, opts)
		if err != nil {
			if tc.errExp == "" || !strings.Contains(err.Error(), tc.errExp) {
				t.Errorf("unexpected error: %v", err)
			}
			continue
		}
		if *cred != tc.creds {
			t.Errorf("expected %#v, but got %#v", tc.creds, *cred)
		}
	}
}
