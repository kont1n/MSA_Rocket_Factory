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

func TestAPI_PayOrder_Success(t *testing.T) {
	// Arrange
	mockService := mocks.NewOrderService(t)
	api := NewAPI(mockService)

	orderUUID := uuid.New()
	transactionUUID := uuid.New()
	userUUID := uuid.New()

	req := &orderV1.PayOrderRequest{
		PaymentMethod: orderV1.PaymentMethodCARD,
	}
	params := orderV1.PayOrderParams{
		OrderUUID: orderUUID,
	}

	expectedOrder := &model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		TransactionUUID: transactionUUID,
		PaymentMethod:   "CARD",
		Status:          "PAID",
		TotalPrice:      1000.0,
	}

	mockService.EXPECT().
		PayOrder(mock.Anything, mock.MatchedBy(func(order *model.Order) bool {
			return order.OrderUUID == orderUUID &&
				order.PaymentMethod == "CARD"
		})).
		Return(expectedOrder, nil).
		Once()

	// Act
	res, err := api.PayOrder(context.Background(), req, params)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, res)

	response, ok := res.(*orderV1.PayOrderResponse)
	assert.True(t, ok)
	assert.Equal(t, orderV1.TransactionUUID(transactionUUID), response.TransactionUUID)

	mockService.AssertExpectations(t)
}

func TestAPI_PayOrder_ServiceNil(t *testing.T) {
	// Arrange
	api := &api{orderService: nil}

	req := &orderV1.PayOrderRequest{
		PaymentMethod: orderV1.PaymentMethodCARD,
	}
	params := orderV1.PayOrderParams{
		OrderUUID: uuid.New(),
	}

	// Act
	res, err := api.PayOrder(context.Background(), req, params)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, res)

	errorResponse, ok := res.(*orderV1.InternalServerError)
	assert.True(t, ok)
	assert.Equal(t, 500, errorResponse.Code)
	assert.Equal(t, "order service not available", errorResponse.Message)
}

func TestAPI_PayOrder_OrderNotFound(t *testing.T) {
	// Arrange
	mockService := mocks.NewOrderService(t)
	api := NewAPI(mockService)

	req := &orderV1.PayOrderRequest{
		PaymentMethod: orderV1.PaymentMethodCARD,
	}
	params := orderV1.PayOrderParams{
		OrderUUID: uuid.New(),
	}

	mockService.EXPECT().
		PayOrder(mock.Anything, mock.Anything).
		Return(nil, model.ErrOrderNotFound).
		Once()

	// Act
	res, err := api.PayOrder(context.Background(), req, params)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, res)

	errorResponse, ok := res.(*orderV1.NotFoundError)
	assert.True(t, ok)
	assert.Equal(t, 404, errorResponse.Code)
	assert.Equal(t, "order not found", errorResponse.Message)

	mockService.AssertExpectations(t)
}

func TestAPI_PayOrder_ConvertFromRepoError(t *testing.T) {
	// Arrange
	mockService := mocks.NewOrderService(t)
	api := NewAPI(mockService)

	req := &orderV1.PayOrderRequest{
		PaymentMethod: orderV1.PaymentMethodCARD,
	}
	params := orderV1.PayOrderParams{
		OrderUUID: uuid.New(),
	}

	mockService.EXPECT().
		PayOrder(mock.Anything, mock.Anything).
		Return(nil, model.ErrConvertFromRepo).
		Once()

	// Act
	res, err := api.PayOrder(context.Background(), req, params)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, res)

	errorResponse, ok := res.(*orderV1.InternalServerError)
	assert.True(t, ok)
	assert.Equal(t, 500, errorResponse.Code)
	assert.Equal(t, "cannot convert order from repository", errorResponse.Message)

	mockService.AssertExpectations(t)
}

func TestAPI_PayOrder_ConvertFromClientError(t *testing.T) {
	// Arrange
	mockService := mocks.NewOrderService(t)
	api := NewAPI(mockService)

	req := &orderV1.PayOrderRequest{
		PaymentMethod: orderV1.PaymentMethodCARD,
	}
	params := orderV1.PayOrderParams{
		OrderUUID: uuid.New(),
	}

	mockService.EXPECT().
		PayOrder(mock.Anything, mock.Anything).
		Return(nil, model.ErrConvertFromClient).
		Once()

	// Act
	res, err := api.PayOrder(context.Background(), req, params)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, res)

	errorResponse, ok := res.(*orderV1.InternalServerError)
	assert.True(t, ok)
	assert.Equal(t, 500, errorResponse.Code)
	assert.Equal(t, "can't parse transaction", errorResponse.Message)

	mockService.AssertExpectations(t)
}

