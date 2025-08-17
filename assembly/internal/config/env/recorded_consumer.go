package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type assemblyRecordedConsumerEnvConfig struct {
	Topic   string `env:"ASSEMBLY_RECORDED_TOPIC_NAME,required"`
	GroupID string `env:"ASSEMBLY_RECORDED_CONSUMER_GROUP_ID,required"`
}

type assemblyRecordedConsumerConfig struct {
	raw assemblyRecordedConsumerEnvConfig
}

func NewAssemblyRecordedConsumerConfig() (*assemblyRecordedConsumerConfig, error) {
	var raw assemblyRecordedConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &assemblyRecordedConsumerConfig{raw: raw}, nil
}

func (cfg *assemblyRecordedConsumerConfig) Topic() string {
	return cfg.raw.Topic
}

func (cfg *assemblyRecordedConsumerConfig) GroupID() string {
	return cfg.raw.GroupID
}

func (cfg *assemblyRecordedConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
