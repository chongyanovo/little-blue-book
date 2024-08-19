//go:build wireinject
// +build wireinject

package main

import (
	"github.com/ChongYanOvO/little-blue-book/core"
	"github.com/ChongYanOvO/little-blue-book/core/bootstrap"
	"github.com/ChongYanOvO/little-blue-book/internal/handler"
	"github.com/ChongYanOvO/little-blue-book/internal/repository"
	"github.com/ChongYanOvO/little-blue-book/internal/repository/cache"
	"github.com/ChongYanOvO/little-blue-book/internal/repository/dao"
	"github.com/ChongYanOvO/little-blue-book/internal/service"
	"github.com/ChongYanOvO/little-blue-book/internal/service/sms"
	"github.com/google/wire"
)

var BaseProvider = wire.NewSet(
	bootstrap.NewViper,
	bootstrap.NewConfig,
	bootstrap.NewMysql,
	bootstrap.NewRedis,
	bootstrap.NewZap,
	bootstrap.NewMiddleware,
	bootstrap.NewServer,
	core.NewApplication,
)

var UserProvider = wire.NewSet(
	cache.NewCodeCache,
	cache.NewRedisUserCache,
	dao.NewUserDao,
	repository.NewCodeRepository,
	repository.NewUserRepository,
	sms.NewMemoryService,
	service.NewCodeService,
	service.NewUserService,
	handler.NewUserHandler,
)

func InitApp() (core.Application, error) {
	wire.Build(
		BaseProvider,
		UserProvider,
	)
	return core.Application{}, nil
}

func InitConfig() (*bootstrap.Config, error) {
	wire.Build(
		bootstrap.NewViper,
		bootstrap.NewConfig,
	)
	return &bootstrap.Config{}, nil
}
