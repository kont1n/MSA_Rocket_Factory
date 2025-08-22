package notification_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	iamMocks "github.com/kont1n/MSA_Rocket_Factory/notification/internal/client/iam/mocks"
	telegramMocks "github.com/kont1n/MSA_Rocket_Factory/notification/internal/client/telegram/mocks"
	"github.com/kont1n/MSA_Rocket_Factory/notification/internal/service"
	"github.com/kont1n/MSA_Rocket_Factory/notification/internal/service/notification"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

type NotificationServiceSuite struct {
	suite.Suite
	service        service.NotificationService
	telegramClient *telegramMocks.TelegramClient
	iamClient      *iamMocks.MockClient
}

func (s *NotificationServiceSuite) SetupSuite() {
	// Инициализируем no-op логгер для тестов
	logger.SetNopLogger()

	s.telegramClient = telegramMocks.NewTelegramClient(s.T())
	s.iamClient = iamMocks.NewMockClient(s.T())

	// Мокаем SetUserRegistrationCallback который вызывается при создании сервиса
	s.telegramClient.On("SetUserRegistrationCallback", mock.AnythingOfType("func(context.Context, string, int64) error")).Return()

	s.service = notification.NewService(context.Background(), s.telegramClient, s.iamClient)
}

func (s *NotificationServiceSuite) SetupTest() {
	// Сбрасываем моки перед каждым тестом
	s.telegramClient.ExpectedCalls = nil
	s.iamClient.ExpectedCalls = nil
}

func (s *NotificationServiceSuite) TearDownSuite() {
}

func TestNotificationServiceSuite(t *testing.T) {
	suite.Run(t, new(NotificationServiceSuite))
}
