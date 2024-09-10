package service

import (
	"context"
	"github.com/ChongYanOvO/little-blue-book/internal/repository"
	"github.com/gin-gonic/gin"
)

type InteractiveService interface {
	IncreaseReadCount(ctx context.Context, biz string, bizId int64) error
	IncreaseLikeCount(ctx *gin.Context, biz string, bizId int64, uid int64) error
}

type InteractiveServiceImpl struct {
	repo repository.InteractiveRepository
}

func (svc *InteractiveServiceImpl) IncreaseLikeCount(ctx *gin.Context, biz string, bizId int64, uid int64) error {
	return svc.repo.IncreaseLikeCount(ctx, biz, bizId, uid)
}

func (svc *InteractiveServiceImpl) IncreaseReadCount(ctx context.Context, biz string, bizId int64) error {
	return svc.repo.IncreaseReadCount(ctx, biz, bizId)
}

func NewInteractiveServiceImpl() *InteractiveServiceImpl {
	return &InteractiveServiceImpl{}
}
