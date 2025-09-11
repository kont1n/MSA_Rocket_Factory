package env

import "github.com/caarlos0/env/v11"

type loggerEnvConfig struct {
	Level        string `env:"LOGGER_LEVEL" envDefault:"info"`
	AsJSON       bool   `env:"LOGGER_AS_JSON" envDefault:"false"`
	Outputs      string `env:"LOG_OUTPUTS" envDefault:"stdout"`
	OtelEndpoint string `env:"OTEL_COLLECTOR_ENDPOINT" envDefault:"otel-collector:4317"`
	ServiceName  string `env:"SERVICE_NAME" envDefault:"notification"`
}

type LoggerConfig struct {
	raw loggerEnvConfig
}

func NewLoggerConfig() (*LoggerConfig, error) {
	var raw loggerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &LoggerConfig{raw: raw}, nil
}

func (cfg *LoggerConfig) Level() string {
	return cfg.raw.Level
}

func (cfg *LoggerConfig) AsJson() bool {
	return cfg.raw.AsJSON
}

func (cfg *LoggerConfig) Outputs() string {
	return cfg.raw.Outputs
}

func (cfg *LoggerConfig) OtelEndpoint() string {
	return cfg.raw.OtelEndpoint
}

func (cfg *LoggerConfig) ServiceName() string {
	return cfg.raw.ServiceName
}
