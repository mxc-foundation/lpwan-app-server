package data

type ExternalAPIStruct struct {
	Bind            string `mapstructure:"bind"`
	TLSCert         string `mapstructure:"tls_cert"`
	TLSKey          string `mapstructure:"tls_key"`
	JWTSecret       string `mapstructure:"jwt_secret"`
	JWTDefaultTTL   int64  `mapstructure:"jwt_default_ttl_sec"`
	OTPSecret       string `mapstructure:"otp_secret"`
	CORSAllowOrigin string `mapstructure:"cors_allow_origin"`
}
