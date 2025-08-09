//go:build integration

package integration

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

// InsertTestOrder — вставляет тестовый заказ в таблицу PostgreSQL и возвращает его UUID
func (env *TestEnvironment) InsertTestOrder(ctx context.Context) (string, error) {
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New(), uuid.New()}

	// Подключаемся к PostgreSQL
	connStr, err := env.Postgres.ConnectionString(ctx)
	if err != nil {
		return "", fmt.Errorf("не удалось получить строку подключения: %w", err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return "", fmt.Errorf("не удалось подключиться к PostgreSQL: %w", err)
	}
	defer pool.Close()

	// Вставляем тестовый заказ
	query := `
		INSERT INTO orders (order_uuid, user_uuid, part_uuid, total_price, transaction_uuid, payment_method, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = pool.Exec(ctx, query,
		orderUUID,
		userUUID,
		partUUIDs,
		15750000.75,
		uuid.New(),
		"credit_card",
		string(model.StatusPendingPayment),
	)
	if err != nil {
		return "", fmt.Errorf("не удалось вставить тестовый заказ: %w", err)
	}

	return orderUUID.String(), nil
}

// InsertTestOrderWithData — вставляет тестовый заказ с заданными данными
func (env *TestEnvironment) InsertTestOrderWithData(ctx context.Context, order *model.Order) (string, error) {
	orderUUID := uuid.New()

	// Подключаемся к PostgreSQL
	connStr, err := env.Postgres.ConnectionString(ctx)
	if err != nil {
		return "", fmt.Errorf("не удалось получить строку подключения: %w", err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return "", fmt.Errorf("не удалось подключиться к PostgreSQL: %w", err)
	}
	defer pool.Close()

	// Вставляем тестовый заказ
	query := `
		INSERT INTO orders (order_uuid, user_uuid, part_uuid, total_price, transaction_uuid, payment_method, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = pool.Exec(ctx, query,
		orderUUID,
		order.UserUUID,
		order.PartUUIDs,
		order.TotalPrice,
		order.TransactionUUID,
		order.PaymentMethod,
		string(order.Status),
	)
	if err != nil {
		return "", fmt.Errorf("не удалось вставить тестовый заказ: %w", err)
	}

	return orderUUID.String(), nil
}

// GetTestOrder — возвращает тестовую информацию о заказе
func (env *TestEnvironment) GetTestOrder() *model.Order {
	return &model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUUIDs:       []uuid.UUID{uuid.New(), uuid.New()},
		TotalPrice:      12500000.00,
		TransactionUUID: uuid.New(),
		PaymentMethod:   "credit_card",
		Status:          model.StatusPendingPayment,
	}
}

