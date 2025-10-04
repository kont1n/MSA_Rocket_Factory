package kafka

import (
	"context"

	"github.com/IBM/sarama"
)

// MessageHandler — обработчик сообщений.
type MessageHandler func(ctx context.Context, msg Message) error

type Consumer interface {
	Consume(ctx context.Context, handler MessageHandler) error
}

type Producer interface {
	Send(ctx context.Context, key, value []byte) error
}

// Middleware — функция middleware для дополнительной обработки сообщений.
type Middleware func(next MessageHandler) MessageHandler

// ConsumerGroupHandler — интерфейс для обработки consumer group
type ConsumerGroupHandler interface {
	Setup(sarama.ConsumerGroupSession) error
	Cleanup(sarama.ConsumerGroupSession) error
	ConsumeClaim(sarama.ConsumerGroupSession, sarama.ConsumerGroupClaim) error
}

// ConsumerGroupHandlerMiddleware — функция middleware для consumer group handler
type ConsumerGroupHandlerMiddleware func(next ConsumerGroupHandler) ConsumerGroupHandler
