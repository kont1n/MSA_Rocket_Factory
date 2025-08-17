package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type assemblyRecordedProducerEnvConfig struct {
	TopicName string `env:"ASSEMBLY_RECORDED_TOPIC_NAME,required"`
}

type assemblyRecordedProducerConfig struct {
	raw assemblyRecordedProducerEnvConfig
}

func NewAssemblyRecordedProducerConfig() (*assemblyRecordedProducerConfig, error) {
	var raw assemblyRecordedProducerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &assemblyRecordedProducerConfig{raw: raw}, nil
}

func (cfg *assemblyRecordedProducerConfig) Topic() string {
	return cfg.raw.TopicName
}

// Config возвращает конфигурацию для sarama consumer
func (cfg *assemblyRecordedProducerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Producer.Return.Successes = true

	return config
}
