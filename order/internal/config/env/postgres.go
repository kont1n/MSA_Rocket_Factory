package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type postgresEnvConfig struct {
	Host          string `env:"POSTGRES_HOST,required"`
	Port          string `env:"POSTGRES_PORT,required"`
	SslMode       string `env:"POSTGRES_SSLMODE,required"`
	Database      string `env:"POSTGRES_DATABASE,required"`
	User          string `env:"POSTGRES_INITDB_ROOT_USERNAME,required"`
	Password      string `env:"POSTGRES_INITDB_ROOT_PASSWORD,required"`
	MigrationsDir string `env:"POSTGRES_MIGRATIONS_DIR,required"`
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

func (cfg *postgresConfig) MigrationsDir() string {
	return cfg.raw.MigrationsDir
}
