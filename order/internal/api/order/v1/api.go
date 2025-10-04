package v1

import (
	"context"

	"go.uber.org/zap"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/service"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

type api struct {
	orderService service.OrderService
}

func NewAPI(service service.OrderService) *api {
	return &api{
		orderService: service,
	}
}

// SetOrderService устанавливает OrderService для API
func (a *api) SetOrderService(service service.OrderService) {
	logger.Info(context.Background(), "🔧 SetOrderService: устанавливаем OrderService в API",
		zap.Bool("service_is_nil", service == nil))
	a.orderService = service
	logger.Info(context.Background(), "✅ SetOrderService: OrderService установлен в API")
}
