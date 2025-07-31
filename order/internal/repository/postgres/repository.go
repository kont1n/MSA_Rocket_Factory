package postgres

import (
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"

	def "github.com/kont1n/MSA_Rocket_Factory/order/internal/repository"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	db   *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		db: pool,
	}
}