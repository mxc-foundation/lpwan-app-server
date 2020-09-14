package store

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/mxc-foundation/lpwan-app-server/internal/pwhash"
)

type UserStore interface {
	CreateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, id int64) (User, error)
	GetUserByExternalID(ctx context.Context, externalID string) (User, error)
	GetUserByUsername(ctx context.Context, userEmail string) (User, error)
	GetUserByEmail(ctx context.Context, userEmail string) (User, error)
	GetUserCount(ctx context.Context) (int, error)
	GetUsers(ctx context.Context, limit, offset int) ([]User, error)
	UpdateUser(ctx context.Context, u *User) error
	DeleteUser(ctx context.Context, id int64) error
	LoginUserByPassword(ctx context.Context, userEmail string, password string) error
	GetProfile(ctx context.Context, id int64) (UserProfile, error)
	GetUserToken(ctx context.Context, u User) (string, error)
	RegisterUser(ctx context.Context, user *User, token string) error
	GetUserByToken(ctx context.Context, token string) (User, error)
	GetTokenByUsername(ctx context.Context, userEmail string) (string, error)
	FinishRegistration(ctx context.Context, userID int64, password string, pwh *pwhash.PasswordHasher) error
	UpdatePassword(ctx context.Context, id int64, newpassword string, pwh *pwhash.PasswordHasher) error
	GetPasswordResetRecord(ctx context.Context, userID int64) (*PasswordResetRecord, error)

	SetOTP(ctx context.Context, pr *PasswordResetRecord) error
	ReduceAttempts(ctx context.Context, pr *PasswordResetRecord) error

	// validator
	CheckActiveUser(ctx context.Context, userEmail string, userID int64) (bool, error)

	CheckCreateUserAcess(ctx context.Context, userEmail string, userID int64) (bool, error)
	CheckListUserAcess(ctx context.Context, userEmail string, userID int64) (bool, error)

	CheckReadUserAccess(ctx context.Context, userEmail string, userID, operatorUserID int64) (bool, error)
	CheckUpdateUserAccess(ctx context.Context, userEmail string, userID, operatorUserID int64) (bool, error)
	CheckDeleteUserAccess(ctx context.Context, userEmail string, userID, operatorUserID int64) (bool, error)
	CheckUpdateProfileUserAccess(ctx context.Context, userEmail string, userID, operatorUserID int64) (bool, error)
	CheckUpdatePasswordUserAccess(ctx context.Context, userEmail string, userID, operatorUserID int64) (bool, error)
}

func (h *Handler) CreateUser(ctx context.Context, user *User) error {
	return h.store.CreateUser(ctx, user)
}
func (h *Handler) GetUser(ctx context.Context, id int64) (User, error) {
	return h.store.GetUser(ctx, id)
}
func (h *Handler) GetUserByExternalID(ctx context.Context, externalID string) (User, error) {
	return h.store.GetUserByExternalID(ctx, externalID)
}
func (h *Handler) GetUserByUsername(ctx context.Context, userEmail string) (User, error) {
	return h.store.GetUserByUsername(ctx, userEmail)
}
func (h *Handler) GetUserByEmail(ctx context.Context, userEmail string) (User, error) {
	return h.store.GetUserByEmail(ctx, userEmail)
}
func (h *Handler) GetUserCount(ctx context.Context) (int, error) {
	return h.store.GetUserCount(ctx)
}
func (h *Handler) GetUsers(ctx context.Context, limit, offset int) ([]User, error) {
	return h.store.GetUsers(ctx, limit, offset)
}
func (h *Handler) UpdateUser(ctx context.Context, u *User) error {
	return h.store.UpdateUser(ctx, u)
}
func (h *Handler) DeleteUser(ctx context.Context, id int64) error {
	return h.store.DeleteUser(ctx, id)
}
func (h *Handler) LoginUserByPassword(ctx context.Context, userEmail string, password string) error {
	return h.store.LoginUserByPassword(ctx, userEmail, password)
}
func (h *Handler) GetProfile(ctx context.Context, id int64) (UserProfile, error) {
	return h.store.GetProfile(ctx, id)
}
func (h *Handler) GetUserToken(ctx context.Context, u User) (string, error) {
	return h.store.GetUserToken(ctx, u)
}
func (h *Handler) RegisterUser(ctx context.Context, user *User, token string) error {
	return h.store.RegisterUser(ctx, user, token)
}
func (h *Handler) GetUserByToken(ctx context.Context, token string) (User, error) {
	return h.store.GetUserByToken(ctx, token)
}
func (h *Handler) GetTokenByUsername(ctx context.Context, userEmail string) (string, error) {
	return h.store.GetTokenByUsername(ctx, userEmail)
}
func (h *Handler) FinishRegistration(ctx context.Context, userID int64, password string, pwh *pwhash.PasswordHasher) error {
	return h.store.FinishRegistration(ctx, userID, password, pwh)
}
func (h *Handler) UpdatePassword(ctx context.Context, id int64, newpassword string, pwh *pwhash.PasswordHasher) error {
	return h.store.UpdatePassword(ctx, id, newpassword, pwh)
}
func (h *Handler) GetPasswordResetRecord(ctx context.Context, userID int64) (*PasswordResetRecord, error) {
	return h.store.GetPasswordResetRecord(ctx, userID)
}

