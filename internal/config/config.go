package config

import (
	"log"

	env "github.com/caarlos0/env/v11"
)

type Config struct {
	Home       string `env:"HOME"`
	SeJWTToken string `env:"SE_JWT_TOKEN"`
	DevicePath string `env:"DEVICE_PATH" envDefault:"/dev/usb/lp0"` // printer device path
	ServerPort string `env:"SERVER_PORT" envDefault:":8082"`        // server port
}

func New() *Config {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	return cfg
}
