package jwt

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

var (
	testJWTKeyEnc     = []byte("BlV5At5TU+LWXSEkiXZVvjuhWy6zBHJzA1jBvDbses4=")
	testExpiredToken  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJsb3JhLWFwcC1zZXJ2ZXIiLCJleHAiOjE1OTExMDUzMTcsImlzcyI6ImxvcmEtYXBwLXNlcnZlciIsIm5iZiI6MTU5MTEwMTcxNywic3ViIjoidXNlciIsInVzZXJuYW1lIjoiYWxpY2VAZXhhbXBsZS5jb20ifQ.A9-adLEdBHMQvc_5XcuOk_Xg_YJkWUUnnx20lvwAJzQ" // nolint: gosec
	testNoExpireToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsibG9yYS1hcHAtc2VydmVyIl0sImlhdCI6MTYwNDQxMTIyMiwidXNlcklkIjo3NywidXNlcm5hbWUiOiJib2JAZXhhbXBsZS5jb20ifQ.6OzYrCKmSNF4qbxfp3q7xyJD70TViUYjggcQg5YR-WM"                                                         // nolint: gosec
)

func TestValidator(t *testing.T) {

	v := NewValidator(jwa.HS256, testJWTKeyEnc, 86400)

	defExpiry, err := v.SignToken(42, "carol@example.com", 0, nil)
	if err != nil {
		t.Fatal(err)
	}

	tok, err := jwt.Parse(strings.NewReader(defExpiry))
	if err != nil {
		t.Fatal(err)
	}
	if tok.Expiration().After(time.Now().Add(24*time.Hour)) ||
		tok.Expiration().Before(time.Now().Add(23*time.Hour)) {
		t.Fatalf("expected the token to expire in 1 day, but it expires in %s",
			time.Until(tok.Expiration()).String())
	}

	bobToken, err := v.SignToken(77, "bob@example.com", 3600, nil)
	if err != nil {
		t.Fatal(err)
	}

	tok, err = jwt.Parse(strings.NewReader(bobToken))
	if err != nil {
		t.Fatal(err)
	}
	if tok.Expiration().After(time.Now().Add(time.Hour)) ||
		tok.Expiration().Before(time.Now().Add(50*time.Minute)) {
		t.Fatalf("expected the token to expire in 1 hour, but it expires in %s",
			time.Until(tok.Expiration()).String())
	}

	v1 := NewValidator(jwa.HS512, testJWTKeyEnc, 86400)
	bobWrongAlgo, err := v1.SignToken(77, "bob@example.com", 3600, nil)
	if err != nil {
		t.Fatal(err)
	}

	bobTestAudience, err := v.SignToken(77, "bob@example.com", 3600, []string{"test", "foo"})
	if err != nil {
		t.Fatal(err)
	}

	expectError := func(s string) func(*Claims, error) error {
		return func(c *Claims, e error) error {
			if e == nil || !strings.Contains(e.Error(), s) {
				return fmt.Errorf("error is not as expected: %v", e)
			}
			return nil
		}
	}
	expectSuccess := func(id int64, s string) func(*Claims, error) error {
		return func(c *Claims, e error) error {
			if e != nil {
				return e
			}
			if c.Username != s {
				return fmt.Errorf("unexpected username %s", c.Username)
			}
			if c.UserID != id {
				return fmt.Errorf("unexpected userId: %d", c.UserID)
			}
			return nil
		}
	}

	testCases := []struct {
		name     string
		token    string
		audience string
		check    func(c *Claims, e error) error
	}{
		{
			name:  "expired token",
			token: testExpiredToken,
			check: expectError(""),
		},
		{
			name:  "valid token",
			token: bobToken,
			check: expectSuccess(77, "bob@example.com"),
		},
		{
			name:  "valid token without expiration",
			token: testNoExpireToken,
			check: expectSuccess(77, "bob@example.com"),
		},
		{
			name:     "valid token, but audience mismatch",
			token:    testNoExpireToken,
			audience: "test",
			check:    expectError(""),
		},
		{
			name:  "token signed with different algo",
			token: bobWrongAlgo,
			check: expectError(""),
		},
		{
			name:  "default audience mismatch",
			token: bobTestAudience,
			check: expectError(""),
		},
	}

	for _, tc := range testCases {
		t.Log(tc.name)
		c, e := v.GetClaims(tc.token, tc.audience)
		if err := tc.check(c, e); err != nil {
			t.Error(err)
		}
	}
}
