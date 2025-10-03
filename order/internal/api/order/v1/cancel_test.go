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

func TestAPI_CancelOrder_Success(t *testing.T) {
	// Arrange
	mockService := mocks.NewOrderService(t)
	api := NewAPI(mockService)

	orderUUID := uuid.New()
	params := orderV1.CancelOrderParams{
		OrderUUID: orderUUID,
	}

	expectedOrder := &model.Order{
		OrderUUID: orderUUID,
		Status:    "CANCELLED",
	}

	mockService.EXPECT().
		CancelOrder(mock.Anything, mock.MatchedBy(func(order *model.Order) bool {
			return order.OrderUUID == orderUUID
		})).
		Return(expectedOrder, nil).
		Once()

	// Act
	res, err := api.CancelOrder(context.Background(), params)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, res)

	_, ok := res.(*orderV1.CancelOrderNoContent)
	assert.True(t, ok)

	mockService.AssertExpectations(t)
}

func TestAPI_CancelOrder_OrderNotFound(t *testing.T) {
	// Arrange
	mockService := mocks.NewOrderService(t)
	api := NewAPI(mockService)

	params := orderV1.CancelOrderParams{
		OrderUUID: uuid.New(),
	}

	mockService.EXPECT().
		CancelOrder(mock.Anything, mock.Anything).
		Return(nil, model.ErrOrderNotFound).
		Once()

	// Act
	res, err := api.CancelOrder(context.Background(), params)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, res)

	errorResponse, ok := res.(*orderV1.NotFoundError)
	assert.True(t, ok)
	assert.Equal(t, 404, errorResponse.Code)
	assert.Equal(t, "order not found", errorResponse.Message)

	mockService.AssertExpectations(t)
}

func TestAPI_CancelOrder_ConvertFromRepoError(t *testing.T) {
	// Arrange
	mockService := mocks.NewOrderService(t)
	api := NewAPI(mockService)

	params := orderV1.CancelOrderParams{
		OrderUUID: uuid.New(),
	}

	mockService.EXPECT().
		CancelOrder(mock.Anything, mock.Anything).
		Return(nil, model.ErrConvertFromRepo).
		Once()

	// Act
	res, err := api.CancelOrder(context.Background(), params)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, res)

	errorResponse, ok := res.(*orderV1.InternalServerError)
	assert.True(t, ok)
	assert.Equal(t, 500, errorResponse.Code)
	assert.Equal(t, "cannot convert order from repository", errorResponse.Message)

	mockService.AssertExpectations(t)
}

func TestAPI_CancelOrder_AlreadyPaid(t *testing.T) {
	// Arrange
	mockService := mocks.NewOrderService(t)
	api := NewAPI(mockService)

	params := orderV1.CancelOrderParams{
		OrderUUID: uuid.New(),
	}

	mockService.EXPECT().
		CancelOrder(mock.Anything, mock.Anything).
		Return(nil, model.ErrPaid).
		Once()

	// Act
	res, err := api.CancelOrder(context.Background(), params)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, res)

	errorResponse, ok := res.(*orderV1.BadRequestError)
	assert.True(t, ok)
	assert.Equal(t, 400, errorResponse.Code)
	assert.Equal(t, "order already paid", errorResponse.Message)

	mockService.AssertExpectations(t)
}

func TestAPI_CancelOrder_AlreadyCancelled(t *testing.T) {
	// Arrange
	mockService := mocks.NewOrderService(t)
	api := NewAPI(mockService)

	params := orderV1.CancelOrderParams{
		OrderUUID: uuid.New(),
	}

	mockService.EXPECT().
		CancelOrder(mock.Anything, mock.Anything).
		Return(nil, model.ErrCancelled).
		Once()

	// Act
	res, err := api.CancelOrder(context.Background(), params)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, res)

	errorResponse, ok := res.(*orderV1.BadRequestError)
	assert.True(t, ok)
	assert.Equal(t, 400, errorResponse.Code)
	assert.Equal(t, "order already cancelled", errorResponse.Message)

	mockService.AssertExpectations(t)
}

func TestAPI_CancelOrder_TableDriven(t *testing.T) {
	tests := []struct {
		name               string
		setupMock          func(*mocks.OrderService)
		expectedStatusCode int
		expectedMessage    string
		expectError        bool
	}{
		{
			name: "успешная отмена заказа",
			setupMock: func(m *mocks.OrderService) {
				m.EXPECT().
					CancelOrder(mock.Anything, mock.Anything).
					Return(&model.Order{
						OrderUUID: uuid.New(),
						Status:    "CANCELLED",
					}, nil).
					Once()
			},
			expectedStatusCode: 204,
			expectError:        false,
		},
		{
			name: "заказ не найден",
			setupMock: func(m *mocks.OrderService) {
				m.EXPECT().
					CancelOrder(mock.Anything, mock.Anything).
					Return(nil, model.ErrOrderNotFound).
					Once()
			},
			expectedStatusCode: 404,
			expectedMessage:    "order not found",
			expectError:        false,
		},
		{
			name: "заказ уже оплачен",
			setupMock: func(m *mocks.OrderService) {
				m.EXPECT().
					CancelOrder(mock.Anything, mock.Anything).
					Return(nil, model.ErrPaid).
					Once()
			},
			expectedStatusCode: 400,
			expectedMessage:    "order already paid",
			expectError:        false,
		},
		{
			name: "заказ уже отменён",
			setupMock: func(m *mocks.OrderService) {
				m.EXPECT().
					CancelOrder(mock.Anything, mock.Anything).
					Return(nil, model.ErrCancelled).
					Once()
			},
			expectedStatusCode: 400,
			expectedMessage:    "order already cancelled",
			expectError:        false,
		},
		{
			name: "ошибка конвертации из репозитория",
			setupMock: func(m *mocks.OrderService) {
				m.EXPECT().
					CancelOrder(mock.Anything, mock.Anything).
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
			tt.setupMock(mockService)
			api := NewAPI(mockService)

			params := orderV1.CancelOrderParams{
				OrderUUID: uuid.New(),
			}

			// Act
			res, err := api.CancelOrder(context.Background(), params)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)

				switch tt.expectedStatusCode {
				case 204: // NoContent
					_, ok := res.(*orderV1.CancelOrderNoContent)
					assert.True(t, ok, "ожидался CancelOrderNoContent")
				case 400:
					errorResponse, ok := res.(*orderV1.BadRequestError)
					assert.True(t, ok, "ожидался BadRequestError")
					assert.Equal(t, tt.expectedStatusCode, errorResponse.Code)
					assert.Equal(t, tt.expectedMessage, errorResponse.Message)
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
