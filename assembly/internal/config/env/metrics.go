package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type metricsEnvConfig struct {
	CollectorEndpoint string        `env:"OTEL_COLLECTOR_ENDPOINT" envDefault:"otel-collector:4317"`
	CollectorInterval time.Duration `env:"METRICS_COLLECTOR_INTERVAL" envDefault:"10s"`
}

type metricsConfig struct {
	raw metricsEnvConfig
}

func NewMetricsConfig() (*metricsConfig, error) {
	var raw metricsEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &metricsConfig{raw: raw}, nil
}

func (cfg *metricsConfig) CollectorEndpoint() string {
	return cfg.raw.CollectorEndpoint
}

func (cfg *metricsConfig) CollectorInterval() time.Duration {
	return cfg.raw.CollectorInterval
}
