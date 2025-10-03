package producer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
	eventsV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/events/v1"
)

func init() {
	_ = logger.InitSimple("error", false)
}

// MockProducer - мок для kafka.Producer
type MockProducer struct {
	mock.Mock
}

func (m *MockProducer) Send(ctx context.Context, key, value []byte) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

// MockKafkaMetrics - мок для KafkaMetrics
type MockKafkaMetrics struct {
	mock.Mock
}

func (m *MockKafkaMetrics) RecordProducerMessage(ctx context.Context, topic string, partition int32, success bool, duration time.Duration) {
	m.Called(ctx, topic, partition, success, duration)
}

func TestNewService(t *testing.T) {
	// Arrange
	mockProducer := &MockProducer{}
	mockMetrics := &MockKafkaMetrics{}

	// Act
	svc := NewService(mockProducer, mockMetrics)

	// Assert
	assert.NotNil(t, svc)
	assert.Equal(t, mockProducer, svc.orderPaidProducer)
	assert.Equal(t, mockMetrics, svc.metrics)
}

func TestProduceOrderPaid_Success(t *testing.T) {
	// Arrange
	mockProducer := &MockProducer{}
	mockMetrics := &MockKafkaMetrics{}
	svc := NewService(mockProducer, mockMetrics)

	eventUUID := uuid.New()
	orderUUID := uuid.New()
	userUUID := uuid.New()
	transactionUUID := uuid.New()

	event := model.OrderPaidEvent{
		EventUUID:       eventUUID,
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PaymentMethod:   "CARD",
		TransactionUUID: transactionUUID,
	}

	// Ожидаем что Send будет вызван с правильными параметрами
	mockProducer.On("Send", mock.Anything, []byte(eventUUID.String()), mock.MatchedBy(func(value []byte) bool {
		// Проверяем что payload можно десериализовать обратно
		msg := &eventsV1.OrderPaid{}
		err := proto.Unmarshal(value, msg)
		if err != nil {
			return false
		}
		return msg.EventUuid == eventUUID.String() &&
			msg.OrderUuid == orderUUID.String() &&
			msg.UserUuid == userUUID.String() &&
			msg.PaymentMethod == "CARD" &&
			msg.TransactionUuid == transactionUUID.String()
	})).Return(nil)

	// Ожидаем запись метрик
	mockMetrics.On("RecordProducerMessage", mock.Anything, "order-paid", int32(0), true, mock.AnythingOfType("time.Duration")).Return()

	// Act
	err := svc.ProduceOrderPaid(context.Background(), event)

	// Assert
	assert.NoError(t, err)
	mockProducer.AssertExpectations(t)
	mockMetrics.AssertExpectations(t)
}

func TestProduceOrderPaid_SendError(t *testing.T) {
	// Arrange
	mockProducer := &MockProducer{}
	mockMetrics := &MockKafkaMetrics{}
	svc := NewService(mockProducer, mockMetrics)

	event := model.OrderPaidEvent{
		EventUUID:       uuid.New(),
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PaymentMethod:   "CARD",
		TransactionUUID: uuid.New(),
	}

	expectedError := errors.New("kafka send error")

	mockProducer.On("Send", mock.Anything, mock.Anything, mock.Anything).Return(expectedError)
	mockMetrics.On("RecordProducerMessage", mock.Anything, "order-paid", int32(0), false, mock.AnythingOfType("time.Duration")).Return()

	// Act
	err := svc.ProduceOrderPaid(context.Background(), event)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockProducer.AssertExpectations(t)
	mockMetrics.AssertExpectations(t)
}

func TestProduceOrderPaid_WithoutMetrics(t *testing.T) {
	// Arrange
	mockProducer := &MockProducer{}
	svc := NewService(mockProducer, nil) // Без метрик

	event := model.OrderPaidEvent{
		EventUUID:       uuid.New(),
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PaymentMethod:   "CARD",
		TransactionUUID: uuid.New(),
	}

	mockProducer.On("Send", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Act
	err := svc.ProduceOrderPaid(context.Background(), event)

	// Assert
	assert.NoError(t, err)
	mockProducer.AssertExpectations(t)
}

func TestProduceOrderPaid_DifferentPaymentMethods(t *testing.T) {
	tests := []struct {
		name          string
		paymentMethod string
	}{
		{
			name:          "оплата картой CARD",
			paymentMethod: "CARD",
		},
		{
			name:          "оплата через SBP",
			paymentMethod: "SBP",
		},
		{
			name:          "оплата через CREDIT_CARD",
			paymentMethod: "CREDIT_CARD",
		},
		{
			name:          "оплата через INVESTOR_MONEY",
			paymentMethod: "INVESTOR_MONEY",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockProducer := &MockProducer{}
			mockMetrics := &MockKafkaMetrics{}
			svc := NewService(mockProducer, mockMetrics)

			event := model.OrderPaidEvent{
				EventUUID:       uuid.New(),
				OrderUUID:       uuid.New(),
				UserUUID:        uuid.New(),
				PaymentMethod:   tt.paymentMethod,
				TransactionUUID: uuid.New(),
			}

			mockProducer.On("Send", mock.Anything, mock.Anything, mock.MatchedBy(func(value []byte) bool {
				msg := &eventsV1.OrderPaid{}
				err := proto.Unmarshal(value, msg)
				return err == nil && msg.PaymentMethod == tt.paymentMethod
			})).Return(nil)

			mockMetrics.On("RecordProducerMessage", mock.Anything, "order-paid", int32(0), true, mock.AnythingOfType("time.Duration")).Return()

			// Act
			err := svc.ProduceOrderPaid(context.Background(), event)

			// Assert
			assert.NoError(t, err)
			mockProducer.AssertExpectations(t)
			mockMetrics.AssertExpectations(t)
		})
	}
}

