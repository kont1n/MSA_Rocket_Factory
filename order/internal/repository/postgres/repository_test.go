package postgres

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/testcontainers/postgres"
)

func init() {
	_ = logger.InitSimple("error", false)
}

// setupTestDB создает тестовую базу данных PostgreSQL используя testcontainers
func setupTestDB(t *testing.T) *pgxpool.Pool {
	ctx := context.Background()

	// Запускаем PostgreSQL контейнер без указания имени (генерируется автоматически)
	container, err := postgres.NewContainer(ctx,
		postgres.WithImageName("postgres:17-alpine"),
		postgres.WithDatabase("test_order_db"),
		postgres.WithAuth("test_user", "test_password"),
		// Не указываем WithContainerName - testcontainers сгенерирует уникальное имя
	)
	require.NoError(t, err)

	// Получаем строку подключения
	connStr, err := container.ConnectionString(ctx)
	require.NoError(t, err)

	// Создаем пул подключений
	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)

	// Выполняем миграции
	repo := &repository{db: pool}
	err = repo.Migrate("../../../migrations")
	require.NoError(t, err)

	// Cleanup function для teardown
	t.Cleanup(func() {
		pool.Close()
		_ = container.Terminate(ctx)
	})

	return pool
}

func TestNewRepository(t *testing.T) {
	pool := setupTestDB(t)

	repo := NewRepository(pool, "../../../migrations")

	assert.NotNil(t, repo)
	assert.NotNil(t, repo.db)
}

func TestRepository_CreateOrder_Success(t *testing.T) {
	pool := setupTestDB(t)
	repo := &repository{db: pool}
	ctx := context.Background()

	// Arrange
	order := &model.Order{
		UserUUID:        uuid.New(),
		PartUUIDs:       []uuid.UUID{uuid.New(), uuid.New()},
		TotalPrice:      1500.50,
		TransactionUUID: uuid.New(),
		PaymentMethod:   "CARD",
		Status:          model.StatusPendingPayment,
	}

	// Act
	createdOrder, err := repo.CreateOrder(ctx, order)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, createdOrder)
	assert.NotEqual(t, uuid.Nil, createdOrder.OrderUUID)
	assert.Equal(t, order.UserUUID, createdOrder.UserUUID)
	assert.Equal(t, order.PartUUIDs, createdOrder.PartUUIDs)
	assert.Equal(t, order.TotalPrice, createdOrder.TotalPrice)
	assert.Equal(t, order.TransactionUUID, createdOrder.TransactionUUID)
	assert.Equal(t, order.PaymentMethod, createdOrder.PaymentMethod)
	assert.Equal(t, order.Status, createdOrder.Status)
}

func TestRepository_CreateOrder_EmptyTransactionUUID(t *testing.T) {
	pool := setupTestDB(t)
	repo := &repository{db: pool}
	ctx := context.Background()

	// Arrange - заказ без TransactionUUID (pending payment)
	order := &model.Order{
		UserUUID:        uuid.New(),
		PartUUIDs:       []uuid.UUID{uuid.New()},
		TotalPrice:      500.0,
		TransactionUUID: uuid.Nil,
		PaymentMethod:   "",
		Status:          model.StatusPendingPayment,
	}

	// Act
	createdOrder, err := repo.CreateOrder(ctx, order)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, createdOrder)
	assert.NotEqual(t, uuid.Nil, createdOrder.OrderUUID)
}

func TestRepository_GetOrder_Success(t *testing.T) {
	pool := setupTestDB(t)
	repo := &repository{db: pool}
	ctx := context.Background()

	// Arrange - создаем заказ
	originalOrder := &model.Order{
		UserUUID:        uuid.New(),
		PartUUIDs:       []uuid.UUID{uuid.New()},
		TotalPrice:      1000.0,
		TransactionUUID: uuid.New(),
		PaymentMethod:   "CARD",
		Status:          model.StatusPaid,
	}
	createdOrder, err := repo.CreateOrder(ctx, originalOrder)
	require.NoError(t, err)

	// Act
	fetchedOrder, err := repo.GetOrder(ctx, createdOrder.OrderUUID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, fetchedOrder)
	assert.Equal(t, createdOrder.OrderUUID, fetchedOrder.OrderUUID)
	assert.Equal(t, createdOrder.UserUUID, fetchedOrder.UserUUID)
	assert.Equal(t, createdOrder.PartUUIDs, fetchedOrder.PartUUIDs)
	assert.Equal(t, createdOrder.TotalPrice, fetchedOrder.TotalPrice)
	assert.Equal(t, createdOrder.PaymentMethod, fetchedOrder.PaymentMethod)
	assert.Equal(t, createdOrder.Status, fetchedOrder.Status)
}

