package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	"github.com/kont1n/MSA_Rocket_Factory/notification/internal/config"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

// client реализация Telegram клиента
type client struct {
	bot    *bot.Bot
	chatID int64
}

// startCommandHandler обрабатывает команду /start
func (c *client) startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	message := "🚀 Добро пожаловать в Rocket Factory!\n\n" +
		"Этот бот уведомляет о важных событиях в системе сборки ракет.\n" +
		"Вы будете получать уведомления о:\n" +
		"• Оплаченных заказах\n" +
		"• Собранных кораблях\n" +
		"• Других важных событиях\n\n" +
		"Бот работает автоматически и не требует дополнительных команд."

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   message,
	})
	if err != nil {
		logger.Error(ctx, "Failed to send start command response",
			zap.Error(err),
			zap.Int64("chat_id", update.Message.Chat.ID))
	}
}

// defaultHandler обрабатывает все остальные сообщения
func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Игнорируем все сообщения, кроме команды /start
	// Этот обработчик нужен для корректной работы бота
}

// NewClient создает новый Telegram клиент
func NewClient(ctx context.Context, cfg config.TelegramConfig) (*client, error) {
	// Валидация токена
	token := cfg.BotToken()
	if token == "" {
		return nil, fmt.Errorf("telegram bot token is empty")
	}

	// Убираем лишние пробелы
	token = strings.TrimSpace(token)

	// Для разработки разрешаем тестовые токены
	if !strings.Contains(token, ":") {
		return nil, fmt.Errorf("invalid telegram bot token format: must contain ':'")
	}

	b, err := bot.New(token, bot.WithDefaultHandler(defaultHandler))
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot: %w", err)
	}

	reportBot := &client{
		bot: b,
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, reportBot.startHandler)

	// Проверяем подключение к Telegram API только если не пропущено
	if !cfg.SkipAPICheck() {
		me, err := b.GetMe(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to telegram API: %w", err)
		}

		logger.Info(ctx, "Telegram bot connected successfully",
			zap.String("bot_username", me.Username),
			zap.Int64("bot_id", me.ID))
	} else {
		logger.Info(ctx, "Telegram bot created (API check skipped)")
	}

	// Парсим ChatID
	chatID, err := strconv.ParseInt(cfg.ChatID(), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid chat ID format: %w", err)
	}

	logger.Info(ctx, "Telegram bot configured",
		zap.Int64("chat_id", chatID))

	return &client{
		bot:    b,
		chatID: chatID,
	}, nil
}

// Start запускает бота и начинает обработку команд
func (c *client) Start(ctx context.Context) error {
	logger.Info(ctx, "Starting Telegram bot...")

	// Запускаем бота в отдельной горутине
	go func() {
		c.bot.Start(ctx)
	}()

	logger.Info(ctx, "Telegram bot started successfully")
	return nil
}

// SendMessage отправляет сообщение в указанный чат
func (c *client) SendMessage(ctx context.Context, chatID int64, message string) error {
	// Используем переданный chatID или дефолтный из конфига
	targetChatID := chatID
	if targetChatID == 0 {
		targetChatID = c.chatID
	}

	_, err := c.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: targetChatID,
		Text:   message,
	})
	if err != nil {
		return fmt.Errorf("failed to send telegram message: %w", err)
	}

	logger.Info(ctx, "Telegram message sent successfully",
		zap.Int64("chat_id", targetChatID))

	return nil
}

// Close закрывает телеграм клиент
func (c *client) Close(ctx context.Context) error {
	if c.bot != nil {
		// Останавливаем бота
		// Игнорируем ошибку закрытия, так как это cleanup операция
		//nolint:gosec // Ошибка закрытия бота не критична для cleanup
		_, _ = c.bot.Close(ctx)
	}
	return nil
}
