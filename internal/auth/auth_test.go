package auth

import (
	"context"
	"fmt"
	"testing"

	"github.com/gofrs/uuid"
)

type tStore struct {
	users    map[string]User
	orgUsers map[int64]map[int64]OrgUser
}

func (ts *tStore) ApplicationOwnedByOrganization(ctx context.Context, orgID, applicationID int64) (bool, error) {
	return false, nil
}

func (ts *tStore) DeviceProfileOwnedByOrganization(ctx context.Context, orgID int64, deviceProfile uuid.UUID) (bool, error) {
	return false, nil
}

func (ts *tStore) AuthGetUser(ctx context.Context, username string) (User, error) {
	u, ok := ts.users[username]
	if !ok {
		return u, fmt.Errorf("not found")
	}
	return u, nil
}

func (ts *tStore) AuthGetOrgUser(ctx context.Context, userID, orgID int64) (OrgUser, error) {
	if orgID == 111 {
		return OrgUser{}, fmt.Errorf("db error")
	}
	return ts.orgUsers[userID][orgID], nil
}

func TestCredentials(t *testing.T) {
	ts := &tStore{
		users: map[string]User{
			"alice@example.com": {
				ID:            13,
				Email:         "alice@example.com",
				IsGlobalAdmin: false,
			},
			"bob@example.com": {
				ID:            17,
				Email:         "bob@example.com",
				IsGlobalAdmin: true,
			},
		},
		orgUsers: map[int64]map[int64]OrgUser{
			13: {
				5: OrgUser{IsOrgUser: true, IsOrgAdmin: true},
				7: OrgUser{IsOrgUser: true, IsDeviceAdmin: true},
			},
			17: {
				5: OrgUser{IsOrgUser: true},
			},
		},
	}
	ctx := context.Background()
	malory, err := NewCredentials(ctx, ts, "malory@example.com", 5, EMAIL, 0, "")
	if err == nil {
		t.Errorf("got credentials for mallory: %v", *malory)
	}
	alice, err := NewCredentials(ctx, ts, "alice@example.com", 111, EMAIL, 0, "")
	if err == nil {
		t.Errorf("expected db error, but got: %v", *alice)
	}
	tests := []struct {
		name     string
		username string
		orgid    int64
		expected Credentials
	}{
		{
			name:     "if user exist but not in the org, all org fields must be false",
			username: "alice@example.com",
			orgid:    1,
			expected: Credentials{
				UserID:     13,
				Username:   "alice@example.com",
				IsExisting: true,
				OrgID:      1,
				Service:    EMAIL,
			},
		},
		{
			name:     "if user is org admin, all org fields must be true",
			username: "alice@example.com",
			orgid:    5,
			expected: Credentials{
				UserID:         13,
				Username:       "alice@example.com",
				IsExisting:     true,
				OrgID:          5,
				IsOrgUser:      true,
				IsOrgAdmin:     true,
				IsDeviceAdmin:  true,
				IsGatewayAdmin: true,
				Service:        EMAIL,
			},
		},
		{
			name:     "if user is org user, the relevalnt org fields must be true",
			username: "alice@example.com",
			orgid:    7,
			expected: Credentials{
				UserID:        13,
				Username:      "alice@example.com",
				IsExisting:    true,
				OrgID:         7,
				IsOrgUser:     true,
				IsDeviceAdmin: true,
				Service:       EMAIL,
			},
		},
		{
			name:     "admin user is a member and admin for every org",
			username: "bob@example.com",
			orgid:    5,
			expected: Credentials{
				UserID:         17,
				Username:       "bob@example.com",
				IsGlobalAdmin:  true,
				IsExisting:     true,
				OrgID:          5,
				IsOrgUser:      true,
				IsOrgAdmin:     true,
				IsDeviceAdmin:  true,
				IsGatewayAdmin: true,
				Service:        EMAIL,
			},
		},
		{
			name:     "admin user is a member and admin for every org",
			username: "bob@example.com",
			orgid:    7,
			expected: Credentials{
				UserID:         17,
				Username:       "bob@example.com",
				IsGlobalAdmin:  true,
				IsExisting:     true,
				OrgID:          7,
				IsOrgUser:      true,
				IsOrgAdmin:     true,
				IsDeviceAdmin:  true,
				IsGatewayAdmin: true,
				Service:        EMAIL,
			},
		},
	}
	for _, tc := range tests {
		t.Logf(tc.name)
		cred, err := NewCredentials(ctx, ts, tc.username, tc.orgid, EMAIL, 0, "")
		if err != nil {
			t.Errorf("couldn't get credentials: %v", err)
		}
		if *cred != tc.expected {
			t.Errorf("expected: %v, got %v", tc.expected, *cred)
		}
	}
}

func TestOptions(t *testing.T) {
	defaults := NewOptions()
	expDefaults := Options{
		Audience: "lora-app-server",
	}
	if *defaults != expDefaults {
		t.Errorf("expected defaults: %v, got %v", expDefaults, *defaults)
	}
	opts := defaults.WithAudience("test").WithAllowNonExisting().WithRequireOTP().WithOrgID(19).WithExternalLimited()
	expOpts := Options{
		Audience:         "test",
		RequireOTP:       true,
		AllowNonExisting: true,
		OrgID:            19,
		ExternalLimited:  true,
	}
	if *opts != expOpts {
		t.Errorf("expected opts: %v, got %v", expOpts, *opts)
	}
}