// InsertMultipleTestOrders — вставляет несколько тестовых заказов для тестирования
func (env *TestEnvironment) InsertMultipleTestOrders(ctx context.Context) ([]string, error) {
	var orderUUIDs []string

	// Подключаемся к PostgreSQL
	connStr, err := env.Postgres.ConnectionString(ctx)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить строку подключения: %w", err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к PostgreSQL: %w", err)
	}
	defer pool.Close()

	// Создаем тестовые заказы
	testOrders := []struct {
		orderUUID       uuid.UUID
		userUUID        uuid.UUID
		partUUIDs       []uuid.UUID
		totalPrice      float32
		transactionUUID uuid.UUID
		paymentMethod   string
		status          string
	}{
		{
			orderUUID:       uuid.New(),
			userUUID:        uuid.New(),
			partUUIDs:       []uuid.UUID{uuid.New()},
			totalPrice:      15000000.50,
			transactionUUID: uuid.New(),
			paymentMethod:   "credit_card",
			status:          string(model.StatusPendingPayment),
		},
		{
			orderUUID:       uuid.New(),
			userUUID:        uuid.New(),
			partUUIDs:       []uuid.UUID{uuid.New(), uuid.New()},
			totalPrice:      2750000.75,
			transactionUUID: uuid.New(),
			paymentMethod:   "bank_transfer",
			status:          string(model.StatusPaid),
		},
		{
			orderUUID:       uuid.New(),
			userUUID:        uuid.New(),
			partUUIDs:       []uuid.UUID{uuid.New()},
			totalPrice:      750000.00,
			transactionUUID: uuid.New(),
			paymentMethod:   "cryptocurrency",
			status:          string(model.StatusCancelled),
		},
	}

	query := `
		INSERT INTO orders (order_uuid, user_uuid, part_uuid, total_price, transaction_uuid, payment_method, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	for _, order := range testOrders {
		_, err := pool.Exec(ctx, query,
			order.orderUUID,
			order.userUUID,
			order.partUUIDs,
			order.totalPrice,
			order.transactionUUID,
			order.paymentMethod,
			order.status,
		)
		if err != nil {
			return nil, fmt.Errorf("не удалось вставить тестовый заказ: %w", err)
		}

		orderUUIDs = append(orderUUIDs, order.orderUUID.String())
	}

	return orderUUIDs, nil
}

// ClearOrdersTable — удаляет все записи из таблицы orders
func (env *TestEnvironment) ClearOrdersTable(ctx context.Context) error {
	// Подключаемся к PostgreSQL
	connStr, err := env.Postgres.ConnectionString(ctx)
	if err != nil {
		return fmt.Errorf("не удалось получить строку подключения: %w", err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к PostgreSQL: %w", err)
	}
	defer pool.Close()

	_, err = pool.Exec(ctx, "DELETE FROM orders")
	if err != nil {
		return fmt.Errorf("не удалось очистить таблицу orders: %w", err)
	}

	return nil
}

// GetOrderByUUID — получает заказ по UUID из базы данных
func (env *TestEnvironment) GetOrderByUUID(ctx context.Context, orderUUID string) (*model.Order, error) {
	// Подключаемся к PostgreSQL
	connStr, err := env.Postgres.ConnectionString(ctx)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить строку подключения: %w", err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к PostgreSQL: %w", err)
	}
	defer pool.Close()

	query := `
		SELECT order_uuid, user_uuid, part_uuid, total_price, transaction_uuid, payment_method, status
		FROM orders WHERE order_uuid = $1
	`

	var order model.Order
	var orderUUIDStr, userUUIDStr, transactionUUIDStr string
	var partUUIDs []uuid.UUID

	err = pool.QueryRow(ctx, query, orderUUID).Scan(
		&orderUUIDStr,
		&userUUIDStr,
		&partUUIDs,
		&order.TotalPrice,
		&transactionUUIDStr,
		&order.PaymentMethod,
		&order.Status,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrOrderNotFound
		}
		return nil, fmt.Errorf("не удалось получить заказ: %w", err)
	}

	// Парсим UUID
	order.OrderUUID, err = uuid.Parse(orderUUIDStr)
	if err != nil {
		return nil, fmt.Errorf("не удалось распарсить order_uuid: %w", err)
	}

	order.UserUUID, err = uuid.Parse(userUUIDStr)
	if err != nil {
		return nil, fmt.Errorf("не удалось распарсить user_uuid: %w", err)
	}

	order.TransactionUUID, err = uuid.Parse(transactionUUIDStr)
	if err != nil {
		return nil, fmt.Errorf("не удалось распарсить transaction_uuid: %w", err)
	}

	order.PartUUIDs = partUUIDs

	return &order, nil
}

// UpdateOrderInDB — обновляет заказ в базе данных (имитирует работу репозитория)
func (env *TestEnvironment) UpdateOrderInDB(ctx context.Context, order *model.Order) error {
	// Подключаемся к PostgreSQL
	connStr, err := env.Postgres.ConnectionString(ctx)
	if err != nil {
		return fmt.Errorf("не удалось получить строку подключения: %w", err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к PostgreSQL: %w", err)
	}
	defer pool.Close()

	query := `
		UPDATE orders 
		SET transaction_uuid = $1, payment_method = $2, status = $3, updated_at = NOW()
		WHERE order_uuid = $4
	`

	_, err = pool.Exec(ctx, query,
		order.TransactionUUID,
		order.PaymentMethod,
		string(order.Status),
		order.OrderUUID,
	)
	if err != nil {
		return fmt.Errorf("не удалось обновить заказ: %w", err)
	}

	return nil
}
