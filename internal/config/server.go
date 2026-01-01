package config

type Keys struct {
	Self   string `yaml:"self"`
	Client string `yaml:"client"`
}

type ServerConfig struct {
	UDPAddr    string `yaml:"udp_addr"`
	DeviceName string `yaml:"device_nane"`
	Keys       Keys   `yaml:"keys"`
}
