package config

type ClientConfig struct {
	SelfKey    string `yaml:"self_key"`
	Nonce      string `yaml:"nonce"`
	DeviceName string `yaml:"device_name"`
	Addrs      Addrs  `yaml:"addrs"`
}
