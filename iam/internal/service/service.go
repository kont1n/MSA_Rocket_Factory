package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/model"
)

// AuthService интерфейс для сервиса аутентификации
type AuthService interface {
	Login(ctx context.Context, login, password string) (*model.Session, error)
	Whoami(ctx context.Context, sessionUUID uuid.UUID) (*model.Session, *model.User, error)
	Logout(ctx context.Context, sessionUUID uuid.UUID) error

	// JWT методы
	JWTLogin(ctx context.Context, login, password string) (*model.TokenPair, error)
	GetAccessToken(ctx context.Context, refreshToken string) (*model.TokenPair, error)
	GetRefreshToken(ctx context.Context, refreshToken string) (*model.TokenPair, error)
}

// UserService интерфейс для сервиса управления пользователями
type UserService interface {
	Register(ctx context.Context, registrationInfo *model.UserRegistrationInfo) (*model.User, error)
	GetUser(ctx context.Context, userUUID uuid.UUID) (*model.User, error)
}

// IAMService объединенный интерфейс для всех сервисов IAM
type IAMService interface {
	AuthService
	UserService
}
