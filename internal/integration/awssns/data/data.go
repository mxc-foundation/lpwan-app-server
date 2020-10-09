package data

// IntegrationAWSSNSConfig holds the AWS SNS integration configuration.
type IntegrationAWSSNSConfig struct {
	Marshaler          string `mapstructure:"marshaler" json:"marshaler"`
	AWSRegion          string `mapstructure:"aws_region" json:"region"`
	AWSAccessKeyID     string `mapstructure:"aws_access_key_id" json:"accessKeyID"`
	AWSSecretAccessKey string `mapstructure:"aws_secret_access_key" json:"secretAccessKey"`
	TopicARN           string `mapstructure:"topic_arn" json:"topicARN"`
}
