package order

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

func init() {
	_ = logger.InitSimple("error", false)
}

// MockOrderRepository - мок для repository.OrderRepository
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) CreateOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	args := m.Called(ctx, order)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderRepository) GetOrder(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderRepository) UpdateOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	args := m.Called(ctx, order)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

// MockOrderPaidProducer - мок для OrderPaidProducer
type MockOrderPaidProducer struct {
	mock.Mock
}

func (m *MockOrderPaidProducer) ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func TestNewService(t *testing.T) {
	// Arrange
	mockRepo := &MockOrderRepository{}
	mockProducer := &MockOrderPaidProducer{}

	// Act
	svc := NewService(mockRepo, nil, nil, mockProducer)

	// Assert
	assert.NotNil(t, svc)
	assert.Equal(t, mockRepo, svc.orderRepository)
	assert.Equal(t, mockProducer, svc.orderPaidProducer)
}

func TestNewService_WithNilProducer(t *testing.T) {
	// Arrange
	mockRepo := &MockOrderRepository{}

	// Act
	svc := NewService(mockRepo, nil, nil, nil)

	// Assert
	assert.NotNil(t, svc)
	assert.Nil(t, svc.orderPaidProducer)
}

func TestUpdateOrderStatus_Success(t *testing.T) {
	// Arrange
	mockRepo := &MockOrderRepository{}
	svc := NewService(mockRepo, nil, nil, nil)

	orderUUID := uuid.New()
	existingOrder := &model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        uuid.New(),
		PartUUIDs:       []uuid.UUID{uuid.New()},
		TotalPrice:      1000.0,
		TransactionUUID: uuid.New(),
		PaymentMethod:   "CARD",
		Status:          model.StatusPendingPayment,
	}

	updatedOrder := &model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        existingOrder.UserUUID,
		PartUUIDs:       existingOrder.PartUUIDs,
		TotalPrice:      existingOrder.TotalPrice,
		TransactionUUID: existingOrder.TransactionUUID,
		PaymentMethod:   existingOrder.PaymentMethod,
		Status:          model.StatusPaid,
	}

	mockRepo.On("GetOrder", mock.Anything, orderUUID).Return(existingOrder, nil)
	mockRepo.On("UpdateOrder", mock.Anything, mock.MatchedBy(func(order *model.Order) bool {
		return order.Status == model.StatusPaid
	})).Return(updatedOrder, nil)

	// Act
	err := svc.UpdateOrderStatus(context.Background(), orderUUID.String(), model.StatusPaid)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateOrderStatus_InvalidUUID(t *testing.T) {
	// Arrange
	mockRepo := &MockOrderRepository{}
	svc := NewService(mockRepo, nil, nil, nil)

	invalidUUID := "invalid-uuid"

	// Act
	err := svc.UpdateOrderStatus(context.Background(), invalidUUID, model.StatusPaid)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid order UUID")
}

func TestUpdateOrderStatus_OrderNotFound(t *testing.T) {
	// Arrange
	mockRepo := &MockOrderRepository{}
	svc := NewService(mockRepo, nil, nil, nil)

	orderUUID := uuid.New()
	expectedError := model.ErrOrderNotFound

	mockRepo.On("GetOrder", mock.Anything, orderUUID).Return(nil, expectedError)

	// Act
	err := svc.UpdateOrderStatus(context.Background(), orderUUID.String(), model.StatusPaid)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get order")
	mockRepo.AssertExpectations(t)
}

