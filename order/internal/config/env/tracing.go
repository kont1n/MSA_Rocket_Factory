package env

import (
	"github.com/caarlos0/env/v11"
)

type tracingEnvConfig struct {
	CollectorEndpoint string `env:"OTEL_COLLECTOR_ENDPOINT" envDefault:"otel-collector:4317"`
	ServiceName       string `env:"SERVICE_NAME" envDefault:"order"`
	Environment       string `env:"ENVIRONMENT" envDefault:"dev"`
	ServiceVersion    string `env:"SERVICE_VERSION" envDefault:"1.0.0"`
}

type tracingConfig struct {
	raw tracingEnvConfig
}

func NewTracingConfig() (*tracingConfig, error) {
	var raw tracingEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &tracingConfig{raw: raw}, nil
}

func (cfg *tracingConfig) CollectorEndpoint() string {
	return cfg.raw.CollectorEndpoint
}

func (cfg *tracingConfig) ServiceName() string {
	return cfg.raw.ServiceName
}

func (cfg *tracingConfig) Environment() string {
	return cfg.raw.Environment
}

func (cfg *tracingConfig) ServiceVersion() string {
	return cfg.raw.ServiceVersion
}
