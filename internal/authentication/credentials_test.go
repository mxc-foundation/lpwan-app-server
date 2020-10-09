package authentication

import (
	"fmt"
	"log"
	"math/rand"
	"testing"

	"google.golang.org/grpc/metadata"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/context"

	"github.com/mxc-foundation/lpwan-app-server/internal/jwt"
	"github.com/mxc-foundation/lpwan-app-server/internal/otp"
)

type testUser struct {
	id            int64
	userEmail     string
	isGlobalAdmin bool
	orgs          map[int64]OrgUser
	jwt           string
	ctx           context.Context
}
type testStore struct {
	users     map[string]testUser
	usersByID map[int64]string
}

func (ts *testStore) GetUser(ctx context.Context, username string) (User, error) {
	u, ok := ts.users[username]
	if !ok {
		return User{}, fmt.Errorf("user not found")
	}
	return User{
		ID:            u.id,
		IsGlobalAdmin: u.isGlobalAdmin,
	}, nil
}

func (ts *testStore) GetOrgUser(ctx context.Context, userID int64, orgID int64) (OrgUser, error) {
	username := ts.usersByID[userID]
	if username == "" {
		return OrgUser{}, fmt.Errorf("user not found")
	}
	u, ok := ts.users[username]
	if !ok {
		panic("inconsistent storage!")
	}
	ou, ok := u.orgs[orgID]
	if !ok {
		return ou, fmt.Errorf("not in org")
	}
	return ou, nil
}

type testOTPStore struct {
	otpByUsername map[string]otp.TOTPInfo
}

func (t testOTPStore) GetTOTPInfo(ctx context.Context, username string) (otp.TOTPInfo, error) {
	return t.otpByUsername[username], nil
}

func (t testOTPStore) Enable(ctx context.Context, username string) error {
	res := t.otpByUsername[username]
	res.Enabled = true
	t.otpByUsername[username] = res
	return nil
}

func (t testOTPStore) Delete(ctx context.Context, username string) error {
	delete(t.otpByUsername, username)
	return nil
}

func (t testOTPStore) StoreNewSecret(ctx context.Context, username, secret string) error {
	res := t.otpByUsername[username]
	res.Secret = secret
	t.otpByUsername[username] = res
	return nil
}

func (t testOTPStore) DeleteRecoveryCode(ctx context.Context, username string, codeID int64) error {
	res := t.otpByUsername[username]
	delete(res.RecoveryCodes, codeID)
	return nil
}

func (t testOTPStore) AddRecoveryCodes(ctx context.Context, username string, codes []string) error {
	res := t.otpByUsername[username]
	if res.RecoveryCodes == nil {
		res.RecoveryCodes = make(map[int64]string)
	}
	for _, v := range codes {
		res.RecoveryCodes[rand.Int63()] = v
	}

	t.otpByUsername[username] = res
	return nil
}

func (t testOTPStore) UpdateLastTimeSlot(ctx context.Context, username string, previousValue, newValue int64) error {
	res := t.otpByUsername[username]
	res.LastTimeSlot = newValue
	t.otpByUsername[username] = res
	return nil
}

type testEnv struct {
	ctx  context.Context
	ts   testStore
	otp  testOTPStore
	cred *Credentials
}

func newTestEnv() *testEnv {
	return &testEnv{
		ctx: context.Background(),
		ts: testStore{
			users:     make(map[string]testUser),
			usersByID: make(map[int64]string),
		},
		otp: testOTPStore{
			otpByUsername: make(map[string]otp.TOTPInfo),
		},
	}
}

