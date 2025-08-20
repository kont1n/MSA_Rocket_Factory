package config

// LoggerConfig интерфейс для конфигурации логгера
type LoggerConfig interface {
	Level() string
	AsJson() bool
}

// GRPCConfig интерфейс для конфигурации gRPC сервера
type GRPCConfig interface {
	Address() string
}

// DBConfig интерфейс для конфигурации базы данных
type DBConfig interface {
	URI() string
	MigrationsDir() string
}
