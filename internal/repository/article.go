package repository

import (
	"context"
	"github.com/ChongYanOvO/little-blue-book/internal/domain"
	"github.com/ChongYanOvO/little-blue-book/internal/repository/cache"
	"github.com/ChongYanOvO/little-blue-book/internal/repository/dao/article"
	"github.com/chongyanovo/zkit/slice"
	"go.uber.org/zap"
)

type ArticleRepository interface {
	Create(ctx context.Context, article *domain.Article) (int64, error)
	Update(ctx context.Context, article *domain.Article) error
	Sync(ctx context.Context, article *domain.Article) (int64, error)
	List(ctx context.Context, offset int, limit int) ([]domain.Article, error)
}

type ArticleRepositoryImpl struct {
	dao    article.ArticleDao
	cache  cache.ArticleCache
	logger *zap.Logger
}

func NewArticleRepository(dao article.ArticleDao, cache cache.ArticleCache, logger *zap.Logger) ArticleRepository {
	return &ArticleRepositoryImpl{
		dao:    dao,
		cache:  cache,
		logger: logger,
	}
}

func (repo *ArticleRepositoryImpl) Create(ctx context.Context, article *domain.Article) (int64, error) {
	defer func() {
		repo.cache.DeleteFirstPage(ctx)
	}()
	return repo.dao.Insert(ctx, domain2entity(article))
}

func (repo *ArticleRepositoryImpl) Update(ctx context.Context, article *domain.Article) error {
	defer func() {
		repo.cache.DeleteFirstPage(ctx)
	}()
	return repo.dao.Update(ctx, domain2entity(article))
}

func (repo *ArticleRepositoryImpl) Sync(ctx context.Context, article *domain.Article) (int64, error) {
	defer func() {
		repo.cache.DeleteFirstPage(ctx)
	}()
	return repo.dao.Sync(ctx, *domain2entity(article))
}

func (repo *ArticleRepositoryImpl) List(ctx context.Context, offset int, limit int) ([]domain.Article, error) {
	if offset == 0 && limit <= 100 {
		data, err := repo.cache.GetFirstPage(ctx)
		if err == nil {
			return data, err
		}
	}
	articles, err := repo.dao.List(ctx, offset, limit)
	if err != nil {
		repo.logger.Error("查询文章列表失败", zap.Error(err))
		return nil, err
	}

	data := slice.Map[article.Article, domain.Article](articles, func(idx int, src article.Article) domain.Article {
		return *entity2domain(&src)
	})
	go func() {
		if err := repo.cache.SetFirstPage(ctx, data); err != nil {
			repo.logger.Error("文章列表缓存回写失败", zap.Error(err))
		}
	}()
	return data, nil
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

func entity2domain(a *article.Article) *domain.Article {
	return &domain.Article{
		Id:      a.Id,
		Title:   a.Title,
		Content: a.Content,
		Author:  domain.Author{Id: a.AuthorId},
		Status:  domain.ArticleStates(a.Status),
	}
}
