package env

import (
	"os"
)

type dbConfig struct {
	uri           string
	migrationsDir string
}

func NewDBConfig() (*dbConfig, error) {
	uri := os.Getenv("DB_URI")
	if uri == "" {
		uri = "postgres://user:password@localhost:5432/order_db?sslmode=disable"
	}

	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		migrationsDir = "./migrations"
	}

	return &dbConfig{
		uri:           uri,
		migrationsDir: migrationsDir,
	}, nil
}

func (c *dbConfig) URI() string {
	return c.uri
}

func (c *dbConfig) MigrationsDir() string {
	return c.migrationsDir
}
