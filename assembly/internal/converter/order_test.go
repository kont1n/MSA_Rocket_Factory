package converter

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kont1n/MSA_Rocket_Factory/assembly/internal/model"
	eventsV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/events/v1"
)

func TestToModelOrder_Success(t *testing.T) {
	// Arrange
	eventUUID := uuid.New()
	orderUUID := uuid.New()
	userUUID := uuid.New()
	transactionUUID := uuid.New()

	protoOrder := &eventsV1.OrderPaid{
		EventUuid:       eventUUID.String(),
		OrderUuid:       orderUUID.String(),
		UserUuid:        userUUID.String(),
		PaymentMethod:   "CARD",
		TransactionUuid: transactionUUID.String(),
	}

	// Act
	result, err := ToModelOrder(protoOrder)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, eventUUID, result.EventUUID)
	assert.Equal(t, orderUUID, result.OrderUUID)
	assert.Equal(t, userUUID, result.UserUUID)
	assert.Equal(t, "CARD", result.PaymentMethod)
	assert.Equal(t, transactionUUID, result.TransactionUUID)
}

func TestToModelOrder_DifferentPaymentMethods(t *testing.T) {
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
			protoOrder := &eventsV1.OrderPaid{
				EventUuid:       uuid.New().String(),
				OrderUuid:       uuid.New().String(),
				UserUuid:        uuid.New().String(),
				PaymentMethod:   tt.paymentMethod,
				TransactionUuid: uuid.New().String(),
			}

			// Act
			result, err := ToModelOrder(protoOrder)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, tt.paymentMethod, result.PaymentMethod)
		})
	}
}

func TestToModelOrder_InvalidEventUUID(t *testing.T) {
	// Arrange
	protoOrder := &eventsV1.OrderPaid{
		EventUuid:       "invalid-uuid",
		OrderUuid:       uuid.New().String(),
		UserUuid:        uuid.New().String(),
		PaymentMethod:   "CARD",
		TransactionUuid: uuid.New().String(),
	}

	// Act
	result, err := ToModelOrder(protoOrder)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, model.ErrConvertFromKafkaEvent, err)
}

func TestToModelOrder_InvalidOrderUUID(t *testing.T) {
	// Arrange
	protoOrder := &eventsV1.OrderPaid{
		EventUuid:       uuid.New().String(),
		OrderUuid:       "not-a-valid-uuid",
		UserUuid:        uuid.New().String(),
		PaymentMethod:   "CARD",
		TransactionUuid: uuid.New().String(),
	}

	// Act
	result, err := ToModelOrder(protoOrder)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, model.ErrConvertFromKafkaEvent, err)
}

func TestToModelOrder_InvalidUserUUID(t *testing.T) {
	// Arrange
	protoOrder := &eventsV1.OrderPaid{
		EventUuid:       uuid.New().String(),
		OrderUuid:       uuid.New().String(),
		UserUuid:        "bad-uuid-format",
		PaymentMethod:   "SBP",
		TransactionUuid: uuid.New().String(),
	}

	// Act
	result, err := ToModelOrder(protoOrder)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, model.ErrConvertFromKafkaEvent, err)
}

func TestToModelOrder_InvalidTransactionUUID(t *testing.T) {
	// Arrange
	protoOrder := &eventsV1.OrderPaid{
		EventUuid:       uuid.New().String(),
		OrderUuid:       uuid.New().String(),
		UserUuid:        uuid.New().String(),
		PaymentMethod:   "CARD",
		TransactionUuid: "invalid-transaction-uuid",
	}

	// Act
	result, err := ToModelOrder(protoOrder)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, model.ErrConvertFromKafkaEvent, err)
}

func TestToModelOrder_EmptyPaymentMethod(t *testing.T) {
	// Arrange
	protoOrder := &eventsV1.OrderPaid{
		EventUuid:       uuid.New().String(),
		OrderUuid:       uuid.New().String(),
		UserUuid:        uuid.New().String(),
		PaymentMethod:   "",
		TransactionUuid: uuid.New().String(),
	}

	// Act
	result, err := ToModelOrder(protoOrder)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "", result.PaymentMethod)
}

func TestToModelOrder_MultipleInvalidUUIDs(t *testing.T) {
	tests := []struct {
		name            string
		eventUUID       string
		orderUUID       string
		userUUID        string
		transactionUUID string
		expectError     bool
	}{
		{
			name:            "все UUID валидны",
			eventUUID:       uuid.New().String(),
			orderUUID:       uuid.New().String(),
			userUUID:        uuid.New().String(),
			transactionUUID: uuid.New().String(),
			expectError:     false,
		},
		{
			name:            "eventUUID невалидный",
			eventUUID:       "bad",
			orderUUID:       uuid.New().String(),
			userUUID:        uuid.New().String(),
			transactionUUID: uuid.New().String(),
			expectError:     true,
		},
		{
			name:            "orderUUID невалидный",
			eventUUID:       uuid.New().String(),
			orderUUID:       "bad",
			userUUID:        uuid.New().String(),
			transactionUUID: uuid.New().String(),
			expectError:     true,
		},
		{
			name:            "userUUID невалидный",
			eventUUID:       uuid.New().String(),
			orderUUID:       uuid.New().String(),
			userUUID:        "bad",
			transactionUUID: uuid.New().String(),
			expectError:     true,
		},
		{
			name:            "transactionUUID невалидный",
			eventUUID:       uuid.New().String(),
			orderUUID:       uuid.New().String(),
			userUUID:        uuid.New().String(),
			transactionUUID: "bad",
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			protoOrder := &eventsV1.OrderPaid{
				EventUuid:       tt.eventUUID,
				OrderUuid:       tt.orderUUID,
				UserUuid:        tt.userUUID,
				PaymentMethod:   "CARD",
				TransactionUuid: tt.transactionUUID,
			}

			// Act
			result, err := ToModelOrder(protoOrder)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, model.ErrConvertFromKafkaEvent, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}
