package consumer

import (
	"context"

	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/kafka"
)

// ConsumerService общий интерфейс для всех consumer сервисов
type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

// MessageHandler обработчик сообщений
type MessageHandler = kafka.MessageHandler
