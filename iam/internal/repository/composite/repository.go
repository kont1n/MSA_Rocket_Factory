package composite

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/repository"
)

var _ repository.IAMRepository = (*compositeRepository)(nil)

// compositeRepository композитный репозиторий, объединяющий PostgreSQL для основной логики и Redis для кеширования
type compositeRepository struct {
	postgres repository.IAMRepository
	cache    repository.SessionCache
}

// NewRepository создает новый композитный репозиторий
func NewRepository(postgres repository.IAMRepository, cache repository.SessionCache) repository.IAMRepository {
	return &compositeRepository{
		postgres: postgres,
		cache:    cache,
	}
}

// Методы для работы с пользователями (делегируются PostgreSQL)
func (r *compositeRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	return r.postgres.CreateUser(ctx, user)
}

func (r *compositeRepository) GetUserByUUID(ctx context.Context, userUUID uuid.UUID) (*model.User, error) {
	return r.postgres.GetUserByUUID(ctx, userUUID)
}

func (r *compositeRepository) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	return r.postgres.GetUserByLogin(ctx, login)
}

func (r *compositeRepository) UpdateUser(ctx context.Context, user *model.User) (*model.User, error) {
	return r.postgres.UpdateUser(ctx, user)
}

// Методы для работы с сессиями в основной БД (делегируются PostgreSQL)
func (r *compositeRepository) CreateSession(ctx context.Context, session *model.Session) (*model.Session, error) {
	return r.postgres.CreateSession(ctx, session)
}

func (r *compositeRepository) GetSessionByUUID(ctx context.Context, sessionUUID uuid.UUID) (*model.Session, error) {
	return r.postgres.GetSessionByUUID(ctx, sessionUUID)
}

func (r *compositeRepository) UpdateSession(ctx context.Context, session *model.Session) (*model.Session, error) {
	return r.postgres.UpdateSession(ctx, session)
}

func (r *compositeRepository) DeleteSession(ctx context.Context, sessionUUID uuid.UUID) error {
	return r.postgres.DeleteSession(ctx, sessionUUID)
}

func (r *compositeRepository) CleanupExpiredSessions(ctx context.Context) error {
	return r.postgres.CleanupExpiredSessions(ctx)
}

// Методы для работы с кешем сессий (делегируются Redis)
func (r *compositeRepository) GetSessionFromCache(ctx context.Context, sessionUUID uuid.UUID) (*model.Session, error) {
	return r.cache.GetSessionByUUID(ctx, sessionUUID)
}

func (r *compositeRepository) Set(ctx context.Context, sessionUUID uuid.UUID, session *model.Session, ttl time.Duration) error {
	return r.cache.Set(ctx, sessionUUID, session, ttl)
}

func (r *compositeRepository) Delete(ctx context.Context, sessionUUID uuid.UUID) error {
	return r.cache.Delete(ctx, sessionUUID)
}
