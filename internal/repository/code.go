package repository

import (
	"context"
	"github.com/ChongYanOvO/little-blue-book/internal/repository/cache"
	"go.uber.org/zap"
)

var (
	ErrCodeSendToMany         = cache.ErrSetCodeTooMany
	ErrCodeVerifyTooManyTimes = cache.ErrCodeVerifyTooManyTimes
)

type CodeRepository interface {
	Store(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type CodeRepositoryImpl struct {
	cache  *cache.CodeCache
	logger *zap.Logger
}

func NewCodeRepository(c *cache.CodeCache, l *zap.Logger) CodeRepository {
	return &CodeRepositoryImpl{
		cache:  c,
		logger: l,
	}
}

func (repo *CodeRepositoryImpl) Store(ctx context.Context, biz, phone, code string) error {
	return repo.cache.Set(ctx, biz, phone, code)
}

func (repo *CodeRepositoryImpl) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return repo.cache.Verify(ctx, biz, phone, inputCode)
}
