package consumer_test

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/kont1n/MSA_Rocket_Factory/assembly/internal/model"
	kafkaPkg "github.com/kont1n/MSA_Rocket_Factory/platform/pkg/kafka"
)

func (s *ConsumerServiceSuite) TestRunConsumer_Success() {
	// Настраиваем мок для успешного запуска consumer
	s.assemblyRecordedConsumer.ConsumeFunc = func(ctx context.Context, handler kafkaPkg.MessageHandler) error {
		// Проверяем, что handler не nil
		assert.NotNil(s.T(), handler)
		return nil
	}

	// Выполняем тест
	ctx := context.Background()
	err := s.service.RunConsumer(ctx)

	// Проверяем результат
	assert.NoError(s.T(), err)
}

func (s *ConsumerServiceSuite) TestRunConsumer_ConsumerError() {
	// Настраиваем мок для ошибки consumer
	expectedError := errors.New("consumer error")
	s.assemblyRecordedConsumer.ConsumeFunc = func(ctx context.Context, handler kafkaPkg.MessageHandler) error {
		return expectedError
	}

	// Выполняем тест
	ctx := context.Background()
	err := s.service.RunConsumer(ctx)

	// Проверяем результат
	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
}

func (s *ConsumerServiceSuite) TestOrderPaidHandler_Success() {
	// Подготавливаем тестовые данные
	event := model.OrderPaidEvent{
		EventUUID:       uuid.New(),
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PaymentMethod:   "CARD",
		TransactionUUID: uuid.New(),
	}

	msg := kafkaPkg.Message{
		Topic: "order.paid",
		Value: []byte("test data"),
	}

	// Настраиваем моки
	s.assemblyRecordedDecoder.DecodeFunc = func(data []byte) (model.OrderPaidEvent, error) {
		return event, nil
	}

	s.assemblyService.AssembleFunc = func(ctx context.Context, event model.OrderPaidEvent) error {
		return nil
	}

	// Выполняем тест через мок consumer
	ctx := context.Background()
	s.assemblyRecordedConsumer.ConsumeFunc = func(ctx context.Context, handler kafkaPkg.MessageHandler) error {
		// Вызываем handler напрямую для тестирования
		return handler(ctx, msg)
	}

	err := s.service.RunConsumer(ctx)

	// Проверяем результат
	assert.NoError(s.T(), err)
}

func (s *ConsumerServiceSuite) TestOrderPaidHandler_DecodeError() {
	// Подготавливаем тестовые данные
	msg := kafkaPkg.Message{
		Topic: "order.paid",
		Value: []byte("invalid data"),
	}

	// Настраиваем мок для ошибки декодирования
	expectedError := errors.New("decode error")
	s.assemblyRecordedDecoder.DecodeFunc = func(data []byte) (model.OrderPaidEvent, error) {
		return model.OrderPaidEvent{}, expectedError
	}

	// Выполняем тест через мок consumer
	ctx := context.Background()
	s.assemblyRecordedConsumer.ConsumeFunc = func(ctx context.Context, handler kafkaPkg.MessageHandler) error {
		// Вызываем handler напрямую для тестирования
		return handler(ctx, msg)
	}

	err := s.service.RunConsumer(ctx)

	// Проверяем результат
	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
}

func (s *ConsumerServiceSuite) TestOrderPaidHandler_AssembleError() {
	// Подготавливаем тестовые данные
	event := model.OrderPaidEvent{
		EventUUID:       uuid.New(),
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PaymentMethod:   "CARD",
		TransactionUUID: uuid.New(),
	}

	msg := kafkaPkg.Message{
		Topic: "order.paid",
		Value: []byte("test data"),
	}

	// Настраиваем моки
	s.assemblyRecordedDecoder.DecodeFunc = func(data []byte) (model.OrderPaidEvent, error) {
		return event, nil
	}

	// Настраиваем мок для ошибки сборки
	expectedError := errors.New("assemble error")
	s.assemblyService.AssembleFunc = func(ctx context.Context, event model.OrderPaidEvent) error {
		return expectedError
	}

	// Выполняем тест через мок consumer
	ctx := context.Background()
	s.assemblyRecordedConsumer.ConsumeFunc = func(ctx context.Context, handler kafkaPkg.MessageHandler) error {
		// Вызываем handler напрямую для тестирования
		return handler(ctx, msg)
	}

	err := s.service.RunConsumer(ctx)

	// Проверяем результат
	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
}
