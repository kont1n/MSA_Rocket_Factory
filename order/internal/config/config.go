package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/config/env"
)

var appConfig *config

type config struct {
	Logger                LoggerConfig
	HTTP                  HTTPConfig
	DB                    DBConfig
	GRPCClient            GRPCClientConfig
	Kafka                 KafkaConfig
	OrderPaidProducer     OrderPaidProducerConfig
	ShipAssembledConsumer ShipAssemblyConsumerConfig
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

	httpCfg, err := env.NewHTTPConfig()
	if err != nil {
		return err
	}

	dbCfg, err := env.NewDBConfig()
	if err != nil {
		return err
	}

	grpcClientCfg, err := env.NewGRPCClientConfig()
	if err != nil {
		return err
	}

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	orderPaidProducerCfg, err := env.NewOrderPaidProducerConfig()
	if err != nil {
		return err
	}

	shipAssembledConsumerCfg, err := env.NewShipAssembledConsumerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:                loggerCfg,
		HTTP:                  httpCfg,
		DB:                    dbCfg,
		GRPCClient:            grpcClientCfg,
		Kafka:                 kafkaCfg,
		OrderPaidProducer:     orderPaidProducerCfg,
		ShipAssembledConsumer: shipAssembledConsumerCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
