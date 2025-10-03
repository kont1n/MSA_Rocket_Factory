package model

import (
	"time"

	"github.com/google/uuid"
)

// UserPostgres модель пользователя для PostgreSQL
type UserPostgres struct {
	UUID         uuid.UUID `db:"user_uuid"`
	Login        string    `db:"login"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// NotificationMethodPostgres модель способа уведомления для PostgreSQL
type NotificationMethodPostgres struct {
	ID           int64     `db:"id"`
	UserUUID     uuid.UUID `db:"user_uuid"`
	ProviderName string    `db:"provider_name"`
	Target       string    `db:"target"`
	CreatedAt    time.Time `db:"created_at"`
}
