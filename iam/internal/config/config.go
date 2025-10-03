package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/config/env"
)

var appConfig *config

type config struct {
	Logger LoggerConfig
	GRPC   GRPCConfig
	DB     DBConfig
	Redis  RedisConfig
	Token  TokenConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	GRPCCfg, err := env.NewGRPCConfig()
	if err != nil {
		return err
	}

	dbCfg, err := env.NewDBConfig()
	if err != nil {
		return err
	}

	redisCfg, err := env.NewRedisConfig()
	if err != nil {
		return err
	}

	jwtCfg, err := env.NewJWTConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger: loggerCfg,
		GRPC:   GRPCCfg,
		DB:     dbCfg,
		Redis:  redisCfg,
		Token:  jwtCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
