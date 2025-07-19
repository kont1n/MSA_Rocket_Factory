package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

type MockOrderRepository struct {
	mock.Mock
}

func NewMockOrderRepository(t mock.TestingT) *MockOrderRepository {
	mock := &MockOrderRepository{}
	mock.Test(t)
	return mock
}

func (m *MockOrderRepository) CreateOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	args := m.Called(ctx, order)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderRepository) GetOrder(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderRepository) UpdateOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	args := m.Called(ctx, order)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}
