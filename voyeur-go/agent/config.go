package agent

type Config struct {
	// DevelMode changes defaults to be development friendly: debug and human
	// readable logging, etc.
	DevelMode bool        `yaml:"devel-mode,omitempty"`
	TribeID   string      `yaml:"tribe-id"`
	REST      *RestConfig `yaml:"rest,omitempty"`
	MQTT      *MqttConfig `yaml:"mqtt,omitempty"`
}

type TlsProfile struct {
	CA   string `yaml:"ca,omitempty"`
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
}

type RestConfig struct {
	Port int         `yaml:"port"`
	TLS  *TlsProfile `yaml:"tls,omitempty"`
}

type MqttConfig struct {
	Port int         `yaml:"port"`
	TLS  *TlsProfile `yaml:"tls,omitempty"`
}
