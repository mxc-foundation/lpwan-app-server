package data

// IntegrationAMQPConfig holds the AMQP integration configuration.
type IntegrationAMQPConfig struct {
	URL                     string `mapstructure:"url"`
	EventRoutingKeyTemplate string `mapstructure:"event_routing_key_template"`
}
