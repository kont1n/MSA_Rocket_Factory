package converter

import (
	"github.com/google/uuid"
	"github.com/kont1n/MSA_Rocket_Factory/payment/internal/model"
	paymentV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/payment/v1"
	"github.com/stretchr/testify/assert"
)

func (s *ConverterSuite) TestToModelOrder_Success() {
	// Подготовка
	orderUUID := uuid.New()
	userUUID := uuid.New()

	req := &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID.String(),
		UserUuid:      userUUID.String(),
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	// Выполнение
	result, err := ToModelOrder(req)

	// Проверка
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), orderUUID, result.OrderUuid)
	assert.Equal(s.T(), userUUID, result.UserUuid)
	assert.Equal(s.T(), "CARD", result.PaymentMethod)
}

func (s *ConverterSuite) TestToModelOrder_InvalidOrderUUID() {
	// Подготовка
	userUUID := uuid.New()

	req := &paymentV1.PayOrderRequest{
		OrderUuid:     "invalid-uuid",
		UserUuid:      userUUID.String(),
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	// Выполнение
	result, err := ToModelOrder(req)

	// Проверка
	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.Order{}, result)
}

func (s *ConverterSuite) TestToModelOrder_InvalidUserUUID() {
	// Подготовка
	orderUUID := uuid.New()

	req := &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID.String(),
		UserUuid:      "invalid-uuid",
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	// Выполнение
	result, err := ToModelOrder(req)

	// Проверка
	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.Order{}, result)
}

func (s *ConverterSuite) TestToModelOrder_EmptyOrderUUID() {
	// Подготовка
	userUUID := uuid.New()

	req := &paymentV1.PayOrderRequest{
		OrderUuid:     "",
		UserUuid:      userUUID.String(),
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	// Выполнение
	result, err := ToModelOrder(req)

	// Проверка
	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.Order{}, result)
}

func (s *ConverterSuite) TestToModelOrder_EmptyUserUUID() {
	// Подготовка
	orderUUID := uuid.New()

	req := &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID.String(),
		UserUuid:      "",
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	// Выполнение
	result, err := ToModelOrder(req)

	// Проверка
	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.Order{}, result)
}

func (s *ConverterSuite) TestToModelOrder_AllPaymentMethods() {
	// Подготовка
	orderUUID := uuid.New()
	userUUID := uuid.New()

	testCases := []struct {
		paymentMethod paymentV1.PaymentMethod
		expected      string
	}{
		{paymentV1.PaymentMethod_PAYMENT_METHOD_CARD, "CARD"},
		{paymentV1.PaymentMethod_PAYMENT_METHOD_SBP, "CARD"},
		{paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD, "CARD"},
		{paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY, "CARD"},
		{paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED, "CARD"},
	}

	for _, tc := range testCases {
		req := &paymentV1.PayOrderRequest{
			OrderUuid:     orderUUID.String(),
			UserUuid:      userUUID.String(),
			PaymentMethod: tc.paymentMethod,
		}

		// Выполнение
		result, err := ToModelOrder(req)

		// Проверка
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), orderUUID, result.OrderUuid)
		assert.Equal(s.T(), userUUID, result.UserUuid)
		assert.Equal(s.T(), tc.expected, result.PaymentMethod)
	}
}
