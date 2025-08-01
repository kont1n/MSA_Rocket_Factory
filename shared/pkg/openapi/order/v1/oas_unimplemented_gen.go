// Code generated by ogen, DO NOT EDIT.

package order_v1

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// CancelOrder implements CancelOrder operation.
//
// Отмена заказа.
//
// POST /api/v1/orders/{order_uuid}/cancel
func (UnimplementedHandler) CancelOrder(ctx context.Context, params CancelOrderParams) (r CancelOrderRes, _ error) {
	return r, ht.ErrNotImplemented
}

// CreateOrder implements CreateOrder operation.
//
// Создание заказа.
//
// POST /api/v1/orders
func (UnimplementedHandler) CreateOrder(ctx context.Context, req *CreateOrderRequest) (r CreateOrderRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetOrderByUUID implements GetOrderByUUID operation.
//
// Получение заказа по UUID.
//
// GET /api/v1/orders/{order_uuid}
func (UnimplementedHandler) GetOrderByUUID(ctx context.Context, params GetOrderByUUIDParams) (r GetOrderByUUIDRes, _ error) {
	return r, ht.ErrNotImplemented
}

// PayOrder implements PayOrder operation.
//
// Оплата заказа.
//
// POST /api/v1/orders/{order_uuid}/pay
func (UnimplementedHandler) PayOrder(ctx context.Context, req *PayOrderRequest, params PayOrderParams) (r PayOrderRes, _ error) {
	return r, ht.ErrNotImplemented
}

// NewError creates *GenericErrorStatusCode from error returned by handler.
//
// Used for common default response.
func (UnimplementedHandler) NewError(ctx context.Context, err error) (r *GenericErrorStatusCode) {
	r = new(GenericErrorStatusCode)
	return r
}