func (h *Handler) SetOTP(ctx context.Context, pr *PasswordResetRecord) error {
	return h.store.SetOTP(ctx, pr)
}
func (h *Handler) ReduceAttempts(ctx context.Context, pr *PasswordResetRecord) error {
	return h.store.ReduceAttempts(ctx, pr)
}

// validator
func (h *Handler) CheckActiveUser(ctx context.Context, userEmail string, userID int64) (bool, error) {
	return h.store.CheckActiveUser(ctx, userEmail, userID)
}

func (h *Handler) CheckCreateUserAcess(ctx context.Context, userEmail string, userID int64) (bool, error) {
	return h.store.CheckCreateUserAcess(ctx, userEmail, userID)
}
func (h *Handler) CheckListUserAcess(ctx context.Context, userEmail string, userID int64) (bool, error) {
	return h.store.CheckListUserAcess(ctx, userEmail, userID)
}

func (h *Handler) CheckReadUserAccess(ctx context.Context, userEmail string, userID, operatorUserID int64) (bool, error) {
	return h.store.CheckReadUserAccess(ctx, userEmail, userID, operatorUserID)
}
func (h *Handler) CheckUpdateUserAccess(ctx context.Context, userEmail string, userID, operatorUserID int64) (bool, error) {
	return h.store.CheckUpdateUserAccess(ctx, userEmail, userID, operatorUserID)
}
func (h *Handler) CheckDeleteUserAccess(ctx context.Context, userEmail string, userID, operatorUserID int64) (bool, error) {
	return h.store.CheckDeleteUserAccess(ctx, userEmail, userID, operatorUserID)
}
func (h *Handler) CheckUpdateProfileUserAccess(ctx context.Context, userEmail string, userID, operatorUserID int64) (bool, error) {
	return h.store.CheckUpdateProfileUserAccess(ctx, userEmail, userID, operatorUserID)
}
func (h *Handler) CheckUpdatePasswordUserAccess(ctx context.Context, userEmail string, userID, operatorUserID int64) (bool, error) {
	return h.store.CheckUpdatePasswordUserAccess(ctx, userEmail, userID, operatorUserID)
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
		return ErrInvalidEmail
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

type PasswordResetRecord struct {
	UserID       int64
	OTP          string
	GeneratedAt  time.Time
	AttemptsLeft int64
	St           UserStore
}

func (pr *PasswordResetRecord) SetOTP(ctx context.Context, otp string) error {
	pr.OTP = otp
	pr.GeneratedAt = time.Now()
	pr.AttemptsLeft = 3

	if err := pr.St.SetOTP(ctx, pr); err != nil {
		return err
	}

	return nil
}

func (pr *PasswordResetRecord) ReduceAttempts(ctx context.Context) error {
	if err := pr.St.ReduceAttempts(ctx, pr); err != nil {
		return err
	}
	return nil
}
