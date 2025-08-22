package v1

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/metadata"

	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/model"
	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
)

// MockInventoryService мок для Inventory сервиса
type MockInventoryService struct {
	mock.Mock
}

func (m *MockInventoryService) ListParts(ctx context.Context, filter *model.Filter) (*[]model.Part, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*[]model.Part), args.Error(1)
}

func (m *MockInventoryService) GetPart(ctx context.Context, uuid uuid.UUID) (*model.Part, error) {
	args := m.Called(ctx, uuid)
	return args.Get(0).(*model.Part), args.Error(1)
}

func TestAPI_ListParts_WithAuth(t *testing.T) {
	// Arrange
	mockService := &MockInventoryService{}
	api := NewAPI(mockService)

	// Создаем контекст с аутентификацией
	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "session-uuid", "test-session-uuid")

	// Мокаем ответ от сервиса
	expectedResponse := &[]model.Part{}
	mockService.On("ListParts", ctx, (*model.Filter)(nil)).Return(expectedResponse, nil)

	// Act
	resp, err := api.ListParts(ctx, &inventoryV1.ListPartsRequest{})

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	mockService.AssertExpectations(t)
}

func TestAPI_ListParts_WithoutAuth(t *testing.T) {
	// Arrange
	mockService := &MockInventoryService{}
	api := NewAPI(mockService)

	// Контекст без аутентификации
	ctx := context.Background()

	// Мокаем ответ от сервиса
	expectedResponse := &[]model.Part{}
	mockService.On("ListParts", ctx, (*model.Filter)(nil)).Return(expectedResponse, nil)

	// Act
	resp, err := api.ListParts(ctx, &inventoryV1.ListPartsRequest{})

	// Assert
	// API должен работать без аутентификации, так как проверка происходит на уровне gRPC интерцептора
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	mockService.AssertExpectations(t)
}
