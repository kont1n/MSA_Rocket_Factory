package kafka

import (
	"context"

	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/kafka"
)

// LagMiddleware middleware для отслеживания lag в Kafka consumer
type LagMiddleware struct {
	lagTracker *LagTracker
}

// NewLagMiddleware создает новый middleware для отслеживания lag
func NewLagMiddleware(lagTracker *LagTracker) *LagMiddleware {
	return &LagMiddleware{
		lagTracker: lagTracker,
	}
}

// Handle обрабатывает сообщение и запускает отслеживание lag
func (m *LagMiddleware) Handle(next kafka.MessageHandler) kafka.MessageHandler {
	return func(ctx context.Context, msg kafka.Message) error {
		// Запускаем отслеживание lag в фоне, если еще не запущено
		// В реальном проекте это должно быть сделано при инициализации consumer'а
		go func() {
			m.lagTracker.Start(ctx)
		}()

		// Передаем управление следующему обработчику
		return next(ctx, msg)
	}
}
