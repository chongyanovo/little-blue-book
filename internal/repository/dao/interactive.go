package dao

import (
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type InteractiveDao interface {
	IncreaseReadCount(ctx context.Context, biz string, bizId int64) error
	IncreaseLikeCount(ctx context.Context, biz string, bizId int64, uid int64) error
	DeletedLike(ctx context.Context, biz string, bizId int64, uid int64) error
}

type InteractiveDaoMysql struct {
	db     *gorm.DB
	logger *zap.Logger
}

func (dao *InteractiveDaoMysql) DeletedLike(ctx context.Context, biz string, bizId int64, uid int64) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&UserLikeBiz{}).
			Where("uid = ? and biz = ? and biz_id = ?", uid, biz, bizId).
			Updates(map[string]any{
				"status":      0,
				"update_time": now,
			}).Error
		if err != nil {
			return err
		}
		return tx.Model(&Interactive{}).
			Where("biz = ? and biz_id = ?", biz, bizId).
			Updates(map[string]any{
				"like_count":  gorm.Expr("like_count - 1"),
				"update_time": now,
			}).Error
	})
}

func (dao *InteractiveDaoMysql) IncreaseLikeCount(ctx context.Context, biz string, bizId int64, uid int64) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			err := tx.Model(&UserLikeBiz{}).Clauses(clause.OnConflict{
				DoUpdates: clause.Assignments(map[string]any{
					"status":      1,
					"update_time": now,
				}),
			}).Create(&UserLikeBiz{
				Uid:        uid,
				Biz:        biz,
				BizId:      bizId,
				Status:     1,
				CreateTime: now,
				UpdateTime: now,
			}).Error
			if err != nil {
				return err
			}
			return tx.Model(&Interactive{}).Clauses(clause.OnConflict{
				DoUpdates: clause.Assignments(map[string]any{
					"like_count": gorm.Expr("like_count + 1"),
				}),
			}).Error
		})
}

func (dao *InteractiveDaoMysql) IncreaseReadCount(ctx context.Context, biz string, bizId int64) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Model(&Interactive{}).
		Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"read_count": gorm.Expr("read_count + 1"),
			}),
		}).
		Create(&Interactive{
			Biz:        biz,
			BizId:      bizId,
			ReadCount:  1,
			CreateTime: now,
			UpdateTime: now,
		}).Error
}

func NewInteractiveDaoMysql(db *gorm.DB, logger *zap.Logger) *InteractiveDaoMysql {
	if err := db.AutoMigrate(&Interactive{}); err != nil {
		logger.Error("初始化点赞收藏表失败", zap.Error(err))
		return nil
	}
	if err := db.AutoMigrate(&UserLikeBiz{}); err != nil {
		logger.Error("初始化用户点赞表失败", zap.Error(err))
		return nil
	}
	return &InteractiveDaoMysql{
		db:     db,
		logger: logger,
	}
}

type Interactive struct {
	Id            int64  `gorm:"primaryKey,autoIncrement"`
	BizId         int64  `gorm:"uniqueIndex=bizId_type"`
	Biz           string `gorm:"uniqueIndex=bizId_type"`
	ReadCount     int64
	LikeCount     int64
	FavoriteCount int64
	CreateTime    int64
	UpdateTime    int64
}

func (i *Interactive) TableName() string {
	return "interactive"
}

type UserLikeBiz struct {
	Id         int64  `gorm:"primaryKey,autoIncrement"`
	Uid        int64  `gorm:"uniqueIndex=uid_biz_bizId"`
	Biz        string `gorm:"uniqueIndex=uid_biz_bizId"`
	BizId      int64  `gorm:"uniqueIndex=uid_biz_bizId"`
	Status     int
	CreateTime int64
	UpdateTime int64
}

func (b *UserLikeBiz) TableName() string {
	return "user_like_biz"
}
