package telegram

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/go-telegram/bot"
	"github.com/kont1n/MSA_Rocket_Factory/notification/internal/config/env"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

type client struct {
	bot    *bot.Bot
	chatID string
}

func NewClient(cfg env.TelegramConfig) (*client, error) {
	b, err := bot.New(cfg.BotToken())
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot: %w", err)
	}

	return &client{
		bot:    b,
		chatID: cfg.ChatID(),
	}, nil
}

func (c *client) SendMessage(ctx context.Context, message string) error {
	_, err := c.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: c.chatID,
		Text:   message,
	})
	if err != nil {
		return fmt.Errorf("failed to send telegram message: %w", err)
	}

	logger.Info(ctx, "Telegram message sent successfully",
		zap.String("chat_id", c.chatID))

	return nil
}
