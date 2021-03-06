package data

type AppserverStruct struct {
	Bind       string `mapstructure:"bind"`
	CACert     string `mapstructure:"ca_cert"`
	TLSCert    string `mapstructure:"tls_cert"`
	TLSKey     string `mapstructure:"tls_key"`
	PublicHost string `mapstructure:"public_host"`
}
