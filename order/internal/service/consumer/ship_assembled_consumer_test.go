package consumer

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/service/mocks"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/kafka"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

func init() {
	_ = logger.InitSimple("error", false)
}

// MockConsumer - мок для kafka.Consumer
type MockConsumer struct {
	mock.Mock
}

func (m *MockConsumer) Consume(ctx context.Context, handler kafka.MessageHandler) error {
	args := m.Called(ctx, handler)
	return args.Error(0)
}

// MockDecoder - мок для ShipAssembledDecoder
type MockDecoder struct {
	mock.Mock
}

func (m *MockDecoder) Decode(data []byte) (*model.ShipAssembledEvent, error) {
	args := m.Called(data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ShipAssembledEvent), args.Error(1)
}

// MockConsumerMetrics - мок для KafkaMetrics
type MockConsumerMetrics struct {
	mock.Mock
}

func (m *MockConsumerMetrics) RecordConsumerMessage(ctx context.Context, topic string, partition int32, groupID string, success bool) {
	m.Called(ctx, topic, partition, groupID, success)
}

func (m *MockConsumerMetrics) RecordConsumerLag(ctx context.Context, topic string, partition int32, groupID string, lagSeconds float64, offsetLag int64) {
	m.Called(ctx, topic, partition, groupID, lagSeconds, offsetLag)
}

func TestNewService(t *testing.T) {
	// Arrange
	mockConsumer := &MockConsumer{}
	mockDecoder := &MockDecoder{}
	mockOrderService := mocks.NewOrderService(t)
	mockMetrics := &MockConsumerMetrics{}

	// Act
	svc := NewService(mockConsumer, mockDecoder, mockOrderService, mockMetrics)

	// Assert
	assert.NotNil(t, svc)
	assert.Equal(t, mockConsumer, svc.shipAssembledConsumer)
	assert.Equal(t, mockDecoder, svc.shipAssembledDecoder)
	assert.Equal(t, mockOrderService, svc.orderService)
	assert.Equal(t, mockMetrics, svc.metrics)
}

func TestSetOrderService(t *testing.T) {
	// Arrange
	mockConsumer := &MockConsumer{}
	mockDecoder := &MockDecoder{}
	mockMetrics := &MockConsumerMetrics{}
	svc := NewService(mockConsumer, mockDecoder, nil, mockMetrics)

	mockOrderService := mocks.NewOrderService(t)

	// Act
	svc.SetOrderService(mockOrderService)

	// Assert
	assert.Equal(t, mockOrderService, svc.orderService)
}

func TestShipAssembledHandler_Success(t *testing.T) {
	// Arrange
	mockConsumer := &MockConsumer{}
	mockDecoder := &MockDecoder{}
	mockOrderService := mocks.NewOrderService(t)
	mockMetrics := &MockConsumerMetrics{}
	svc := NewService(mockConsumer, mockDecoder, mockOrderService, mockMetrics)

	eventUUID := uuid.New()
	orderUUID := uuid.New()
	userUUID := uuid.New()

	event := &model.ShipAssembledEvent{
		EventUUID: eventUUID,
		OrderUUID: orderUUID,
		UserUUID:  userUUID,
		BuildTime: 3600,
	}

	messageData := []byte("test-message-data")
	msg := kafka.Message{
		Topic:     "ship-assembled",
		Partition: 0,
		Offset:    100,
		Key:       []byte(eventUUID.String()),
		Value:     messageData,
	}

	// Ожидаем вызовы
	mockMetrics.On("RecordConsumerMessage", mock.Anything, "ship-assembled", int32(0), "order-ship-assembled-group", true).Return()
	mockDecoder.On("Decode", messageData).Return(event, nil)
	mockOrderService.EXPECT().UpdateOrderStatus(mock.Anything, orderUUID.String(), model.StatusAssembled).Return(nil).Once()

	// Act
	err := svc.ShipAssembledHandler(context.Background(), msg)

	// Assert
	assert.NoError(t, err)
	mockDecoder.AssertExpectations(t)
	mockMetrics.AssertExpectations(t)
	mockOrderService.AssertExpectations(t)
}