func TestRepository_GetOrder_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	repo := &repository{db: pool}
	ctx := context.Background()

	// Arrange
	nonExistentUUID := uuid.New()

	// Act
	order, err := repo.GetOrder(ctx, nonExistentUUID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, order)
	assert.Equal(t, model.ErrOrderNotFound, err)
}

func TestRepository_UpdateOrder_Success(t *testing.T) {
	pool := setupTestDB(t)
	repo := &repository{db: pool}
	ctx := context.Background()

	// Arrange - создаем заказ
	originalOrder := &model.Order{
		UserUUID:        uuid.New(),
		PartUUIDs:       []uuid.UUID{uuid.New()},
		TotalPrice:      1500.0,
		TransactionUUID: uuid.Nil,
		PaymentMethod:   "",
		Status:          model.StatusPendingPayment,
	}
	createdOrder, err := repo.CreateOrder(ctx, originalOrder)
	require.NoError(t, err)

	// Act - обновляем заказ (оплачиваем)
	createdOrder.Status = model.StatusPaid
	createdOrder.PaymentMethod = "CARD"
	createdOrder.TransactionUUID = uuid.New()

	updatedOrder, err := repo.UpdateOrder(ctx, createdOrder)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, updatedOrder)
	assert.Equal(t, model.StatusPaid, updatedOrder.Status)
	assert.Equal(t, "CARD", updatedOrder.PaymentMethod)
	assert.NotEqual(t, uuid.Nil, updatedOrder.TransactionUUID)

	// Проверяем, что изменения сохранились в БД
	fetchedOrder, err := repo.GetOrder(ctx, createdOrder.OrderUUID)
	require.NoError(t, err)
	assert.Equal(t, model.StatusPaid, fetchedOrder.Status)
	assert.Equal(t, "CARD", fetchedOrder.PaymentMethod)
	assert.Equal(t, createdOrder.TransactionUUID, fetchedOrder.TransactionUUID)
}

func TestRepository_UpdateOrder_ChangeStatus(t *testing.T) {
	pool := setupTestDB(t)
	repo := &repository{db: pool}
	ctx := context.Background()

	// Arrange - создаем оплаченный заказ
	order := &model.Order{
		UserUUID:        uuid.New(),
		PartUUIDs:       []uuid.UUID{uuid.New()},
		TotalPrice:      2000.0,
		TransactionUUID: uuid.New(),
		PaymentMethod:   "CARD",
		Status:          model.StatusPaid,
	}
	createdOrder, err := repo.CreateOrder(ctx, order)
	require.NoError(t, err)

	// Act - отменяем заказ
	createdOrder.Status = model.StatusCancelled

	updatedOrder, err := repo.UpdateOrder(ctx, createdOrder)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, model.StatusCancelled, updatedOrder.Status)

	// Проверяем в БД
	fetchedOrder, err := repo.GetOrder(ctx, createdOrder.OrderUUID)
	require.NoError(t, err)
	assert.Equal(t, model.StatusCancelled, fetchedOrder.Status)
}

func TestRepository_CreateOrder_MultipleParts(t *testing.T) {
	pool := setupTestDB(t)
	repo := &repository{db: pool}
	ctx := context.Background()

	// Arrange - заказ с несколькими деталями
	partUUIDs := []uuid.UUID{
		uuid.New(),
		uuid.New(),
		uuid.New(),
		uuid.New(),
	}

	order := &model.Order{
		UserUUID:        uuid.New(),
		PartUUIDs:       partUUIDs,
		TotalPrice:      5000.0,
		TransactionUUID: uuid.New(),
		PaymentMethod:   "CARD",
		Status:          model.StatusPaid,
	}

	// Act
	createdOrder, err := repo.CreateOrder(ctx, order)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, createdOrder)
	assert.Len(t, createdOrder.PartUUIDs, 4)
	assert.Equal(t, partUUIDs, createdOrder.PartUUIDs)

	// Проверяем через GetOrder
	fetchedOrder, err := repo.GetOrder(ctx, createdOrder.OrderUUID)
	require.NoError(t, err)
	assert.Len(t, fetchedOrder.PartUUIDs, 4)
	assert.Equal(t, partUUIDs, fetchedOrder.PartUUIDs)
}

