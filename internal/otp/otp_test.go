package otp

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
)

type testStore struct {
	m       map[string]*TOTPInfo
	codeCnt int64
}

func newTestStore() *testStore {
	return &testStore{
		m: make(map[string]*TOTPInfo),
	}
}

func (tos *testStore) GetTOTPInfo(ctx context.Context, username string) (TOTPInfo, error) {
	var r TOTPInfo
	if u, ok := tos.m[username]; ok {
		r = *u
		rc := make(map[int64]string)
		for k, v := range u.RecoveryCodes {
			rc[k] = v
		}
		r.RecoveryCodes = rc
	}
	return r, nil
}

func (tos *testStore) StoreNewSecret(ctx context.Context, username, secret string) error {
	var r TOTPInfo
	r.Secret = secret
	tos.m[username] = &r
	return nil
}

func (tos *testStore) Enable(ctx context.Context, username string) error {
	u, ok := tos.m[username]
	if !ok {
		return fmt.Errorf("user does not have configuration")
	}
	if u.Enabled {
		return fmt.Errorf("already enabled")
	}
	tos.m[username].Enabled = true
	return nil
}

func (tos *testStore) Delete(ctx context.Context, username string) error {
	delete(tos.m, username)
	return nil
}

func (tos *testStore) DeleteRecoveryCode(ctx context.Context, username string, codeID int64) error {
	u, ok := tos.m[username]
	if !ok {
		return fmt.Errorf("user %s doesn't exist", username)
	}
	_, ok = u.RecoveryCodes[codeID]
	if !ok {
		return fmt.Errorf("key %d doesn't exist for user %s", codeID, username)
	}
	delete(u.RecoveryCodes, codeID)
	return nil
}

func (tos *testStore) AddRecoveryCodes(ctx context.Context, username string, codes []string) error {
	u, ok := tos.m[username]
	if !ok {
		return fmt.Errorf("user %s doesn't exist", username)
	}
	if u.RecoveryCodes == nil {
		u.RecoveryCodes = make(map[int64]string)
	}
	for _, c := range codes {
		n := atomic.AddInt64(&tos.codeCnt, 1)
		u.RecoveryCodes[n] = c
	}
	return nil
}

func (tos *testStore) UpdateLastTimeSlot(ctx context.Context, username string, previousValue, newValue int64) error {
	u, ok := tos.m[username]
	if !ok {
		return fmt.Errorf("user %s doesn't exist", username)
	}
	if u.LastTimeSlot != previousValue {
		return fmt.Errorf("user %s has been changed", username)
	}
	u.LastTimeSlot = newValue
	return nil
}

var testKey = "ee1227f349578e44c8e362e18ee2d05e"

func TestOTPConfigure(t *testing.T) {
	ctx := context.Background()
	ts := newTestStore()
	if ti, err := ts.GetTOTPInfo(ctx, "alice"); err != nil || ti.Enabled || ti.Secret != "" || len(ti.RecoveryCodes) > 0 {
		t.Fatal("expected empty config for alice")
	}

	tv, err := NewValidator("test", testKey, ts, authcus.NewJWTValidator(storage.DB(), "HS256", config.C.ApplicationServer.ExternalAPI.JWTSecret))
	if err != nil {
		t.Fatal(err)
	}
	conf, err := tv.NewConfiguration(ctx, "alice")
	if err != nil {
		t.Fatal(err)
	}

	if ok, err := tv.IsEnabled(ctx, "alice"); err != nil || ok {
		t.Fatal("expected TOTP for user to be disabled")
	}
	if ti, err := ts.GetTOTPInfo(ctx, "alice"); err != nil || ti.Enabled || ti.Secret == "" || len(ti.RecoveryCodes) != 10 {
		t.Fatalf("user config after configuration is not as expected %#v %v", ti, err)
	}

	key, err := otp.NewKeyFromURL(conf.URL)
	if err != nil {
		t.Fatal(err)
	}
	if key.Secret() != conf.Secret {
		t.Fatal("secret in url doesn't match plain text secret")
	}
	code, err := totp.GenerateCodeCustom(key.Secret(), time.Now(), totpOptions)
	if err != nil {
		t.Fatal(err)
	}

	if err := tv.Enable(ctx, "alice", "000000"); err == nil {
		t.Fatal("enabled user TOTP using code 000000, bug or you've got very lucky")
	}
	if err := tv.Enable(ctx, "alice", code); err != nil {
		t.Fatalf("couldn't enable user using generated code: %v", err)
	}
	if err := tv.Validate(ctx, "alice", code); err == nil {
		t.Fatal("accepted the same code second time")
	}
	if ok, err := tv.IsEnabled(ctx, "alice"); !ok || err != nil {
		t.Fatal("totp hasn't been enabled")
	}

	if err := tv.Disable(ctx, "alice"); err != nil {
		t.Fatalf("couldn't disable TOTP for alice: %v", err)
	}
	if ok, err := tv.IsEnabled(ctx, "alice"); ok || err != nil {
		t.Fatal("totp is still enabled for alice")
	}
}

