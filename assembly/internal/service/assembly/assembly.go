package assembly

import (
	"context"
	"time"

	"github.com/kont1n/MSA_Rocket_Factory/assembly/internal/model"
)

const delayTime = 10

func (s *service) Assemble(ctx context.Context, event model.OrderPaidEvent) error {
	time.Sleep(delayTime * time.Second)

	err := s.assemblyProducerService.ProduceAssembly(ctx, model.ShipAssembledEvent{
		EventUUID: event.EventUUID,
		OrderUUID: event.OrderUUID,
		UserUUID:  event.UserUUID,
		BuildTime: delayTime,
	})
	if err != nil {
		return err
	}

	return nil
}
