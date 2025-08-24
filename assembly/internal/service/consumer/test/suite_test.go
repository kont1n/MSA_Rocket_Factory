package consumer_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kont1n/MSA_Rocket_Factory/assembly/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/assembly/internal/service"
	consumerPkg "github.com/kont1n/MSA_Rocket_Factory/assembly/internal/service/consumer"
	kafkaPkg "github.com/kont1n/MSA_Rocket_Factory/platform/pkg/kafka"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

type ConsumerServiceSuite struct {
	suite.Suite
	assemblyRecordedConsumer *MockConsumer
	assemblyRecordedDecoder  *MockAssemblyRecordedDecoder
	assemblyService          *MockAssemblyService
	service                  service.ConsumerService
}

type MockConsumer struct {
	ConsumeFunc func(ctx context.Context, handler kafkaPkg.MessageHandler) error
}

func (m *MockConsumer) Consume(ctx context.Context, handler kafkaPkg.MessageHandler) error {
	if m.ConsumeFunc != nil {
		return m.ConsumeFunc(ctx, handler)
	}
	return nil
}

type MockAssemblyRecordedDecoder struct {
	DecodeFunc func(data []byte) (model.OrderPaidEvent, error)
}

func (m *MockAssemblyRecordedDecoder) Decode(data []byte) (model.OrderPaidEvent, error) {
	if m.DecodeFunc != nil {
		return m.DecodeFunc(data)
	}
	return model.OrderPaidEvent{}, nil
}

type MockAssemblyService struct {
	AssembleFunc func(ctx context.Context, event model.OrderPaidEvent) error
}

func (m *MockAssemblyService) Assemble(ctx context.Context, event model.OrderPaidEvent) error {
	if m.AssembleFunc != nil {
		return m.AssembleFunc(ctx, event)
	}
	return nil
}

func (s *ConsumerServiceSuite) SetupSuite() {
	// Инициализируем logger для тестов
	if err := logger.Init("debug", false); err != nil {
		panic(err)
	}

	s.assemblyRecordedConsumer = &MockConsumer{}
	s.assemblyRecordedDecoder = &MockAssemblyRecordedDecoder{}
	s.assemblyService = &MockAssemblyService{}
	s.service = consumerPkg.NewService(s.assemblyRecordedConsumer, s.assemblyRecordedDecoder, s.assemblyService)
}

func (s *ConsumerServiceSuite) SetupTest() {
	// Сбрасываем моки перед каждым тестом
	s.assemblyRecordedConsumer.ConsumeFunc = nil
	s.assemblyRecordedDecoder.DecodeFunc = nil
	s.assemblyService.AssembleFunc = nil
}

func (s *ConsumerServiceSuite) TearDownSuite() {
}

func TestConsumerService(t *testing.T) {
	suite.Run(t, new(ConsumerServiceSuite))
}
