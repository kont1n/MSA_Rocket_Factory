package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type assemblyProducerEnvConfig struct {
	TopicName string `env:"PRODUCER_TOPIC_NAME,required"`
}

type assemblyProducerConfig struct {
	raw assemblyProducerEnvConfig
}

func NewAssemblyRecordedProducerConfig() (*assemblyProducerConfig, error) {
	var raw assemblyProducerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &assemblyProducerConfig{raw: raw}, nil
}

func (cfg *assemblyProducerConfig) Topic() string {
	return cfg.raw.TopicName
}

// Config возвращает конфигурацию для sarama consumer
func (cfg *assemblyProducerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Producer.Return.Successes = true

	return config
}
