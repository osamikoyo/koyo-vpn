package config

type ClientConfig struct {
	SelfKey       string `yaml:"self_key"`
	DeviceName    string `yaml:"device_name"`
	ServerUDPAddr string `yaml:"server_udp_addr"`
}
