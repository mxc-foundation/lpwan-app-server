package types

import (
	"github.com/mxc-foundation/lpwan-app-server/internal/grpccli"
)

// ProvisioningServerStruct defines credentails to connect to provisioning-server
type ProvisioningServerStruct struct {
	ServerConifig  grpccli.ConnectionOpts `mapstructure:"grpc_connection"`
	UpdateSchedule string                 `mapstructure:"update_schedule"`
}