func (te *testEnv) newTestUser(user testUser, is2FA bool) (err error) {
	var audience []string
	ttl := int64(100000)
	te.ts.usersByID[user.id] = user.userEmail

	totpConfig, err := te.cred.NewConfiguration(te.ctx, user.userEmail)
	if err != nil {
		return err
	}
	totpInfo := te.otp.otpByUsername[user.userEmail]
	totpInfo.Secret = totpConfig.Secret
	totpInfo.Enabled = is2FA
	te.otp.otpByUsername[user.userEmail] = totpInfo

	if is2FA {
		ttl = 600
		audience = []string{"login-2fa"}
	}
	jwt, err := te.cred.SignJWToken(user.userEmail, ttl, audience)
	if err != nil {
		return err
	}
	ctx := metadata.NewIncomingContext(te.ctx, map[string][]string{
		"authorization": {fmt.Sprintf("bearer %s", jwt)},
	})
	te.ts.users[user.userEmail] = testUser{
		id:            user.id,
		userEmail:     user.userEmail,
		isGlobalAdmin: user.isGlobalAdmin,
		orgs:          user.orgs,
		jwt:           jwt,
		ctx:           ctx,
	}

	return nil
}

type keyType string

var aliceAdminDisable2FA = testUser{
	userEmail:     "aliceAdminDisable2FA",
	id:            1001,
	isGlobalAdmin: true,
	orgs:          map[int64]OrgUser{1001: {}},
	jwt:           "",
}
var bobAdminEnable2FA = testUser{
	userEmail:     "bobAdminEnable2FA",
	id:            1001,
	isGlobalAdmin: true,
	orgs: map[int64]OrgUser{
		1002: {IsOrgAdmin: true},
		2000: {IsGatewayAdmin: true},
		2001: {IsOrgAdmin: true, IsGatewayAdmin: true},
		2002: {IsDeviceAdmin: true},
	},
	jwt: "",
}