func TestShipAssembledHandler_DecodeError(t *testing.T) {
	// Arrange
	mockConsumer := &MockConsumer{}
	mockDecoder := &MockDecoder{}
	mockOrderService := mocks.NewOrderService(t)
	mockMetrics := &MockConsumerMetrics{}
	svc := NewService(mockConsumer, mockDecoder, mockOrderService, mockMetrics)

	messageData := []byte("invalid-data")
	msg := kafka.Message{
		Topic:     "ship-assembled",
		Partition: 0,
		Offset:    100,
		Value:     messageData,
	}

	expectedError := errors.New("decode error")

	mockMetrics.On("RecordConsumerMessage", mock.Anything, "ship-assembled", int32(0), "order-ship-assembled-group", true).Return()
	mockMetrics.On("RecordConsumerMessage", mock.Anything, "ship-assembled", int32(0), "order-ship-assembled-group", false).Return()
	mockDecoder.On("Decode", messageData).Return(nil, expectedError)

	// Act
	err := svc.ShipAssembledHandler(context.Background(), msg)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockDecoder.AssertExpectations(t)
	mockMetrics.AssertExpectations(t)
}

func TestShipAssembledHandler_UpdateStatusError(t *testing.T) {
	// Arrange
	mockConsumer := &MockConsumer{}
	mockDecoder := &MockDecoder{}
	mockOrderService := mocks.NewOrderService(t)
	mockMetrics := &MockConsumerMetrics{}
	svc := NewService(mockConsumer, mockDecoder, mockOrderService, mockMetrics)

	orderUUID := uuid.New()
	event := &model.ShipAssembledEvent{
		EventUUID: uuid.New(),
		OrderUUID: orderUUID,
		UserUUID:  uuid.New(),
		BuildTime: 3600,
	}

	messageData := []byte("test-message-data")
	msg := kafka.Message{
		Topic:     "ship-assembled",
		Partition: 0,
		Offset:    100,
		Value:     messageData,
	}

	expectedError := errors.New("update status error")

	mockMetrics.On("RecordConsumerMessage", mock.Anything, "ship-assembled", int32(0), "order-ship-assembled-group", true).Return()
	mockMetrics.On("RecordConsumerMessage", mock.Anything, "ship-assembled", int32(0), "order-ship-assembled-group", false).Return()
	mockDecoder.On("Decode", messageData).Return(event, nil)
	mockOrderService.EXPECT().UpdateOrderStatus(mock.Anything, orderUUID.String(), model.StatusAssembled).Return(expectedError).Once()

	// Act
	err := svc.ShipAssembledHandler(context.Background(), msg)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockDecoder.AssertExpectations(t)
	mockMetrics.AssertExpectations(t)
	mockOrderService.AssertExpectations(t)
}

func TestShipAssembledHandler_WithoutOrderService(t *testing.T) {
	// Arrange
	mockConsumer := &MockConsumer{}
	mockDecoder := &MockDecoder{}
	mockMetrics := &MockConsumerMetrics{}
	svc := NewService(mockConsumer, mockDecoder, nil, mockMetrics) // Без OrderService

	event := &model.ShipAssembledEvent{
		EventUUID: uuid.New(),
		OrderUUID: uuid.New(),
		UserUUID:  uuid.New(),
		BuildTime: 3600,
	}

	messageData := []byte("test-message-data")
	msg := kafka.Message{
		Topic:     "ship-assembled",
		Partition: 0,
		Offset:    100,
		Value:     messageData,
	}

	mockMetrics.On("RecordConsumerMessage", mock.Anything, "ship-assembled", int32(0), "order-ship-assembled-group", true).Return()
	mockDecoder.On("Decode", messageData).Return(event, nil)

	// Act
	err := svc.ShipAssembledHandler(context.Background(), msg)

	// Assert
	assert.NoError(t, err) // Успешно, хотя статус не обновляется
	mockDecoder.AssertExpectations(t)
	mockMetrics.AssertExpectations(t)
}

func TestShipAssembledHandler_WithoutMetrics(t *testing.T) {
	// Arrange
	mockConsumer := &MockConsumer{}
	mockDecoder := &MockDecoder{}
	mockOrderService := mocks.NewOrderService(t)
	svc := NewService(mockConsumer, mockDecoder, mockOrderService, nil) // Без метрик

	orderUUID := uuid.New()
	event := &model.ShipAssembledEvent{
		EventUUID: uuid.New(),
		OrderUUID: orderUUID,
		UserUUID:  uuid.New(),
		BuildTime: 3600,
	}

	messageData := []byte("test-message-data")
	msg := kafka.Message{
		Topic:     "ship-assembled",
		Partition: 0,
		Offset:    100,
		Value:     messageData,
	}

	mockDecoder.On("Decode", messageData).Return(event, nil)
	mockOrderService.EXPECT().UpdateOrderStatus(mock.Anything, orderUUID.String(), model.StatusAssembled).Return(nil).Once()

	// Act
	err := svc.ShipAssembledHandler(context.Background(), msg)

	// Assert
	assert.NoError(t, err)
	mockDecoder.AssertExpectations(t)
	mockOrderService.AssertExpectations(t)
}

