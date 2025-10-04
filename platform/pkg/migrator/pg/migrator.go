package pg

import (
	"database/sql"

	"github.com/pressly/goose/v3"

	def "github.com/kont1n/MSA_Rocket_Factory/platform/pkg/migrator"
)

var _ def.Migrator = (*migrator)(nil)

type migrator struct {
	db            *sql.DB
	migrationsDir string
}

// NewMigrator создает новый экземпляр мигратора для PostgreSQL
func NewMigrator(db *sql.DB, migrationsDir string) *migrator {
	return &migrator{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

// Up применяет все невыполненные миграции
func (m *migrator) Up() error {
	return goose.Up(m.db, m.migrationsDir)
}
