package repository

import (
	"context"
	"github.com/ChongYanOvO/little-blue-book/internal/domain"
	"github.com/ChongYanOvO/little-blue-book/internal/repository/cache"
	cachemock "github.com/ChongYanOvO/little-blue-book/internal/repository/cache/mock"
	"github.com/ChongYanOvO/little-blue-book/internal/repository/dao"
	daomock "github.com/ChongYanOvO/little-blue-book/internal/repository/dao/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestUserRepositoryImpl_FindById(t *testing.T) {
	//now := time.Now().UnixMilli()
	testCases := []struct {
		name     string
		mock     func(ctl *gomock.Controller) (dao.UserDao, cache.UserCache)
		email    string
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "缓存未命中,查询成功",
			mock: func(ctl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				uc := cachemock.NewMockUserCache(ctl)
				uc.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{}, cache.ErrKeyNotExist)

				ud := daomock.NewMockUserDao(ctl)
				ud.EXPECT().FindById(gomock.Any(), int64(123)).
					Return(gomock.Any(), nil)
				return ud, uc
			},
			email: "test@example.com",
			wantUser: domain.User{
				Id:       123,
				Email:    "test@example.com",
				Password: "passwd",
				Phone:    "1234567890",
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			ud, uc := tc.mock(ctl)
			repo := NewUserRepository(ud, uc, nil)
			u, err := repo.FindByEmail(context.Background(), tc.email)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)
		})
	}
}
