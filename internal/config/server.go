package config

type Keys struct {
	Self   string `yaml:"self"`
	Remote string `yaml:"remote"`
}

type ServerConfig struct {
	DeviceName string `yaml:"device_name"`

	Addrs Addrs `yaml:"addrs"`
	Keys  Keys  `yaml:"keys"`
}
