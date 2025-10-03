package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type jwtEnvConfig struct {
	AccessTokenSecret  string        `env:"JWT_ACCESS_TOKEN_SECRET,required"`
	RefreshTokenSecret string        `env:"JWT_REFRESH_TOKEN_SECRET,required"`
	AccessTokenTTL     time.Duration `env:"JWT_ACCESS_TOKEN_TTL" envDefault:"15m"`
	RefreshTokenTTL    time.Duration `env:"JWT_REFRESH_TOKEN_TTL" envDefault:"24h"`
}

// JWTConfig содержит конфигурацию для JWT токенов
type JWTConfig interface {
	AccessTokenSecret() string
	RefreshTokenSecret() string
	AccessTokenTTL() time.Duration
	RefreshTokenTTL() time.Duration
}

type jwtConfig struct {
	raw jwtEnvConfig
}

// NewJWTConfig создает новую конфигурацию JWT
func NewJWTConfig() (JWTConfig, error) {
	var raw jwtEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &jwtConfig{raw: raw}, nil
}

func (cfg *jwtConfig) AccessTokenSecret() string {
	return cfg.raw.AccessTokenSecret
}

func (cfg *jwtConfig) RefreshTokenSecret() string {
	return cfg.raw.RefreshTokenSecret
}

func (cfg *jwtConfig) AccessTokenTTL() time.Duration {
	return cfg.raw.AccessTokenTTL
}

func (cfg *jwtConfig) RefreshTokenTTL() time.Duration {
	return cfg.raw.RefreshTokenTTL
}
