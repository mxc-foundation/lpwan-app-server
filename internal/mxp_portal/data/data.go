package data

// MxprotocolServerStruct defines credentails to connect to mxprotocol-server
type MxprotocolServerStruct struct {
	Server  string `mapstructure:"m2m_server"`
	CACert  string `mapstructure:"ca_cert"`
	TLSCert string `mapstructure:"tls_cert"`
	TLSKey  string `mapstructure:"tls_key"`
}

type MxprotocolClientStruct struct {
	Bind    string `mapstructure:"bind"`
	CACert  string `mapstructure:"ca_cert"`
	TLSCert string `mapstructure:"tls_cert"`
	TLSKey  string `mapstructure:"tls_key"`
}
