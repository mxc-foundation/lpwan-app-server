package default_gateway_config

import (
	"github.com/jmoiron/sqlx"

	"github.com/mxc-foundation/lpwan-app-server/internal/storage/postgresql"
)

type DefaultGatewayConfigTable interface {
	AddNewDefaultGatewayConfig(db sqlx.Execer, defaultConfig *DefaultGatewayConfig) error
	UpdateDefaultGatewayConfig(db sqlx.Execer, defaultConfig *DefaultGatewayConfig) error
	GetDefaultGatewayConfig(db sqlx.Queryer, defaultConfig *DefaultGatewayConfig) error
}

var DefaultGatewayConfigDB = DefaultGatewayConfigTable(&postgresql.DefaultGatewayConfigTable)

type GatewayTable interface {
	//TODO
}

var Gateway = GatewayTable(&struct {
	//TODO
}{})
