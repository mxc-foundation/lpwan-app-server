package storage

import (
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

// Action defines the action type.
type Action store.Action

// Possible actions
const (
	Select = Action(store.Select)
	Insert = Action(store.Insert)
	Update = Action(store.Update)
	Delete = Action(store.Delete)
	Scan   = Action(store.Scan)
)

// errors
var (
	ErrAlreadyExists                   = store.ErrAlreadyExists
	ErrDoesNotExist                    = store.ErrDoesNotExist
	ErrUsedByOtherObjects              = store.ErrUsedByOtherObjects
	ErrApplicationInvalidName          = store.ErrApplicationInvalidName
	ErrNodeInvalidName                 = store.ErrNodeInvalidName
	ErrNodeMaxRXDelay                  = store.ErrNodeMaxRXDelay
	ErrCFListTooManyChannels           = store.ErrCFListTooManyChannels
	ErrUserInvalidUsername             = store.ErrUserInvalidUsername
	ErrUserPasswordLength              = store.ErrUserPasswordLength
	ErrInvalidUsernameOrPassword       = store.ErrInvalidUsernameOrPassword
	ErrOrganizationInvalidName         = store.ErrOrganizationInvalidName
	ErrGatewayInvalidName              = store.ErrGatewayInvalidName
	ErrInvalidEmail                    = store.ErrInvalidEmail
	ErrInvalidGatewayDiscoveryInterval = store.ErrInvalidGatewayDiscoveryInterval
	ErrDeviceProfileInvalidName        = store.ErrDeviceProfileInvalidName
	ErrServiceProfileInvalidName       = store.ErrServiceProfileInvalidName
	ErrFUOTADeploymentInvalidName      = store.ErrFUOTADeploymentInvalidName
	ErrFUOTADeploymentNullPayload      = store.ErrFUOTADeploymentNullPayload
	ErrMulticastGroupInvalidName       = store.ErrMulticastGroupInvalidName
	ErrOrganizationMaxDeviceCount      = store.ErrOrganizationMaxDeviceCount
	ErrOrganizationMaxGatewayCount     = store.ErrOrganizationMaxGatewayCount
	ErrNetworkServerInvalidName        = store.ErrNetworkServerInvalidName
)

func handlePSQLError(action Action, err error, description string) error {
	return store.HandlePSQLError(store.Action(action), err, description)
}
