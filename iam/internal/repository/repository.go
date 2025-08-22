package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/model"
)

// IAMRepository объединенный интерфейс для всех репозиториев IAM
type IAMRepository interface {
	UserRepository
	SessionRepository
	SessionCache
}

// UserRepository интерфейс для работы с пользователями
type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUserByUUID(ctx context.Context, userUUID uuid.UUID) (*model.User, error)
	GetUserByLogin(ctx context.Context, login string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) (*model.User, error)
}

// SessionRepository интерфейс для работы с сессиями
type SessionRepository interface {
	CreateSession(ctx context.Context, session *model.Session) (*model.Session, error)
	GetSessionByUUID(ctx context.Context, sessionUUID uuid.UUID) (*model.Session, error)
	UpdateSession(ctx context.Context, session *model.Session) (*model.Session, error)
	DeleteSession(ctx context.Context, sessionUUID uuid.UUID) error
	CleanupExpiredSessions(ctx context.Context) error
}

type SessionCache interface {
	GetSessionByUUID(ctx context.Context, sessionUUID uuid.UUID) (*model.Session, error)
	Set(ctx context.Context, sessionUUID uuid.UUID, *model.Session, ttl time.Duration) error
	Delete(ctx context.Context, sessionUUID uuid.UUID) error
}