func TestCredentials(t *testing.T) {
	te := newTestEnv()

	jwtSecret := "ToQDKH22aAJ/raMSgKtjWsveS9GY9VdlnX/YSBq5JGg="
	jwtValidator := jwt.NewJWTValidator("HS256", []byte(jwtSecret))
	otpValidator, err := otp.NewValidator("lpwan-app-server", "3592961e33fb5225b6479b2593dc47be", &te.otp)
	if err != nil {
		log.Fatal()
	}

	SetupCred(&te.ts, jwtValidator, otpValidator)
	te.cred = NewCredentials()

	err = te.newTestUser(aliceAdminDisable2FA, false)
	if err != nil {
		t.Fatal("Failed to get new test user alice")
	}

	err = te.newTestUser(bobAdminEnable2FA, true)
	if err != nil {
		t.Fatal("Failed to get new test user bob")
	}

	Convey("Check Is2FAEnabled", t, func() {

		res, err := te.cred.Is2FAEnabled(te.ts.users[aliceAdminDisable2FA.userEmail].ctx, aliceAdminDisable2FA.userEmail)
		So(err, ShouldBeNil)
		So(res, ShouldBeFalse)

		res, err = te.cred.Is2FAEnabled(te.ts.users[bobAdminEnable2FA.userEmail].ctx, bobAdminEnable2FA.userEmail)
		So(err, ShouldBeNil)
		So(res, ShouldBeTrue)
	})

	Convey("Check GetUser", t, func() {
		aliceRaw := te.ts.users[aliceAdminDisable2FA.userEmail]
		alice, err := te.cred.GetUser(te.ts.users[aliceAdminDisable2FA.userEmail].ctx)
		So(err, ShouldBeNil)
		So(alice.Email, ShouldEqual, aliceRaw.userEmail)
		So(alice.IsGlobalAdmin, ShouldBeTrue)
		So(alice.ID, ShouldEqual, aliceRaw.id)

		bobRaw := te.ts.users[bobAdminEnable2FA.userEmail]
		bob, err := te.cred.GetUser(te.ts.users[bobAdminEnable2FA.userEmail].ctx)
		So(err, ShouldBeNil)
		So(bob.Email, ShouldEqual, bobRaw.userEmail)
		So(bob.IsGlobalAdmin, ShouldBeTrue)
		So(bob.ID, ShouldEqual, bobRaw.id)
	})

	Convey("Check GetUserPermissionWithOrgID", t, func() {

	})

	Convey("Check Username", t, func() {

	})

	Convey("Check UserID", t, func() {

	})

	Convey("Check IsGlobalAdmin", t, func() {

	})

	Convey("Check IsOrgUser", t, func() {

	})

	Convey("Check IsOrgAdmin", t, func() {

	})

	Convey("Check IsDeviceAdmin", t, func() {

	})

	Convey("Check IsGatewayAdmin", t, func() {

	})

	Convey("Check SignJWToken", t, func() {

	})

	Convey("Check GetOTP", t, func() {

	})

	Convey("Check EnableOTP", t, func() {

	})

	Convey("Check DisableOTP", t, func() {

	})

	Convey("Check OTPGetRecoveryCodes", t, func() {

	})
	/*
		// admin has no limits, they belong to any org, they admins in any org
		if c, err := GetCredentials(ctx, ts, "alice"); err != nil {
			t.Errorf("couldn't get Credentials for alice: %v", err)
		} else {
			if c.Username() != "alice" {
				t.Errorf("alice's username is %s", c.Username())
			}
			if c.UserID() != 1001 {
				t.Errorf("alice's user id is %d", c.UserID())
			}
			expectNoError(t, c.IsGlobalAdmin(ctx), "alice is not a global admin")
			expectNoError(t, c.IsOrgUser(ctx, 1001), "alice is not orgUser 1001")
			expectNoError(t, c.IsOrgUser(ctx, 1002), "alice is not orgUser 1002")
			expectNoError(t, c.IsOrgAdmin(ctx, 1002), "alice is not orgAdmin 1002")
			expectNoError(t, c.IsDeviceAdmin(ctx, 2000), "alice is not deviceAdmin 2000")
			expectNoError(t, c.IsGatewayAdmin(ctx, 2000), "alice is not gatewayAdmin 2000")
		}

		if c, err := GetCredentials(ctx, ts, "bob"); err != nil {
			t.Errorf("couldn't get Credentials for bob: %v", err)
		} else {
			expectError(t, c.IsGlobalAdmin(ctx), "bob is a global admin")
			expectError(t, c.IsOrgUser(ctx, 1001), "bob is orgUser 1001")
			expectNoError(t, c.IsOrgUser(ctx, 1002), "bob is not orgUser 1002")
			expectNoError(t, c.IsOrgAdmin(ctx, 1002), "bob is not orgAdmin 1002")
			expectNoError(t, c.IsDeviceAdmin(ctx, 1002), "bob is not devAdmin 1002")
			expectNoError(t, c.IsGatewayAdmin(ctx, 1002), "bob is not gwAdmin 1002")

			expectNoError(t, c.IsOrgUser(ctx, 2000), "bob is not orgUser 1002")
			expectError(t, c.IsOrgAdmin(ctx, 2000), "bob is orgAdmin 2000")
			expectNoError(t, c.IsGatewayAdmin(ctx, 2000), "bob is not gwAdmin 2000")
			expectError(t, c.IsDeviceAdmin(ctx, 2000), "bob is devAdmin 2000")

			expectNoError(t, c.IsGatewayAdmin(ctx, 2001), "bob is not gwAdmin 2001")
			expectNoError(t, c.IsDeviceAdmin(ctx, 2001), "bob is not devAdmin 2001")
			expectNoError(t, c.IsOrgAdmin(ctx, 2001), "bob is not orgAdmin 2001")

			expectError(t, c.IsGatewayAdmin(ctx, 2002), "bob is gwAdmin 2002")
			expectNoError(t, c.IsDeviceAdmin(ctx, 2002), "bob is not devAdmin 2002")
			expectError(t, c.IsOrgAdmin(ctx, 2002), "bob is orgAdmin 2002")
		}*/
}

func expectNoError(t *testing.T, err error, message string) {
	t.Helper()
	if err != nil {
		t.Error(message, " ", err)
	}
}

func expectError(t *testing.T, err error, message string) {
	t.Helper()
	if err == nil {
		t.Error(message)
	}
}
