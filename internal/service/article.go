package service

import (
	"context"
	"github.com/ChongYanOvO/little-blue-book/internal/domain"
	"github.com/ChongYanOvO/little-blue-book/internal/repository"
	"go.uber.org/zap"
)

type ArticleService interface {
	Save(ctx context.Context, article *domain.Article) (int64, error)
	Create(ctx context.Context, article *domain.Article) (int64, error)
	Update(ctx context.Context, article *domain.Article) error
	Publish(ctx context.Context, article *domain.Article) (int64, error)
}

type ArticleServiceImpl struct {
	repo   repository.ArticleRepository
	logger *zap.Logger
}

func NewArticleService(repo repository.ArticleRepository, l *zap.Logger) ArticleService {
	return &ArticleServiceImpl{
		repo:   repo,
		logger: l,
	}
}

func (svc *ArticleServiceImpl) Create(ctx context.Context, article *domain.Article) (int64, error) {
	return svc.repo.Create(ctx, article)
}

func (svc *ArticleServiceImpl) Update(ctx context.Context, article *domain.Article) error {
	return svc.repo.Update(ctx, article)
}

func (svc *ArticleServiceImpl) Save(ctx context.Context, article *domain.Article) (int64, error) {
	article.Status = domain.ArticleStatusUnpublished
	if article.Id > 0 {
		if err := svc.repo.Update(ctx, article); err != nil {
			svc.logger.Error("更新文章失败", zap.Error(err))
		}
		return article.Id, nil
	}
	return svc.repo.Create(ctx, article)
}

func (svc *ArticleServiceImpl) Publish(ctx context.Context, article *domain.Article) (int64, error) {
	article.Status = domain.ArticleStatusPublished
	return svc.repo.Sync(ctx, article)
}
