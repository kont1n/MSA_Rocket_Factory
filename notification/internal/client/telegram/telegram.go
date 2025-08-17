package telegram

import "context"

type TelegramClient interface {
	SendMessage(ctx context.Context, message string) error
}
