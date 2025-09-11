package assembly

import (
	"context"
	"time"

	"github.com/kont1n/MSA_Rocket_Factory/assembly/internal/model"
)

const delayTime = 10

func (s *service) Assemble(ctx context.Context, event model.OrderPaidEvent) error {
	startTime := time.Now()
	rocketType := "standard" // По умолчанию стандартный тип ракеты

	// Записываем метрики начала сборки
	if s.metrics != nil {
		s.metrics.recordAssemblyStart(ctx, rocketType)
	}

	// Используем таймер вместо time.Sleep для корректной работы с контекстом
	timer := time.NewTimer(delayTime * time.Second)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		// Записываем метрики ошибки при отмене контекста
		if s.metrics != nil {
			s.metrics.recordAssemblyError(ctx, rocketType, "context_cancelled")
		}
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
		// Записываем метрики ошибки при неудачной отправке события
		if s.metrics != nil {
			s.metrics.recordAssemblyError(ctx, rocketType, "producer_error")
		}
		return err
	}

	// Записываем метрики успешного завершения сборки
	if s.metrics != nil {
		s.metrics.recordAssemblyComplete(ctx, rocketType, time.Since(startTime))
	}

	return nil
}
