package tests

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/kont1n/MSA_Rocket_Factory/payment/internal/model"
)

func (s *ServiceSuite) TestPaySuccess() {
	// Тестовые данные
	orderUUID := uuid.New()
	userUUID := uuid.New()
	testOrder := model.Order{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: "CARD",
		TransactionId: uuid.Nil,
	}

	// Вызов метода
	transactionUUID, err := s.service.Pay(s.ctx, testOrder)

	// Проверка результата
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), transactionUUID)
	assert.NotEqual(s.T(), uuid.Nil, transactionUUID)
	assert.NotEqual(s.T(), orderUUID, transactionUUID) // UUID транзакции должен отличаться от UUID заказа
}

func (s *ServiceSuite) TestPayWithDifferentPaymentMethods() {
	// Тестируем различные методы оплаты
	paymentMethods := []string{"CARD", "SBP", "CREDIT_CARD", "INVESTOR_MONEY"}

	for _, method := range paymentMethods {
		s.T().Run("PaymentMethod_"+method, func(t *testing.T) {
			order := model.Order{
				OrderUuid:     uuid.New(),
				UserUuid:      uuid.New(),
				PaymentMethod: method,
				TransactionId: uuid.Nil,
			}

			transactionUUID, err := s.service.Pay(s.ctx, order)

			assert.NoError(t, err)
			assert.NotNil(t, transactionUUID)
			assert.NotEqual(t, uuid.Nil, transactionUUID)
		})
	}
}

func (s *ServiceSuite) TestPayWithContext() {
	// Тестируем с контекстом с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	order := model.Order{
		OrderUuid:     uuid.New(),
		UserUuid:      uuid.New(),
		PaymentMethod: "CARD",
		TransactionId: uuid.Nil,
	}

	transactionUUID, err := s.service.Pay(ctx, order)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), transactionUUID)
}

func (s *ServiceSuite) TestPayGeneratesUniqueUUIDs() {
	// Проверяем, что каждый вызов генерирует уникальный UUID
	order := model.Order{
		OrderUuid:     uuid.New(),
		UserUuid:      uuid.New(),
		PaymentMethod: "CARD",
		TransactionId: uuid.Nil,
	}

	uuid1, err1 := s.service.Pay(s.ctx, order)
	uuid2, err2 := s.service.Pay(s.ctx, order)

	assert.NoError(s.T(), err1)
	assert.NoError(s.T(), err2)
	assert.NotEqual(s.T(), uuid1, uuid2) // UUID должны быть разными
}

func (s *ServiceSuite) TestPayWithNilContext() {
	// Тестируем с nil контекстом
	order := model.Order{
		OrderUuid:     uuid.New(),
		UserUuid:      uuid.New(),
		PaymentMethod: "CARD",
		TransactionId: uuid.Nil,
	}

	transactionUUID, err := s.service.Pay(nil, order)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), transactionUUID)
}
