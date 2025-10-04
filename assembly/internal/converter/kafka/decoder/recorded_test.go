package decoder

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/kont1n/MSA_Rocket_Factory/assembly/internal/model"
	eventsV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/events/v1"
)

func TestDecoder_Decode_Success(t *testing.T) {
	// Arrange
	eventUUID := uuid.New()
	orderUUID := uuid.New()
	userUUID := uuid.New()
	transactionUUID := uuid.New()

	protoEvent := &eventsV1.OrderPaid{
		EventUuid:       eventUUID.String(),
		OrderUuid:       orderUUID.String(),
		UserUuid:        userUUID.String(),
		PaymentMethod:   "CARD",
		TransactionUuid: transactionUUID.String(),
	}

	data, err := proto.Marshal(protoEvent)
	require.NoError(t, err)

	decoder := NewAssemblyRecordedDecoder()

	// Act
	result, err := decoder.Decode(data)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, eventUUID, result.EventUUID)
	assert.Equal(t, orderUUID, result.OrderUUID)
	assert.Equal(t, userUUID, result.UserUUID)
	assert.Equal(t, "CARD", result.PaymentMethod)
	assert.Equal(t, transactionUUID, result.TransactionUUID)
}

func TestDecoder_Decode_DifferentPaymentMethods(t *testing.T) {
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
			protoEvent := &eventsV1.OrderPaid{
				EventUuid:       uuid.New().String(),
				OrderUuid:       uuid.New().String(),
				UserUuid:        uuid.New().String(),
				PaymentMethod:   tt.paymentMethod,
				TransactionUuid: uuid.New().String(),
			}

			data, err := proto.Marshal(protoEvent)
			require.NoError(t, err)

			decoder := NewAssemblyRecordedDecoder()

			// Act
			result, err := decoder.Decode(data)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, tt.paymentMethod, result.PaymentMethod)
		})
	}
}

func TestDecoder_Decode_InvalidProtobuf(t *testing.T) {
	// Arrange
	invalidData := []byte("this is not valid protobuf data")
	decoder := NewAssemblyRecordedDecoder()

	// Act
	result, err := decoder.Decode(invalidData)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal protobuf")
	assert.Equal(t, model.OrderPaidEvent{}, result)
}

func TestDecoder_Decode_EmptyData(t *testing.T) {
	// Arrange
	emptyData := []byte{}
	decoder := NewAssemblyRecordedDecoder()

	// Act
	result, err := decoder.Decode(emptyData)

	// Assert
	// Пустой protobuf содержит пустые строки для UUID, что вызовет ошибку парсинга
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to convert protobuf to model")
	assert.Equal(t, model.OrderPaidEvent{}, result)
}

func TestDecoder_Decode_InvalidEventUUID(t *testing.T) {
	// Arrange
	protoEvent := &eventsV1.OrderPaid{
		EventUuid:       "invalid-uuid",
		OrderUuid:       uuid.New().String(),
		UserUuid:        uuid.New().String(),
		PaymentMethod:   "CARD",
		TransactionUuid: uuid.New().String(),
	}

	data, err := proto.Marshal(protoEvent)
	require.NoError(t, err)

	decoder := NewAssemblyRecordedDecoder()

	// Act
	result, err := decoder.Decode(data)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to convert protobuf to model")
	assert.Equal(t, model.OrderPaidEvent{}, result)
}

func TestDecoder_Decode_InvalidOrderUUID(t *testing.T) {
	// Arrange
	protoEvent := &eventsV1.OrderPaid{
		EventUuid:       uuid.New().String(),
		OrderUuid:       "not-a-valid-uuid",
		UserUuid:        uuid.New().String(),
		PaymentMethod:   "SBP",
		TransactionUuid: uuid.New().String(),
	}

	data, err := proto.Marshal(protoEvent)
	require.NoError(t, err)

	decoder := NewAssemblyRecordedDecoder()

	// Act
	result, err := decoder.Decode(data)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to convert protobuf to model")
	assert.Equal(t, model.OrderPaidEvent{}, result)
}

func TestDecoder_Decode_InvalidUserUUID(t *testing.T) {
	// Arrange
	protoEvent := &eventsV1.OrderPaid{
		EventUuid:       uuid.New().String(),
		OrderUuid:       uuid.New().String(),
		UserUuid:        "bad-uuid-format",
		PaymentMethod:   "CARD",
		TransactionUuid: uuid.New().String(),
	}

	data, err := proto.Marshal(protoEvent)
	require.NoError(t, err)

	decoder := NewAssemblyRecordedDecoder()

	// Act
	result, err := decoder.Decode(data)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to convert protobuf to model")
	assert.Equal(t, model.OrderPaidEvent{}, result)
}

func TestDecoder_Decode_InvalidTransactionUUID(t *testing.T) {
	// Arrange
	protoEvent := &eventsV1.OrderPaid{
		EventUuid:       uuid.New().String(),
		OrderUuid:       uuid.New().String(),
		UserUuid:        uuid.New().String(),
		PaymentMethod:   "CARD",
		TransactionUuid: "invalid-transaction",
	}

	data, err := proto.Marshal(protoEvent)
	require.NoError(t, err)

	decoder := NewAssemblyRecordedDecoder()

	// Act
	result, err := decoder.Decode(data)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to convert protobuf to model")
	assert.Equal(t, model.OrderPaidEvent{}, result)
}

func TestDecoder_Decode_NilData(t *testing.T) {
	// Arrange
	decoder := NewAssemblyRecordedDecoder()

	// Act
	result, err := decoder.Decode(nil)

	// Assert
	// nil данные приводят к пустым строкам UUID, что вызывает ошибку
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to convert protobuf to model")
	assert.Equal(t, model.OrderPaidEvent{}, result)
}

func TestDecoder_Decode_EmptyPaymentMethod(t *testing.T) {
	// Arrange
	protoEvent := &eventsV1.OrderPaid{
		EventUuid:       uuid.New().String(),
		OrderUuid:       uuid.New().String(),
		UserUuid:        uuid.New().String(),
		PaymentMethod:   "",
		TransactionUuid: uuid.New().String(),
	}

	data, err := proto.Marshal(protoEvent)
	require.NoError(t, err)

	decoder := NewAssemblyRecordedDecoder()

	// Act
	result, err := decoder.Decode(data)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "", result.PaymentMethod)
}

func TestNewAssemblyRecordedDecoder(t *testing.T) {
	// Act
	decoder := NewAssemblyRecordedDecoder()

	// Assert
	assert.NotNil(t, decoder)
}
