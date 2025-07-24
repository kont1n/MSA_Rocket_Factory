package v1

import "github.com/kont1n/MSA_Rocket_Factory/order/internal/service"

type api struct {
	orderService service.OrderService
}

func NewAPI(service service.OrderService) *api {
	return &api{
		orderService: service,
	}
}
