package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
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

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (d *UserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := d.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (d *UserDao) Insert(ctx context.Context, u User) error {
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

func (d *UserDao) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := d.db.WithContext(ctx).Where("`id` = ?", id).First(&u).Error
	return u, err
}

func (d UserDao) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := d.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
	return u, err
}
