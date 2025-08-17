package config

import "github.com/IBM/sarama"

// LoggerConfig интерфейс для конфигурации логгера
type LoggerConfig interface {
	Level() string
	AsJson() bool
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

// GRPCClientConfig интерфейс для конфигурации gRPC клиентов
type GRPCClientConfig interface {
	InventoryAddress() string
	PaymentAddress() string
}

// KafkaConfig интерфейс для конфигурации Kafka
type KafkaConfig interface {
	Brokers() []string
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
