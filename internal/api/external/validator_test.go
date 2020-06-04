package external

import (
	"context"
	"fmt"

	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/auth"
)

type TestValidator struct {
	ctx            context.Context
	validatorFuncs []auth.ValidatorFunc
	returnError    error
	returnUsername string
	returnIsAdmin  bool
}

func (v *TestValidator) Validate(ctx context.Context, funcs ...auth.ValidatorFunc) error {
	v.ctx = ctx
	v.validatorFuncs = funcs
	return v.returnError
}

func (v *TestValidator) GetUsername(ctx context.Context) (string, error) {
	return v.returnUsername, v.returnError
}

func (v *TestValidator) GetOTP(ctx context.Context) string {
	return ""
}

func (v *TestValidator) GetIsAdmin(ctx context.Context) (bool, error) {
	return v.returnIsAdmin, v.returnError
}

func (v *TestValidator) GetCredentials(ctx context.Context, opts ...auth.Option) (auth.Credentials, error) {
	return nil, fmt.Errorf("not implemented")
}

func (v *TestValidator) SignToken(username string, ttl int64, audience []string) (string, error) {
	return "foo", nil
}
