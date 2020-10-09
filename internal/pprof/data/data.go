package data

type Config struct {
	Enabled bool   `mapstructure:"enabled"`
	Bind    string `mapstructure:"bind"`
}
