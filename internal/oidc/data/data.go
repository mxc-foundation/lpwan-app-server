package data

type UserAuthenticationStruct struct {
	OpenIDConnect struct {
		Enabled                 bool   `mapstructure:"enabled"`
		RegistrationEnabled     bool   `mapstructure:"registration_enabled"`
		RegistrationCallbackURL string `mapstructure:"registration_callback_url"`
		ProviderURL             string `mapstructure:"provider_url"`
		ClientID                string `mapstructure:"client_id"`
		ClientSecret            string `mapstructure:"client_secret"`
		RedirectURL             string `mapstructure:"redirect_url"`
		LogoutURL               string `mapstructure:"logout_url"`
		LoginLabel              string `mapstructure:"login_label"`
	} `mapstructure:"openid_connect"`
}
