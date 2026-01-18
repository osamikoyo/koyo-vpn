package config

import (
	"fmt"
	"os"

	"koyo-vpn/pkg/logger"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type Keys struct {
	Self   string `yaml:"self"`
	Remote string `yaml:"remote"`
}

type ServerConfig struct {
	DeviceName string `yaml:"device_name"`

	Addrs Addrs `yaml:"addrs"`
	Keys  Keys  `yaml:"keys"`
}

func LoadServerConfig(logger *logger.Logger, path string) (*ServerConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		logger.Error("failed open config path",
			zap.String("path", path),
			zap.Error(err))

		return nil, fmt.Errorf("failed open config file: %s: %w", path, err)
	}

	var cfg ServerConfig

	if err = yaml.NewDecoder(file).Decode(&cfg); err != nil {
		logger.Error("failed load server config",
			zap.String("path", path),
			zap.Error(err))

		return nil, fmt.Errorf("failed load server config: %w", err)
	}

	return nil, nil
}
