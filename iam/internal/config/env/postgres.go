package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type postgresEnvConfig struct {
	Host          string `env:"POSTGRES_HOST,required"`
	Port          string `env:"POSTGRES_PORT,required"`
	SslMode       string `env:"POSTGRES_SSL_MODE,required"`
	Database      string `env:"POSTGRES_DB,required"`
	User          string `env:"POSTGRES_USER,required"`
	Password      string `env:"POSTGRES_PASSWORD,required"`
	MigrationsDir string `env:"MIGRATION_DIRECTORY,required"`
}

type postgresConfig struct {
	raw postgresEnvConfig
}

func NewDBConfig() (*postgresConfig, error) {
	var raw postgresEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &postgresConfig{raw: raw}, nil
}

func (cfg *postgresConfig) URI() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.raw.User,
		cfg.raw.Password,
		cfg.raw.Host,
		cfg.raw.Port,
		cfg.raw.Database,
		cfg.raw.SslMode,
	)
}

// SafeURI возвращает строку подключения с замаскированным паролем для логирования
func (cfg *postgresConfig) SafeURI() string {
	return fmt.Sprintf(
		"postgres://%s:***@%s:%s/%s?sslmode=%s",
		cfg.raw.User,
		cfg.raw.Host,
		cfg.raw.Port,
		cfg.raw.Database,
		cfg.raw.SslMode,
	)
}

func (cfg *postgresConfig) MigrationsDir() string {
	return cfg.raw.MigrationsDir
}
