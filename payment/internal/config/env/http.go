package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type HttpEnvConfig struct {
	Host string `env:"HTTP_HOST,required"`
	Port string `env:"HTTP_PORT,required"`
}

type HttpConfig struct {
	raw HttpEnvConfig
}

func NewHttpConfig() (*HttpConfig, error) {
	var raw HttpEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &HttpConfig{raw: raw}, nil
}

func (cfg *HttpConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
