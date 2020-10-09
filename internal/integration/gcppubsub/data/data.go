package data

// IntegrationGCPConfig holds the GCP Pub/Sub integration configuration.
type IntegrationGCPConfig struct {
	Marshaler            string `mapstructure:"marshaler" json:"marshaler"`
	CredentialsFile      string `mapstructure:"credentials_file" json:"-"`
	CredentialsFileBytes []byte `mapstructure:"-" json:"credentialsFile"`
	ProjectID            string `mapstructure:"project_id" json:"projectID"`
	TopicName            string `mapstructure:"topic_name" json:"topicName"`
}
