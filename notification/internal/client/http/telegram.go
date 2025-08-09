package http

import "context"

type TelegramClient interface {
	SendMessage(ctx context.Context, chatID int64, text string) error
}
