package consumer_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	decoderMocks "github.com/kont1n/MSA_Rocket_Factory/notification/internal/converter/kafka/mocks"
	"github.com/kont1n/MSA_Rocket_Factory/notification/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/notification/internal/service"
	"github.com/kont1n/MSA_Rocket_Factory/notification/internal/service/consumer"
	consumerMocks "github.com/kont1n/MSA_Rocket_Factory/notification/internal/service/consumer/mocks"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/kafka"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

type ShipAssembledConsumerTestSuite struct {
	suite.Suite
	shipAssembledConsumer *consumerMocks.Consumer
	shipAssembledDecoder  *decoderMocks.ShipAssembledDecoder
	notificationService   *mockNotificationService
	service               service.ShipAssembledConsumerService
}

func (s *ShipAssembledConsumerTestSuite) SetupSuite() {
	// Инициализируем no-op логгер для тестов
	logger.SetNopLogger()

	s.shipAssembledConsumer = consumerMocks.NewConsumer(s.T())
	s.shipAssembledDecoder = decoderMocks.NewShipAssembledDecoder(s.T())
	s.notificationService = &mockNotificationService{}

	s.service = consumer.NewShipAssembledService(
		s.shipAssembledConsumer,
		s.shipAssembledDecoder,
		s.notificationService,
	)
}

func (s *ShipAssembledConsumerTestSuite) SetupTest() {
	// Сбрасываем моки перед каждым тестом
	s.shipAssembledConsumer.ExpectedCalls = nil
	s.shipAssembledDecoder.ExpectedCalls = nil
	s.notificationService.reset()
}

func (s *ShipAssembledConsumerTestSuite) TearDownSuite() {
}

func TestShipAssembledConsumerTestSuite(t *testing.T) {
	suite.Run(t, new(ShipAssembledConsumerTestSuite))
}

func (s *ShipAssembledConsumerTestSuite) TestShipAssembledHandler_Success() {
	// Подготавливаем тестовые данные
	event := &model.ShipAssembledEvent{
		EventUUID: uuid.New(),
		OrderUUID: uuid.New(),
		UserUUID:  uuid.New(),
		BuildTime: 120,
	}

	msg := kafka.Message{
		Topic:     "ship.assembled",
		Partition: 0,
		Offset:    1,
		Value:     []byte("test message"),
	}

	// Настраиваем моки
	s.shipAssembledDecoder.On("Decode", msg.Value).Return(event, nil)
	s.notificationService.On("NotifyShipAssembled", mock.Anything, event).Return(nil)

	// Выполняем тест
	err := s.service.(interface {
		ShipAssembledHandler(ctx context.Context, msg kafka.Message) error
	}).ShipAssembledHandler(context.Background(), msg)

	// Проверяем результат
	s.NoError(err)
	s.shipAssembledDecoder.AssertExpectations(s.T())
	s.notificationService.AssertExpectations(s.T())
}

func (s *ShipAssembledConsumerTestSuite) TestShipAssembledHandler_DecodeError() {
	// Подготавливаем тестовые данные
	msg := kafka.Message{
		Topic:     "ship.assembled",
		Partition: 0,
		Offset:    1,
		Value:     []byte("invalid message"),
	}

	// Настраиваем мок для ошибки декодирования
	expectedError := errors.New("decode error")
	s.shipAssembledDecoder.On("Decode", msg.Value).Return(nil, expectedError)

	// Выполняем тест
	err := s.service.(interface {
		ShipAssembledHandler(ctx context.Context, msg kafka.Message) error
	}).ShipAssembledHandler(context.Background(), msg)

	// Проверяем результат
	s.Error(err)
	s.Equal(expectedError, err)
	s.shipAssembledDecoder.AssertExpectations(s.T())
}

func (s *ShipAssembledConsumerTestSuite) TestShipAssembledHandler_NotificationError() {
	// Подготавливаем тестовые данные
	event := &model.ShipAssembledEvent{
		EventUUID: uuid.New(),
		OrderUUID: uuid.New(),
		UserUUID:  uuid.New(),
		BuildTime: 120,
	}

	msg := kafka.Message{
		Topic:     "ship.assembled",
		Partition: 0,
		Offset:    1,
		Value:     []byte("test message"),
	}

	// Настраиваем моки
	s.shipAssembledDecoder.On("Decode", msg.Value).Return(event, nil)
	expectedError := errors.New("notification error")
	s.notificationService.On("NotifyShipAssembled", mock.Anything, event).Return(expectedError)

	// Выполняем тест
	err := s.service.(interface {
		ShipAssembledHandler(ctx context.Context, msg kafka.Message) error
	}).ShipAssembledHandler(context.Background(), msg)

	// Проверяем результат
	s.Error(err)
	s.Equal(expectedError, err)
	s.shipAssembledDecoder.AssertExpectations(s.T())
	s.notificationService.AssertExpectations(s.T())
}

func (s *ShipAssembledConsumerTestSuite) TestRunConsumer_Success() {
	// Настраиваем мок для успешного запуска consumer'а
	s.shipAssembledConsumer.On("Consume", mock.Anything, mock.AnythingOfType("kafka.MessageHandler")).Return(nil)

	// Выполняем тест
	err := s.service.RunConsumer(context.Background())

	// Проверяем результат
	s.NoError(err)
	s.shipAssembledConsumer.AssertExpectations(s.T())
}

func (s *ShipAssembledConsumerTestSuite) TestRunConsumer_Error() {
	// Настраиваем мок для ошибки запуска consumer'а
	expectedError := errors.New("consumer error")
	s.shipAssembledConsumer.On("Consume", mock.Anything, mock.AnythingOfType("kafka.MessageHandler")).Return(expectedError)

	// Выполняем тест
	err := s.service.RunConsumer(context.Background())

	// Проверяем результат
	s.Error(err)
	s.Equal(expectedError, err)
	s.shipAssembledConsumer.AssertExpectations(s.T())
}
