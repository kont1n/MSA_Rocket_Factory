//go:build integration

package postgres_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

func (s *PostgresRepositorySuite) TestCreateOrder_Success() {
	// Подготавливаем тестовые данные
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New(), uuid.New()}
	testOrder := &model.Order{
		UserUUID:      userUUID,
		PartUUIDs:     partUUIDs,
		TotalPrice:    1500.50,
		PaymentMethod: "credit_card",
		Status:        model.StatusPendingPayment,
	}

	// Вызываем метод репозитория
	result, err := s.repository.CreateOrder(context.Background(), testOrder)

	// Проверяем результат
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.NotEqual(s.T(), uuid.Nil, result.OrderUUID) // UUID должен быть сгенерирован
	assert.Equal(s.T(), userUUID, result.UserUUID)
	assert.Equal(s.T(), partUUIDs, result.PartUUIDs)
	assert.Equal(s.T(), float32(1500.50), result.TotalPrice)
	assert.Equal(s.T(), "credit_card", result.PaymentMethod)
	assert.Equal(s.T(), model.StatusPendingPayment, result.Status)

	// Проверяем, что заказ действительно сохранился в базе
	savedOrder, err := s.repository.GetOrder(context.Background(), result.OrderUUID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), result.OrderUUID, savedOrder.OrderUUID)
	assert.Equal(s.T(), userUUID, savedOrder.UserUUID)
}

func (s *PostgresRepositorySuite) TestCreateOrder_WithTransactionUUID() {
	// Подготавливаем тестовые данные с транзакцией
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New()}
	transactionUUID := uuid.New()
	testOrder := &model.Order{
		UserUUID:        userUUID,
		PartUUIDs:       partUUIDs,
		TotalPrice:      500.0,
		TransactionUUID: transactionUUID,
		PaymentMethod:   "bank_transfer",
		Status:          model.StatusPaid,
	}

	// Вызываем метод репозитория
	result, err := s.repository.CreateOrder(context.Background(), testOrder)

	// Проверяем результат
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), transactionUUID, result.TransactionUUID)
	assert.Equal(s.T(), "bank_transfer", result.PaymentMethod)
	assert.Equal(s.T(), model.StatusPaid, result.Status)
}

func (s *PostgresRepositorySuite) TestCreateOrder_MultipleOrders() {
	// Создаем несколько заказов для проверки независимости
	userUUID1 := uuid.New()
	userUUID2 := uuid.New()

	order1 := &model.Order{
		UserUUID:      userUUID1,
		PartUUIDs:     []uuid.UUID{uuid.New()},
		TotalPrice:    100.0,
		PaymentMethod: "cash",
		Status:        model.StatusPendingPayment,
	}

	order2 := &model.Order{
		UserUUID:      userUUID2,
		PartUUIDs:     []uuid.UUID{uuid.New(), uuid.New()},
		TotalPrice:    200.0,
		PaymentMethod: "credit_card",
		Status:        model.StatusPendingPayment,
	}

	// Создаем оба заказа
	result1, err1 := s.repository.CreateOrder(context.Background(), order1)
	result2, err2 := s.repository.CreateOrder(context.Background(), order2)

	// Проверяем результаты
	assert.NoError(s.T(), err1)
	assert.NoError(s.T(), err2)
	assert.NotEqual(s.T(), result1.OrderUUID, result2.OrderUUID) // UUID должны быть разными
	assert.Equal(s.T(), userUUID1, result1.UserUUID)
	assert.Equal(s.T(), userUUID2, result2.UserUUID)
}
