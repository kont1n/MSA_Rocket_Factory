package redis

import (
	"context"
	"fmt"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/model"
	repoConverter "github.com/kont1n/MSA_Rocket_Factory/iam/internal/repository/converter"
	repoModel "github.com/kont1n/MSA_Rocket_Factory/iam/internal/repository/model"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/cache"
)

const (
	cacheKeyPrefix = "iam:session:"
)

type repository struct {
	cache cache.RedisClient
}

func NewRepository(cache cache.RedisClient) *repository {
	return &repository{
		cache: cache,
	}
}

// getCacheKey формирует ключ для кеширования сессии в Redis
func (r *repository) getCacheKey(sessionUUID uuid.UUID) string {
	return fmt.Sprintf("%s%s", cacheKeyPrefix, sessionUUID.String())
}

// GetSessionByUUID получает сессию из кеша Redis по UUID
func (r *repository) GetSessionByUUID(ctx context.Context, sessionUUID uuid.UUID) (*model.Session, error) {
	cacheKey := r.getCacheKey(sessionUUID)

	values, err := r.cache.HGetAll(ctx, cacheKey)
	if err != nil {
		if errors.Is(err, redigo.ErrNil) {
			return nil, model.ErrSessionNotFound
		}
		return nil, err
	}

	if len(values) == 0 {
		return nil, model.ErrSessionNotFound
	}

	var sessionRedis repoModel.SessionRedis
	err = redigo.ScanStruct(values, &sessionRedis)
	if err != nil {
		return nil, err
	}

	session, err := repoConverter.ToModelSessionFromRedis(&sessionRedis)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// Set сохраняет сессию в кеше Redis с указанным TTL
func (r *repository) Set(ctx context.Context, sessionUUID uuid.UUID, session *model.Session, ttl time.Duration) error {
	cacheKey := r.getCacheKey(sessionUUID)

	redisSession := repoConverter.ToRepoSessionRedis(session)

	err := r.cache.HashSet(ctx, cacheKey, redisSession)
	if err != nil {
		return err
	}

	return r.cache.Expire(ctx, cacheKey, ttl)
}

// Delete удаляет сессию из кеша Redis
func (r *repository) Delete(ctx context.Context, sessionUUID uuid.UUID) error {
	cacheKey := r.getCacheKey(sessionUUID)
	return r.cache.Del(ctx, cacheKey)
}
