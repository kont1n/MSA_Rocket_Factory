//go:build integration

package postgres_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

func (s *PostgresRepositorySuite) TestGetOrder_Success() {
	// Сначала создаем заказ
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New(), uuid.New()}
	testOrder := &model.Order{
		UserUUID:      userUUID,
		PartUUIDs:     partUUIDs,
		TotalPrice:    1500.50,
		PaymentMethod: "credit_card",
		Status:        model.StatusPendingPayment,
	}

	createdOrder, err := s.repository.CreateOrder(context.Background(), testOrder)
	s.Require().NoError(err)

	// Теперь получаем заказ по ID
	result, err := s.repository.GetOrder(context.Background(), createdOrder.OrderUUID)

	// Проверяем результат
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), createdOrder.OrderUUID, result.OrderUUID)
	assert.Equal(s.T(), userUUID, result.UserUUID)
	assert.Equal(s.T(), partUUIDs, result.PartUUIDs)
	assert.Equal(s.T(), float32(1500.50), result.TotalPrice)
	assert.Equal(s.T(), "credit_card", result.PaymentMethod)
	assert.Equal(s.T(), model.StatusPendingPayment, result.Status)
}

func (s *PostgresRepositorySuite) TestGetOrder_NotFound() {
	// Используем несуществующий UUID
	nonExistentUUID := uuid.New()

	// Вызываем метод репозитория
	result, err := s.repository.GetOrder(context.Background(), nonExistentUUID)

	// Проверяем результат
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.ErrorIs(s.T(), err, model.ErrOrderNotFound)
}

func (s *PostgresRepositorySuite) TestGetOrder_WithTransactionData() {
	// Создаем заказ с данными транзакции
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New()}
	transactionUUID := uuid.New()
	testOrder := &model.Order{
		UserUUID:        userUUID,
		PartUUIDs:       partUUIDs,
		TotalPrice:      750.25,
		TransactionUUID: transactionUUID,
		PaymentMethod:   "bank_transfer",
		Status:          model.StatusPaid,
	}

	createdOrder, err := s.repository.CreateOrder(context.Background(), testOrder)
	s.Require().NoError(err)

	// Получаем заказ
	result, err := s.repository.GetOrder(context.Background(), createdOrder.OrderUUID)

	// Проверяем результат
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), transactionUUID, result.TransactionUUID)
	assert.Equal(s.T(), "bank_transfer", result.PaymentMethod)
	assert.Equal(s.T(), model.StatusPaid, result.Status)
}

func (s *PostgresRepositorySuite) TestGetOrder_DatabaseError() {
	// Закрываем соединение для симуляции ошибки базы данных
	s.pool.Close()

	orderUUID := uuid.New()

	// Вызываем метод репозитория
	result, err := s.repository.GetOrder(context.Background(), orderUUID)

	// Проверяем результат
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.NotErrorIs(s.T(), err, model.ErrOrderNotFound)

	// Восстанавливаем соединение для последующих тестов
	s.SetupSuite()
}
