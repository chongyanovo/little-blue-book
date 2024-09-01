package service

import (
	"context"
	"errors"
	"github.com/ChongYanOvO/little-blue-book/internal/domain"
	"github.com/ChongYanOvO/little-blue-book/internal/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail
	ErrUserNotFound       = repository.ErrUserNotFound
	ErrInvalidUserOrEmail = errors.New("邮箱或密码不对")
)

type UserService interface {
	Login(ctx context.Context, email, password string) (domain.User, error)
	SignUp(ctx context.Context, u domain.User) error
	Profile(ctx context.Context, id int64) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
}

type UserServiceImpl struct {
	repo   repository.UserRepository
	logger *zap.Logger
}

func NewUserService(repo repository.UserRepository, l *zap.Logger) UserService {
	return &UserServiceImpl{
		repo:   repo,
		logger: l,
	}
}

func (svc *UserServiceImpl) Login(ctx context.Context, email, password string) (domain.User, error) {
	// 先找用户
	u, err := svc.repo.FindByEmail(ctx, email)
	if errors.Is(err, ErrUserNotFound) {
		return domain.User{}, err
	}

	if err != nil {
		return domain.User{}, err
	}
	// 比较密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		// DEBUG
		return domain.User{}, ErrInvalidUserOrEmail
	}
	return u, nil
}

func (svc *UserServiceImpl) SignUp(ctx context.Context, u domain.User) error {
	// 你要考虑加密放在哪里
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hash)
	// 然后就是存起来
	return svc.repo.Create(ctx, u)
}

func (svc *UserServiceImpl) Profile(ctx context.Context, id int64) (domain.User, error) {
	u, err := svc.repo.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return u, err
}

func (svc UserServiceImpl) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {

	if user, err := svc.repo.FindByPhone(ctx, phone); err == nil {
		return user, err
	} else if err := svc.repo.Create(ctx, domain.User{
		Phone: phone,
	}); err != nil {
		return domain.User{}, err
	}
	return svc.repo.FindByPhone(ctx, phone)
}
