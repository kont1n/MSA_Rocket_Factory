package app

import (
	"context"
	"fmt"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v5/pgxpool"

	authV1API "github.com/kont1n/MSA_Rocket_Factory/iam/internal/api/auth/v1"
	jwtV1API "github.com/kont1n/MSA_Rocket_Factory/iam/internal/api/jwt/v1"
	userV1API "github.com/kont1n/MSA_Rocket_Factory/iam/internal/api/user/v1"
	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/config"
	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/repository"
	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/repository/composite"
	postgresRepository "github.com/kont1n/MSA_Rocket_Factory/iam/internal/repository/postgres"
	redisRepository "github.com/kont1n/MSA_Rocket_Factory/iam/internal/repository/redis"
	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/service"
	iamService "github.com/kont1n/MSA_Rocket_Factory/iam/internal/service/iam"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/cache"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/cache/redis"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/closer"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
	iamV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/iam/v1"
	jwtV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/jwt/v1"
)

type diContainer struct {
	authAPIv1     iamV1.AuthServiceServer
	jwtAPIv1      jwtV1.JWTServiceServer
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

func (d *diContainer) JWTV1API(ctx context.Context) jwtV1.JWTServiceServer {
	if d.jwtAPIv1 == nil {
		d.jwtAPIv1 = jwtV1API.NewAPI(d.IAMService(ctx))
	}
	return d.jwtAPIv1
}

func (d *diContainer) UserV1API(ctx context.Context) iamV1.UserServiceServer {
	if d.userAPIv1 == nil {
		d.userAPIv1 = userV1API.NewAPI(d.IAMService(ctx))
	}
	return d.userAPIv1
}

func (d *diContainer) IAMRepository(ctx context.Context) repository.IAMRepository {
	if d.iamRepository == nil {
		// Создаем PostgreSQL репозиторий
		pgRepo := postgresRepository.NewRepository(
			d.DBPool(ctx),
			config.AppConfig().DB.MigrationsDir(),
		)

		// Создаем Redis репозиторий для кеширования
		redisRepo := redisRepository.NewRepository(d.RedisClient())

		// Создаем композитный репозиторий
		d.iamRepository = composite.NewRepository(pgRepo, redisRepo)
	}
	return d.iamRepository
}

func (d *diContainer) IAMService(ctx context.Context) service.IAMService {
	if d.iamService == nil {
		// Убеждаемся что репозиторий создан перед сервисом
		d.IAMRepository(ctx)

		d.iamService = iamService.NewService(d.IAMRepository(ctx), config.AppConfig().Token)
	}
	return d.iamService
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
		redisConfig := config.AppConfig().Redis
		d.redisPool = &redigo.Pool{
			MaxIdle:     redisConfig.MaxIdle(),
			MaxActive:   100, // Добавляем максимальное количество активных соединений
			IdleTimeout: redisConfig.IdleTimeout(),
			Wait:        true, // Ждать доступного соединения вместо возврата ошибки
			TestOnBorrow: func(c redigo.Conn, t time.Time) error {
				// Проверяем соединение каждые 30 секунд
				if time.Since(t) < time.Minute {
					return nil
				}
				_, err := c.Do("PING")
				return err
			},
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", redisConfig.Address(),
					redigo.DialConnectTimeout(redisConfig.ConnectionTimeout()),
					redigo.DialReadTimeout(10*time.Second),
					redigo.DialWriteTimeout(10*time.Second),
				)
			},
		}

		// Добавляем graceful shutdown для Redis pool
		closer.AddNamed("Redis pool", func(ctx context.Context) error {
			if err := d.redisPool.Close(); err != nil {
				return fmt.Errorf("failed to close redis pool: %w", err)
			}
			return nil
		})
	}

	return d.redisPool
}

func (d *diContainer) RedisClient() cache.RedisClient {
	if d.redisClient == nil {
		d.redisClient = redis.NewClient(d.RedisPool(), logger.Logger(), config.AppConfig().Redis.ConnectionTimeout())
	}

	return d.redisClient
}
