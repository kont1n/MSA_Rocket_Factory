//go:build integration

package mongo_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/repository"
)

type MongoRepositorySuite struct {
	suite.Suite
	repository repository.InventoryRepository
	client     *mongo.Client
	db         *mongo.Database
}

func (s *MongoRepositorySuite) SetupSuite() {
	// Для интеграционных тестов с MongoDB потребуется настройка testcontainers
	// В текущей реализации используем заглушку - тесты будут пропущены
	// TODO: настроить testcontainers для MongoDB
}

func (s *MongoRepositorySuite) SetupTest() {
	// Заглушка для тестов
}

func (s *MongoRepositorySuite) TearDownSuite() {
	// Заглушка для тестов
}

func TestMongoRepositoryIntegration(t *testing.T) {
	suite.Run(t, new(MongoRepositorySuite))
}
