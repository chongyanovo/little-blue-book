package service

import (
	"context"
	"errors"
	"github.com/ChongYanOvO/little-blue-book/internal/domain"
	"github.com/ChongYanOvO/little-blue-book/internal/repository"
	repomock "github.com/ChongYanOvO/little-blue-book/internal/repository/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestUserServiceImpl_Login(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctl *gomock.Controller) repository.UserRepository
		email    string
		password string
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登录成功",
			mock: func(ctl *gomock.Controller) repository.UserRepository {
				ur := repomock.NewMockUserRepository(ctl)
				ur.EXPECT().FindByEmail(gomock.Any(), "test@gmail.com").
					Return(domain.User{
						Email:    "test@gmail.com",
						Password: "$2a$10$DEFY1AeFZidKeHuKVleFSueNUOP9mjiNq7YmCmyXA/Miwqyrk.1Ze",
						Phone:    "1234567890",
					}, nil)
				return ur
			},
			email:    "test@gmail.com",
			password: "1qaz@WSX",
			wantUser: domain.User{
				Email:    "test@gmail.com",
				Password: "$2a$10$DEFY1AeFZidKeHuKVleFSueNUOP9mjiNq7YmCmyXA/Miwqyrk.1Ze",
				Phone:    "1234567890",
			},
			wantErr: nil,
		},
		{
			name: "用户不存在",
			mock: func(ctl *gomock.Controller) repository.UserRepository {
				ur := repomock.NewMockUserRepository(ctl)
				ur.EXPECT().FindByEmail(gomock.Any(), "test@gmail.com").
					Return(domain.User{}, repository.ErrUserNotFound)
				return ur
			},
			email:    "test@gmail.com",
			password: "1qaz@WSX",
			wantUser: domain.User{},
			wantErr:  ErrUserNotFound,
		},
		{
			name: "数据库错误",
			mock: func(ctl *gomock.Controller) repository.UserRepository {
				ur := repomock.NewMockUserRepository(ctl)
				ur.EXPECT().FindByEmail(gomock.Any(), "test@gmail.com").
					Return(domain.User{}, errors.New("数据库错误"))
				return ur
			},
			email:    "test@gmail.com",
			password: "1qaz@WSX2",
			wantUser: domain.User{},
			wantErr:  errors.New("数据库错误"),
		},
		{
			name: "密码错误",
			mock: func(ctl *gomock.Controller) repository.UserRepository {
				ur := repomock.NewMockUserRepository(ctl)
				ur.EXPECT().FindByEmail(gomock.Any(), "test@gmail.com").
					Return(domain.User{
						Email:    "test@gmail.com",
						Password: "$2a$10$DEFY1AeFZidKeHuKVleFSueNUOP9mjiNq7YmCmyXA/Miwqyrk.1Ze",
						Phone:    "1234567890",
					}, nil)
				return ur
			},
			email:    "test@gmail.com",
			password: "1222qaz@WSX",
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrEmail,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			ur := tc.mock(ctl)
			svc := NewUserService(ur, nil)
			u, err := svc.Login(context.Background(), tc.email, tc.password)
			assert.Equal(t, tc.wantUser, u)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
