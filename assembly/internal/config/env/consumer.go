package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type assemblyConsumerEnvConfig struct {
	Topic   string `env:"CONSUMER_TOPIC_NAME,required"`
	GroupID string `env:"CONSUMER_GROUP_ID,required"`
}

type assemblyConsumerConfig struct {
	raw assemblyConsumerEnvConfig
}

func NewAssemblyRecordedConsumerConfig() (*assemblyConsumerConfig, error) {
	var raw assemblyConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &assemblyConsumerConfig{raw: raw}, nil
}

func (cfg *assemblyConsumerConfig) Topic() string {
	return cfg.raw.Topic
}

func (cfg *assemblyConsumerConfig) GroupID() string {
	return cfg.raw.GroupID
}

func (cfg *assemblyConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
