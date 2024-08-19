package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

// User 用户数据库对象
type User struct {
	Id         int64          `gorm:"primaryKey,autoIncrement"` // 用户ID
	Email      sql.NullString `gorm:"unique"`                   // 用户邮箱
	Password   string         // 用户密码
	Phone      sql.NullString `gorm:"unique"`
	CreateTime int64          // 创建时间 毫秒数
	UpdateTime int64          // 更新时间 毫秒数
}

type UserDao interface {
	FindByEmail(ctx context.Context, email string) (User, error)
	Insert(ctx context.Context, u User) error
	FindById(ctx context.Context, id int64) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
}

type UserDaoImpl struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewUserDao(db *gorm.DB, logger *zap.Logger) UserDao {
	if err := db.AutoMigrate(&User{}); err != nil {
		logger.Error("自动建表失败", zap.Error(err))
	}
	return &UserDaoImpl{
		db:     db,
		logger: logger,
	}
}

func (d *UserDaoImpl) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := d.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (d *UserDaoImpl) Insert(ctx context.Context, u User) error {
	// 存毫秒数
	now := time.Now().UnixMilli()

	u.CreateTime = now
	u.UpdateTime = now
	err := d.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictErrNo {
			// 邮箱冲突
			return ErrUserDuplicateEmail
		}
	}
	return err
}

func (d *UserDaoImpl) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := d.db.WithContext(ctx).Where("`id` = ?", id).First(&u).Error
	return u, err
}

func (d UserDaoImpl) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := d.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
	return u, err
}
