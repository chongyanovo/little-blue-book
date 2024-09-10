package cache

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type InteractiveCache interface {
	IncreaseReadCountIfPresent(ctx context.Context, biz string, bizId int64) error
}

type RedisInteractiveCache struct {
	redis  redis.Cmdable
	logger *zap.Logger
}

//go:embed lua/interactive_increase.lua
var luaInteractiveIncrease string

func (r *RedisInteractiveCache) IncreaseReadCountIfPresent(ctx context.Context, biz string, bizId int64) error {
	return r.redis.Eval(ctx, luaInteractiveIncrease, []string{r.key(biz, bizId)}, "read_count", 1).Err()
}

func (r *RedisInteractiveCache) key(biz string, bizId int64) string {
	return fmt.Sprintf("%s:%d", biz, bizId)
}

func NewRedisInteractiveCache(redis redis.Cmdable, logger *zap.Logger) *RedisInteractiveCache {
	return &RedisInteractiveCache{redis: redis, logger: logger}
}
