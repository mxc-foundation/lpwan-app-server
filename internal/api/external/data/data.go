package data

type ExternalAPIStruct struct {
	Bind            string
	TLSCert         string `mapstructure:"tls_cert"`
	TLSKey          string `mapstructure:"tls_key"`
	JWTSecret       string `mapstructure:"jwt_secret"`
	OTPSecret       string `mapstructure:"otp_secret"`
	CORSAllowOrigin string `mapstructure:"cors_allow_origin"`
}
