package article

import (
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ArticleDao interface {
	Insert(context.Context, *Article) (int64, error)
	Update(context.Context, *Article) error
	Sync(context.Context, Article) (int64, error)
	Upsert(context.Context, *PublishedArticle) error
	List(ctx context.Context, offset int, limit int) ([]Article, error)
}

type ArticleDaoImpl struct {
	db     *gorm.DB
	logger *zap.Logger
}

func (dao *ArticleDaoImpl) Insert(ctx context.Context, article *Article) (int64, error) {
	now := time.Now().UnixMilli()
	article.CreateTime = now
	article.UpdateTime = now
	return article.Id, dao.db.WithContext(ctx).Model(&Article{}).Create(&article).Error
}

func (dao *ArticleDaoImpl) Update(ctx context.Context, a *Article) error {
	a.UpdateTime = time.Now().UnixMilli()
	return dao.db.WithContext(ctx).
		Model(&Article{}).
		Where("id=? and author_id=?", a.Id, a.AuthorId).
		Updates(map[string]any{
			"title":       a.Title,
			"content":     a.Content,
			"status":      a.Status,
			"update_time": a.UpdateTime,
		}).Error
}

func (dao *ArticleDaoImpl) Sync(ctx context.Context, article Article) (int64, error) {
	var (
		id  = article.Id
		err error
	)
	err = dao.db.Transaction(func(tx *gorm.DB) error {
		txDao := NewArticleDao(tx, dao.logger)
		if id > 0 {
			err = txDao.Update(ctx, &article)
		} else {
			id, err = txDao.Insert(ctx, &article)
		}
		if err != nil {
			return err
		}
		publishedArticle := PublishedArticle(article)
		publishedArticle.Id = id
		return txDao.Upsert(ctx, &publishedArticle)
	})
	return id, err
}

func (dao *ArticleDaoImpl) Upsert(ctx context.Context, article *PublishedArticle) error {
	now := time.Now().UnixMilli()
	article.CreateTime = now
	article.UpdateTime = now
	return dao.db.WithContext(ctx).Model(&PublishedArticle{}).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.Assignments(map[string]any{
				"title":       article.Title,
				"content":     article.Content,
				"status":      article.Status,
				"update_time": article.UpdateTime,
			}),
		}).Create(&article).Error
}

func (dao *ArticleDaoImpl) List(ctx context.Context, offset int, limit int) ([]Article, error) {
	articles := []Article{}
	err := dao.db.WithContext(ctx).Model(&Article{}).
		Offset(offset).Limit(limit).
		Order("update_time desc").Find(&articles).Error
	return articles, err
}

func NewArticleDao(db *gorm.DB, l *zap.Logger) ArticleDao {
	if err := db.AutoMigrate(&Article{}); err != nil {
		l.Error("初始化制作库失败", zap.Error(err))
		return nil
	}
	if err := db.AutoMigrate(&PublishedArticle{}); err != nil {
		l.Error("初始化线上库失败", zap.Error(err))
		return nil
	}
	return &ArticleDaoImpl{
		db:     db,
		logger: l,
	}
}

//
//type MongoArticleDao struct {
//	collection *mongo.Collection
//	logger     *zap.Logger
//}
//
//func NewMongoArticleDao(collection *mongo.Collection, l *zap.Logger) ArticleDao {
//	return &MongoArticleDao{
//		collection: collection,
//		logger:     l,
//	}
//}
//
//func (dao *MongoArticleDao) Insert(ctx context.Context, article *Article) (int64, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (dao *MongoArticleDao) Update(ctx context.Context, article *Article) error {
//	filter := bson.M{"id": article.Id, "author_id": article.AuthorId}
//	update := bson.D{bson.E{"$set", bson.M{
//		"title":       article.Title,
//		"content":     article.Content,
//		"status":      article.Status,
//		"update_time": article.UpdateTime,
//	}}}
//	updateResult, err := dao.collection.UpdateOne(ctx, filter, update)
//	if err != nil {
//		dao.logger.Error("mongo:Update:updateResult", zap.Error(err))
//		return err
//	}
//	if updateResult.ModifiedCount == 0 {
//		return errors.New("更新文档失败")
//	}
//	return nil
//}
//
//func (dao *MongoArticleDao) Sync(ctx context.Context, article *Article) error {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (dao *MongoArticleDao) Upsert(ctx context.Context, article *Article) error {
//	//TODO implement me
//	panic("implement me")
//}
