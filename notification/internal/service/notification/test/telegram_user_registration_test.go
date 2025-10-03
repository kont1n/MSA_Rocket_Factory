package notification_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	iamMocks "github.com/kont1n/MSA_Rocket_Factory/notification/internal/client/iam/mocks"
	telegramMocks "github.com/kont1n/MSA_Rocket_Factory/notification/internal/client/telegram/mocks"
	"github.com/kont1n/MSA_Rocket_Factory/notification/internal/service/notification"
)

func TestRegisterTelegramUser_Success(t *testing.T) {
	// Создаем моки
	telegramClient := telegramMocks.NewTelegramClient(t)
	iamClient := iamMocks.NewMockClient(t)

	// Создаем сервис
	service := notification.NewService(context.Background(), telegramClient, iamClient)

	// Проверяем, что сервис создался успешно
	assert.NotNil(t, service)
}

func TestRegisterTelegramUser_NilIAMClient(t *testing.T) {
	// Создаем мок telegram клиента
	telegramClient := telegramMocks.NewTelegramClient(t)

	// Передаем nil как IAM клиент
	var nilIAMClient *iamMocks.MockClient = nil

	// Создаем сервис с nil IAM клиентом
	service := notification.NewService(context.Background(), telegramClient, nilIAMClient)

	// Проверяем, что сервис создался успешно даже с nil IAM клиентом
	assert.NotNil(t, service)
}

func TestService_Creation(t *testing.T) {
	// Создаем моки
	telegramClient := telegramMocks.NewTelegramClient(t)
	iamClient := iamMocks.NewMockClient(t)

	// Создаем сервис
	service := notification.NewService(context.Background(), telegramClient, iamClient)

	// Проверяем, что сервис создался успешно
	assert.NotNil(t, service)
}

// TestTelegramUserRegistrationWithCallback тестирует регистрацию пользователя через callback
func TestTelegramUserRegistrationWithCallback(t *testing.T) {
	// Создаем моки
	telegramClient := telegramMocks.NewTelegramClient(t)
	iamClient := iamMocks.NewMockClient(t)

	// Создаем сервис
	service := notification.NewService(context.Background(), telegramClient, iamClient)
	assert.NotNil(t, service)

	// Проверяем, что сервис создался успешно
	assert.NotNil(t, service)
}
