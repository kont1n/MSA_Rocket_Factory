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

type HttpConfig interface {
	Address() string
}

type TracingConfig interface {
	CollectorEndpoint() string
	ServiceName() string
	Environment() string
	ServiceVersion() string
}
