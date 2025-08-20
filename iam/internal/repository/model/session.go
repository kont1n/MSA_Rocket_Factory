package model

import (
	"time"

	"github.com/google/uuid"
)

// SessionPostgres модель сессии для PostgreSQL
type SessionPostgres struct {
	UUID      uuid.UUID `db:"session_uuid"`
	UserUUID  uuid.UUID `db:"user_uuid"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	ExpiresAt time.Time `db:"expires_at"`
}
