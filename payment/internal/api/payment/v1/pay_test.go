package v1

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kont1n/MSA_Rocket_Factory/payment/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
	paymentV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/payment/v1"
)

func init() {
	// Инициализация логгера для всех тестов
	_ = logger.InitSimple("info", false)
}

func TestPayOrder_Success(t *testing.T) {
	// Arrange
	mockService := new(MockPaymentService)
	api := NewAPI(mockService)

	orderUUID := uuid.New()
	userUUID := uuid.New()
	transactionUUID := uuid.New()

	req := &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID.String(),
		UserUuid:      userUUID.String(),
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	mockService.On("Pay", mock.Anything, mock.MatchedBy(func(o model.Order) bool {
		return o.OrderUuid == orderUUID && o.UserUuid == userUUID
	})).Return(transactionUUID, nil)

	ctx := context.Background()

	// Act
	resp, err := api.PayOrder(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, transactionUUID.String(), resp.TransactionUuid)
	mockService.AssertExpectations(t)
}

func TestPayOrder_InvalidOrderUUID(t *testing.T) {
	// Arrange
	mockService := new(MockPaymentService)
	api := NewAPI(mockService)

	req := &paymentV1.PayOrderRequest{
		OrderUuid:     "invalid-uuid",
		UserUuid:      uuid.New().String(),
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	ctx := context.Background()

	// Act
	resp, err := api.PayOrder(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestPayOrder_InvalidUserUUID(t *testing.T) {
	// Arrange
	mockService := new(MockPaymentService)
	api := NewAPI(mockService)

	req := &paymentV1.PayOrderRequest{
		OrderUuid:     uuid.New().String(),
		UserUuid:      "invalid-uuid",
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	ctx := context.Background()

	// Act
	resp, err := api.PayOrder(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestPayOrder_EmptyOrderUUID(t *testing.T) {
	// Arrange
	mockService := new(MockPaymentService)
	api := NewAPI(mockService)

	req := &paymentV1.PayOrderRequest{
		OrderUuid:     "",
		UserUuid:      uuid.New().String(),
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	ctx := context.Background()

	// Act
	resp, err := api.PayOrder(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestPayOrder_EmptyUserUUID(t *testing.T) {
	// Arrange
	mockService := new(MockPaymentService)
	api := NewAPI(mockService)

	req := &paymentV1.PayOrderRequest{
		OrderUuid:     uuid.New().String(),
		UserUuid:      "",
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	ctx := context.Background()

	// Act
	resp, err := api.PayOrder(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestPayOrder_PaymentServiceError(t *testing.T) {
	// Arrange
	mockService := new(MockPaymentService)
	api := NewAPI(mockService)

	orderUUID := uuid.New()
	userUUID := uuid.New()

	req := &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID.String(),
		UserUuid:      userUUID.String(),
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	mockService.On("Pay", mock.Anything, mock.Anything).
		Return(uuid.Nil, errors.New("payment processing error"))

	ctx := context.Background()

	// Act
	resp, err := api.PayOrder(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	mockService.AssertExpectations(t)
}

func TestPayOrder_DifferentPaymentMethods(t *testing.T) {
	tests := []struct {
		name          string
		paymentMethod paymentV1.PaymentMethod
	}{
		{
			name:          "CARD",
			paymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
		},
		{
			name:          "SBP",
			paymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_SBP,
		},
		{
			name:          "CREDIT_CARD",
			paymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
		},
		{
			name:          "INVESTOR_MONEY",
			paymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockPaymentService)
			api := NewAPI(mockService)

			orderUUID := uuid.New()
			userUUID := uuid.New()
			transactionUUID := uuid.New()

			req := &paymentV1.PayOrderRequest{
				OrderUuid:     orderUUID.String(),
				UserUuid:      userUUID.String(),
				PaymentMethod: tt.paymentMethod,
			}

			mockService.On("Pay", mock.Anything, mock.Anything).
				Return(transactionUUID, nil)

			ctx := context.Background()

			// Act
			resp, err := api.PayOrder(ctx, req)

			// Assert
			require.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, transactionUUID.String(), resp.TransactionUuid)
			mockService.AssertExpectations(t)
		})
	}
}

type testContextKey string

func TestPayOrder_WithContext(t *testing.T) {
	// Arrange
	mockService := new(MockPaymentService)
	api := NewAPI(mockService)

	orderUUID := uuid.New()
	userUUID := uuid.New()
	transactionUUID := uuid.New()

	req := &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID.String(),
		UserUuid:      userUUID.String(),
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	mockService.On("Pay", mock.Anything, mock.Anything).
		Return(transactionUUID, nil)

	// Создаем контекст со значением
	ctx := context.WithValue(context.Background(), testContextKey("test_key"), "test_value")

	// Act
	resp, err := api.PayOrder(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	mockService.AssertExpectations(t)
}

func TestPayOrder_TransactionUUIDFormat(t *testing.T) {
	// Arrange
	mockService := new(MockPaymentService)
	api := NewAPI(mockService)

	orderUUID := uuid.New()
	userUUID := uuid.New()
	transactionUUID := uuid.New()

	req := &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID.String(),
		UserUuid:      userUUID.String(),
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	mockService.On("Pay", mock.Anything, mock.Anything).
		Return(transactionUUID, nil)

	ctx := context.Background()

	// Act
	resp, err := api.PayOrder(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)

	// Проверяем, что transactionUUID является валидным UUID
	parsedUUID, err := uuid.Parse(resp.TransactionUuid)
	require.NoError(t, err)
	assert.Equal(t, transactionUUID, parsedUUID)
}

func TestPayOrder_TableDriven(t *testing.T) {
	tests := []struct {
		name          string
		orderUUID     string
		userUUID      string
		paymentMethod paymentV1.PaymentMethod
		serviceError  error
		expectError   bool
		expectedCode  codes.Code
	}{
		{
			name:          "успешная оплата",
			orderUUID:     uuid.New().String(),
			userUUID:      uuid.New().String(),
			paymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
			serviceError:  nil,
			expectError:   false,
		},
		{
			name:          "невалидный order UUID",
			orderUUID:     "invalid",
			userUUID:      uuid.New().String(),
			paymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
			serviceError:  nil,
			expectError:   true,
			expectedCode:  codes.InvalidArgument,
		},
		{
			name:          "невалидный user UUID",
			orderUUID:     uuid.New().String(),
			userUUID:      "invalid",
			paymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
			serviceError:  nil,
			expectError:   true,
			expectedCode:  codes.InvalidArgument,
		},
		{
			name:          "ошибка сервиса оплаты",
			orderUUID:     uuid.New().String(),
			userUUID:      uuid.New().String(),
			paymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_SBP,
			serviceError:  errors.New("service error"),
			expectError:   true,
			expectedCode:  codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockPaymentService)
			api := NewAPI(mockService)

			req := &paymentV1.PayOrderRequest{
				OrderUuid:     tt.orderUUID,
				UserUuid:      tt.userUUID,
				PaymentMethod: tt.paymentMethod,
			}

			if tt.serviceError != nil {
				mockService.On("Pay", mock.Anything, mock.Anything).
					Return(uuid.Nil, tt.serviceError)
			} else if !tt.expectError {
				mockService.On("Pay", mock.Anything, mock.Anything).
					Return(uuid.New(), nil)
			}

			ctx := context.Background()

			// Act
			resp, err := api.PayOrder(ctx, req)

			// Assert
			if tt.expectError {
				require.Error(t, err)
				assert.Nil(t, resp)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tt.expectedCode, st.Code())
			} else {
				require.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}
