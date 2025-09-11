package config

import (
	"time"

	"github.com/IBM/sarama"
)

type LoggerConfig interface {
	Level() string
	AsJson() bool
	Outputs() string
	OtelEndpoint() string
	ServiceName() string
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

type MetricsConfig interface {
	CollectorEndpoint() string
	CollectorInterval() time.Duration
}
