package inmemory_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

func (s *InMemoryOrderRepositorySuite) TestGetOrder_Success() {
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

func (s *InMemoryOrderRepositorySuite) TestGetOrder_NotFound() {
	// Используем несуществующий UUID
	nonExistentUUID := uuid.New()

	// Вызываем метод репозитория
	result, err := s.repository.GetOrder(context.Background(), nonExistentUUID)

	// Проверяем результат
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.ErrorIs(s.T(), err, model.ErrOrderNotFound)
}

func (s *InMemoryOrderRepositorySuite) TestGetOrder_WithTransactionData() {
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

func (s *InMemoryOrderRepositorySuite) TestGetOrder_MultipleOrders() {
	// Создаем несколько заказов и проверяем, что получаем правильный
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
		Status:        model.StatusPaid,
	}

	// Создаем оба заказа
	createdOrder1, err1 := s.repository.CreateOrder(context.Background(), order1)
	createdOrder2, err2 := s.repository.CreateOrder(context.Background(), order2)
	s.Require().NoError(err1)
	s.Require().NoError(err2)

	// Получаем первый заказ
	result1, err := s.repository.GetOrder(context.Background(), createdOrder1.OrderUUID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), createdOrder1.OrderUUID, result1.OrderUUID)
	assert.Equal(s.T(), userUUID1, result1.UserUUID)
	assert.Equal(s.T(), float32(100.0), result1.TotalPrice)
	assert.Equal(s.T(), "cash", result1.PaymentMethod)

	// Получаем второй заказ
	result2, err := s.repository.GetOrder(context.Background(), createdOrder2.OrderUUID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), createdOrder2.OrderUUID, result2.OrderUUID)
	assert.Equal(s.T(), userUUID2, result2.UserUUID)
	assert.Equal(s.T(), float32(200.0), result2.TotalPrice)
	assert.Equal(s.T(), "credit_card", result2.PaymentMethod)
}
