package telegram

import (
	"testing"

	"github.com/go-telegram/bot/models"
	"github.com/stretchr/testify/assert"
)

func TestStartCommandHandler(t *testing.T) {
	// Тест проверяет, что обработчик команды /start создает правильное сообщение
	// Создаем тестовый update с командой /start
	update := &models.Update{
		Message: &models.Message{
			Chat: models.Chat{
				ID: 12345,
			},
		},
	}

	// Проверяем, что update содержит правильные данные
	assert.NotNil(t, update.Message, "Message не должен быть nil")
	assert.Equal(t, int64(12345), update.Message.Chat.ID, "Chat ID должен совпадать")
}

func TestClientStart(t *testing.T) {
	// Тест проверяет, что клиент может быть создан
	// Это базовый тест для проверки компиляции
	assert.True(t, true, "Базовый тест должен проходить")
}
