package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type orderPaidConsumerEnvConfig struct {
	Topic   string `env:"CONSUMER_ORDER_PAID_TOPIC_NAME,required"`
	GroupID string `env:"CONSUMER_ORDER_PAID_GROUP_ID,required"`
}

type OrderPaidConsumerConfig struct {
	raw orderPaidConsumerEnvConfig
}

func NewOrderPaidConsumerConfig() (*OrderPaidConsumerConfig, error) {
	var raw orderPaidConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &OrderPaidConsumerConfig{raw: raw}, nil
}

func (cfg *OrderPaidConsumerConfig) Topic() string {
	return cfg.raw.Topic
}

func (cfg *OrderPaidConsumerConfig) GroupID() string {
	return cfg.raw.GroupID
}

func (cfg *OrderPaidConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
