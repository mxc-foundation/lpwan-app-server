package data

// IntegrationPostgreSQLConfig holds the PostgreSQL integration configuration.
type IntegrationPostgreSQLConfig struct {
	DSN                string `json:"dsn"`
	MaxOpenConnections int    `mapstructure:"max_open_connections"`
	MaxIdleConnections int    `mapstructure:"max_idle_connections"`
}