func TestTOTP(t *testing.T) {
	ctx := context.Background()
	ts := newTestStore()
	tv, err := NewValidator("test", testKey, ts, authcus.NewJWTValidator(storage.DB(), "HS256", config.C.ApplicationServer.ExternalAPI.JWTSecret))
	if err != nil {
		t.Fatal(err)
	}
	secrets := make(map[string]string)

	enable := func(username string) {
		conf, err := tv.NewConfiguration(ctx, username)
		if err != nil {
			t.Fatal(err)
		}
		if err := ts.Enable(ctx, username); err != nil {
			t.Fatalf("couldn't enable totp for %s: %v", username, err)
		}
		secrets[username] = conf.Secret
	}

	getCode := func(username string, seconds int64) string {
		secret, ok := secrets[username]
		if !ok {
			t.Fatalf("user %s hasn't been configured", username)
		}
		code, err := totp.GenerateCodeCustom(secret, time.Now().Add(time.Duration(seconds)*time.Second), totpOptions)
		if err != nil {
			t.Fatalf("couldn't generate code for %s: %v", username, err)
		}
		return code
	}

	testCases := []struct {
		username    string
		description string
		tries       []int64
		results     []string
	}{
		{
			username:    "alice",
			description: "accept codes at least to a minute back, but not from the future",
			tries:       []int64{-60, -30, 40},
			results:     []string{"", "", errOTPNotValid},
		},
		{
			username:    "bob",
			description: "do not accept the same code twice",
			tries:       []int64{-10, -10},
			results:     []string{"", errOTPNotValid},
		},
		{
			username:    "carol",
			description: "do not accept codes for earlier intervals than the last one",
			tries:       []int64{0, -40},
			results:     []string{"", errOTPNotValid},
		},
		{
			username:    "dan",
			description: "after two failures we still accept the valid code",
			tries:       []int64{-100, -100, 0},
			results:     []string{errOTPNotValid, errOTPNotValid, ""},
		},
		{
			username:    "mallory",
			description: "after 3 unsuccessful attempts we lock user out",
			tries:       []int64{-100, -100, -100, 0, 0},
			results:     []string{errOTPNotValid, errOTPNotValid, errOTPNotValid, errOTPNotValid, errOTPLockedOut},
		},
	}

	for _, tc := range testCases {
		t.Logf("%s: %s", tc.username, tc.description)
		enable(tc.username)
		for i, sec := range tc.tries {
			code := getCode(tc.username, sec)
			err := tv.Validate(ctx, tc.username, code)
			if (tc.results[i] == "" && err != nil) ||
				(tc.results[i] != "" && !strings.Contains(err.Error(), tc.results[i])) {

				t.Errorf("try %d, expected error '%s', got: %s", i, tc.results[i], err.Error())
			}
		}
	}
}

