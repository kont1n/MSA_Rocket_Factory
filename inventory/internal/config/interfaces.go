package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type GRPCConfig interface {
	Address() string
}

type MongoConfig interface {
	URI() string
	DatabaseName() string
}

// GRPCClientConfig интерфейс для конфигурации gRPC клиентов
type GRPCClientConfig interface {
	IAMAddress() string
}
