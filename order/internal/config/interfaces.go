package config

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
