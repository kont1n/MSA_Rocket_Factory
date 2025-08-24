package postgres

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"

	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/model"
	def "github.com/kont1n/MSA_Rocket_Factory/iam/internal/repository"
)

var _ def.IAMRepository = (*repository)(nil)

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool, migrationsDir string) *repository {
	repo := repository{
		db: pool,
	}
	err := repo.Migrate(migrationsDir)
	if err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	return &repo
}

// Migrate выполняет миграции базы данных
func (r *repository) Migrate(migrationsDir string) error {
	sqlDB := stdlib.OpenDBFromPool(r.db)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	if err := goose.Up(sqlDB, migrationsDir); err != nil {
		return err
	}

	log.Println("✅ Миграции IAM сервиса успешно применены.")
	return nil
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
