package config

type Keys struct {
	Self   string `yaml:"self"`
	Remote string `yaml:"remote"`
}

type ServerConfig struct {
	Nonce string `yaml:"nonce"`
	DeviceName string `yaml:"device_nane"`
	
	Addrs Addrs `yaml:"addr"`
	Keys       Keys   `yaml:"keys"`
}
