package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type orderPaidProducerEnvConfig struct {
	TopicName string `env:"PRODUCER_TOPIC_NAME,required"`
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
