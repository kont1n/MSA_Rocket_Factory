package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type GRPCEnvConfig struct {
	Host     string `env:"GRPC_HOST,required"`
	Port     string `env:"GRPC_PORT,required"`
	TLSCert  string `env:"GRPC_TLS_CERT_FILE"`
	TLSKey   string `env:"GRPC_TLS_KEY_FILE"`
	Insecure bool   `env:"GRPC_INSECURE" envDefault:"true"` // По умолчанию true для dev среды
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

// TLS конфигурация
func (cfg *GRPCConfig) TLSCertFile() string {
	return cfg.raw.TLSCert
}

func (cfg *GRPCConfig) TLSKeyFile() string {
	return cfg.raw.TLSKey
}

func (cfg *GRPCConfig) IsInsecure() bool {
	return cfg.raw.Insecure
}
