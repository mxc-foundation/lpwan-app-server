package storage

import (
	"context"
	"time"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

// saltSize defines the salt size
const saltSize = 16

// defaultSessionTTL defines the default session TTL
const defaultSessionTTL = time.Hour * 24

// User defines the user structure.
type User store.User

// Validate validates the user data.
func (u User) Validate() error {
	return store.User(u).Validate()
}

// UserProfile contains the profile of the user.
type UserProfile store.UserProfile

// UserProfileUser contains the user information of the profile.
type UserProfileUser store.UserProfileUser

// UserProfileOrganization contains the organizations to which the user
// is linked.
type UserProfileOrganization store.UserProfileOrganization

// CreateUser creates the given user.
func CreateUser(ctx context.Context, handler *store.Handler, user *User) error {
	return handler.CreateUser(ctx, (*store.User)(user))
}

// GetUser returns the User for the given id.
func GetUser(ctx context.Context, handler *store.Handler, id int64) (User, error) {
	res, err := handler.GetUser(ctx, id)
	return User(res), err
}

// GetUserByExternalID returns the User for the given ext. ID.
func GetUserByExternalID(ctx context.Context, handler *store.Handler, externalID string) (User, error) {
	res, err := handler.GetUserByExternalID(ctx, externalID)
	return User(res), err
}

// GetUserByEmail returns the User for the given email.
func GetUserByEmail(ctx context.Context, handler *store.Handler, email string) (User, error) {
	res, err := handler.GetUserByEmail(ctx, email)
	return User(res), err
}

// GetUserCount returns the total number of users.
func GetUserCount(ctx context.Context, handler *store.Handler) (int, error) {
	return handler.GetUserCount(ctx)
}

// GetUsers returns a slice of users, respecting the given limit and offset.
func GetUsers(ctx context.Context, handler *store.Handler, limit, offset int) ([]User, error) {
	res, err := handler.GetUsers(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	var userList []User
	for _, v := range res {
		userItem := User(v)
		userList = append(userList, userItem)
	}

	return userList, nil
}

// UpdateUser updates the given User.
func UpdateUser(ctx context.Context, handler *store.Handler, u *User) error {
	return handler.UpdateUser(ctx, (*store.User)(u))
}

// DeleteUser deletes the User record matching the given ID.
func DeleteUser(ctx context.Context, handler *store.Handler, id int64) error {
	return handler.DeleteUser(ctx, id)
}

// LoginUserByPassword returns a JWT token for the user matching the given email
// and password combination.
func LoginUserByPassword(ctx context.Context, handler *store.Handler, email string, password string) error {
	return handler.LoginUserByPassword(ctx, email, password)
}

// GetProfile returns the user profile (user, applications and organizations
// to which the user is linked).
func GetProfile(ctx context.Context, handler *store.Handler, id int64) (UserProfile, error) {
	res, err := handler.GetProfile(ctx, id)
	return UserProfile(res), err
}
