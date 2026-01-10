package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Addrs struct{
	Self string `yaml:"self"`
	Remote string `yaml:"remote"`
}

type GenericConfig interface {
	*ClientConfig | *ServerConfig
}

func NewConfig[T GenericConfig](cType, path string) (T, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed open config file: %w", err)
	}
	defer file.Close()

	var zero T

	switch any(zero).(type) {
	case *ServerConfig:
		if cType != "server" {
			return zero, fmt.Errorf("expected server config")
		}
		cfg := new(ServerConfig)
		if err := yaml.NewDecoder(file).Decode(cfg); err != nil {
			return zero, fmt.Errorf("failed decode: %w", err)
		}
		return any(cfg).(T), nil

	case *ClientConfig:
		if cType != "client" {
			return zero, fmt.Errorf("expected client config")
		}
		cfg := new(ClientConfig)
		if err := yaml.NewDecoder(file).Decode(cfg); err != nil {
			return zero, fmt.Errorf("failed decode: %w", err)
		}
		return any(cfg).(T), nil

	default:
		return zero, fmt.Errorf("unsupported config type")
	}
}
