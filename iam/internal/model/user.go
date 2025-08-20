package model

import (
	"time"

	"github.com/google/uuid"
)

// User представляет пользователя в системе
type User struct {
	UUID                uuid.UUID
	Login               string
	Email               string
	PasswordHash        string
	NotificationMethods []NotificationMethod
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// NotificationMethod представляет способ уведомления пользователя
type NotificationMethod struct {
	ProviderName string // telegram, email, push и т.д.
	Target       string // email адрес, telegram chat id и т.д.
}

// UserRegistrationInfo содержит данные для регистрации нового пользователя
type UserRegistrationInfo struct {
	Login               string
	Email               string
	Password            string
	NotificationMethods []NotificationMethod
}
