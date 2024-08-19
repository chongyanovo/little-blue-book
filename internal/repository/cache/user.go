package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ChongYanOvO/little-blue-book/internal/domain"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

var ErrKeyNotExist = redis.Nil

type UserCache interface {
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, u domain.User) error
}

type RedisUserCache struct {
	redis      redis.Cmdable
	expiration time.Duration
	logger     *zap.Logger
}

func NewRedisUserCache(r redis.Cmdable, l *zap.Logger) UserCache {
	return &RedisUserCache{
		redis: r,
		//expiration: time.Duration(config.CacheConfig.UserExpiration) * time.Minute,
		expiration: 10 * time.Minute,
		logger:     l,
	}
}

// Get 只有 err 为 nil，就认为 u 是一定在的
// 如果没有数据，返回一个特定的 error
func (cache *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.generateKey(id)
	val, err := cache.redis.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, ErrKeyNotExist
	}
	var u domain.User
	err = json.Unmarshal(val, &u)
	return u, err
}

func (cache *RedisUserCache) Set(ctx context.Context, u domain.User) error {
	val, err := json.Marshal(u)
	if err != nil {
		cache.logger.Error("failed to marshal user", zap.Error(err))
		return err
	}
	key := cache.generateKey(u.Id)
	return cache.redis.Set(ctx, key, val, cache.expiration).Err()
}

func (cache *RedisUserCache) generateKey(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
