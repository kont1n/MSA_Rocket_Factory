package producer_test

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/kont1n/MSA_Rocket_Factory/assembly/internal/model"
)

func (s *ProducerServiceSuite) TestProduceAssembly_Success() {
	// Подготавливаем тестовые данные
	event := model.ShipAssembledEvent{
		EventUUID: uuid.New(),
		OrderUUID: uuid.New(),
		UserUUID:  uuid.New(),
		BuildTime: 10,
	}

	// Настраиваем мок для успешной отправки
	s.assemblyProducer.SendFunc = func(ctx context.Context, key, value []byte) error {
		// Проверяем, что ключ соответствует EventUUID
		assert.Equal(s.T(), event.EventUUID.String(), string(key))
		// Проверяем, что value не пустой (protobuf сообщение)
		assert.NotEmpty(s.T(), value)
		return nil
	}

	// Выполняем тест
	ctx := context.Background()
	err := s.service.ProduceAssembly(ctx, event)

	// Проверяем результат
	assert.NoError(s.T(), err)
}

func (s *ProducerServiceSuite) TestProduceAssembly_ProducerError() {
	// Подготавливаем тестовые данные
	event := model.ShipAssembledEvent{
		EventUUID: uuid.New(),
		OrderUUID: uuid.New(),
		UserUUID:  uuid.New(),
		BuildTime: 10,
	}

	// Настраиваем мок для ошибки отправки
	expectedError := errors.New("send error")
	s.assemblyProducer.SendFunc = func(ctx context.Context, key, value []byte) error {
		return expectedError
	}

	// Выполняем тест
	ctx := context.Background()
	err := s.service.ProduceAssembly(ctx, event)

	// Проверяем результат
	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
}

func (s *ProducerServiceSuite) TestProduceAssembly_ContextCancelled() {
	// Подготавливаем тестовые данные
	event := model.ShipAssembledEvent{
		EventUUID: uuid.New(),
		OrderUUID: uuid.New(),
		UserUUID:  uuid.New(),
		BuildTime: 10,
	}

	// Создаем отмененный контекст
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Выполняем тест
	err := s.service.ProduceAssembly(ctx, event)

	// Проверяем результат - должен быть успех, так как контекст проверяется только в kafka producer
	assert.NoError(s.T(), err)
}
