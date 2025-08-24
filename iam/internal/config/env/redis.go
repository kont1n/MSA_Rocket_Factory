package env

import (
	"net"
	"time"

	"github.com/caarlos0/env/v11"
)

type redisEnvConfig struct {
	Host              string        `env:"REDIS_HOST,required"`
	Port              string        `env:"REDIS_PORT,required"`
	ConnectionTimeout time.Duration `env:"REDIS_CONNECTION_TIMEOUT,required"`
	MaxIdle           int           `env:"REDIS_MAX_IDLE,required"`
	IdleTimeout       time.Duration `env:"REDIS_IDLE_TIMEOUT,required"`
	CacheTTL          time.Duration `env:"SESSION_TTL,required"`
}

type redisConfig struct {
	raw redisEnvConfig
}

func NewRedisConfig() (*redisConfig, error) {
	var raw redisEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &redisConfig{raw: raw}, nil
}

func (cfg *redisConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}

func (cfg *redisConfig) ConnectionTimeout() time.Duration {
	return cfg.raw.ConnectionTimeout
}

func (cfg *redisConfig) MaxIdle() int {
	return cfg.raw.MaxIdle
}

func (cfg *redisConfig) IdleTimeout() time.Duration {
	return cfg.raw.IdleTimeout
}

func (cfg *redisConfig) CacheTTL() time.Duration {
	return cfg.raw.CacheTTL
}
