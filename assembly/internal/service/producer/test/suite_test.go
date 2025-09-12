package producer_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/kont1n/MSA_Rocket_Factory/assembly/internal/service"
	"github.com/kont1n/MSA_Rocket_Factory/assembly/internal/service/producer"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

type ProducerServiceSuite struct {
	suite.Suite
	assemblyProducer *MockProducer
	mockMetrics      *MockKafkaMetrics
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

type MockKafkaMetrics struct {
	RecordProducerMessageFunc func(ctx context.Context, topic string, partition int32, success bool, duration time.Duration)
}

func (m *MockKafkaMetrics) RecordProducerMessage(ctx context.Context, topic string, partition int32, success bool, duration time.Duration) {
	if m.RecordProducerMessageFunc != nil {
		m.RecordProducerMessageFunc(ctx, topic, partition, success, duration)
	}
}

func (s *ProducerServiceSuite) SetupSuite() {
	// Инициализируем logger для тестов
	if err := logger.Init(context.Background(), "debug", false, "stdout", "", "assembly-test", "test"); err != nil {
		panic(err)
	}

	s.assemblyProducer = &MockProducer{}
	s.mockMetrics = &MockKafkaMetrics{}
	s.service = producer.NewService(s.assemblyProducer, s.mockMetrics)
}

func (s *ProducerServiceSuite) SetupTest() {
	// Сбрасываем моки перед каждым тестом
	s.assemblyProducer.SendFunc = nil
	s.mockMetrics.RecordProducerMessageFunc = nil
}

func (s *ProducerServiceSuite) TearDownSuite() {
}

func TestProducerService(t *testing.T) {
	suite.Run(t, new(ProducerServiceSuite))
}
