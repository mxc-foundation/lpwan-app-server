package data

// SMTPStruct defines smtp service settings
type SMTPStruct struct {
	Email    string `mapstructure:"email"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	AuthType string `mapstructure:"auth_type"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
}

// OperatorStruct defines basic settings of operator of this supernode
type OperatorStruct struct {
	Operator           string `mapstructure:"name"`
	PrimaryColor       string `mapstructure:"primary_color"`
	SecondaryColor     string `mapstructure:"secondary_color"`
	DownloadAppStore   string `mapstructure:"download_appstore"`
	DownloadGoogle     string `mapstructure:"download_google"`
	DownloadTestFlight string `mapstructure:"download_testflight"`
	DownloadAPK        string `mapstructure:"download_apk"`
	OperatorAddress    string `mapstructure:"operator_address"`
	OperatorLegal      string `mapstructure:"operator_legal_name"`
	OperatorLogo       string `mapstructure:"operator_logo"`
	OperatorContact    string `mapstructure:"operator_contact"`
	OperatorSupport    string `mapstructure:"operator_support"`
}

// ServerInfoStruct defines general settings of the server
type ServerInfoStruct struct {
	ServerAddr      string
	DefaultLanguage string
}
