package config

import (
	"time"

	"github.com/IBM/sarama"
)

// LoggerConfig интерфейс для конфигурации логгера
type LoggerConfig interface {
	Level() string
	AsJson() bool
	Outputs() string
	OtelEndpoint() string
	ServiceName() string
}

// HTTPConfig интерфейс для конфигурации HTTP сервера
type HTTPConfig interface {
	Address() string
	ReadHeaderTimeout() int
	ShutdownTimeout() int
}

// DBConfig интерфейс для конфигурации базы данных
type DBConfig interface {
	URI() string
	MigrationsDir() string
}

// KafkaConfig интерфейс для конфигурации Kafka
type KafkaConfig interface {
	Brokers() []string
	Config() *sarama.Config
}

// GRPCClientConfig интерфейс для конфигурации gRPC клиентов
type GRPCClientConfig interface {
	InventoryAddress() string
	PaymentAddress() string
	IAMAddress() string
}

// OrderPaidProducerConfig интерфейс для конфигурации Kafka producer
type OrderPaidProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}

// ShipAssemblyConsumerConfig интерфейс для конфигурации Kafka consumer
type ShipAssemblyConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}

// MetricsConfig интерфейс для конфигурации метрик
type MetricsConfig interface {
	CollectorEndpoint() string
	CollectorInterval() time.Duration
}

// TracingConfig интерфейс для конфигурации трассировки
type TracingConfig interface {
	CollectorEndpoint() string
	ServiceName() string
	Environment() string
	ServiceVersion() string
}
