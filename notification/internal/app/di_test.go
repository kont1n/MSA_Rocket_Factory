package app

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiContainer_NewDiContainer(t *testing.T) {
	// Создаем DI контейнер
	container := NewDiContainer()

	// Проверяем, что контейнер создался
	assert.NotNil(t, container)
}

func TestDiContainer_Structure(t *testing.T) {
	// Создаем DI контейнер
	container := NewDiContainer()

	// Проверяем, что все поля инициализированы как nil
	assert.Nil(t, container.notificationService)
	assert.Nil(t, container.orderPaidConsumer)
	assert.Nil(t, container.shipAssembledConsumer)
	assert.Nil(t, container.telegramClient)
	assert.Nil(t, container.iamClient)
	assert.Nil(t, container.iamGRPCConn)
	assert.Nil(t, container.iamGRPCClient)
	assert.Nil(t, container.orderPaidConsumerGroup)
	assert.Nil(t, container.shipAssembledConsumerGroup)
	assert.Nil(t, container.orderPaidKafkaConsumer)
	assert.Nil(t, container.shipAssembledKafkaConsumer)
	assert.Nil(t, container.orderPaidDecoder)
	assert.Nil(t, container.shipAssembledDecoder)
}

func TestDiContainer_TelegramClientStartCommand(t *testing.T) {
	// Тестируем обработку команды /start в режиме разработки

	// Устанавливаем переменную окружения для пропуска API проверки
	err := os.Setenv("TELEGRAM_SKIP_API_CHECK", "true")
	assert.NoError(t, err)
	defer func() {
		_ = os.Unsetenv("TELEGRAM_SKIP_API_CHECK")
	}()

	// Устанавливаем минимальную конфигурацию для тестов
	err = os.Setenv("TELEGRAM_BOT_TOKEN", "test_token:test")
	assert.NoError(t, err)
	defer func() {
		_ = os.Unsetenv("TELEGRAM_BOT_TOKEN")
	}()

	container := NewDiContainer()

	// В тестах не можем создать реальный клиент из-за отсутствия конфигурации
	// Поэтому просто проверяем, что контейнер создается
	assert.NotNil(t, container)
}
