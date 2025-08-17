package env

import "github.com/caarlos0/env/v11"

type loggerEnvConfig struct {
	Level  string `env:"LOGGER_LEVEL" envDefault:"info"`
	AsJSON bool   `env:"LOGGER_AS_JSON" envDefault:"false"`
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
