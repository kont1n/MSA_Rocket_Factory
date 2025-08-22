package app

import (
	"context"
	"fmt"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/cache/redis"

	authV1API "github.com/kont1n/MSA_Rocket_Factory/iam/internal/api/auth/v1"
	userV1API "github.com/kont1n/MSA_Rocket_Factory/iam/internal/api/user/v1"
	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/config"
	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/repository"
	iamRepository "github.com/kont1n/MSA_Rocket_Factory/iam/internal/repository/postgres"
	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/service"
	iamService "github.com/kont1n/MSA_Rocket_Factory/iam/internal/service/iam"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/cache"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/closer"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
	iamV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/iam/v1"
)

type diContainer struct {
	authAPIv1     iamV1.AuthServiceServer
	userAPIv1     iamV1.UserServiceServer
	iamService    service.IAMService
	iamRepository repository.IAMRepository
	dbPool        *pgxpool.Pool
	redisPool     *redigo.Pool
	redisClient   cache.RedisClient
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) AuthV1API(ctx context.Context) iamV1.AuthServiceServer {
	if d.authAPIv1 == nil {
		d.authAPIv1 = authV1API.NewAPI(d.IAMService(ctx))
	}
	return d.authAPIv1
}

func (d *diContainer) UserV1API(ctx context.Context) iamV1.UserServiceServer {
	if d.userAPIv1 == nil {
		d.userAPIv1 = userV1API.NewAPI(d.IAMService(ctx))
	}
	return d.userAPIv1
}

func (d *diContainer) IAMService(ctx context.Context) service.IAMService {
	if d.iamService == nil {
		d.iamService = iamService.NewService(d.IAMRepository(ctx))
	}
	return d.iamService
}

func (d *diContainer) IAMRepository(ctx context.Context) repository.IAMRepository {
	if d.iamRepository == nil {
		d.iamRepository = iamRepository.NewRepository(
			d.DBPool(ctx),
			config.AppConfig().DB.MigrationsDir(),
		)
	}
	return d.iamRepository
}

func (d *diContainer) DBPool(ctx context.Context) *pgxpool.Pool {
	if d.dbPool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().DB.URI())
		if err != nil {
			panic(fmt.Sprintf("failed to connect to database: %v", err))
		}

		closer.AddNamed("DB pool", func(ctx context.Context) error {
			d.dbPool.Close()
			return nil
		})

		d.dbPool = pool
	}
	return d.dbPool
}

func (d *diContainer) RedisPool() *redigo.Pool {
	if d.redisPool == nil {
		d.redisPool = &redigo.Pool{
			MaxIdle:     config.AppConfig().Redis.MaxIdle(),
			IdleTimeout: config.AppConfig().Redis.IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", config.AppConfig().Redis.Address())
			},
		}
	}

	return d.redisPool
}

func (d *diContainer) RedisClient() cache.RedisClient {
	if d.redisClient == nil {
		d.redisClient = redis.NewClient(d.RedisPool(), logger.Logger(), config.AppConfig().Redis.ConnectionTimeout())
	}

	return d.redisClient
}
