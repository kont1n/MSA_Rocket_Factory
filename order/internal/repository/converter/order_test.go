package converter

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	repoModel "github.com/kont1n/MSA_Rocket_Factory/order/internal/repository/model"
)

func TestToRepoOrderPostgres(t *testing.T) {
	// Arrange
	userUUID := uuid.New()
	partUUID1 := uuid.New()
	partUUID2 := uuid.New()
	transactionUUID := uuid.New()

	order := &model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        userUUID,
		PartUUIDs:       []uuid.UUID{partUUID1, partUUID2},
		TotalPrice:      1500.50,
		TransactionUUID: transactionUUID,
		PaymentMethod:   "CARD",
		Status:          model.StatusPaid,
	}

	// Act
	repoOrder := ToRepoOrderPostgres(order)

	// Assert
	assert.NotNil(t, repoOrder)
	assert.Equal(t, userUUID, repoOrder.UserUUID)
	assert.Equal(t, []uuid.UUID{partUUID1, partUUID2}, repoOrder.PartUUIDs)
	assert.Equal(t, float32(1500.50), repoOrder.TotalPrice)
	assert.Equal(t, transactionUUID, repoOrder.TransactionUUID)
	assert.Equal(t, "CARD", repoOrder.PaymentMethod)
	assert.Equal(t, "PAID", repoOrder.Status)
}

func TestToRepoOrderPostgres_WithNilTransactionUUID(t *testing.T) {
	// Arrange - заказ без TransactionUUID (pending payment)
	order := &model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUUIDs:       []uuid.UUID{uuid.New()},
		TotalPrice:      500.0,
		TransactionUUID: uuid.Nil,
		PaymentMethod:   "",
		Status:          model.StatusPendingPayment,
	}

	// Act
	repoOrder := ToRepoOrderPostgres(order)

	// Assert
	assert.NotNil(t, repoOrder)
	assert.Equal(t, uuid.Nil, repoOrder.TransactionUUID)
	assert.Equal(t, "", repoOrder.PaymentMethod)
	assert.Equal(t, "PENDING_PAYMENT", repoOrder.Status)
}

func TestToRepoOrderPostgres_WithMultipleParts(t *testing.T) {
	// Arrange - заказ с несколькими деталями
	partUUIDs := []uuid.UUID{
		uuid.New(),
		uuid.New(),
		uuid.New(),
		uuid.New(),
	}

	order := &model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUUIDs:       partUUIDs,
		TotalPrice:      5000.0,
		TransactionUUID: uuid.New(),
		PaymentMethod:   "SBP",
		Status:          model.StatusPaid,
	}

	// Act
	repoOrder := ToRepoOrderPostgres(order)

	// Assert
	assert.NotNil(t, repoOrder)
	assert.Len(t, repoOrder.PartUUIDs, 4)
	assert.Equal(t, partUUIDs, repoOrder.PartUUIDs)
}

func TestToRepoOrderPostgres_DifferentStatuses(t *testing.T) {
	tests := []struct {
		name   string
		status model.OrderStatus
	}{
		{
			name:   "pending payment статус",
			status: model.StatusPendingPayment,
		},
		{
			name:   "paid статус",
			status: model.StatusPaid,
		},
		{
			name:   "cancelled статус",
			status: model.StatusCancelled,
		},
		{
			name:   "assembled статус",
			status: model.StatusAssembled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			order := &model.Order{
				OrderUUID:       uuid.New(),
				UserUUID:        uuid.New(),
				PartUUIDs:       []uuid.UUID{uuid.New()},
				TotalPrice:      1000.0,
				TransactionUUID: uuid.New(),
				PaymentMethod:   "CARD",
				Status:          tt.status,
			}

			// Act
			repoOrder := ToRepoOrderPostgres(order)

			// Assert
			assert.NotNil(t, repoOrder)
			assert.Equal(t, string(tt.status), repoOrder.Status)
		})
	}
}

func TestToModelOrderFromPostgres(t *testing.T) {
	// Arrange
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUID1 := uuid.New()
	partUUID2 := uuid.New()
	transactionUUID := uuid.New()

	repoOrder := &repoModel.OrderPostgres{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUUIDs:       []uuid.UUID{partUUID1, partUUID2},
		TotalPrice:      2500.75,
		TransactionUUID: transactionUUID,
		PaymentMethod:   "CREDIT_CARD",
		Status:          "PAID",
	}

	// Act
	order, err := ToModelOrderFromPostgres(repoOrder)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, orderUUID, order.OrderUUID)
	assert.Equal(t, userUUID, order.UserUUID)
	assert.Equal(t, []uuid.UUID{partUUID1, partUUID2}, order.PartUUIDs)
	assert.Equal(t, float32(2500.75), order.TotalPrice)
	assert.Equal(t, transactionUUID, order.TransactionUUID)
	assert.Equal(t, "CREDIT_CARD", order.PaymentMethod)
	assert.Equal(t, model.StatusPaid, order.Status)
}

