package data

// ProvisioningServerStruct defines credentails to connect to provisioning-server
type ProvisioningServerStruct struct {
	Server         string `mapstructure:"provision_server"`
	CACert         string `mapstructure:"ca_cert"`
	TLSCert        string `mapstructure:"tls_cert"`
	TLSKey         string `mapstructure:"tls_key"`
	UpdateSchedule string `mapstructure:"update_schedule"`
}
