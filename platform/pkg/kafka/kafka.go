package kafka

import (
	"context"
)

// MessageHandler — обработчик сообщений.
type MessageHandler func(ctx context.Context, msg Message) error

type Consumer interface {
	Consume(ctx context.Context, handler MessageHandler) error
}

type Producer interface {
	Send(ctx context.Context, key, value []byte) error
}