func TestRepository_TableDriven_DifferentStatuses(t *testing.T) {
	pool := setupTestDB(t)
	repo := &repository{db: pool}
	ctx := context.Background()

	tests := []struct {
		name          string
		status        model.OrderStatus
		paymentMethod string
		hasTxUUID     bool
	}{
		{
			name:          "pending payment заказ",
			status:        model.StatusPendingPayment,
			paymentMethod: "",
			hasTxUUID:     false,
		},
		{
			name:          "оплаченный заказ картой",
			status:        model.StatusPaid,
			paymentMethod: "CARD",
			hasTxUUID:     true,
		},
		{
			name:          "оплаченный заказ через SBP",
			status:        model.StatusPaid,
			paymentMethod: "SBP",
			hasTxUUID:     true,
		},
		{
			name:          "отменённый заказ",
			status:        model.StatusCancelled,
			paymentMethod: "CARD",
			hasTxUUID:     true,
		},
		{
			name:          "заказ в статусе assembled",
			status:        model.StatusAssembled,
			paymentMethod: "CARD",
			hasTxUUID:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			order := &model.Order{
				UserUUID:      uuid.New(),
				PartUUIDs:     []uuid.UUID{uuid.New()},
				TotalPrice:    1000.0,
				PaymentMethod: tt.paymentMethod,
				Status:        tt.status,
			}

			if tt.hasTxUUID {
				order.TransactionUUID = uuid.New()
			} else {
				order.TransactionUUID = uuid.Nil
			}

			// Act
			createdOrder, err := repo.CreateOrder(ctx, order)

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, createdOrder)
			assert.Equal(t, tt.status, createdOrder.Status)
			assert.Equal(t, tt.paymentMethod, createdOrder.PaymentMethod)

			// Verify via GetOrder
			fetchedOrder, err := repo.GetOrder(ctx, createdOrder.OrderUUID)
			require.NoError(t, err)
			assert.Equal(t, tt.status, fetchedOrder.Status)
		})
	}
}

func TestRepository_FullLifecycle(t *testing.T) {
	pool := setupTestDB(t)
	repo := &repository{db: pool}
	ctx := context.Background()

	// Шаг 1: Создаем заказ в статусе pending
	order := &model.Order{
		UserUUID:        uuid.New(),
		PartUUIDs:       []uuid.UUID{uuid.New()},
		TotalPrice:      1200.0,
		TransactionUUID: uuid.Nil,
		PaymentMethod:   "",
		Status:          model.StatusPendingPayment,
	}

	createdOrder, err := repo.CreateOrder(ctx, order)
	require.NoError(t, err)
	assert.Equal(t, model.StatusPendingPayment, createdOrder.Status)

	// Шаг 2: Получаем заказ
	fetchedOrder, err := repo.GetOrder(ctx, createdOrder.OrderUUID)
	require.NoError(t, err)
	assert.Equal(t, model.StatusPendingPayment, fetchedOrder.Status)

	// Шаг 3: Оплачиваем заказ
	fetchedOrder.Status = model.StatusPaid
	fetchedOrder.PaymentMethod = "CARD"
	fetchedOrder.TransactionUUID = uuid.New()

	updatedOrder, err := repo.UpdateOrder(ctx, fetchedOrder)
	require.NoError(t, err)
	assert.Equal(t, model.StatusPaid, updatedOrder.Status)

	// Шаг 4: Проверяем что статус обновился
	paidOrder, err := repo.GetOrder(ctx, createdOrder.OrderUUID)
	require.NoError(t, err)
	assert.Equal(t, model.StatusPaid, paidOrder.Status)
	assert.Equal(t, "CARD", paidOrder.PaymentMethod)
	assert.NotEqual(t, uuid.Nil, paidOrder.TransactionUUID)

	// Шаг 5: Отменяем заказ
	paidOrder.Status = model.StatusCancelled

	cancelledOrder, err := repo.UpdateOrder(ctx, paidOrder)
	require.NoError(t, err)
	assert.Equal(t, model.StatusCancelled, cancelledOrder.Status)

	// Шаг 6: Финальная проверка
	finalOrder, err := repo.GetOrder(ctx, createdOrder.OrderUUID)
	require.NoError(t, err)
	assert.Equal(t, model.StatusCancelled, finalOrder.Status)
}
