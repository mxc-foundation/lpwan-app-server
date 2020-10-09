package data

import "time"

// IntegrationMQTTConfig holds the configuration for the MQTT integration.
type IntegrationMQTTConfig struct {
	Server               string        `mapstructure:"server"`
	Username             string        `mapstructure:"username"`
	Password             string        `mapstructure:"password"`
	MaxReconnectInterval time.Duration `mapstructure:"max_reconnect_interval"`
	QOS                  uint8         `mapstructure:"qos"`
	CleanSession         bool          `mapstructure:"clean_session"`
	ClientID             string        `mapstructure:"client_id"`
	CACert               string        `mapstructure:"ca_cert"`
	TLSCert              string        `mapstructure:"tls_cert"`
	TLSKey               string        `mapstructure:"tls_key"`
	EventTopicTemplate   string        `mapstructure:"event_topic_template"`
	CommandTopicTemplate string        `mapstructure:"command_topic_template"`
	RetainEvents         bool          `mapstructure:"retain_events"`

	// For backards compatibility
	UplinkTopicTemplate        string `mapstructure:"uplink_topic_template"`
	DownlinkTopicTemplate      string `mapstructure:"downlink_topic_template"`
	JoinTopicTemplate          string `mapstructure:"join_topic_template"`
	AckTopicTemplate           string `mapstructure:"ack_topic_template"`
	ErrorTopicTemplate         string `mapstructure:"error_topic_template"`
	StatusTopicTemplate        string `mapstructure:"status_topic_template"`
	LocationTopicTemplate      string `mapstructure:"location_topic_template"`
	TxAckTopicTemplate         string `mapstructure:"tx_ack_topic_template"`
	IntegrationTopicTemplate   string `mapstructure:"integration_topic_template"`
	UplinkRetainedMessage      bool   `mapstructure:"uplink_retained_message"`
	JoinRetainedMessage        bool   `mapstructure:"join_retained_message"`
	AckRetainedMessage         bool   `mapstructure:"ack_retained_message"`
	ErrorRetainedMessage       bool   `mapstructure:"error_retained_message"`
	StatusRetainedMessage      bool   `mapstructure:"status_retained_message"`
	LocationRetainedMessage    bool   `mapstructure:"location_retained_message"`
	TxAckRetainedMessage       bool   `mapstructure:"tx_ack_retained_message"`
	IntegrationRetainedMessage bool   `mapstructure:"integration_retained_message"`
}
