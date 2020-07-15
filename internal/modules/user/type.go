package user

import (
	"context"
	"regexp"
	"time"

	"github.com/pkg/errors"
)

// UserProfile contains the profile of the user.
type UserProfile struct {
	User          UserProfileUser
	Organizations []UserProfileOrganization
}

// UserProfileUser contains the user information of the profile.
type UserProfileUser struct {
	ID         int64     `db:"id"`
	Username   string    `db:"username"`
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

// User defines the user structure.
type User struct {
	ID            int64     `db:"id"`
	IsAdmin       bool      `db:"is_admin"`
	IsActive      bool      `db:"is_active"`
	SessionTTL    int32     `db:"session_ttl"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
	Password      string
	PasswordHash  string  `db:"password_hash"`
	Email         string  `db:"email"`
	EmailVerified bool    `db:"email_verified"`
	EmailOld      string  `db:"email_old"`
	Note          string  `db:"note"`
	ExternalID    *string `db:"external_id"` // must be pointer for unique index
	SecurityToken *string `db:"security_token"`
	Username      string  `db:"username"`
}

const (
	// saltSize defines the salt size
	saltSize = 16
	//  defines the default session TTL
	DefaultSessionTTL = time.Hour * 24
)

var (
	// Any upper, lower, digit characters, at least 6 characters.
	usernameValidator = regexp.MustCompile(`^[[:alnum:]]+$`)

	// Any printable characters, at least 6 characters.
	passwordValidator = regexp.MustCompile(`^.{6,}$`)

	// Must contain @ (this is far from perfect)
	emailValidator = regexp.MustCompile(`.+@.+`)
)

// Validate validates the user data.
func (u *User) Validate() error {
	if u.Email == "" || !emailValidator.MatchString(u.Email) {
		return errors.New("invalid email")
	}

	if u.Password != "" && !passwordValidator.MatchString(u.Password) {
		return errors.New("invalid password")
	}

	if u.Username == "" || !usernameValidator.MatchString(u.Username) {
		return errors.New("invalid username")
	}

	return nil
}

func ValidatePassword(password string) error {
	if !passwordValidator.MatchString(password) {
		return errors.New("invalid password")
	}
	return nil
}

type PasswordResetRecord struct {
	UserID       int64
	OTP          string
	GeneratedAt  time.Time
	AttemptsLeft int64
}

func (pr *PasswordResetRecord) SetOTP(ctx context.Context, otp string) error {
	pr.OTP = otp
	pr.GeneratedAt = time.Now()
	pr.AttemptsLeft = 3

	if err := Service.St.SetOTP(ctx, pr); err != nil {
		return err
	}

	return nil
}

func (pr *PasswordResetRecord) ReduceAttempts(ctx context.Context) error {
	if err := Service.St.ReduceAttempts(ctx, pr); err != nil {
		return err
	}
	return nil
}
