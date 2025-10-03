package config

import "time"

// Константы для конфигурации производительности
const (
	// Cache TTL
	DefaultSessionCacheTTL = 24 * time.Hour
	UserCacheTTL           = 30 * time.Minute

	// Connection pools
	DefaultMaxDBConnections     = 25
	DefaultMaxIdleDBConnections = 5
	DefaultDBConnectionTimeout  = 30 * time.Second

	// Redis
	DefaultRedisMaxIdle        = 10
	DefaultRedisMaxActive      = 100
	DefaultRedisIdleTimeout    = 240 * time.Second
	DefaultRedisConnectTimeout = 10 * time.Second
	DefaultRedisReadTimeout    = 10 * time.Second
	DefaultRedisWriteTimeout   = 10 * time.Second

	// Rate limiting
	DefaultRateLimitRequests = 5
	DefaultRateLimitWindow   = 1 * time.Minute

	// Sessions
	DefaultSessionCleanupInterval = 1 * time.Hour

	// Argon2 parameters - оптимизированы для безопасности и производительности
	Argon2Time    = 3         // Увеличено для лучшей безопасности
	Argon2Memory  = 64 * 1024 // 64 MB
	Argon2Threads = 4
	Argon2KeyLen  = 32
	Argon2SaltLen = 16

	// JWT
	JWTSigningMethod = "HS256"

	// Batch sizes
	DefaultBatchSize = 100
)

// Performance settings
var (
	// Timeouts
	DefaultRequestTimeout  = 30 * time.Second
	DefaultShutdownTimeout = 5 * time.Second
)
