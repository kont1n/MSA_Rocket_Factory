package telegram

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"go.uber.org/zap"

	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

type client struct {
	bot    *bot.Bot
	chatID int64
}

// NewClient создает новый HTTP клиент для Telegram
func NewClient(botToken string, chatID int64) (*client, error) {
	b, err := bot.New(botToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot: %w", err)
	}

	return &client{
		bot:    b,
		chatID: chatID,
	}, nil
}

// SendMessage отправляет сообщение в указанный чат
func (c *client) SendMessage(ctx context.Context, message string) error {
	_, err := c.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: c.chatID,
		Text:   message,
	})
	if err != nil {
		return fmt.Errorf("failed to send telegram message: %w", err)
	}

	logger.Info(ctx, "Telegram message sent successfully",
		zap.Int64("chat_id", c.chatID))

	return nil
}
