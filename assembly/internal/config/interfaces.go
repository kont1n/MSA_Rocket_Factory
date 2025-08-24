package config

import "github.com/IBM/sarama"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type KafkaConfig interface {
	Brokers() []string
}

type AssemblyProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}

type AssemblyConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}
