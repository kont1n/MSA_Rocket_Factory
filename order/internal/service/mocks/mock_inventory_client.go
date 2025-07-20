package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

type MockInventoryClient struct {
	mock.Mock
}

func NewMockInventoryClient(t mock.TestingT) *MockInventoryClient {
	mock := &MockInventoryClient{}
	mock.Test(t)
	return mock
}

func (m *MockInventoryClient) ListParts(ctx context.Context, filter *model.Filter) (*[]model.Part, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]model.Part), args.Error(1)
}
