package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type shipAssembledConsumerEnvConfig struct {
	Topic   string `env:"CONSUMER_SHIP_ASSEMBLED_TOPIC_NAME,required"`
	GroupID string `env:"CONSUMER_SHIP_ASSEMBLED_GROUP_ID,required"`
}

type ShipAssembledConsumerConfig struct {
	raw shipAssembledConsumerEnvConfig
}

func NewShipAssembledConsumerConfig() (*ShipAssembledConsumerConfig, error) {
	var raw shipAssembledConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &ShipAssembledConsumerConfig{raw: raw}, nil
}

func (cfg *ShipAssembledConsumerConfig) Topic() string {
	return cfg.raw.Topic
}

func (cfg *ShipAssembledConsumerConfig) GroupID() string {
	return cfg.raw.GroupID
}

func (cfg *ShipAssembledConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}