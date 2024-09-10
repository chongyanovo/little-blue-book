package repository

import (
	"context"
	"github.com/ChongYanOvO/little-blue-book/internal/repository/cache"
	"github.com/ChongYanOvO/little-blue-book/internal/repository/dao"
)

type InteractiveRepository interface {
	IncreaseReadCount(ctx context.Context, biz string, bizId int64) error
	IncreaseLikeCount(ctx context.Context, biz string, id int64, uid int64) error
}

type InteractiveRepositoryImpl struct {
	dao   dao.InteractiveDao
	cache cache.InteractiveCache
}

func (repo *InteractiveRepositoryImpl) IncreaseLikeCount(ctx context.Context, biz string, id int64, uid int64) error {
	repo.dao.IncreaseLikeCount(ctx, biz, id, uid)
	//TODO implement me
	panic("implement me")
}

func (repo *InteractiveRepositoryImpl) IncreaseReadCount(ctx context.Context, biz string, bizId int64) error {
	if err := repo.dao.IncreaseReadCount(ctx, biz, bizId); err != nil {
		return err
	}
	return repo.cache.IncreaseReadCountIfPresent(ctx, biz, bizId)
}

func NewInteractiveRepositoryImpl() *InteractiveRepositoryImpl {
	return &InteractiveRepositoryImpl{}
}
