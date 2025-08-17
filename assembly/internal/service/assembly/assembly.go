package assembly

import (
	"context"
	"time"

	"github.com/kont1n/MSA_Rocket_Factory/assembly/internal/model"
)

const delayTime = 10

func (s *service) Assemble(ctx context.Context, event model.OrderPaidEvent) error {
	// Используем таймер вместо time.Sleep для корректной работы с контекстом
	timer := time.NewTimer(delayTime * time.Second)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		// Продолжаем выполнение после истечения таймера
	}

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
