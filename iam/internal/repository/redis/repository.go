package redis

import (
	"context"
	"fmt"
	"time"

	redigo "github.com/gomodule/redigo/redis"
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

// api -> service -> cache repo -> redis client (обертка наша) -> redis
func (r *repository) getCacheKey(uuid string) string {
	return fmt.Sprintf("%s%s", cacheKeyPrefix, uuid)
}

func (r *repository) Get(ctx context.Context, uuid string) (model.Sighting, error) {
	cacheKey := r.getCacheKey(uuid)

	values, err := r.cache.HGetAll(ctx, cacheKey)
	if err != nil {
		if errors.Is(err, redigo.ErrNil) {
			return model.Sighting{}, model.ErrSightingNotFound
		}
		return model.Sighting{}, err
	}

	if len(values) == 0 {
		return model.Sighting{}, model.ErrSightingNotFound
	}

	var sightingRedisView repoModel.SightingRedisView
	err = redigo.ScanStruct(values, &sightingRedisView)
	if err != nil {
		return model.Sighting{}, err
	}

	return repoConverter.SightingFromRedisView(sightingRedisView), nil
}

func (r *repository) Set(ctx context.Context, uuid string, sighting model.Sighting, ttl time.Duration) error {
	cacheKey := r.getCacheKey(uuid)

	redisView := repoConverter.SightingToRedisView(sighting)

	err := r.cache.HashSet(ctx, cacheKey, redisView)
	if err != nil {
		return err
	}

	return r.cache.Expire(ctx, cacheKey, ttl)
}

func (r *repository) Delete(ctx context.Context, uuid string) error {
	cacheKey := r.getCacheKey(uuid)
	return r.cache.Del(ctx, cacheKey)
}
