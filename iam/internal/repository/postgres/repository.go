package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/model"
	def "github.com/kont1n/MSA_Rocket_Factory/iam/internal/repository"
)

var _ def.IAMRepository = (*repository)(nil)

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *repository {
	repo := repository{
		db: pool,
	}

	return &repo
}

// Заглушки для SessionCache методов (не используются в PostgreSQL, только в Redis)
func (r *repository) Set(ctx context.Context, sessionUUID uuid.UUID, session *model.Session, ttl time.Duration) error {
	// Не реализован для PostgreSQL
	return model.ErrFailedToCreateSession
}

func (r *repository) Delete(ctx context.Context, sessionUUID uuid.UUID) error {
	// Используем реализованный метод DeleteSession
	return r.DeleteSession(ctx, sessionUUID)
}

func (r *repository) GetSessionFromCache(ctx context.Context, sessionUUID uuid.UUID) (*model.Session, error) {
	// Не реализован для PostgreSQL
	return nil, model.ErrSessionNotFound
}
