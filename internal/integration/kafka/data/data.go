package data

// IntegrationKafkaConfig holds the Kafka integration configuration.
type IntegrationKafkaConfig struct {
	Brokers          []string `mapstructure:"brokers"`
	TLS              bool     `mapstructure:"tls"`
	Topic            string   `mapstructure:"topic"`
	EventKeyTemplate string   `mapstructure:"event_key_template"`
	Username         string   `mapstructure:"username"`
	Password         string   `mapstructure:"password"`
}
