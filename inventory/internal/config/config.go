package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/config/env"
)

var appConfig *config

type config struct {
	Logger     LoggerConfig
	GRPC       GRPCConfig
	Mongo      MongoConfig
	GRPCClient GRPCClientConfig
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

	mongoCfg, err := env.NewMongoConfig()
	if err != nil {
		return err
	}

	GRPCClientCfg, err := env.NewGRPCClientConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:     loggerCfg,
		GRPC:       GRPCCfg,
		Mongo:      mongoCfg,
		GRPCClient: GRPCClientCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
