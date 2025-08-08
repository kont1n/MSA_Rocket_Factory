//go:build integration

package postgres_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

func (s *PostgresRepositorySuite) TestUpdateOrder_Success() {
	// Сначала создаем заказ
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New()}
	testOrder := &model.Order{
		UserUUID:      userUUID,
		PartUUIDs:     partUUIDs,
		TotalPrice:    1000.0,
		PaymentMethod: "",
		Status:        model.StatusPendingPayment,
	}

	createdOrder, err := s.repository.CreateOrder(context.Background(), testOrder)
	s.Require().NoError(err)

	// Обновляем заказ с данными платежа
	transactionUUID := uuid.New()
	createdOrder.TransactionUUID = transactionUUID
	createdOrder.PaymentMethod = "credit_card"
	createdOrder.Status = model.StatusPaid

	// Вызываем метод обновления
	result, err := s.repository.UpdateOrder(context.Background(), createdOrder)

	// Проверяем результат
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), createdOrder.OrderUUID, result.OrderUUID)
	assert.Equal(s.T(), transactionUUID, result.TransactionUUID)
	assert.Equal(s.T(), "credit_card", result.PaymentMethod)
	assert.Equal(s.T(), model.StatusPaid, result.Status)

	// Проверяем, что изменения сохранились в базе
	savedOrder, err := s.repository.GetOrder(context.Background(), createdOrder.OrderUUID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), transactionUUID, savedOrder.TransactionUUID)
	assert.Equal(s.T(), "credit_card", savedOrder.PaymentMethod)
	assert.Equal(s.T(), model.StatusPaid, savedOrder.Status)
}

func (s *PostgresRepositorySuite) TestUpdateOrder_ChangeStatus() {
	// Создаем заказ
	userUUID := uuid.New()
	testOrder := &model.Order{
		UserUUID:      userUUID,
		PartUUIDs:     []uuid.UUID{uuid.New()},
		TotalPrice:    500.0,
		PaymentMethod: "cash",
		Status:        model.StatusPendingPayment,
	}

	createdOrder, err := s.repository.CreateOrder(context.Background(), testOrder)
	s.Require().NoError(err)

	// Обновляем статус на отмененный
	createdOrder.Status = model.StatusCancelled

	result, err := s.repository.UpdateOrder(context.Background(), createdOrder)

	// Проверяем результат
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), model.StatusCancelled, result.Status)

	// Проверяем в базе данных
	savedOrder, err := s.repository.GetOrder(context.Background(), createdOrder.OrderUUID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), model.StatusCancelled, savedOrder.Status)
}

func (s *PostgresRepositorySuite) TestUpdateOrder_NonExistentOrder() {
	// Пытаемся обновить несуществующий заказ
	nonExistentOrder := &model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUUIDs:       []uuid.UUID{uuid.New()},
		TotalPrice:      100.0,
		TransactionUUID: uuid.New(),
		PaymentMethod:   "credit_card",
		Status:          model.StatusPaid,
	}

	// Вызываем метод обновления
	result, err := s.repository.UpdateOrder(context.Background(), nonExistentOrder)

	// Операция должна завершиться без ошибки, но ничего не обновить
	// В PostgreSQL UPDATE не возвращает ошибку, если строка не найдена
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), nonExistentOrder.OrderUUID, result.OrderUUID)
}

func (s *PostgresRepositorySuite) TestUpdateOrder_MultipleUpdates() {
	// Создаем заказ
	userUUID := uuid.New()
	testOrder := &model.Order{
		UserUUID:      userUUID,
		PartUUIDs:     []uuid.UUID{uuid.New()},
		TotalPrice:    750.0,
		PaymentMethod: "",
		Status:        model.StatusPendingPayment,
	}

	createdOrder, err := s.repository.CreateOrder(context.Background(), testOrder)
	s.Require().NoError(err)

	// Первое обновление - добавляем метод платежа
	createdOrder.PaymentMethod = "credit_card"
	result1, err1 := s.repository.UpdateOrder(context.Background(), createdOrder)
	assert.NoError(s.T(), err1)
	assert.Equal(s.T(), "credit_card", result1.PaymentMethod)

	// Второе обновление - добавляем транзакцию и меняем статус
	transactionUUID := uuid.New()
	createdOrder.TransactionUUID = transactionUUID
	createdOrder.Status = model.StatusPaid
	result2, err2 := s.repository.UpdateOrder(context.Background(), createdOrder)
	assert.NoError(s.T(), err2)
	assert.Equal(s.T(), transactionUUID, result2.TransactionUUID)
	assert.Equal(s.T(), model.StatusPaid, result2.Status)

	// Проверяем финальное состояние в базе
	finalOrder, err := s.repository.GetOrder(context.Background(), createdOrder.OrderUUID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "credit_card", finalOrder.PaymentMethod)
	assert.Equal(s.T(), transactionUUID, finalOrder.TransactionUUID)
	assert.Equal(s.T(), model.StatusPaid, finalOrder.Status)
}
