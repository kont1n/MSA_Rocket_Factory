package converter

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	repoModel "github.com/kont1n/MSA_Rocket_Factory/order/internal/repository/model"
)

func (s *ConverterSuite) TestModelToRepo_Success() {
	// Подготовка
	orderUUID := uuid.New()
	userUUID := uuid.New()
	transactionUUID := uuid.New()
	partUUID1 := uuid.New()
	partUUID2 := uuid.New()

	order := &model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUUIDs:       []uuid.UUID{partUUID1, partUUID2},
		TotalPrice:      1500.75,
		TransactionUUID: transactionUUID,
		PaymentMethod:   "CARD",
		Status:          model.StatusPaid,
	}

	// Выполнение
	result := ModelToRepo(order)

	// Проверка
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), orderUUID.String(), result.OrderUUID)
	assert.Equal(s.T(), userUUID.String(), result.UserUUID)
	assert.Len(s.T(), result.PartUUIDs, 2)
	assert.Contains(s.T(), result.PartUUIDs, partUUID1.String())
	assert.Contains(s.T(), result.PartUUIDs, partUUID2.String())
	assert.Equal(s.T(), float32(1500.75), result.TotalPrice)
	assert.Equal(s.T(), transactionUUID.String(), result.TransactionUUID)
	assert.Equal(s.T(), "CARD", result.PaymentMethod)
	assert.Equal(s.T(), string(model.StatusPaid), result.Status)
}

func (s *ConverterSuite) TestModelToRepo_EmptyParts() {
	// Подготовка
	orderUUID := uuid.New()
	userUUID := uuid.New()
	transactionUUID := uuid.New()

	order := &model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUUIDs:       []uuid.UUID{},
		TotalPrice:      0,
		TransactionUUID: transactionUUID,
		PaymentMethod:   "",
		Status:          model.StatusPendingPayment,
	}

	// Выполнение
	result := ModelToRepo(order)

	// Проверка
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), orderUUID.String(), result.OrderUUID)
	assert.Equal(s.T(), userUUID.String(), result.UserUUID)
	assert.Empty(s.T(), result.PartUUIDs)
	assert.Equal(s.T(), float32(0), result.TotalPrice)
	assert.Equal(s.T(), transactionUUID.String(), result.TransactionUUID)
	assert.Equal(s.T(), "", result.PaymentMethod)
	assert.Equal(s.T(), string(model.StatusPendingPayment), result.Status)
}

func (s *ConverterSuite) TestRepoToModel_Success() {
	// Подготовка
	orderUUID := uuid.New()
	userUUID := uuid.New()
	transactionUUID := uuid.New()
	partUUID1 := uuid.New()
	partUUID2 := uuid.New()

	repoOrder := &repoModel.Order{
		OrderUUID:       orderUUID.String(),
		UserUUID:        userUUID.String(),
		PartUUIDs:       []string{partUUID1.String(), partUUID2.String()},
		TotalPrice:      1500.75,
		TransactionUUID: transactionUUID.String(),
		PaymentMethod:   "CARD",
		Status:          string(model.StatusPaid),
	}

	// Выполнение
	result, err := RepoToModel(repoOrder)

	// Проверка
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), orderUUID, result.OrderUUID)
	assert.Equal(s.T(), userUUID, result.UserUUID)
	assert.Len(s.T(), result.PartUUIDs, 2)
	assert.Contains(s.T(), result.PartUUIDs, partUUID1)
	assert.Contains(s.T(), result.PartUUIDs, partUUID2)
	assert.Equal(s.T(), 1500.75, result.TotalPrice)
	assert.Equal(s.T(), transactionUUID, result.TransactionUUID)
	assert.Equal(s.T(), "CARD", result.PaymentMethod)
	assert.Equal(s.T(), model.StatusPaid, result.Status)
}

func (s *ConverterSuite) TestRepoToModel_InvalidOrderUUID() {
	// Подготовка
	userUUID := uuid.New()
	transactionUUID := uuid.New()
	partUUID := uuid.New()

	repoOrder := &repoModel.Order{
		OrderUUID:       "invalid-uuid",
		UserUUID:        userUUID.String(),
		PartUUIDs:       []string{partUUID.String()},
		TotalPrice:      1500.75,
		TransactionUUID: transactionUUID.String(),
		PaymentMethod:   "CARD",
		Status:          string(model.StatusPaid),
	}

	// Выполнение
	result, err := RepoToModel(repoOrder)

	// Проверка
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), model.ErrConvertFromRepo, err)
}

