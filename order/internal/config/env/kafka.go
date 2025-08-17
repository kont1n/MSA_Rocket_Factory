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
	config.Version = sarama.V4_0_0_0
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	return config
}

type orderPaidProducerEnvConfig struct {
	TopicName string `env:"KAFKA_ORDER_PAID_TOPIC,required"`
}

type OrderPaidProducerConfig struct {
	raw orderPaidProducerEnvConfig
}

func NewOrderPaidProducerConfig() (*OrderPaidProducerConfig, error) {
	var raw orderPaidProducerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &OrderPaidProducerConfig{raw: raw}, nil
}

func (cfg *OrderPaidProducerConfig) Topic() string {
	return cfg.raw.TopicName
}

func (cfg *OrderPaidProducerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	return config
}

type shipAssembledConsumerEnvConfig struct {
	Topic   string `env:"KAFKA_SHIP_ASSEMBLED_TOPIC,required"`
	GroupID string `env:"KAFKA_SHIP_ASSEMBLED_GROUP_ID,required"`
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