func TestAPI_PayOrder_PaymentClientError(t *testing.T) {
	// Arrange
	mockService := mocks.NewOrderService(t)
	api := NewAPI(mockService)

	req := &orderV1.PayOrderRequest{
		PaymentMethod: orderV1.PaymentMethodCARD,
	}
	params := orderV1.PayOrderParams{
		OrderUUID: uuid.New(),
	}

	mockService.EXPECT().
		PayOrder(mock.Anything, mock.Anything).
		Return(nil, model.ErrPaymentClient).
		Once()

	// Act
	res, err := api.PayOrder(context.Background(), req, params)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, res)

	errorResponse, ok := res.(*orderV1.InternalServerError)
	assert.True(t, ok)
	assert.Equal(t, 424, errorResponse.Code) // HTTP_FAILED_DEPENDENCY
	assert.Equal(t, "payment client error", errorResponse.Message)

	mockService.AssertExpectations(t)
}

func TestAPI_PayOrder_DifferentPaymentMethods(t *testing.T) {
	tests := []struct {
		name          string
		paymentMethod orderV1.PaymentMethod
	}{
		{
			name:          "оплата картой CARD",
			paymentMethod: orderV1.PaymentMethodCARD,
		},
		{
			name:          "оплата через SBP",
			paymentMethod: orderV1.PaymentMethodSBP,
		},
		{
			name:          "оплата через CREDIT_CARD",
			paymentMethod: orderV1.PaymentMethodCREDITCARD,
		},
		{
			name:          "оплата через INVESTOR_MONEY",
			paymentMethod: orderV1.PaymentMethodINVESTORMONEY,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := mocks.NewOrderService(t)
			api := NewAPI(mockService)

			orderUUID := uuid.New()
			transactionUUID := uuid.New()

			req := &orderV1.PayOrderRequest{
				PaymentMethod: tt.paymentMethod,
			}
			params := orderV1.PayOrderParams{
				OrderUUID: orderUUID,
			}

			expectedOrder := &model.Order{
				OrderUUID:       orderUUID,
				TransactionUUID: transactionUUID,
				PaymentMethod:   string(tt.paymentMethod),
			}

			mockService.EXPECT().
				PayOrder(mock.Anything, mock.MatchedBy(func(order *model.Order) bool {
					return order.PaymentMethod == string(tt.paymentMethod)
				})).
				Return(expectedOrder, nil).
				Once()

			// Act
			res, err := api.PayOrder(context.Background(), req, params)

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, res)

			response, ok := res.(*orderV1.PayOrderResponse)
			assert.True(t, ok)
			assert.Equal(t, orderV1.TransactionUUID(transactionUUID), response.TransactionUUID)

			mockService.AssertExpectations(t)
		})
	}
}

func TestAPI_PayOrder_TableDriven(t *testing.T) {
	tests := []struct {
		name               string
		setupMock          func(*mocks.OrderService)
		paymentMethod      orderV1.PaymentMethod
		expectedStatusCode int
		expectedMessage    string
		expectError        bool
	}{
		{
			name: "успешная оплата",
			setupMock: func(m *mocks.OrderService) {
				m.EXPECT().
					PayOrder(mock.Anything, mock.Anything).
					Return(&model.Order{
						OrderUUID:       uuid.New(),
						TransactionUUID: uuid.New(),
					}, nil).
					Once()
			},
			paymentMethod:      orderV1.PaymentMethodCARD,
			expectedStatusCode: 0,
			expectError:        false,
		},
		{
			name: "заказ не найден",
			setupMock: func(m *mocks.OrderService) {
				m.EXPECT().
					PayOrder(mock.Anything, mock.Anything).
					Return(nil, model.ErrOrderNotFound).
					Once()
			},
			paymentMethod:      orderV1.PaymentMethodCARD,
			expectedStatusCode: 404,
			expectedMessage:    "order not found",
			expectError:        false,
		},
		{
			name: "недоступен payment service",
			setupMock: func(m *mocks.OrderService) {
				m.EXPECT().
					PayOrder(mock.Anything, mock.Anything).
					Return(nil, model.ErrPaymentClient).
					Once()
			},
			paymentMethod:      orderV1.PaymentMethodCARD,
			expectedStatusCode: 424,
			expectedMessage:    "payment client error",
			expectError:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := mocks.NewOrderService(t)
			tt.setupMock(mockService)
			api := NewAPI(mockService)

			req := &orderV1.PayOrderRequest{
				PaymentMethod: tt.paymentMethod,
			}
			params := orderV1.PayOrderParams{
				OrderUUID: uuid.New(),
			}

			// Act
			res, err := api.PayOrder(context.Background(), req, params)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)

				switch tt.expectedStatusCode {
				case 0: // Success
					_, ok := res.(*orderV1.PayOrderResponse)
					assert.True(t, ok, "ожидался PayOrderResponse")
				case 404:
					errorResponse, ok := res.(*orderV1.NotFoundError)
					assert.True(t, ok, "ожидался NotFoundError")
					assert.Equal(t, tt.expectedStatusCode, errorResponse.Code)
					assert.Equal(t, tt.expectedMessage, errorResponse.Message)
				case 424, 500:
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
