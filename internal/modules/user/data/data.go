package data

import (
	"errors"
	"regexp"
	"time"

	"github.com/brocaar/lorawan"

	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
)

type RecaptchaStruct struct {
	HostServer string `mapstructure:"host_server"`
	Secret     string `mapstructure:"secret"`
}

type Config struct {
	Recaptcha      RecaptchaStruct
	Enable2FALogin bool
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
}

// Validate validates the user data.
func (u User) Validate() error {
	if !emailValidator.MatchString(u.Email) {
		return errHandler.ErrInvalidEmail
	}

	return nil
}

const (
	// saltSize defines the salt size
	saltSize = 16
	//  defines the default session TTL
	DefaultSessionTTL = time.Hour * 24
)

var (
	// Any printable characters, at least 6 characters.
	passwordValidator = regexp.MustCompile(`^.{6,}$`)

	// Must contain @ (this is far from perfect)
	emailValidator = regexp.MustCompile(`.+@.+`)
)

func ValidatePassword(password string) error {
	if !passwordValidator.MatchString(password) {
		return errors.New("invalid password")
	}
	return nil
}

// SearchResult defines a search result.
type SearchResult struct {
	Kind             string         `db:"kind"`
	Score            float64        `db:"score"`
	OrganizationID   *int64         `db:"organization_id"`
	OrganizationName *string        `db:"organization_name"`
	ApplicationID    *int64         `db:"application_id"`
	ApplicationName  *string        `db:"application_name"`
	DeviceDevEUI     *lorawan.EUI64 `db:"device_dev_eui"`
	DeviceName       *string        `db:"device_name"`
	GatewayMAC       *lorawan.EUI64 `db:"gateway_mac"`
	GatewayName      *string        `db:"gateway_name"`
}