func TestProduceOrderPaid_EventUUIDAsKey(t *testing.T) {
	// Arrange
	mockProducer := &MockProducer{}
	mockMetrics := &MockKafkaMetrics{}
	svc := NewService(mockProducer, mockMetrics)

	eventUUID := uuid.New()
	event := model.OrderPaidEvent{
		EventUUID:       eventUUID,
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PaymentMethod:   "CARD",
		TransactionUUID: uuid.New(),
	}

	// Проверяем что ключ сообщения == EventUUID
	mockProducer.On("Send", mock.Anything, []byte(eventUUID.String()), mock.Anything).Return(nil)
	mockMetrics.On("RecordProducerMessage", mock.Anything, "order-paid", int32(0), true, mock.AnythingOfType("time.Duration")).Return()

	// Act
	err := svc.ProduceOrderPaid(context.Background(), event)

	// Assert
	assert.NoError(t, err)
	mockProducer.AssertExpectations(t)
}

func TestProduceOrderPaid_MetricsRecordDuration(t *testing.T) {
	// Arrange
	mockProducer := &MockProducer{}
	mockMetrics := &MockKafkaMetrics{}
	svc := NewService(mockProducer, mockMetrics)

	event := model.OrderPaidEvent{
		EventUUID:       uuid.New(),
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PaymentMethod:   "CARD",
		TransactionUUID: uuid.New(),
	}

	mockProducer.On("Send", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Проверяем что метрики записывают длительность
	var recordedDuration time.Duration
	mockMetrics.On("RecordProducerMessage", mock.Anything, "order-paid", int32(0), true, mock.AnythingOfType("time.Duration")).
		Run(func(args mock.Arguments) {
			recordedDuration = args.Get(4).(time.Duration)
		}).Return()

	// Act
	err := svc.ProduceOrderPaid(context.Background(), event)

	// Assert
	assert.NoError(t, err)
	assert.Greater(t, recordedDuration, time.Duration(0))
	mockMetrics.AssertExpectations(t)
}

func TestProduceOrderPaid_TableDriven(t *testing.T) {
	tests := []struct {
		name          string
		event         model.OrderPaidEvent
		sendError     error
		expectedError bool
	}{
		{
			name: "успешная отправка события",
			event: model.OrderPaidEvent{
				EventUUID:       uuid.New(),
				OrderUUID:       uuid.New(),
				UserUUID:        uuid.New(),
				PaymentMethod:   "CARD",
				TransactionUUID: uuid.New(),
			},
			sendError:     nil,
			expectedError: false,
		},
		{
			name: "ошибка подключения к Kafka",
			event: model.OrderPaidEvent{
				EventUUID:       uuid.New(),
				OrderUUID:       uuid.New(),
				UserUUID:        uuid.New(),
				PaymentMethod:   "SBP",
				TransactionUUID: uuid.New(),
			},
			sendError:     errors.New("connection refused"),
			expectedError: true,
		},
		{
			name: "ошибка таймаута Kafka",
			event: model.OrderPaidEvent{
				EventUUID:       uuid.New(),
				OrderUUID:       uuid.New(),
				UserUUID:        uuid.New(),
				PaymentMethod:   "CREDIT_CARD",
				TransactionUUID: uuid.New(),
			},
			sendError:     errors.New("request timeout"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockProducer := &MockProducer{}
			mockMetrics := &MockKafkaMetrics{}
			svc := NewService(mockProducer, mockMetrics)

			mockProducer.On("Send", mock.Anything, mock.Anything, mock.Anything).Return(tt.sendError)
			mockMetrics.On("RecordProducerMessage", mock.Anything, "order-paid", int32(0), tt.sendError == nil, mock.AnythingOfType("time.Duration")).Return()

			// Act
			err := svc.ProduceOrderPaid(context.Background(), tt.event)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
				assert.Equal(t, tt.sendError, err)
			} else {
				assert.NoError(t, err)
			}

			mockProducer.AssertExpectations(t)
			mockMetrics.AssertExpectations(t)
		})
	}
}

func TestProduceOrderPaid_ContextCancellation(t *testing.T) {
	// Arrange
	mockProducer := &MockProducer{}
	mockMetrics := &MockKafkaMetrics{}
	svc := NewService(mockProducer, mockMetrics)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Отменяем контекст сразу

	event := model.OrderPaidEvent{
		EventUUID:       uuid.New(),
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PaymentMethod:   "CARD",
		TransactionUUID: uuid.New(),
	}

	// Producer может вернуть ошибку из-за отмененного контекста
	mockProducer.On("Send", mock.Anything, mock.Anything, mock.Anything).Return(context.Canceled)
	mockMetrics.On("RecordProducerMessage", mock.Anything, "order-paid", int32(0), false, mock.AnythingOfType("time.Duration")).Return()

	// Act
	err := svc.ProduceOrderPaid(ctx, event)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
	mockProducer.AssertExpectations(t)
}
