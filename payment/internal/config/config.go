package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/kont1n/MSA_Rocket_Factory/payment/internal/config/env"
)

var appConfig *config

type config struct {
	Logger LoggerConfig
	GRPC   GRPCConfig
	Http   HttpConfig
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

	HttpCfg, err := env.NewHttpConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger: loggerCfg,
		GRPC:   GRPCCfg,
		Http:   HttpCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
