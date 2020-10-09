package helpers

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
)

var errToCode = map[error]codes.Code{
	errHandler.ErrAlreadyExists:                   codes.AlreadyExists,
	errHandler.ErrDoesNotExist:                    codes.NotFound,
	errHandler.ErrUsedByOtherObjects:              codes.FailedPrecondition,
	errHandler.ErrApplicationInvalidName:          codes.InvalidArgument,
	errHandler.ErrNodeInvalidName:                 codes.InvalidArgument,
	errHandler.ErrNodeMaxRXDelay:                  codes.InvalidArgument,
	errHandler.ErrCFListTooManyChannels:           codes.InvalidArgument,
	errHandler.ErrUserInvalidUsername:             codes.InvalidArgument,
	errHandler.ErrUserPasswordLength:              codes.InvalidArgument,
	errHandler.ErrInvalidUsernameOrPassword:       codes.Unauthenticated,
	errHandler.ErrInvalidEmail:                    codes.InvalidArgument,
	errHandler.ErrInvalidGatewayDiscoveryInterval: codes.InvalidArgument,
	errHandler.ErrDeviceProfileInvalidName:        codes.InvalidArgument,
	errHandler.ErrServiceProfileInvalidName:       codes.InvalidArgument,
	errHandler.ErrMulticastGroupInvalidName:       codes.InvalidArgument,
	errHandler.ErrOrganizationMaxDeviceCount:      codes.FailedPrecondition,
	errHandler.ErrOrganizationMaxGatewayCount:     codes.FailedPrecondition,
	errHandler.ErrNetworkServerInvalidName:        codes.InvalidArgument,
	errHandler.ErrFUOTADeploymentInvalidName:      codes.InvalidArgument,
	errHandler.ErrFUOTADeploymentNullPayload:      codes.InvalidArgument,
	errHandler.ErrInvalidHeaderName:               codes.InvalidArgument,
	errHandler.ErrInvalidPrecision:                codes.InvalidArgument,
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
