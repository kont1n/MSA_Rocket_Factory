package producer_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kont1n/MSA_Rocket_Factory/assembly/internal/service"
	"github.com/kont1n/MSA_Rocket_Factory/assembly/internal/service/producer"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

type ProducerServiceSuite struct {
	suite.Suite
	assemblyProducer *MockProducer
	service          service.ProducerService
}

type MockProducer struct {
	SendFunc func(ctx context.Context, key, value []byte) error
}

func (m *MockProducer) Send(ctx context.Context, key, value []byte) error {
	if m.SendFunc != nil {
		return m.SendFunc(ctx, key, value)
	}
	return nil
}

func (s *ProducerServiceSuite) SetupSuite() {
	// Инициализируем logger для тестов
	if err := logger.Init("debug", false); err != nil {
		panic(err)
	}

	s.assemblyProducer = &MockProducer{}
	s.service = producer.NewService(s.assemblyProducer)
}

func (s *ProducerServiceSuite) SetupTest() {
	// Сбрасываем моки перед каждым тестом
	s.assemblyProducer.SendFunc = nil
}

func (s *ProducerServiceSuite) TearDownSuite() {
}

func TestProducerService(t *testing.T) {
	suite.Run(t, new(ProducerServiceSuite))
}
