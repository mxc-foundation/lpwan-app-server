package authentication

import (
	"fmt"
	"testing"

	"golang.org/x/net/context"
)

type testUser struct {
	id      int64
	isAdmin bool
	orgs    map[int64]OrgUser
}

type testStore struct {
	users     map[string]testUser
	usersByID map[int64]string
}

func (ts *testStore) init() {
	ts.usersByID = make(map[int64]string)
	for name, u := range ts.users {
		ts.usersByID[u.id] = name
	}
}

func (ts *testStore) GetUser(ctx context.Context, username string) (*User, error) {
	u, ok := ts.users[username]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return &User{
		ID:            u.id,
		IsGlobalAdmin: u.isAdmin,
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

func TestCredentials(t *testing.T) {
	ts := &testStore{
		users: map[string]testUser{
			"alice": {
				id:      1001,
				isAdmin: true,
				orgs: map[int64]OrgUser{
					1001: {},
				},
			},
			"bob": {
				id: 1002,
				orgs: map[int64]OrgUser{
					1002: {IsOrgAdmin: true},
					2000: {IsGatewayAdmin: true},
					2001: {IsOrgAdmin: true, IsGatewayAdmin: true},
					2002: {IsDeviceAdmin: true},
				},
			},
		},
	}
	ts.init()
	ctx := context.Background()

	if c, err := GetCredentials(ctx, ts, "mallory"); err == nil {
		t.Errorf("got Credentials for mallory: %#v", c)
	}

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
	}
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
