package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type mongoEnvConfig struct {
	Host     string `env:"MONGO_HOST,required"`
	Port     string `env:"MONGO_PORT,required"`
	Database string `env:"MONGO_DATABASE,required"`
	User     string `env:"MONGO_INITDB_ROOT_USERNAME,required"`
	Password string `env:"MONGO_INITDB_ROOT_PASSWORD,required"`
	AuthDB   string `env:"MONGO_AUTH_DB,required"`
}

type mongoConfig struct {
	raw mongoEnvConfig
}

func NewMongoConfig() (*mongoConfig, error) {
	var raw mongoEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &mongoConfig{raw: raw}, nil
}

func (cfg *mongoConfig) URI() string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?authSource=%s",
		cfg.raw.User,
		cfg.raw.Password,
		cfg.raw.Host,
		cfg.raw.Port,
		cfg.raw.Database,
		cfg.raw.AuthDB,
	)
}

func (cfg *mongoConfig) DatabaseName() string {
	return cfg.raw.Database
}
