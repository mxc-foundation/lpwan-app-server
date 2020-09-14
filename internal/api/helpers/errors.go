package helpers

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	err "github.com/mxc-foundation/lpwan-app-server/internal/errors"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

var errToCode = map[error]codes.Code{
	store.ErrAlreadyExists:                   codes.AlreadyExists,
	store.ErrDoesNotExist:                    codes.NotFound,
	store.ErrUsedByOtherObjects:              codes.FailedPrecondition,
	store.ErrApplicationInvalidName:          codes.InvalidArgument,
	store.ErrNodeInvalidName:                 codes.InvalidArgument,
	store.ErrNodeMaxRXDelay:                  codes.InvalidArgument,
	store.ErrCFListTooManyChannels:           codes.InvalidArgument,
	store.ErrUserInvalidUsername:             codes.InvalidArgument,
	store.ErrUserPasswordLength:              codes.InvalidArgument,
	store.ErrInvalidUsernameOrPassword:       codes.Unauthenticated,
	store.ErrInvalidEmail:                    codes.InvalidArgument,
	store.ErrInvalidGatewayDiscoveryInterval: codes.InvalidArgument,
	store.ErrDeviceProfileInvalidName:        codes.InvalidArgument,
	store.ErrServiceProfileInvalidName:       codes.InvalidArgument,
	store.ErrMulticastGroupInvalidName:       codes.InvalidArgument,
	store.ErrOrganizationMaxDeviceCount:      codes.FailedPrecondition,
	store.ErrOrganizationMaxGatewayCount:     codes.FailedPrecondition,
	store.ErrNetworkServerInvalidName:        codes.InvalidArgument,
	store.ErrFUOTADeploymentInvalidName:      codes.InvalidArgument,
	store.ErrFUOTADeploymentNullPayload:      codes.InvalidArgument,
	err.ErrInvalidHeaderName:                 codes.InvalidArgument,
	err.ErrInvalidPrecision:                  codes.InvalidArgument,
}

// ErrToRPCError converts the given error into a gRPC error.
func ErrToRPCError(err error) error {
	cause := errors.Cause(err)

	// if the err has already a gRPC status (it is a gRPC error), just
	// return the error.
	if code := status.Code(cause); code != codes.Unknown {
		return cause
	}

	code, ok := errToCode[cause]
	if !ok {
		code = codes.Unknown
	}
	return grpc.Errorf(code, cause.Error())
}
