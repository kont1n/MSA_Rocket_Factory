package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type kafkaEnvConfig struct {
	Brokers []string `env:"KAFKA_BROKERS,required"`
}

type KafkaConfig struct {
	raw kafkaEnvConfig
}

func NewKafkaConfig() (*KafkaConfig, error) {
	var raw kafkaEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &KafkaConfig{raw: raw}, nil
}

func (cfg *KafkaConfig) Brokers() []string {
	return cfg.raw.Brokers
}

func (cfg *KafkaConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	return config
}
