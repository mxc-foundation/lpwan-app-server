package store

import (
	"context"

	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"

	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/user/data"
)

func NewStore(pg pgstore.PgStore) *usrs {
	return &usrs{
		pg: pg,
	}
}

type usrs struct {
	pg pgstore.UserPgStore
}

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
	FinishRegistration(ctx context.Context, userID int64, password string) error
	UpdatePassword(ctx context.Context, id int64, newpassword string) error
	RequestPasswordReset(ctx context.Context, userID int64, otp string) (string, error)
	ConfirmPasswordReset(ctx context.Context, userID int64, otp string, newPassword string) error
	GlobalSearch(ctx context.Context, userID int64, globalAdmin bool, search string, limit, offset int) ([]SearchResult, error)

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

func (h *usrs) GlobalSearch(ctx context.Context, userID int64, globalAdmin bool, search string, limit, offset int) ([]SearchResult, error) {
	return h.pg.GlobalSearch(ctx, userID, globalAdmin, search, limit, offset)
}

func (h *usrs) RequestPasswordReset(ctx context.Context, userID int64, otp string) (string, error) {
	return h.pg.ResetPassword(ctx, userID, otp)
}

func (h *usrs) ConfirmPasswordReset(ctx context.Context, userID int64, otp string, newPassword string) error {
	return h.pg.ConfirmPasswordReset(ctx, userID, otp, newPassword)
}

func (h *usrs) CreateUser(ctx context.Context, user *User) error {
	return h.pg.CreateUser(ctx, user)
}
func (h *usrs) GetUser(ctx context.Context, id int64) (User, error) {
	return h.pg.GetUser(ctx, id)
}
func (h *usrs) GetUserByExternalID(ctx context.Context, externalID string) (User, error) {
	return h.pg.GetUserByExternalID(ctx, externalID)
}
func (h *usrs) GetUserByUsername(ctx context.Context, userEmail string) (User, error) {
	return h.pg.GetUserByUsername(ctx, userEmail)
}
func (h *usrs) GetUserByEmail(ctx context.Context, userEmail string) (User, error) {
	return h.pg.GetUserByEmail(ctx, userEmail)
}
func (h *usrs) GetUserCount(ctx context.Context) (int, error) {
	return h.pg.GetUserCount(ctx)
}
func (h *usrs) GetUsers(ctx context.Context, limit, offset int) ([]User, error) {
	return h.pg.GetUsers(ctx, limit, offset)
}
func (h *usrs) UpdateUser(ctx context.Context, u *User) error {
	return h.pg.UpdateUser(ctx, u)
}
func (h *usrs) DeleteUser(ctx context.Context, id int64) error {
	return h.pg.DeleteUser(ctx, id)
}
func (h *usrs) LoginUserByPassword(ctx context.Context, userEmail string, password string) error {
	return h.pg.LoginUserByPassword(ctx, userEmail, password)
}
func (h *usrs) GetProfile(ctx context.Context, id int64) (UserProfile, error) {
	return h.pg.GetProfile(ctx, id)
}
func (h *usrs) GetUserToken(ctx context.Context, u User) (string, error) {
	return h.pg.GetUserToken(ctx, u)
}
func (h *usrs) RegisterUser(ctx context.Context, user *User, token string) error {
	return h.pg.RegisterUser(ctx, user, token)
}
func (h *usrs) GetUserByToken(ctx context.Context, token string) (User, error) {
	return h.pg.GetUserByToken(ctx, token)
}
func (h *usrs) GetTokenByUsername(ctx context.Context, userEmail string) (string, error) {
	return h.pg.GetTokenByUsername(ctx, userEmail)
}
func (h *usrs) FinishRegistration(ctx context.Context, userID int64, password string) error {
	return h.pg.FinishRegistration(ctx, userID, password)
}
func (h *usrs) UpdatePassword(ctx context.Context, id int64, newpassword string) error {
	return h.pg.UpdatePassword(ctx, id, newpassword)
}

// validator
func (h *usrs) CheckActiveUser(ctx context.Context, userEmail string, userID int64) (bool, error) {
	return h.pg.CheckActiveUser(ctx, userEmail, userID)
}

func (h *usrs) CheckCreateUserAcess(ctx context.Context, userEmail string, userID int64) (bool, error) {
	return h.pg.CheckCreateUserAcess(ctx, userEmail, userID)
}
func (h *usrs) CheckListUserAcess(ctx context.Context, userEmail string, userID int64) (bool, error) {
	return h.pg.CheckListUserAcess(ctx, userEmail, userID)
}

func (h *usrs) CheckReadUserAccess(ctx context.Context, userEmail string, userID, operatorUserID int64) (bool, error) {
	return h.pg.CheckReadUserAccess(ctx, userEmail, userID, operatorUserID)
}
func (h *usrs) CheckUpdateUserAccess(ctx context.Context, userEmail string, userID, operatorUserID int64) (bool, error) {
	return h.pg.CheckUpdateUserAccess(ctx, userEmail, userID, operatorUserID)
}
func (h *usrs) CheckDeleteUserAccess(ctx context.Context, userEmail string, userID, operatorUserID int64) (bool, error) {
	return h.pg.CheckDeleteUserAccess(ctx, userEmail, userID, operatorUserID)
}
func (h *usrs) CheckUpdateProfileUserAccess(ctx context.Context, userEmail string, userID, operatorUserID int64) (bool, error) {
	return h.pg.CheckUpdateProfileUserAccess(ctx, userEmail, userID, operatorUserID)
}
func (h *usrs) CheckUpdatePasswordUserAccess(ctx context.Context, userEmail string, userID, operatorUserID int64) (bool, error) {
	return h.pg.CheckUpdatePasswordUserAccess(ctx, userEmail, userID, operatorUserID)
}
