package storage

import (
	"time"
)

// User defines the user structure.
type User struct {
	ID            int64     `db:"id"`
	IsAdmin       bool      `db:"is_admin"`
	IsActive      bool      `db:"is_active"`
	SessionTTL    int32     `db:"session_ttl"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
	PasswordHash  string    `db:"password_hash"`
	Email         string    `db:"email"`
	EmailVerified bool      `db:"email_verified"`
	EmailOld      string    `db:"email_old"`
	Note          string    `db:"note"`
	ExternalID    *string   `db:"external_id"` // must be pointer for unique index
	SecurityToken *string   `db:"security_token"`
}

// UserProfile contains the profile of the user.
type UserProfile struct {
	User          UserProfileUser
	Organizations []UserProfileOrganization
}

// UserProfileUser contains the user information of the profile.
type UserProfileUser struct {
	ID         int64     `db:"id"`
	Email      string    `db:"email"`
	IsAdmin    bool      `db:"is_admin"`
	IsActive   bool      `db:"is_active"`
	SessionTTL int32     `db:"session_ttl"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

// UserProfileOrganization contains the organizations to which the user
// is linked.
type UserProfileOrganization struct {
	ID             int64     `db:"organization_id"`
	Name           string    `db:"organization_name"`
	IsAdmin        bool      `db:"is_admin"`
	IsDeviceAdmin  bool      `db:"is_device_admin"`
	IsGatewayAdmin bool      `db:"is_gateway_admin"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
