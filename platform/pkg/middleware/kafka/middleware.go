package kafka

import (
	"context"

	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/kafka"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

// MetricsMiddleware создает middleware для сбора Kafka метрик
func MetricsMiddleware(metrics *KafkaMetrics, groupID string) kafka.Middleware {
	return func(next kafka.MessageHandler) kafka.MessageHandler {
		return func(ctx context.Context, message kafka.Message) error {
			// Выполняем обработку сообщения
			err := next(ctx, message)

			// Определяем статус обработки
			success := err == nil

			// Записываем метрики consumer'а
			if metrics != nil {
				metrics.RecordConsumerMessage(
					ctx,
					message.Topic,
					message.Partition,
					groupID,
					success,
				)
			}

			return err
		}
	}
}

// RebalancingMiddleware создает middleware для отслеживания ребалансировок
func RebalancingMiddleware(metrics *KafkaMetrics, groupID string) kafka.ConsumerGroupHandlerMiddleware {
	return func(next kafka.ConsumerGroupHandler) kafka.ConsumerGroupHandler {
		return &rebalancingHandler{
			next:    next,
			metrics: metrics,
			groupID: groupID,
		}
	}
}

// rebalancingHandler обертка для отслеживания ребалансировок
type rebalancingHandler struct {
	next    kafka.ConsumerGroupHandler
	metrics *KafkaMetrics
	groupID string
}

func (h *rebalancingHandler) Setup(session sarama.ConsumerGroupSession) error {
	// Записываем метрику ребалансировки
	if h.metrics != nil {
		h.metrics.RecordConsumerRebalancing(session.Context(), h.groupID)
	}

	logger.Info(session.Context(), "Consumer group setup",
		zap.String("member_id", session.MemberID()),
		zap.String("group_id", h.groupID),
	)

	return h.next.Setup(session)
}

func (h *rebalancingHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	logger.Info(session.Context(), "Consumer group cleanup",
		zap.String("member_id", session.MemberID()),
		zap.String("group_id", h.groupID),
	)

	return h.next.Cleanup(session)
}

func (h *rebalancingHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	return h.next.ConsumeClaim(session, claim)
}
