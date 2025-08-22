package telegram

import "context"

// TelegramClient интерфейс для работы с Telegram Bot API
type TelegramClient interface {
	// Start запускает бота и начинает обработку команд
	Start(ctx context.Context) error
	// SendMessage отправляет сообщение в указанный чат
	SendMessage(ctx context.Context, chatID int64, message string) error
	// Close закрывает клиент и освобождает ресурсы
	Close(ctx context.Context) error
	// SetUserRegistrationCallback устанавливает callback для регистрации пользователей
	SetUserRegistrationCallback(callback func(ctx context.Context, username string, chatID int64) error)
	// HandleStartCommand обрабатывает команду /start (для тестирования и разработки)
	HandleStartCommand(ctx context.Context, username string, chatID int64) error
}
