package consumer

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/kafka"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type consumer struct {
	group                    sarama.ConsumerGroup
	topics                   []string
	logger                   Logger
	middlewares              []Middleware
	consumerGroupMiddlewares []kafka.ConsumerGroupHandlerMiddleware
}

// NewConsumer — создаёт новый consumer.
func NewConsumer(group sarama.ConsumerGroup, topics []string, logger Logger, middlewares ...Middleware) *consumer {
	return &consumer{
		group:       group,
		topics:      topics,
		logger:      logger,
		middlewares: middlewares,
	}
}

// NewConsumerWithConsumerGroupMiddleware — создаёт новый consumer с consumer group middleware.
func NewConsumerWithConsumerGroupMiddleware(group sarama.ConsumerGroup, topics []string, logger Logger, messageMiddlewares []Middleware, consumerGroupMiddlewares []kafka.ConsumerGroupHandlerMiddleware) *consumer {
	return &consumer{
		group:                    group,
		topics:                   topics,
		logger:                   logger,
		middlewares:              messageMiddlewares,
		consumerGroupMiddlewares: consumerGroupMiddlewares,
	}
}

// Consume запускает консьюмер для списка топиков.
func (c *consumer) Consume(ctx context.Context, handler kafka.MessageHandler) error {
	var newGroupHandler kafka.ConsumerGroupHandler

	if len(c.consumerGroupMiddlewares) > 0 {
		newGroupHandler = NewGroupHandlerWithConsumerGroupMiddleware(handler, c.logger, c.middlewares, c.consumerGroupMiddlewares)
	} else {
		newGroupHandler = NewGroupHandler(handler, c.logger, c.middlewares...)
	}

	for {
		if err := c.group.Consume(ctx, c.topics, newGroupHandler); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return nil
			}

			c.logger.Error(ctx, "Kafka consume error", zap.Error(err))
			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		c.logger.Info(ctx, "Kafka consumer group rebalancing...")
	}
}
