//go:build integration

package postgres_test

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/repository"
)

type PostgresRepositorySuite struct {
	suite.Suite
	repository repository.OrderRepository
	pool       *pgxpool.Pool
}

func (s *PostgresRepositorySuite) SetupSuite() {
	// Для интеграционных тестов с PostgreSQL потребуется настройка testcontainers
	// В текущей реализации используем заглушку - тесты будут пропущены
	// TODO: настроить testcontainers для PostgreSQL
}

func (s *PostgresRepositorySuite) SetupTest() {
	// Заглушка для тестов
}

func (s *PostgresRepositorySuite) TearDownSuite() {
	// Заглушка для тестов
}

func TestPostgresRepositoryIntegration(t *testing.T) {
	suite.Run(t, new(PostgresRepositorySuite))
}
