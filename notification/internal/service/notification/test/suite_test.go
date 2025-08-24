package notification_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	telegramMocks "github.com/kont1n/MSA_Rocket_Factory/notification/internal/client/telegram/mocks"
	"github.com/kont1n/MSA_Rocket_Factory/notification/internal/service"
	"github.com/kont1n/MSA_Rocket_Factory/notification/internal/service/notification"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

type NotificationServiceSuite struct {
	suite.Suite
	service        service.NotificationService
	telegramClient *telegramMocks.TelegramClient
}

func (s *NotificationServiceSuite) SetupSuite() {
	// Инициализируем no-op логгер для тестов
	logger.SetNopLogger()

	s.telegramClient = telegramMocks.NewTelegramClient(s.T())
	s.service = notification.NewService(s.telegramClient)
}

func (s *NotificationServiceSuite) SetupTest() {
	// Сбрасываем моки перед каждым тестом
	s.telegramClient.ExpectedCalls = nil
}

func (s *NotificationServiceSuite) TearDownSuite() {
}

func TestNotificationServiceSuite(t *testing.T) {
	suite.Run(t, new(NotificationServiceSuite))
}