func TestShipAssembledHandler_TableDriven(t *testing.T) {
	tests := []struct {
		name          string
		decodeError   error
		updateError   error
		expectedError bool
	}{
		{
			name:          "успешная обработка события",
			decodeError:   nil,
			updateError:   nil,
			expectedError: false,
		},
		{
			name:          "ошибка декодирования",
			decodeError:   errors.New("invalid protobuf"),
			updateError:   nil,
			expectedError: true,
		},
		{
			name:          "ошибка обновления статуса",
			decodeError:   nil,
			updateError:   errors.New("database error"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockConsumer := &MockConsumer{}
			mockDecoder := &MockDecoder{}
			mockOrderService := mocks.NewOrderService(t)
			mockMetrics := &MockConsumerMetrics{}
			svc := NewService(mockConsumer, mockDecoder, mockOrderService, mockMetrics)

			orderUUID := uuid.New()
			event := &model.ShipAssembledEvent{
				EventUUID: uuid.New(),
				OrderUUID: orderUUID,
				UserUUID:  uuid.New(),
				BuildTime: 3600,
			}

			messageData := []byte("test-message-data")
			msg := kafka.Message{
				Topic:     "ship-assembled",
				Partition: 0,
				Offset:    100,
				Value:     messageData,
			}

			mockMetrics.On("RecordConsumerMessage", mock.Anything, "ship-assembled", int32(0), "order-ship-assembled-group", true).Return()

			if tt.decodeError != nil {
				mockDecoder.On("Decode", messageData).Return(nil, tt.decodeError)
				mockMetrics.On("RecordConsumerMessage", mock.Anything, "ship-assembled", int32(0), "order-ship-assembled-group", false).Return()
			} else {
				mockDecoder.On("Decode", messageData).Return(event, nil)
				if tt.updateError != nil {
					mockOrderService.EXPECT().UpdateOrderStatus(mock.Anything, orderUUID.String(), model.StatusAssembled).Return(tt.updateError).Once()
					mockMetrics.On("RecordConsumerMessage", mock.Anything, "ship-assembled", int32(0), "order-ship-assembled-group", false).Return()
				} else {
					mockOrderService.EXPECT().UpdateOrderStatus(mock.Anything, orderUUID.String(), model.StatusAssembled).Return(nil).Once()
				}
			}

			// Act
			err := svc.ShipAssembledHandler(context.Background(), msg)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockDecoder.AssertExpectations(t)
			mockMetrics.AssertExpectations(t)
			if tt.decodeError == nil {
				mockOrderService.AssertExpectations(t)
			}
		})
	}
}

func TestRunConsumer(t *testing.T) {
	// Arrange
	mockConsumer := &MockConsumer{}
	mockDecoder := &MockDecoder{}
	mockOrderService := mocks.NewOrderService(t)
	mockMetrics := &MockConsumerMetrics{}
	svc := NewService(mockConsumer, mockDecoder, mockOrderService, mockMetrics)

	mockConsumer.On("Consume", mock.Anything, mock.AnythingOfType("kafka.MessageHandler")).Return(nil)

	// Act
	err := svc.RunConsumer(context.Background())

	// Assert
	assert.NoError(t, err)
	mockConsumer.AssertExpectations(t)
}

func TestRunConsumer_Error(t *testing.T) {
	// Arrange
	mockConsumer := &MockConsumer{}
	mockDecoder := &MockDecoder{}
	mockOrderService := mocks.NewOrderService(t)
	mockMetrics := &MockConsumerMetrics{}
	svc := NewService(mockConsumer, mockDecoder, mockOrderService, mockMetrics)

	expectedError := errors.New("consumer error")
	mockConsumer.On("Consume", mock.Anything, mock.AnythingOfType("kafka.MessageHandler")).Return(expectedError)

	// Act
	err := svc.RunConsumer(context.Background())

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockConsumer.AssertExpectations(t)
}
