package data

// AzurePublishMode defines the publish-mode type.
type AzurePublishMode string

// IntegrationAzureConfig holds the Azure Service-Bus integration configuration.
type IntegrationAzureConfig struct {
	Marshaler        string           `mapstructure:"marshaler" json:"marshaler"`
	ConnectionString string           `mapstructure:"connection_string" json:"connectionString"`
	PublishMode      AzurePublishMode `mapstructure:"publish_mode" json:"-"`
	PublishName      string           `mapstructure:"publish_name" json:"publishName"`
}
