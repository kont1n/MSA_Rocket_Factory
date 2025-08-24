package kafka

import (
	"context"

	"go.uber.org/zap"

	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/kafka"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/kafka/consumer"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
}

func Logging(logger Logger) consumer.Middleware {
	return func(next kafka.MessageHandler) kafka.MessageHandler {
		return func(ctx context.Context, msg kafka.Message) error {
			logger.Info(ctx, "Kafka msg received", zap.String("topic", msg.Topic))
			return next(ctx, msg)
		}
	}
}
