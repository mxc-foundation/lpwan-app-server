package data

type GeneralSettingsStruct struct {
	LogLevel               int    `mapstructure:"log_level"`
	LogToSyslog            bool   `mapstructure:"log_to_syslog"`
	PasswordHashIterations int    `mapstructure:"password_hash_iterations"`
	Enable2FALogin         bool   `mapstructure:"enable_2fa_login"`
	DefaultLanguage        string `mapstructure:"defualt_language"`
	ServerAddr             string `mapstructure:"server_addr"`
	ServerRegion           string `mapstructure:"server_region"`
	EnableSTC              bool   `mapstructure:"enable_stc"`
}
