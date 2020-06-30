package auth

import (
	"fmt"
	"strings"
	"testing"

	"github.com/lestrrat-go/jwx/jwa"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

var (
	testJWTKeyEnc     = []byte("BlV5At5TU+LWXSEkiXZVvjuhWy6zBHJzA1jBvDbses4=")
	testExpiredToken  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJsb3JhLWFwcC1zZXJ2ZXIiLCJleHAiOjE1OTExMDUzMTcsImlzcyI6ImxvcmEtYXBwLXNlcnZlciIsIm5iZiI6MTU5MTEwMTcxNywic3ViIjoidXNlciIsInVzZXJuYW1lIjoiYWxpY2VAZXhhbXBsZS5jb20ifQ.A9-adLEdBHMQvc_5XcuOk_Xg_YJkWUUnnx20lvwAJzQ" // nolint: gosec
	testNoExpireToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJsb3JhLWFwcC1zZXJ2ZXIiLCJpc3MiOiJsb3JhLWFwcC1zZXJ2ZXIiLCJuYmYiOjE1OTExMDE2NDgsInN1YiI6InVzZXIiLCJ1c2VybmFtZSI6ImFsaWNlQGV4YW1wbGUuY29tIn0.pFsgDyepoi0hAbxk-mgCOKk_BQtXHbKZyP5isb9gV_M"                        // nolint: gosec
)

type testOTPValidator struct {
	users map[string]string
}

func (tov *testOTPValidator) IsEnabled(ctx context.Context, username string) (bool, error) {
	if _, ok := tov.users[username]; ok {
		return true, nil
	}
	return false, nil
}

func (tov *testOTPValidator) Validate(ctx context.Context, username, otp string) error {
	eotp, ok := tov.users[username]
	if !ok {
		return fmt.Errorf("otp is not enabled for %s", username)
	}
	if otp != eotp {
		return fmt.Errorf("otp mismatch")
	}
	return nil
}

func TestGetCredentials(t *testing.T) {
	ts := &testStore{
		users: map[string]testUser{
			"alice@example.com": {
				id:      1001,
				isAdmin: true,
				orgs: map[int64]OrgUser{
					1001: {IsOrgAdmin: true},
				},
			},
			"bob@example.com": {
				id: 1002,
				orgs: map[int64]OrgUser{
					1002: {IsOrgAdmin: true},
					2000: {},
				},
			},
		},
	}
	ts.init()

	tov := &testOTPValidator{
		users: map[string]string{"bob@example.com": "123456"},
	}
	v := NewJWTValidator(nil, jwa.HS256, testJWTKeyEnc, tov, ts)

	bobToken, err := v.SignToken("bob@example.com", 3600, nil)
	if err != nil {
		t.Fatal(err)
	}

	v1 := NewJWTValidator(nil, jwa.HS512, testJWTKeyEnc, tov, ts)
	bobWrongAlgo, err := v1.SignToken("bob@example.com", 3600, nil)
	if err != nil {
		t.Fatal(err)
	}

	bobTestAudience, err := v.SignToken("bob@example.com", 3600, []string{"test", "foo"})
	if err != nil {
		t.Fatal(err)
	}

	malloryToken, err := v.SignToken("mallory@example.com", 3600, nil)
	if err != nil {
		t.Fatal(err)
	}

	expectError := func(s string) func(Credentials, error) error {
		return func(c Credentials, e error) error {
			if e == nil || !strings.Contains(e.Error(), s) {
				return fmt.Errorf("error is not as expected: %v", e)
			}
			return nil
		}
	}
	expectSuccess := func(s string) func(Credentials, error) error {
		return func(c Credentials, e error) error {
			if e != nil {
				return e
			}
			if u := c.Username(); u != s {
				return fmt.Errorf("unexpected username %s", u)
			}
			return nil
		}
	}

	ctx := context.Background()
	testCases := []struct {
		name    string
		token   string
		options []Option
		otp     string
		check   func(cred Credentials, e error) error
	}{
		{
			name:  "expired token",
			token: testExpiredToken,
			check: expectError(""),
		},
		{
			name:  "valid token without expiration",
			token: testNoExpireToken,
			check: func(cred Credentials, e error) error {
				if e != nil {
					return fmt.Errorf("token has been rejected: %v", e)
				}
				if u := cred.Username(); u != "alice@example.com" {
					return fmt.Errorf("unexpected username: %s", u)
				}
				if err := cred.IsGlobalAdmin(ctx); err != nil {
					return fmt.Errorf("expected alice to be a global admin: %v", err)
				}
				return nil
			},
		},
		{
			name:    "valid token, but OTP required and is not enabled",
			token:   testNoExpireToken,
			options: []Option{WithValidOTP()},
			otp:     "111111",
			check:   expectError("not enabled"),
		},
		{
			name:    "valid token, but audience mismatch",
			token:   testNoExpireToken,
			options: []Option{WithAudience("test")},
			check:   expectError(""),
		},
		{
			name:  "token signed with different algo",
			token: bobWrongAlgo,
			check: expectError(""),
		},
		{
			name:    "invalid OTP",
			token:   bobToken,
			otp:     "111111",
			options: []Option{WithValidOTP()},
			check:   expectError("not valid"),
		},
		{
			name:    "missing OTP",
			token:   bobToken,
			options: []Option{WithValidOTP()},
			check:   expectError("is required"),
		},
		{
			name:    "valid OTP",
			token:   bobToken,
			otp:     "123456",
			options: []Option{WithValidOTP()},
			check:   expectSuccess("bob@example.com"),
		},
		{
			name:    "valid OTP but not audience",
			token:   bobToken,
			otp:     "123456",
			options: []Option{WithValidOTP(), WithAudience("test")},
			check:   expectError(""),
		},
		{
			name:    "valid OTP and audience",
			token:   bobTestAudience,
			otp:     "123456",
			options: []Option{WithAudience("test"), WithValidOTP()},
			check:   expectSuccess("bob@example.com"),
		},
		{
			name:  "default audience mismatch",
			token: bobTestAudience,
			check: expectError(""),
		},
		{
			name:  "non-existing user",
			token: malloryToken,
			check: expectError(""),
		},
	}

	for _, tc := range testCases {
		t.Log(tc.name)
		md := make(metadata.MD)
		md.Set("authorization", "Bearer "+tc.token)
		if tc.otp != "" {
			md.Set("x-otp", tc.otp)
		}
		ctx := metadata.NewIncomingContext(context.Background(), md)
		cred, e := v.GetCredentials(ctx, tc.options...)
		if err := tc.check(cred, e); err != nil {
			t.Error(err)
		}
	}
}