func TestToModelOrderFromPostgres_WithNilTransactionUUID(t *testing.T) {
	// Arrange - заказ без TransactionUUID
	repoOrder := &repoModel.OrderPostgres{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUUIDs:       []uuid.UUID{uuid.New()},
		TotalPrice:      500.0,
		TransactionUUID: uuid.Nil,
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	// Act
	order, err := ToModelOrderFromPostgres(repoOrder)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, uuid.Nil, order.TransactionUUID)
	assert.Equal(t, "", order.PaymentMethod)
	assert.Equal(t, model.StatusPendingPayment, order.Status)
}

func TestToModelOrderFromPostgres_WithMultipleParts(t *testing.T) {
	// Arrange - заказ с несколькими деталями
	partUUIDs := []uuid.UUID{
		uuid.New(),
		uuid.New(),
		uuid.New(),
	}

	repoOrder := &repoModel.OrderPostgres{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUUIDs:       partUUIDs,
		TotalPrice:      3000.0,
		TransactionUUID: uuid.New(),
		PaymentMethod:   "SBP",
		Status:          "PAID",
	}

	// Act
	order, err := ToModelOrderFromPostgres(repoOrder)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Len(t, order.PartUUIDs, 3)
	assert.Equal(t, partUUIDs, order.PartUUIDs)
}

func TestToModelOrderFromPostgres_DifferentStatuses(t *testing.T) {
	tests := []struct {
		name           string
		repoStatus     string
		expectedStatus model.OrderStatus
	}{
		{
			name:           "pending payment статус",
			repoStatus:     "PENDING_PAYMENT",
			expectedStatus: model.StatusPendingPayment,
		},
		{
			name:           "paid статус",
			repoStatus:     "PAID",
			expectedStatus: model.StatusPaid,
		},
		{
			name:           "cancelled статус",
			repoStatus:     "CANCELLED",
			expectedStatus: model.StatusCancelled,
		},
		{
			name:           "assembled статус",
			repoStatus:     "ASSEMBLED",
			expectedStatus: model.StatusAssembled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			repoOrder := &repoModel.OrderPostgres{
				OrderUUID:       uuid.New(),
				UserUUID:        uuid.New(),
				PartUUIDs:       []uuid.UUID{uuid.New()},
				TotalPrice:      1000.0,
				TransactionUUID: uuid.New(),
				PaymentMethod:   "CARD",
				Status:          tt.repoStatus,
			}

			// Act
			order, err := ToModelOrderFromPostgres(repoOrder)

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, order)
			assert.Equal(t, tt.expectedStatus, order.Status)
		})
	}
}

func TestToModelOrderFromPostgres_EmptyPartsList(t *testing.T) {
	// Arrange - заказ с пустым списком деталей
	repoOrder := &repoModel.OrderPostgres{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUUIDs:       []uuid.UUID{},
		TotalPrice:      0.0,
		TransactionUUID: uuid.Nil,
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	// Act
	order, err := ToModelOrderFromPostgres(repoOrder)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Empty(t, order.PartUUIDs)
}

func TestToModelOrderFromPostgres_DifferentPaymentMethods(t *testing.T) {
	tests := []struct {
		name          string
		paymentMethod string
	}{
		{
			name:          "оплата картой CARD",
			paymentMethod: "CARD",
		},
		{
			name:          "оплата через SBP",
			paymentMethod: "SBP",
		},
		{
			name:          "оплата через CREDIT_CARD",
			paymentMethod: "CREDIT_CARD",
		},
		{
			name:          "оплата через INVESTOR_MONEY",
			paymentMethod: "INVESTOR_MONEY",
		},
		{
			name:          "пустой payment method",
			paymentMethod: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			repoOrder := &repoModel.OrderPostgres{
				OrderUUID:       uuid.New(),
				UserUUID:        uuid.New(),
				PartUUIDs:       []uuid.UUID{uuid.New()},
				TotalPrice:      1000.0,
				TransactionUUID: uuid.New(),
				PaymentMethod:   tt.paymentMethod,
				Status:          "PAID",
			}

			// Act
			order, err := ToModelOrderFromPostgres(repoOrder)

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, order)
			assert.Equal(t, tt.paymentMethod, order.PaymentMethod)
		})
	}
}

