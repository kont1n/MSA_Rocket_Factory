package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type GRPCConfig interface {
	Address() string
	TLSCertFile() string
	TLSKeyFile() string
	IsInsecure() bool
}

type DBConfig interface {
	URI() string
	SafeURI() string
	MigrationsDir() string
}

type RedisConfig interface {
	Address() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
	CacheTTL() time.Duration
}

type TokenConfig interface {
	AccessTokenSecret() string
	RefreshTokenSecret() string
	AccessTokenTTL() time.Duration
	RefreshTokenTTL() time.Duration
}
