package postgres

import "github.com/kont1n/MSA_Rocket_Factory/platform/pkg/testcontainers"

// Config представляет конфигурацию PostgreSQL контейнера
type Config struct {
	ContainerName string
	ImageName     string
	Database      string
	Username      string
	Password      string
	Port          string
}

// NewConfig создает новую конфигурацию PostgreSQL с значениями по умолчанию
func NewConfig() *Config {
	return &Config{
		ContainerName: testcontainers.PostgresContainerName,
		ImageName:     testcontainers.PostgresImageName,
		Database:      testcontainers.PostgresDatabase,
		Username:      testcontainers.PostgresUsername,
		Password:      testcontainers.PostgresPassword,
		Port:          testcontainers.PostgresPort,
	}
}
