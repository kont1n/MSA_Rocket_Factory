package telegram

import (
	"bytes"
	"context"
	"embed"
	"text/template"
	"time"

	"go.uber.org/zap"

	"github.com/kont1n/MSA_Rocket_Factory/notification/internal/client/http"
	"github.com/kont1n/MSA_Rocket_Factory/notification/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

const chatID = 234586218

//go:embed templates/notification.tmpl
var templateFS embed.FS

type ufoTemplateData struct {
	UUID            string
	Location        string
	Description     string
	ObservedAt      *time.Time
	Color           *string
	Sound           *soundValue
	DurationSeconds *int32
	RegisteredAt    time.Time
}

type soundValue struct {
	Value bool
}

var ufoTemplate = template.Must(template.ParseFS(templateFS, "templates/ufo_notification.tmpl"))

type service struct {
	telegramClient http.TelegramClient
}

// NewService создает новый Telegram сервис
func NewService(telegramClient http.TelegramClient) *service {
	return &service{
		telegramClient: telegramClient,
	}
}

// SendUFONotification отправляет уведомление о новом наблюдении UFO
func (s *service) SendUFONotification(ctx context.Context, uuid string, sighting model.SightingInfo) error {
	message, err := s.buildUFOMessage(uuid, sighting)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, chatID, message)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message sent to chat", zap.Int("chat_id", chatID), zap.String("message", message))
	return nil
}

// buildUFOMessage создает сообщение о наблюдении UFO из шаблона
func (s *service) buildUFOMessage(uuid string, sighting model.SightingInfo) (string, error) {
	data := ufoTemplateData{
		UUID:         uuid,
		Location:     sighting.Location,
		Description:  sighting.Description,
		ObservedAt:   sighting.ObservedAt,
		Color:        sighting.Color,
		RegisteredAt: time.Now(),
	}

	if sighting.Sound != nil {
		data.Sound = &soundValue{Value: *sighting.Sound}
	}

	if sighting.DurationSeconds != nil {
		data.DurationSeconds = sighting.DurationSeconds
	}

	var buf bytes.Buffer
	err := ufoTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
