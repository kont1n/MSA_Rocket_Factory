package consumer

import (
	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/kafka"
)

// Middleware — функция middleware для дополнительной обработки.
type Middleware func(next kafka.MessageHandler) kafka.MessageHandler

// groupHandler — обёртка для sarama.ConsumerGroupHandler
type groupHandler struct {
	handler kafka.MessageHandler
	logger  Logger
}

// Реализуем интерфейс kafka.ConsumerGroupHandler
var _ kafka.ConsumerGroupHandler = (*groupHandler)(nil)

// NewGroupHandler создаёт новый groupHandler с middleware цепочкой.
func NewGroupHandler(handler kafka.MessageHandler, logger Logger, middlewares ...Middleware) *groupHandler {
	// Применяем middleware цепочку
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return &groupHandler{
		handler: handler,
		logger:  logger,
	}
}

// NewGroupHandlerWithConsumerGroupMiddleware создаёт новый groupHandler с middleware цепочкой и consumer group middleware.
func NewGroupHandlerWithConsumerGroupMiddleware(handler kafka.MessageHandler, logger Logger, messageMiddlewares []Middleware, consumerGroupMiddlewares []kafka.ConsumerGroupHandlerMiddleware) kafka.ConsumerGroupHandler {
	// Применяем middleware цепочку для сообщений
	for i := len(messageMiddlewares) - 1; i >= 0; i-- {
		handler = messageMiddlewares[i](handler)
	}

	// Создаем базовый handler
	baseHandler := &groupHandler{
		handler: handler,
		logger:  logger,
	}

	// Применяем consumer group middleware цепочку
	result := kafka.ConsumerGroupHandler(baseHandler)
	for i := len(consumerGroupMiddlewares) - 1; i >= 0; i-- {
		result = consumerGroupMiddlewares[i](result)
	}

	return result
}

func (g *groupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (g *groupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (g *groupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				g.logger.Info(session.Context(), "Kafka message channel closed")
				return nil
			}

			msg := kafka.Message{
				Key:            message.Key,
				Value:          message.Value,
				Topic:          message.Topic,
				Partition:      message.Partition,
				Offset:         message.Offset,
				Timestamp:      message.Timestamp,
				BlockTimestamp: message.BlockTimestamp,
				Headers:        extractHeaders(message.Headers),
			}

			if err := g.handler(session.Context(), msg); err != nil {
				g.logger.Error(session.Context(), "Kafka handler error", zap.Error(err))
				continue
			}

			session.MarkMessage(message, "")

		case <-session.Context().Done():
			g.logger.Info(session.Context(), "Kafka session context done")
			return nil
		}
	}
}

func extractHeaders(headers []*sarama.RecordHeader) map[string][]byte {
	result := make(map[string][]byte)
	for _, h := range headers {
		if h != nil && h.Key != nil {
			result[string(h.Key)] = h.Value
		}
	}

	return result
}