func TestConverters_RoundTrip(t *testing.T) {
	// Arrange - создаем исходный заказ
	originalOrder := &model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUUIDs:       []uuid.UUID{uuid.New(), uuid.New()},
		TotalPrice:      3500.99,
		TransactionUUID: uuid.New(),
		PaymentMethod:   "CARD",
		Status:          model.StatusPaid,
	}

	// Act - конвертируем туда и обратно
	repoOrder := ToRepoOrderPostgres(originalOrder)
	repoOrder.OrderUUID = originalOrder.OrderUUID // Добавляем OrderUUID, который не копируется в ToRepoOrderPostgres
	convertedOrder, err := ToModelOrderFromPostgres(repoOrder)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, convertedOrder)
	assert.Equal(t, originalOrder.OrderUUID, convertedOrder.OrderUUID)
	assert.Equal(t, originalOrder.UserUUID, convertedOrder.UserUUID)
	assert.Equal(t, originalOrder.PartUUIDs, convertedOrder.PartUUIDs)
	assert.Equal(t, originalOrder.TotalPrice, convertedOrder.TotalPrice)
	assert.Equal(t, originalOrder.TransactionUUID, convertedOrder.TransactionUUID)
	assert.Equal(t, originalOrder.PaymentMethod, convertedOrder.PaymentMethod)
	assert.Equal(t, originalOrder.Status, convertedOrder.Status)
}

func TestConverters_RoundTrip_PendingPayment(t *testing.T) {
	// Arrange - создаем заказ в статусе pending payment
	originalOrder := &model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUUIDs:       []uuid.UUID{uuid.New()},
		TotalPrice:      1000.0,
		TransactionUUID: uuid.Nil,
		PaymentMethod:   "",
		Status:          model.StatusPendingPayment,
	}

	// Act - конвертируем туда и обратно
	repoOrder := ToRepoOrderPostgres(originalOrder)
	repoOrder.OrderUUID = originalOrder.OrderUUID
	convertedOrder, err := ToModelOrderFromPostgres(repoOrder)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, convertedOrder)
	assert.Equal(t, originalOrder.OrderUUID, convertedOrder.OrderUUID)
	assert.Equal(t, originalOrder.UserUUID, convertedOrder.UserUUID)
	assert.Equal(t, originalOrder.Status, convertedOrder.Status)
	assert.Equal(t, uuid.Nil, convertedOrder.TransactionUUID)
	assert.Equal(t, "", convertedOrder.PaymentMethod)
}

func TestToRepoOrderPostgres_PriceTypes(t *testing.T) {
	tests := []struct {
		name  string
		price float32
	}{
		{
			name:  "минимальная цена",
			price: 0.01,
		},
		{
			name:  "средняя цена",
			price: 1500.50,
		},
		{
			name:  "большая цена",
			price: 999999.99,
		},
		{
			name:  "нулевая цена",
			price: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			order := &model.Order{
				OrderUUID:       uuid.New(),
				UserUUID:        uuid.New(),
				PartUUIDs:       []uuid.UUID{uuid.New()},
				TotalPrice:      tt.price,
				TransactionUUID: uuid.New(),
				PaymentMethod:   "CARD",
				Status:          model.StatusPaid,
			}

			// Act
			repoOrder := ToRepoOrderPostgres(order)

			// Assert
			assert.NotNil(t, repoOrder)
			assert.Equal(t, tt.price, repoOrder.TotalPrice)
		})
	}
}

func TestToModelOrderFromPostgres_PriceTypes(t *testing.T) {
	tests := []struct {
		name  string
		price float32
	}{
		{
			name:  "минимальная цена",
			price: 0.01,
		},
		{
			name:  "средняя цена",
			price: 2500.75,
		},
		{
			name:  "большая цена",
			price: 1000000.00,
		},
		{
			name:  "нулевая цена",
			price: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			repoOrder := &repoModel.OrderPostgres{
				OrderUUID:       uuid.New(),
				UserUUID:        uuid.New(),
				PartUUIDs:       []uuid.UUID{uuid.New()},
				TotalPrice:      tt.price,
				TransactionUUID: uuid.New(),
				PaymentMethod:   "CARD",
				Status:          "PAID",
			}

			// Act
			order, err := ToModelOrderFromPostgres(repoOrder)

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, order)
			assert.Equal(t, tt.price, order.TotalPrice)
		})
	}
}
