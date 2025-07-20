package converter

import (
	"github.com/stretchr/testify/assert"

	paymentV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/payment/v1"
)

func (s *ConverterSuite) TestPaymentToProto_Card() {
	// Подготовка
	paymentMethod := "CARD"

	// Выполнение
	result := PaymentToProto(paymentMethod)

	// Проверка
	assert.Equal(s.T(), paymentV1.PaymentMethod_PAYMENT_METHOD_CARD, result)
}

func (s *ConverterSuite) TestPaymentToProto_Sbp() {
	// Подготовка
	paymentMethod := "SBP"

	// Выполнение
	result := PaymentToProto(paymentMethod)

	// Проверка
	assert.Equal(s.T(), paymentV1.PaymentMethod_PAYMENT_METHOD_SBP, result)
}

func (s *ConverterSuite) TestPaymentToProto_CreditCard() {
	// Подготовка
	paymentMethod := "CREDIT_CARD"

	// Выполнение
	result := PaymentToProto(paymentMethod)

	// Проверка
	assert.Equal(s.T(), paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD, result)
}

func (s *ConverterSuite) TestPaymentToProto_InvestorMoney() {
	// Подготовка
	paymentMethod := "INVESTOR_MONEY"

	// Выполнение
	result := PaymentToProto(paymentMethod)

	// Проверка
	assert.Equal(s.T(), paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY, result)
}

func (s *ConverterSuite) TestPaymentToProto_Unknown() {
	// Подготовка
	paymentMethod := "UNKNOWN"

	// Выполнение
	result := PaymentToProto(paymentMethod)

	// Проверка
	assert.Equal(s.T(), paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED, result)
}

func (s *ConverterSuite) TestPaymentToProto_InvalidMethod() {
	// Подготовка
	paymentMethod := "INVALID_METHOD"

	// Выполнение
	result := PaymentToProto(paymentMethod)

	// Проверка
	assert.Equal(s.T(), paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED, result)
}

func (s *ConverterSuite) TestPaymentToProto_EmptyString() {
	// Подготовка
	paymentMethod := ""

	// Выполнение
	result := PaymentToProto(paymentMethod)

	// Проверка
	assert.Equal(s.T(), paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED, result)
}

func (s *ConverterSuite) TestPaymentToProto_Lowercase() {
	// Подготовка
	paymentMethod := "card"

	// Выполнение
	result := PaymentToProto(paymentMethod)

	// Проверка
	assert.Equal(s.T(), paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED, result)
}
