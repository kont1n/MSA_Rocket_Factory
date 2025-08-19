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

type OrderPaidConsumerTestSuite struct {
	suite.Suite
	orderPaidConsumer   *consumerMocks.Consumer
	orderPaidDecoder    *decoderMocks.OrderPaidDecoder
	notificationService *mockNotificationService
	service             service.OrderPaidConsumerService
}

func (s *OrderPaidConsumerTestSuite) SetupSuite() {
	// Инициализируем no-op логгер для тестов
	logger.SetNopLogger()

	s.orderPaidConsumer = consumerMocks.NewConsumer(s.T())
	s.orderPaidDecoder = decoderMocks.NewOrderPaidDecoder(s.T())
	s.notificationService = &mockNotificationService{}

	s.service = consumer.NewOrderPaidService(
		s.orderPaidConsumer,
		s.orderPaidDecoder,
		s.notificationService,
	)
}

func (s *OrderPaidConsumerTestSuite) SetupTest() {
	// Сбрасываем моки перед каждым тестом
	s.orderPaidConsumer.ExpectedCalls = nil
	s.orderPaidDecoder.ExpectedCalls = nil
	s.notificationService.reset()
}

func (s *OrderPaidConsumerTestSuite) TearDownSuite() {
}

func TestOrderPaidConsumerTestSuite(t *testing.T) {
	suite.Run(t, new(OrderPaidConsumerTestSuite))
}

func (s *OrderPaidConsumerTestSuite) TestOrderPaidHandler_Success() {
	// Подготавливаем тестовые данные
	event := &model.OrderPaidEvent{
		EventUUID:       uuid.New(),
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PaymentMethod:   "card",
		TransactionUUID: uuid.New(),
	}

	msg := kafka.Message{
		Topic:     "order.paid",
		Partition: 0,
		Offset:    1,
		Value:     []byte("test message"),
	}

	// Настраиваем моки
	s.orderPaidDecoder.On("Decode", msg.Value).Return(event, nil)
	s.notificationService.On("NotifyOrderPaid", mock.Anything, event).Return(nil)

	// Выполняем тест
	err := s.service.(interface {
		OrderPaidHandler(ctx context.Context, msg kafka.Message) error
	}).OrderPaidHandler(context.Background(), msg)

	// Проверяем результат
	s.NoError(err)
	s.orderPaidDecoder.AssertExpectations(s.T())
	s.notificationService.AssertExpectations(s.T())
}

func (s *OrderPaidConsumerTestSuite) TestOrderPaidHandler_DecodeError() {
	// Подготавливаем тестовые данные
	msg := kafka.Message{
		Topic:     "order.paid",
		Partition: 0,
		Offset:    1,
		Value:     []byte("invalid message"),
	}

	// Настраиваем мок для ошибки декодирования
	expectedError := errors.New("decode error")
	s.orderPaidDecoder.On("Decode", msg.Value).Return(nil, expectedError)

	// Выполняем тест
	err := s.service.(interface {
		OrderPaidHandler(ctx context.Context, msg kafka.Message) error
	}).OrderPaidHandler(context.Background(), msg)

	// Проверяем результат
	s.Error(err)
	s.Equal(expectedError, err)
	s.orderPaidDecoder.AssertExpectations(s.T())
}

func (s *OrderPaidConsumerTestSuite) TestOrderPaidHandler_NotificationError() {
	// Подготавливаем тестовые данные
	event := &model.OrderPaidEvent{
		EventUUID:       uuid.New(),
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PaymentMethod:   "card",
		TransactionUUID: uuid.New(),
	}

	msg := kafka.Message{
		Topic:     "order.paid",
		Partition: 0,
		Offset:    1,
		Value:     []byte("test message"),
	}

	// Настраиваем моки
	s.orderPaidDecoder.On("Decode", msg.Value).Return(event, nil)
	expectedError := errors.New("notification error")
	s.notificationService.On("NotifyOrderPaid", mock.Anything, event).Return(expectedError)

	// Выполняем тест
	err := s.service.(interface {
		OrderPaidHandler(ctx context.Context, msg kafka.Message) error
	}).OrderPaidHandler(context.Background(), msg)

	// Проверяем результат
	s.Error(err)
	s.Equal(expectedError, err)
	s.orderPaidDecoder.AssertExpectations(s.T())
	s.notificationService.AssertExpectations(s.T())
}

func (s *OrderPaidConsumerTestSuite) TestRunConsumer_Success() {
	// Настраиваем мок для успешного запуска consumer'а
	s.orderPaidConsumer.On("Consume", mock.Anything, mock.AnythingOfType("kafka.MessageHandler")).Return(nil)

	// Выполняем тест
	err := s.service.RunConsumer(context.Background())

	// Проверяем результат
	s.NoError(err)
	s.orderPaidConsumer.AssertExpectations(s.T())
}

func (s *OrderPaidConsumerTestSuite) TestRunConsumer_Error() {
	// Настраиваем мок для ошибки запуска consumer'а
	expectedError := errors.New("consumer error")
	s.orderPaidConsumer.On("Consume", mock.Anything, mock.AnythingOfType("kafka.MessageHandler")).Return(expectedError)

	// Выполняем тест
	err := s.service.RunConsumer(context.Background())

	// Проверяем результат
	s.Error(err)
	s.Equal(expectedError, err)
	s.orderPaidConsumer.AssertExpectations(s.T())
}