func TestUpdateOrderStatus_RepositoryUpdateError(t *testing.T) {
	// Arrange
	mockRepo := &MockOrderRepository{}
	svc := NewService(mockRepo, nil, nil, nil)

	orderUUID := uuid.New()
	existingOrder := &model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        uuid.New(),
		PartUUIDs:       []uuid.UUID{uuid.New()},
		TotalPrice:      1000.0,
		TransactionUUID: uuid.New(),
		PaymentMethod:   "CARD",
		Status:          model.StatusPendingPayment,
	}

	expectedError := errors.New("database error")

	mockRepo.On("GetOrder", mock.Anything, orderUUID).Return(existingOrder, nil)
	mockRepo.On("UpdateOrder", mock.Anything, mock.Anything).Return(nil, expectedError)

	// Act
	err := svc.UpdateOrderStatus(context.Background(), orderUUID.String(), model.StatusPaid)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update order status")
	mockRepo.AssertExpectations(t)
}

func TestUpdateOrderStatus_DifferentStatuses(t *testing.T) {
	tests := []struct {
		name          string
		initialStatus model.OrderStatus
		newStatus     model.OrderStatus
	}{
		{
			name:          "pending to paid",
			initialStatus: model.StatusPendingPayment,
			newStatus:     model.StatusPaid,
		},
		{
			name:          "pending to cancelled",
			initialStatus: model.StatusPendingPayment,
			newStatus:     model.StatusCancelled,
		},
		{
			name:          "paid to assembled",
			initialStatus: model.StatusPaid,
			newStatus:     model.StatusAssembled,
		},
		{
			name:          "paid to cancelled",
			initialStatus: model.StatusPaid,
			newStatus:     model.StatusCancelled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockOrderRepository{}
			svc := NewService(mockRepo, nil, nil, nil)

			orderUUID := uuid.New()
			existingOrder := &model.Order{
				OrderUUID:       orderUUID,
				UserUUID:        uuid.New(),
				PartUUIDs:       []uuid.UUID{uuid.New()},
				TotalPrice:      1000.0,
				TransactionUUID: uuid.New(),
				PaymentMethod:   "CARD",
				Status:          tt.initialStatus,
			}

			updatedOrder := &model.Order{
				OrderUUID:       orderUUID,
				UserUUID:        existingOrder.UserUUID,
				PartUUIDs:       existingOrder.PartUUIDs,
				TotalPrice:      existingOrder.TotalPrice,
				TransactionUUID: existingOrder.TransactionUUID,
				PaymentMethod:   existingOrder.PaymentMethod,
				Status:          tt.newStatus,
			}

			mockRepo.On("GetOrder", mock.Anything, orderUUID).Return(existingOrder, nil)
			mockRepo.On("UpdateOrder", mock.Anything, mock.MatchedBy(func(order *model.Order) bool {
				return order.Status == tt.newStatus
			})).Return(updatedOrder, nil)

			// Act
			err := svc.UpdateOrderStatus(context.Background(), orderUUID.String(), tt.newStatus)

			// Assert
			assert.NoError(t, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSetOrderPaidProducer_Success(t *testing.T) {
	// Arrange
	mockRepo := &MockOrderRepository{}
	svc := NewService(mockRepo, nil, nil, nil)
	mockProducer := &MockOrderPaidProducer{}

	// Act
	svc.SetOrderPaidProducer(mockProducer)

	// Assert
	assert.NotNil(t, svc.orderPaidProducer)
	assert.Equal(t, mockProducer, svc.orderPaidProducer)
}

func TestSetOrderPaidProducer_ReplaceExisting(t *testing.T) {
	// Arrange
	oldProducer := &MockOrderPaidProducer{}
	mockRepo := &MockOrderRepository{}
	svc := NewService(mockRepo, nil, nil, oldProducer)

	assert.NotNil(t, svc.orderPaidProducer)
	assert.Equal(t, oldProducer, svc.orderPaidProducer)

	newProducer := &MockOrderPaidProducer{}

	// Act
	svc.SetOrderPaidProducer(newProducer)

	// Assert - проверяем что producer был заменен
	assert.NotNil(t, svc.orderPaidProducer)
	assert.Equal(t, newProducer, svc.orderPaidProducer)
}

func TestSetOrderPaidProducer_SetToNil(t *testing.T) {
	// Arrange
	oldProducer := &MockOrderPaidProducer{}
	mockRepo := &MockOrderRepository{}
	svc := NewService(mockRepo, nil, nil, oldProducer)

	// Act
	svc.SetOrderPaidProducer(nil)

	// Assert
	assert.Nil(t, svc.orderPaidProducer)
}

func TestUpdateOrderStatus_WithContext(t *testing.T) {
	// Arrange
	mockRepo := &MockOrderRepository{}
	svc := NewService(mockRepo, nil, nil, nil)

	orderUUID := uuid.New()
	existingOrder := &model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        uuid.New(),
		PartUUIDs:       []uuid.UUID{uuid.New()},
		TotalPrice:      1000.0,
		TransactionUUID: uuid.New(),
		PaymentMethod:   "CARD",
		Status:          model.StatusPendingPayment,
	}

	updatedOrder := &model.Order{
		OrderUUID: orderUUID,
		Status:    model.StatusPaid,
	}

	ctx := context.Background()

	mockRepo.On("GetOrder", ctx, orderUUID).Return(existingOrder, nil)
	mockRepo.On("UpdateOrder", ctx, mock.Anything).Return(updatedOrder, nil)

	// Act
	err := svc.UpdateOrderStatus(ctx, orderUUID.String(), model.StatusPaid)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateOrderStatus_TableDriven(t *testing.T) {
	tests := []struct {
		name          string
		orderUUID     string
		newStatus     model.OrderStatus
		getOrderError error
		updateError   error
		expectedError bool
		errorContains string
	}{
		{
			name:          "успешное обновление статуса",
			orderUUID:     uuid.New().String(),
			newStatus:     model.StatusPaid,
			getOrderError: nil,
			updateError:   nil,
			expectedError: false,
		},
		{
			name:          "невалидный UUID",
			orderUUID:     "invalid-uuid",
			newStatus:     model.StatusPaid,
			expectedError: true,
			errorContains: "invalid order UUID",
		},
		{
			name:          "заказ не найден",
			orderUUID:     uuid.New().String(),
			newStatus:     model.StatusPaid,
			getOrderError: model.ErrOrderNotFound,
			expectedError: true,
			errorContains: "failed to get order",
		},
		{
			name:          "ошибка обновления в БД",
			orderUUID:     uuid.New().String(),
			newStatus:     model.StatusPaid,
			getOrderError: nil,
			updateError:   errors.New("database error"),
			expectedError: true,
			errorContains: "failed to update order status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockOrderRepository{}
			svc := NewService(mockRepo, nil, nil, nil)

			if tt.orderUUID != "invalid-uuid" {
				orderUUID, _ := uuid.Parse(tt.orderUUID)
				existingOrder := &model.Order{
					OrderUUID:       orderUUID,
					UserUUID:        uuid.New(),
					PartUUIDs:       []uuid.UUID{uuid.New()},
					TotalPrice:      1000.0,
					TransactionUUID: uuid.New(),
					PaymentMethod:   "CARD",
					Status:          model.StatusPendingPayment,
				}

				if tt.getOrderError != nil {
					mockRepo.On("GetOrder", mock.Anything, orderUUID).Return(nil, tt.getOrderError)
				} else {
					mockRepo.On("GetOrder", mock.Anything, orderUUID).Return(existingOrder, nil)

					updatedOrder := &model.Order{
						OrderUUID: orderUUID,
						Status:    tt.newStatus,
					}

					if tt.updateError != nil {
						mockRepo.On("UpdateOrder", mock.Anything, mock.Anything).Return(nil, tt.updateError)
					} else {
						mockRepo.On("UpdateOrder", mock.Anything, mock.Anything).Return(updatedOrder, nil)
					}
				}
			}

			// Act
			err := svc.UpdateOrderStatus(context.Background(), tt.orderUUID, tt.newStatus)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}

			if tt.orderUUID != "invalid-uuid" {
				mockRepo.AssertExpectations(t)
			}
		})
	}
}
