package data

import (
	"github.com/gofrs/uuid"

	"github.com/mxc-foundation/lpwan-app-server/internal/pwhash"
)

type Config struct {
	ApplicationServerID         uuid.UUID
	JWTSecret                   string
	ApplicationServerPublicHost string
	PWH                         *pwhash.PasswordHasher
}

type PostgreSQLStruct struct {
	DSN                string `mapstructure:"dsn"`
	Automigrate        bool
	MaxOpenConnections int `mapstructure:"max_open_connections"`
	MaxIdleConnections int `mapstructure:"max_idle_connections"`
}
