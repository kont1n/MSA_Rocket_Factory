package assembly_test

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/kont1n/MSA_Rocket_Factory/assembly/internal/model"
)

func (s *AssemblyServiceSuite) TestAssemble_Success() {
	// Подготавливаем тестовые данные
	event := model.OrderPaidEvent{
		EventUUID:       uuid.New(),
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PaymentMethod:   "CARD",
		TransactionUUID: uuid.New(),
	}

	// Настраиваем мок для успешного производства
	s.assemblyProducerService.ProduceAssemblyFunc = func(ctx context.Context, event model.ShipAssembledEvent) error {
		// Проверяем, что переданы правильные данные
		assert.Equal(s.T(), event.OrderUUID, event.OrderUUID)
		assert.Equal(s.T(), event.UserUUID, event.UserUUID)
		assert.Equal(s.T(), int64(10), event.BuildTime) // delayTime = 10
		return nil
	}

	// Выполняем тест
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := s.service.Assemble(ctx, event)

	// Проверяем результат
	assert.NoError(s.T(), err)
}

func (s *AssemblyServiceSuite) TestAssemble_ContextCancelled() {
	// Подготавливаем тестовые данные
	event := model.OrderPaidEvent{
		EventUUID:       uuid.New(),
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PaymentMethod:   "CARD",
		TransactionUUID: uuid.New(),
	}

	// Создаем контекст, который будет отменен сразу
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Отменяем контекст сразу

	err := s.service.Assemble(ctx, event)

	// Проверяем, что получили ошибку отмены контекста
	assert.Error(s.T(), err)
	assert.Equal(s.T(), context.Canceled, err)
}

func (s *AssemblyServiceSuite) TestAssemble_ProducerError() {
	// Подготавливаем тестовые данные
	event := model.OrderPaidEvent{
		EventUUID:       uuid.New(),
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PaymentMethod:   "CARD",
		TransactionUUID: uuid.New(),
	}

	// Настраиваем мок для ошибки производства
	expectedError := assert.AnError
	s.assemblyProducerService.ProduceAssemblyFunc = func(ctx context.Context, event model.ShipAssembledEvent) error {
		return expectedError
	}

	// Выполняем тест
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := s.service.Assemble(ctx, event)

	// Проверяем, что получили ошибку от producer
	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
}

func (s *AssemblyServiceSuite) TestAssemble_Timeout() {
	// Подготавливаем тестовые данные
	event := model.OrderPaidEvent{
		EventUUID:       uuid.New(),
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PaymentMethod:   "CARD",
		TransactionUUID: uuid.New(),
	}

	// Создаем контекст с очень коротким таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	err := s.service.Assemble(ctx, event)

	// Проверяем, что получили ошибку таймаута
	assert.Error(s.T(), err)
	assert.Equal(s.T(), context.DeadlineExceeded, err)
}
