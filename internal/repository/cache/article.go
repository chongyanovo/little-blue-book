package cache

import (
	"context"
	"encoding/json"
	"github.com/ChongYanOvO/little-blue-book/internal/domain"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type ArticleCache interface {
	GetFirstPage(ctx context.Context) ([]domain.Article, error)
	SetFirstPage(ctx context.Context, articles []domain.Article) error
	DeleteFirstPage(ctx context.Context)
}

type RedisArticleCache struct {
	redis  redis.Cmdable
	logger *zap.Logger
}

func (ac *RedisArticleCache) DeleteFirstPage(ctx context.Context) {
	ac.redis.Del(ctx, ac.key())
}

func (ac *RedisArticleCache) GetFirstPage(ctx context.Context) (articles []domain.Article, err error) {
	articleJson, err := ac.redis.Get(ctx, ac.key()).Bytes()
	err = json.Unmarshal(articleJson, &articles)
	return articles, err
}

func (ac *RedisArticleCache) SetFirstPage(ctx context.Context, articles []domain.Article) error {
	for _, article := range articles {
		article.Content = article.Abstract()
	}
	articleJson, err := json.Marshal(articles)
	if err != nil {
		ac.logger.Error("格式化json字符串失败", zap.Error(err))
		return err
	}
	ac.redis.Set(ctx, ac.key(), articleJson, 10*time.Minute)
	return nil
}

func (ac *RedisArticleCache) key() string {
	return "article:first_page"
}

func NewRedisArticleCache(r redis.Cmdable, l *zap.Logger) *RedisArticleCache {
	return &RedisArticleCache{
		redis:  r,
		logger: l,
	}
}
