package model

import (
	"time"

	"github.com/google/uuid"
)

// Session представляет пользовательскую сессию
type Session struct {
	UUID      uuid.UUID
	UserUUID  uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}

// IsExpired проверяет, истекла ли сессия
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// ExtendSession продлевает сессию на указанную длительность
func (s *Session) ExtendSession(duration time.Duration) {
	s.ExpiresAt = time.Now().Add(duration)
	s.UpdatedAt = time.Now()
}
