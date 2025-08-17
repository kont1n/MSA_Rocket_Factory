package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/kont1n/MSA_Rocket_Factory/assembly/internal/config/env"
)

var appConfig *config

type config struct {
	Logger                   LoggerConfig
	Kafka                    KafkaConfig
	AssemblyRecordedProducer AssemblyRecordedProducerConfig
	AssemblyRecordedConsumer AssemblyRecordedConsumerConfig
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

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	assemblyRecordedProducerCfg, err := env.NewAssemblyRecordedProducerConfig()
	if err != nil {
		return err
	}

	assemblyRecordedConsumerCfg, err := env.NewAssemblyRecordedConsumerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:                   loggerCfg,
		Kafka:                    kafkaCfg,
		AssemblyRecordedProducer: assemblyRecordedProducerCfg,
		AssemblyRecordedConsumer: assemblyRecordedConsumerCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
