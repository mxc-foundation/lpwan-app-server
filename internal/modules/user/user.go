package user

import (
	"context"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
)

type UserStore interface {
	CreateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, id int64) (User, error)
	GetUserByExternalID(ctx context.Context, externalID string) (User, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserCount(ctx context.Context) (int, error)
	GetUsers(ctx context.Context, limit, offset int) ([]User, error)
	UpdateUser(ctx context.Context, u *User) error
	DeleteUser(ctx context.Context, id int64) error
	LoginUserByPassword(ctx context.Context, email string, password string) error
	GetProfile(ctx context.Context, id int64) (UserProfile, error)
	GetUserToken(u User) (string, error)
	RegisterUser(user *User, token string) error
	GetUserByToken(token string) (User, error)
	GetTokenByUsername(ctx context.Context, username string) (string, error)
	FinishRegistration(userID int64, password string) error
	UpdatePassword(ctx context.Context, id int64, newpassword string) error
	GetPasswordResetRecord(userID int64) (*PasswordResetRecord, error)

	SetOTP(ctx context.Context, pr *PasswordResetRecord) error
	ReduceAttempts(ctx context.Context, pr *PasswordResetRecord) error

	// validator
	CheckActiveUser(username string, userID int64) (bool, error)

	CheckCreateUserAcess(username string, userID int64) (bool, error)
	CheckListUserAcess(username string, userID int64) (bool, error)

	CheckReadUserAccess(username string, userID, operatorUserID int64) (bool, error)
	CheckUpdateDeleteUserAccess(username string, userID, operatorUserID int64) (bool, error)
	CheckUpdateProfileUserAccess(username string, userID, operatorUserID int64) (bool, error)
	CheckUpdatePasswordUserAccess(username string, userID, operatorUserID int64) (bool, error)
}

type Controller struct {
	St        UserStore
	Validator Validator
}

var Service *Controller

func Setup() error {
	Service.St = store.New(storage.DB().DB)
	return nil
}
