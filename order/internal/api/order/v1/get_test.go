package v1

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/service/mocks"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
	orderV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/openapi/order/v1"
)

func init() {
	_ = logger.InitSimple("error", false)
}

func TestAPI_GetOrderByUUID_Success(t *testing.T) {
	// Arrange
	mockService := mocks.NewOrderService(t)
	api := NewAPI(mockService)

	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUID1 := uuid.New()
	partUUID2 := uuid.New()
	transactionUUID := uuid.New()

	params := orderV1.GetOrderByUUIDParams{
		OrderUUID: orderUUID,
	}

	expectedOrder := &model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUUIDs:       []uuid.UUID{partUUID1, partUUID2},
		TransactionUUID: transactionUUID,
		PaymentMethod:   "CARD",
		Status:          "PAID",
		TotalPrice:      1500.50,
	}

	mockService.EXPECT().
		GetOrder(mock.Anything, orderUUID).
		Return(expectedOrder, nil).
		Once()

	// Act
	res, err := api.GetOrderByUUID(context.Background(), params)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, res)

	response, ok := res.(*orderV1.OrderDto)
	assert.True(t, ok)
	assert.Equal(t, orderUUID, response.OrderUUID)
	assert.Equal(t, userUUID, response.UserUUID)
	assert.Equal(t, 2, len(response.PartUuids))
	assert.True(t, response.TotalPrice.Set)
	assert.Equal(t, float32(1500.50), response.TotalPrice.Value)
	assert.True(t, response.TransactionUUID.Set)
	assert.Equal(t, transactionUUID, response.TransactionUUID.Value)
	assert.True(t, response.PaymentMethod.Set)
	assert.Equal(t, orderV1.PaymentMethod("CARD"), response.PaymentMethod.Value)
	assert.Equal(t, orderV1.OrderStatus("PAID"), response.Status)

	mockService.AssertExpectations(t)
}

func TestAPI_GetOrderByUUID_ServiceNil(t *testing.T) {
	// Arrange
	api := &api{orderService: nil}

	params := orderV1.GetOrderByUUIDParams{
		OrderUUID: uuid.New(),
	}

	// Act
	res, err := api.GetOrderByUUID(context.Background(), params)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, res)

	errorResponse, ok := res.(*orderV1.InternalServerError)
	assert.True(t, ok)
	assert.Equal(t, 500, errorResponse.Code)
	assert.Equal(t, "order service not available", errorResponse.Message)
}

func TestAPI_GetOrderByUUID_OrderNotFound(t *testing.T) {
	// Arrange
	mockService := mocks.NewOrderService(t)
	api := NewAPI(mockService)

	orderUUID := uuid.New()
	params := orderV1.GetOrderByUUIDParams{
		OrderUUID: orderUUID,
	}

	mockService.EXPECT().
		GetOrder(mock.Anything, orderUUID).
		Return(nil, model.ErrOrderNotFound).
		Once()

	// Act
	res, err := api.GetOrderByUUID(context.Background(), params)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, res)

	errorResponse, ok := res.(*orderV1.NotFoundError)
	assert.True(t, ok)
	assert.Equal(t, 404, errorResponse.Code)
	assert.Equal(t, "order not found", errorResponse.Message)

	mockService.AssertExpectations(t)
}

func TestAPI_GetOrderByUUID_ConvertFromRepoError(t *testing.T) {
	// Arrange
	mockService := mocks.NewOrderService(t)
	api := NewAPI(mockService)

	orderUUID := uuid.New()
	params := orderV1.GetOrderByUUIDParams{
		OrderUUID: orderUUID,
	}

	mockService.EXPECT().
		GetOrder(mock.Anything, orderUUID).
		Return(nil, model.ErrConvertFromRepo).
		Once()

	// Act
	res, err := api.GetOrderByUUID(context.Background(), params)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, res)

	errorResponse, ok := res.(*orderV1.InternalServerError)
	assert.True(t, ok)
	assert.Equal(t, 500, errorResponse.Code)
	assert.Equal(t, "cannot convert order from repository", errorResponse.Message)

	mockService.AssertExpectations(t)
}

