package postgres

import (
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"

	def "github.com/kont1n/MSA_Rocket_Factory/iam/internal/repository"
)

var _ def.IAMRepository = (*repository)(nil)

type repository struct {
	db *pgxpool.Pool
}

func (r *repository) Set() {
	//TODO implement me
	panic("implement me")
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
