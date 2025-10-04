package env

import (
	"github.com/caarlos0/env/v11"
)

type loggerEnvConfig struct {
	Level        string `env:"LOGGER_LEVEL,required"`
	AsJson       bool   `env:"LOGGER_AS_JSON,required"`
	Outputs      string `env:"LOG_OUTPUTS" envDefault:"stdout"`
	OtelEndpoint string `env:"OTEL_COLLECTOR_ENDPOINT" envDefault:"otel-collector:4317"`
	ServiceName  string `env:"SERVICE_NAME" envDefault:"iam"`
}

type loggerConfig struct {
	raw loggerEnvConfig
}

func NewLoggerConfig() (*loggerConfig, error) {
	var raw loggerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &loggerConfig{raw: raw}, nil
}

func (cfg *loggerConfig) Level() string {
	return cfg.raw.Level
}

func (cfg *loggerConfig) AsJson() bool {
	return cfg.raw.AsJson
}

func (cfg *loggerConfig) Outputs() string {
	return cfg.raw.Outputs
}

func (cfg *loggerConfig) OtelEndpoint() string {
	return cfg.raw.OtelEndpoint
}

func (cfg *loggerConfig) ServiceName() string {
	return cfg.raw.ServiceName
}