func TestAPI_GetOrderByUUID_DifferentStatuses(t *testing.T) {
	tests := []struct {
		name          string
		orderStatus   string
		paymentMethod string
	}{
		{
			name:          "заказ в ожидании оплаты",
			orderStatus:   "PENDING_PAYMENT",
			paymentMethod: "",
		},
		{
			name:          "заказ оплачен",
			orderStatus:   "PAID",
			paymentMethod: "CARD",
		},
		{
			name:          "заказ собран",
			orderStatus:   "ASSEMBLED",
			paymentMethod: "CARD",
		},
		{
			name:          "заказ отменён",
			orderStatus:   "CANCELLED",
			paymentMethod: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := mocks.NewOrderService(t)
			api := NewAPI(mockService)

			orderUUID := uuid.New()
			params := orderV1.GetOrderByUUIDParams{
				OrderUUID: orderUUID,
			}

			expectedOrder := &model.Order{
				OrderUUID:       orderUUID,
				UserUUID:        uuid.New(),
				Status:          model.OrderStatus(tt.orderStatus),
				PaymentMethod:   tt.paymentMethod,
				TotalPrice:      1000.0,
				TransactionUUID: uuid.New(),
			}

			mockService.EXPECT().
				GetOrder(mock.Anything, orderUUID).
				Return(expectedOrder, nil).
				Once()

			// Act
			res, err := api.GetOrderByUUID(context.Background(), params)

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, res)

			response, ok := res.(*orderV1.OrderDto)
			assert.True(t, ok)
			assert.Equal(t, orderV1.OrderStatus(tt.orderStatus), response.Status)

			mockService.AssertExpectations(t)
		})
	}
}

func TestAPI_GetOrderByUUID_TableDriven(t *testing.T) {
	tests := []struct {
		name               string
		setupMock          func(*mocks.OrderService, uuid.UUID)
		expectedStatusCode int
		expectedMessage    string
		expectError        bool
	}{
		{
			name: "успешное получение заказа",
			setupMock: func(m *mocks.OrderService, orderUUID uuid.UUID) {
				m.EXPECT().
					GetOrder(mock.Anything, orderUUID).
					Return(&model.Order{
						OrderUUID:  orderUUID,
						UserUUID:   uuid.New(),
						Status:     "PAID",
						TotalPrice: 1000.0,
					}, nil).
					Once()
			},
			expectedStatusCode: 0,
			expectError:        false,
		},
		{
			name: "заказ не найден",
			setupMock: func(m *mocks.OrderService, orderUUID uuid.UUID) {
				m.EXPECT().
					GetOrder(mock.Anything, orderUUID).
					Return(nil, model.ErrOrderNotFound).
					Once()
			},
			expectedStatusCode: 404,
			expectedMessage:    "order not found",
			expectError:        false,
		},
		{
			name: "ошибка конвертации из репозитория",
			setupMock: func(m *mocks.OrderService, orderUUID uuid.UUID) {
				m.EXPECT().
					GetOrder(mock.Anything, orderUUID).
					Return(nil, model.ErrConvertFromRepo).
					Once()
			},
			expectedStatusCode: 500,
			expectedMessage:    "cannot convert order from repository",
			expectError:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := mocks.NewOrderService(t)
			orderUUID := uuid.New()
			tt.setupMock(mockService, orderUUID)
			api := NewAPI(mockService)

			params := orderV1.GetOrderByUUIDParams{
				OrderUUID: orderUUID,
			}

			// Act
			res, err := api.GetOrderByUUID(context.Background(), params)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)

				switch tt.expectedStatusCode {
				case 0: // Success
					response, ok := res.(*orderV1.OrderDto)
					assert.True(t, ok, "ожидался OrderDto")
					assert.Equal(t, orderUUID, response.OrderUUID)
				case 404:
					errorResponse, ok := res.(*orderV1.NotFoundError)
					assert.True(t, ok, "ожидался NotFoundError")
					assert.Equal(t, tt.expectedStatusCode, errorResponse.Code)
					assert.Equal(t, tt.expectedMessage, errorResponse.Message)
				case 500:
					errorResponse, ok := res.(*orderV1.InternalServerError)
					assert.True(t, ok, "ожидался InternalServerError")
					assert.Equal(t, tt.expectedStatusCode, errorResponse.Code)
					assert.Equal(t, tt.expectedMessage, errorResponse.Message)
				}
			}

			mockService.AssertExpectations(t)
		})
	}
}
