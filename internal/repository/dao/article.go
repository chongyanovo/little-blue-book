package dao

import (
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Article struct {
	Id         int64  `gorm:"primaryKey,autoIncrement"`
	Title      string `gorm:"type=varchar(1024)"`
	Content    string `gorm:"type=BLOB"`
	AuthorId   int64  `gorm:"index=aid_ctime"`
	CreateTime int64  `gorm:"index=aid_ctime"`
	UpdateTime int64
}

type Author struct {
	Id   int64  `gorm:"primaryKey,autoIncrement"`
	Name string `gorm:"type=varchar(1024)"`
}

func (a *Article) TableName() string {
	return "articles"
}

type ArticleDao interface {
	Insert(context.Context, *Article) (int64, error)
	UpdateById(context.Context, *Article) error
	Sync(context.Context, *Article) error
	Upsert(context.Context, *Article) error
}

type ArticleDaoImpl struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewArticleDao(db *gorm.DB, l *zap.Logger) ArticleDao {
	return &ArticleDaoImpl{
		db:     db,
		logger: l,
	}
}

func (dao *ArticleDaoImpl) Insert(ctx context.Context, article *Article) (int64, error) {
	now := time.Now().UnixMilli()
	article.CreateTime = now
	article.UpdateTime = now
	return article.Id, dao.db.WithContext(ctx).Create(&article).Error
}

func (dao *ArticleDaoImpl) UpdateById(ctx context.Context, article *Article) error {
	article.UpdateTime = time.Now().UnixMilli()
	return dao.db.WithContext(ctx).
		Model(&Article{}).
		Where("id=? and author_id=?", article.Id, article.AuthorId).
		Updates(map[string]any{
			"title":       article.Title,
			"content":     article.Content,
			"update_time": article.UpdateTime,
		}).Error
}

func (dao *ArticleDaoImpl) Sync(ctx context.Context, article *Article) error {
	err := dao.db.Transaction(func(tx *gorm.DB) error {
		var (
			id  = article.Id
			err error
		)
		txDao := NewArticleDao(tx, dao.logger)
		if id > 0 {
			err = txDao.UpdateById(ctx, article)
		} else {
			id, err = txDao.Insert(ctx, article)
		}
		if err != nil {
			return err
		}
		return txDao.Upsert(ctx, article)
	})
	return err
}

func (dao *ArticleDaoImpl) Upsert(ctx context.Context, article *Article) error {
	now := time.Now().UnixMilli()
	article.CreateTime = now
	article.UpdateTime = now
	return dao.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]any{
			"title":       article.Title,
			"content":     article.Content,
			"update_time": now,
		}),
	}).Create(&article).Error
}
