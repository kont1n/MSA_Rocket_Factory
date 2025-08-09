package telegram

import (
	"context"

	"github.com/go-telegram/bot"
)

type client struct {
	bot *bot.Bot
}

// NewClient создает новый клиент для Telegram Bot API
func NewClient(bot *bot.Bot) *client {
	return &client{
		bot: bot,
	}
}

// SendMessage отправляет сообщение в указанный чат
func (c *client) SendMessage(ctx context.Context, chatID int64, text string) error {
	_, err := c.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      text,
		ParseMode: "Markdown",
	})

	if err != nil {
		return err
	}

	return nil
}
