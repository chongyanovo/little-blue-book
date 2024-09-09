package repository

import (
	"context"
	"github.com/ChongYanOvO/little-blue-book/internal/domain"
	"github.com/ChongYanOvO/little-blue-book/internal/repository/dao/article"
	"go.uber.org/zap"
)

type ArticleRepository interface {
	Create(ctx context.Context, article *domain.Article) (int64, error)
	Update(ctx context.Context, article *domain.Article) error
	Sync(ctx context.Context, article *domain.Article) (int64, error)
}

type ArticleRepositoryImpl struct {
	dao    article.ArticleDao
	logger *zap.Logger
}

func NewArticleRepository(dao article.ArticleDao, l *zap.Logger) ArticleRepository {
	return &ArticleRepositoryImpl{
		dao:    dao,
		logger: l,
	}
}

func (repo *ArticleRepositoryImpl) Create(ctx context.Context, article *domain.Article) (int64, error) {
	return repo.dao.Insert(ctx, domain2entity(article))
}

func (repo *ArticleRepositoryImpl) Update(ctx context.Context, article *domain.Article) error {
	return repo.dao.Update(ctx, domain2entity(article))
}

func (repo *ArticleRepositoryImpl) Sync(ctx context.Context, article *domain.Article) (int64, error) {
	return repo.dao.Sync(ctx, *domain2entity(article))
}

func domain2entity(a *domain.Article) *article.Article {
	return &article.Article{
		Id:       a.Id,
		Title:    a.Title,
		Content:  a.Content,
		AuthorId: a.Author.Id,
		Status:   a.Status.ToUint8(),
	}
}