func TestRecoveryCodes(t *testing.T) {
	ctx := context.Background()
	ts := newTestStore()
	tv, err := NewValidator("test", testKey, ts, authcus.NewJWTValidator(storage.DB(), "HS256", config.C.ApplicationServer.ExternalAPI.JWTSecret))
	if err != nil {
		t.Fatal(err)
	}
	secrets := make(map[string][]string)

	enable := func(username string) {
		conf, err := tv.NewConfiguration(ctx, username)
		if err != nil {
			t.Fatal(err)
		}
		if err := ts.Enable(ctx, username); err != nil {
			t.Fatalf("couldn't enable totp for %s: %v", username, err)
		}
		secrets[username] = conf.RecoveryCodes
	}

	getRecoveryCodes := func(username string) []string {
		codes, err := tv.GetRecoveryCodes(ctx, username, false)
		if err != nil {
			t.Fatalf("couldn't get recovery codes for %s: %v", username, err)
		}
		return codes
	}

	enable("alice")
	codes := getRecoveryCodes("alice")
	if len(codes) != 10 {
		t.Fatalf("expected to get 10 recovery codes, but got %d: %#v", len(codes), codes)
	}
	codeMap := make(map[string]int)
	for _, c := range codes {
		codeMap[c] = 0
	}
	if len(codeMap) != 10 {
		t.Fatalf("expected to get 10 different recovery codes, but got %d: %#v", len(codeMap), codes)
	}

	code := codes[0]
	delete(codeMap, code)
	if err := tv.Validate(ctx, "alice", code); err != nil {
		t.Fatalf("recovery code %s wasn't accepted: %v", code, err)
	}
	re := regexp.MustCompile("^[0-9a-f]{5}-[0-9a-f]{5}$")
	if !re.MatchString(code) {
		t.Fatalf("recovery code format is wrong: %s", code)
	}

	t.Log("recovery code can be used only once, and then is deleted")
	t.Log("if user downloads recovery codes, they should get 9 old codes and one new")
	codes = getRecoveryCodes("alice")
	for _, c := range codes {
		if c == code {
			t.Fatal("the used recovery code hasn't been removed")
		}
		codeMap[c]++
		if codeMap[c] > 1 {
			t.Fatalf("got %s more than once: %#v", c, codes)
		}
	}
	if len(codeMap) != 10 {
		t.Fatalf("expected 10 different codes, 9 old ones and 1 new, got: %#v", codeMap)
	}

	t.Log("after we used a recovery code we can't immediately use another one")
	if err := tv.Validate(ctx, "alice", codes[0]); err == nil {
		t.Fatal("was able to validate the user with another recovery code immediately")
	}

	t.Log("we should allow three attempts to enter recovery code")
	enable("bob")
	codes = getRecoveryCodes("bob")
	codeMap = make(map[string]int)
	for _, c := range codes {
		codeMap[c]++
	}
	code = codes[0]
	var randCodes []string
	for i := 0; i < 3; i++ {
		wCode, err := generateRecoveryCode()
		if err != nil {
			t.Fatal(err)
		}
		randCodes = append(randCodes, wCode)
	}
	for i := 0; i < 2; i++ {
		if err := tv.Validate(ctx, "bob", randCodes[i]); err == nil {
			t.Fatalf("validated random recovery code")
		}
	}
	if err := tv.Validate(ctx, "bob", code); err != nil {
		t.Fatalf("validation failed: %v", err)
	}

	t.Log("if user regenerated recovery codes, the old ones all should be replaced")
	codes, err = tv.GetRecoveryCodes(ctx, "bob", true)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range codes {
		codeMap[c]++
	}
	if len(codeMap) != 20 {
		t.Fatalf("expected to get 10 new recovery codes: %#v", codeMap)
	}

	t.Log("after three unsuccessful attemps user should be locked out")
	enable("mallory")
	code = getRecoveryCodes("mallory")[0]
	for i := 0; i < 3; i++ {
		if err := tv.Validate(ctx, "mallory", randCodes[i]); err == nil {
			t.Fatalf("validated random recovery code")
		}
	}
	if err := tv.Validate(ctx, "mallory", code); err == nil {
		t.Fatal("validation succeeded")
	}
}
