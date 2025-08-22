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

// TokenPair - пара токенов
type TokenPair struct {
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

// Claims - кастомные claims для JWT
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
}
