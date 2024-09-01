package repository

import (
	"context"
	"database/sql"
	"github.com/ChongYanOvO/little-blue-book/internal/domain"
	"github.com/ChongYanOvO/little-blue-book/internal/repository/cache"
	"github.com/ChongYanOvO/little-blue-book/internal/repository/dao"
	"go.uber.org/zap"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	Create(ctx context.Context, user domain.User) error
	FindById(ctx context.Context, id int64) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
}

type UserRepositoryImpl struct {
	dao    dao.UserDao
	cache  cache.UserCache
	logger *zap.Logger
}

func NewUserRepository(d dao.UserDao, c cache.UserCache, l *zap.Logger) UserRepository {
	return &UserRepositoryImpl{
		dao:    d,
		cache:  c,
		logger: l,
	}
}

func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.Entity2Domain(u), err
}

func (r *UserRepositoryImpl) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.Domain2Entity(u))
}

// FindById
// 缺点：只要缓存返回了 error，就直接取数据库查询。
//
//	回写缓存的时候，忽略掉了错误
func (r *UserRepositoryImpl) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.cache.Get(ctx, id)

	switch err {
	case nil:
		return u, err
	case cache.ErrKeyNotExist:
		ue, err := r.dao.FindById(ctx, id)
		if err != nil {
			return domain.User{}, err
		}
		u = r.Entity2Domain(ue)
		_ = r.cache.Set(ctx, u)
		return u, nil
	default:
		return domain.User{}, err
	}
}

func (r *UserRepositoryImpl) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return r.Entity2Domain(u), err
}

func (r *UserRepositoryImpl) Entity2Domain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
		Phone:    u.Phone.String,
	}
}

func (r *UserRepositoryImpl) Domain2Entity(u domain.User) dao.User {
	return dao.User{
		Email:    sql.NullString{String: u.Email, Valid: u.Email != ""},
		Password: u.Password,
		Phone:    sql.NullString{String: u.Phone, Valid: u.Phone != ""},
	}
}
