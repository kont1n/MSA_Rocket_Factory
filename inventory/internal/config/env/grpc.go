package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type GRPCEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type GRPCConfig struct {
	raw GRPCEnvConfig
}

func NewGRPCConfig() (*GRPCConfig, error) {
	var raw GRPCEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &GRPCConfig{raw: raw}, nil
}

func (cfg *GRPCConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
