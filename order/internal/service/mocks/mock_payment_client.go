package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

type MockPaymentClient struct {
	mock.Mock
}

func NewMockPaymentClient(t mock.TestingT) *MockPaymentClient {
	mock := &MockPaymentClient{}
	mock.Test(t)
	return mock
}

func (m *MockPaymentClient) CreatePayment(ctx context.Context, order *model.Order) (*model.Order, error) {
	args := m.Called(ctx, order)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}
