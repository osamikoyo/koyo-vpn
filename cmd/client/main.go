package main

import (
	"koyo-vpn/internal/config"
	"koyo-vpn/internal/core"
	"koyo-vpn/pkg/logger"
	"os"
)

func main() {
	logger.Init(logger.Config{
		AppName:   "vpn-client",
		AddCaller: false,
		LogFile:   "vpn-client.log",
		LogLevel:  "debug",
	})

	logger := logger.Get()

	logger.Info("setup client...")

	cfg, err := getCfg(logger)
	if err != nil{
		return
	}

	

}

func getCfg(logger *logger.Logger) (*config.ClientConfig, error) {
	path := "client.yaml"

	for i, arg := range os.Args {
		if arg == "--config" {
			path = os.Args[i+1]
		}
	}

	cfg, err := config.NewConfig[*config.ClientConfig]("client", path)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
