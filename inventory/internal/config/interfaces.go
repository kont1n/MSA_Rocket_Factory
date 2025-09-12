package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
	Outputs() string
	OtelEndpoint() string
	ServiceName() string
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

type TracingConfig interface {
	CollectorEndpoint() string
	ServiceName() string
	Environment() string
	ServiceVersion() string
}