func (s *ConverterSuite) TestRepoToModel_InvalidUserUUID() {
	// Подготовка
	orderUUID := uuid.New()
	transactionUUID := uuid.New()
	partUUID := uuid.New()

	repoOrder := &repoModel.Order{
		OrderUUID:       orderUUID.String(),
		UserUUID:        "invalid-uuid",
		PartUUIDs:       []string{partUUID.String()},
		TotalPrice:      1500.75,
		TransactionUUID: transactionUUID.String(),
		PaymentMethod:   "CARD",
		Status:          string(model.StatusPaid),
	}

	// Выполнение
	result, err := RepoToModel(repoOrder)

	// Проверка
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), model.ErrConvertFromRepo, err)
}

func (s *ConverterSuite) TestRepoToModel_InvalidTransactionUUID() {
	// Подготовка
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUID := uuid.New()

	repoOrder := &repoModel.Order{
		OrderUUID:       orderUUID.String(),
		UserUUID:        userUUID.String(),
		PartUUIDs:       []string{partUUID.String()},
		TotalPrice:      1500.75,
		TransactionUUID: "invalid-uuid",
		PaymentMethod:   "CARD",
		Status:          string(model.StatusPaid),
	}

	// Выполнение
	result, err := RepoToModel(repoOrder)

	// Проверка
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), model.ErrConvertFromRepo, err)
}

func (s *ConverterSuite) TestRepoToModel_InvalidPartUUID() {
	// Подготовка
	orderUUID := uuid.New()
	userUUID := uuid.New()
	transactionUUID := uuid.New()

	repoOrder := &repoModel.Order{
		OrderUUID:       orderUUID.String(),
		UserUUID:        userUUID.String(),
		PartUUIDs:       []string{"invalid-uuid"},
		TotalPrice:      1500.75,
		TransactionUUID: transactionUUID.String(),
		PaymentMethod:   "CARD",
		Status:          string(model.StatusPaid),
	}

	// Выполнение
	result, err := RepoToModel(repoOrder)

	// Проверка
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), model.ErrConvertFromRepo, err)
}

func (s *ConverterSuite) TestRepoToModel_EmptyParts() {
	// Подготовка
	orderUUID := uuid.New()
	userUUID := uuid.New()
	transactionUUID := uuid.New()

	repoOrder := &repoModel.Order{
		OrderUUID:       orderUUID.String(),
		UserUUID:        userUUID.String(),
		PartUUIDs:       []string{},
		TotalPrice:      0,
		TransactionUUID: transactionUUID.String(),
		PaymentMethod:   "",
		Status:          string(model.StatusPendingPayment),
	}

	// Выполнение
	result, err := RepoToModel(repoOrder)

	// Проверка
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), orderUUID, result.OrderUUID)
	assert.Equal(s.T(), userUUID, result.UserUUID)
	assert.Empty(s.T(), result.PartUUIDs)
	assert.Equal(s.T(), float64(0), result.TotalPrice)
	assert.Equal(s.T(), transactionUUID, result.TransactionUUID)
	assert.Equal(s.T(), "", result.PaymentMethod)
	assert.Equal(s.T(), model.StatusPendingPayment, result.Status)
}

func (s *ConverterSuite) TestRepoToModel_AllStatuses() {
	// Подготовка
	orderUUID := uuid.New()
	userUUID := uuid.New()
	transactionUUID := uuid.New()

	testCases := []struct {
		status     string
		expected   model.OrderStatus
		shouldPass bool
	}{
		{string(model.StatusPendingPayment), model.StatusPendingPayment, true},
		{string(model.StatusPaid), model.StatusPaid, true},
		{string(model.StatusCancelled), model.StatusCancelled, true},
		{"INVALID_STATUS", model.OrderStatus("INVALID_STATUS"), true}, // строка просто копируется
	}

	for _, tc := range testCases {
		repoOrder := &repoModel.Order{
			OrderUUID:       orderUUID.String(),
			UserUUID:        userUUID.String(),
			PartUUIDs:       []string{},
			TotalPrice:      0,
			TransactionUUID: transactionUUID.String(),
			PaymentMethod:   "",
			Status:          tc.status,
		}

		// Выполнение
		result, err := RepoToModel(repoOrder)

		// Проверка
		if tc.shouldPass {
			assert.NoError(s.T(), err)
			assert.NotNil(s.T(), result)
			assert.Equal(s.T(), tc.expected, result.Status)
		} else {
			assert.Error(s.T(), err)
			assert.Nil(s.T(), result)
		}
	}
}
