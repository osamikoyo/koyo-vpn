package config

type Keys struct {
	Self   string `yaml:"self"`
	Remote string `yaml:"remote"`
}

type ServerConfig struct {
	SelfUDPAddr   string `yaml:"self_udp_addr"`
	RemoteUDPAddr string `yaml:"remote_udp_addr"`

	Nonce string `yaml:"nonce"`

	DeviceName string `yaml:"device_nane"`
	Keys       Keys   `yaml:"keys"`
}
