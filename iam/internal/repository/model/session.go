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

// SessionRedis модель для сессии в Redis hash map
type SessionRedis struct {
	UUID      string `redis:"uuid"`
	UserUUID  string `redis:"user_uuid"`
	CreatedAt int64  `redis:"created_at"`
	UpdatedAt *int64 `redis:"updated_at,omitempty"`
	ExpiresAt int64  `redis:"expires_at"`
}
