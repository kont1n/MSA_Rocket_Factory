package converter

import (
	paymentV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/payment/v1"
	"github.com/stretchr/testify/assert"
)

func (s *ConverterSuite) TestToProtoPaymentMethod_Card() {
	// Подготовка
	paymentMethod := "CARD"

	// Выполнение
	result := ToProtoPaymentMethod(paymentMethod)

	// Проверка
	assert.Equal(s.T(), paymentV1.PaymentMethod_PAYMENT_METHOD_CARD, result)
}

func (s *ConverterSuite) TestToProtoPaymentMethod_Sbp() {
	// Подготовка
	paymentMethod := "SBP"

	// Выполнение
	result := ToProtoPaymentMethod(paymentMethod)

	// Проверка
	assert.Equal(s.T(), paymentV1.PaymentMethod_PAYMENT_METHOD_SBP, result)
}

func (s *ConverterSuite) TestToProtoPaymentMethod_CreditCard() {
	// Подготовка
	paymentMethod := "CREDIT_CARD"

	// Выполнение
	result := ToProtoPaymentMethod(paymentMethod)

	// Проверка
	assert.Equal(s.T(), paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD, result)
}

func (s *ConverterSuite) TestToProtoPaymentMethod_InvestorMoney() {
	// Подготовка
	paymentMethod := "INVESTOR_MONEY"

	// Выполнение
	result := ToProtoPaymentMethod(paymentMethod)

	// Проверка
	assert.Equal(s.T(), paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY, result)
}

func (s *ConverterSuite) TestToProtoPaymentMethod_Unknown() {
	// Подготовка
	paymentMethod := "UNKNOWN"

	// Выполнение
	result := ToProtoPaymentMethod(paymentMethod)

	// Проверка
	assert.Equal(s.T(), paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED, result)
}

func (s *ConverterSuite) TestToProtoPaymentMethod_InvalidMethod() {
	// Подготовка
	paymentMethod := "INVALID_METHOD"

	// Выполнение
	result := ToProtoPaymentMethod(paymentMethod)

	// Проверка
	assert.Equal(s.T(), paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED, result)
}

func (s *ConverterSuite) TestToProtoPaymentMethod_EmptyString() {
	// Подготовка
	paymentMethod := ""

	// Выполнение
	result := ToProtoPaymentMethod(paymentMethod)

	// Проверка
	assert.Equal(s.T(), paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED, result)
}

func (s *ConverterSuite) TestToProtoPaymentMethod_Lowercase() {
	// Подготовка
	paymentMethod := "card"

	// Выполнение
	result := ToProtoPaymentMethod(paymentMethod)

	// Проверка
	assert.Equal(s.T(), paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED, result)
}
