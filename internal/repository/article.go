package repository

import (
	"context"
	"github.com/ChongYanOvO/little-blue-book/internal/domain"
	"github.com/ChongYanOvO/little-blue-book/internal/repository/dao"
	"go.uber.org/zap"
)

type ArticleRepository interface {
	Create(ctx context.Context, article *domain.Article) (int64, error)
	Update(ctx context.Context, article *domain.Article) error
	Sync(ctx context.Context, article *domain.Article) error
}

type ArticleRepositoryImpl struct {
	dao    dao.ArticleDao
	logger *zap.Logger
}

func NewArticleRepository(dao dao.ArticleDao, l *zap.Logger) ArticleRepository {
	return &ArticleRepositoryImpl{
		dao:    dao,
		logger: l,
	}
}

func (repo *ArticleRepositoryImpl) Create(ctx context.Context, article *domain.Article) (int64, error) {
	return repo.dao.Insert(ctx, domain2entity(article))
}

func (repo *ArticleRepositoryImpl) Update(ctx context.Context, article *domain.Article) error {
	return repo.dao.UpdateById(ctx, domain2entity(article))
}

func (repo *ArticleRepositoryImpl) Sync(ctx context.Context, article *domain.Article) error {
	return repo.dao.Sync(ctx, domain2entity(article))
}

func domain2entity(article *domain.Article) *dao.Article {
	return &dao.Article{
		Id:       article.Id,
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.Author.Id,
	}
}
