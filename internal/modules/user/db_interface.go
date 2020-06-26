package user

import (
	"context"

	"github.com/jmoiron/sqlx"

	pg "github.com/mxc-foundation/lpwan-app-server/internal/storage/postgresql"
)

type UserTable interface {
	CreateUser(ctx context.Context, db sqlx.Queryer, user *User) error
	GetUser(ctx context.Context, db sqlx.Queryer, id int64) (User, error)
	GetUserByExternalID(ctx context.Context, db sqlx.Queryer, externalID string) (User, error)
	GetUserByEmail(ctx context.Context, db sqlx.Queryer, email string) (User, error)
	GetUserCount(ctx context.Context, db sqlx.Queryer) (int, error)
	GetUsers(ctx context.Context, db sqlx.Queryer, limit, offset int) ([]User, error)
	UpdateUser(ctx context.Context, db sqlx.Execer, u *User) error
	DeleteUser(ctx context.Context, db sqlx.Execer, id int64) error
	LoginUserByPassword(ctx context.Context, db sqlx.Queryer, email string, password string) (string, error)
	GetProfile(ctx context.Context, db sqlx.Queryer, id int64) (UserProfile, error)
	GetUserToken(u User) (string, error)
	RegisterUser(db sqlx.Queryer, user *User, token string) error
	GetUserByToken(db sqlx.Queryer, token string) (User, error)
	GetTokenByUsername(ctx context.Context, db sqlx.Queryer, username string) (string, error)
	FinishRegistration(db sqlx.Execer, userID int64, newPwd string) error
}

var UserDB = UserTable(&pg.UserTable)
