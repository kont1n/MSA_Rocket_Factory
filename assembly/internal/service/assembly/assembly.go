package assembly

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

	"github.com/kont1n/MSA_Rocket_Factory/assembly/internal/model"
)

func (s *service) Assemble(ctx context.Context, event model.OrderPaidEvent) error {
	startTime := time.Now()

	// Генерируем случайное время сборки от 5 до 15 секунд
	randomNum, err := rand.Int(rand.Reader, big.NewInt(11)) // Генерируем число от 0 до 10
	if err != nil {
		// В случае ошибки используем значение по умолчанию
		randomNum = big.NewInt(5)
	}
	assemblyTime := int64(5 + randomNum.Int64()) // Добавляем 5 для диапазона 5-15

	// Записываем метрики начала сборки
	if s.metrics != nil {
		s.metrics.recordAssemblyStart(ctx)
	}

	// Используем таймер вместо time.Sleep для корректной работы с контекстом
	timer := time.NewTimer(time.Duration(assemblyTime) * time.Second)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		// Записываем метрики ошибки при отмене контекста
		if s.metrics != nil {
			s.metrics.recordAssemblyError(ctx, "context_cancelled")
		}
		return ctx.Err()
	case <-timer.C:
		// Продолжаем выполнение после истечения таймера
	}

	err = s.assemblyProducerService.ProduceAssembly(ctx, model.ShipAssembledEvent{
		EventUUID: event.EventUUID,
		OrderUUID: event.OrderUUID,
		UserUUID:  event.UserUUID,
		BuildTime: assemblyTime,
	})
	if err != nil {
		// Записываем метрики ошибки при неудачной отправке события
		if s.metrics != nil {
			s.metrics.recordAssemblyError(ctx, "producer_error")
		}
		return err
	}

	// Записываем метрики успешного завершения сборки
	if s.metrics != nil {
		s.metrics.recordAssemblyComplete(ctx, time.Since(startTime))
	}

	return nil
}
