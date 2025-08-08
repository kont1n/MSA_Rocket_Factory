package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type GRPCConfig interface {
	Address() string
}

type HttpConfig interface {
	Address() string
}
