package migrator

// Migrator интерфейс для выполнения миграций базы данных
type Migrator interface {
	// Up применяет все невыполненные миграции
	Up() error
}
