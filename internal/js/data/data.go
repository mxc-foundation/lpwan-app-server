package data

type JoinServerStruct struct {
	Bind    string
	CACert  string `mapstructure:"ca_cert"`
	TLSCert string `mapstructure:"tls_cert"`
	TLSKey  string `mapstructure:"tls_key"`

	KEK struct {
		ASKEKLabel string `mapstructure:"as_kek_label"`

		Set []struct {
			Label string `mapstructure:"label"`
			KEK   string `mapstructure:"kek"`
		}
	} `mapstructure:"kek"`
}
