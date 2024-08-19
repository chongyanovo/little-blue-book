//go:build wireinject

package wire

import (
	"github.com/ChongYanOvO/little-blue-book/bootstrap"
	"github.com/ChongYanOvO/little-blue-book/internal/repository"
	"github.com/ChongYanOvO/little-blue-book/internal/repository/cache"
	"github.com/ChongYanOvO/little-blue-book/internal/repository/dao"
	"github.com/ChongYanOvO/little-blue-book/internal/service"
	"github.com/ChongYanOvO/little-blue-book/internal/service/sms"
	"github.com/ChongYanOvO/little-blue-book/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var BootstrapProviderSet = wire.NewSet(
	bootstrap.NewViper,
	bootstrap.NewConfig,
	bootstrap.NewMysql,
	bootstrap.NewZap,
	bootstrap.NewRedis,
	bootstrap.NewWeb,
)

var CacheProviderSet = wire.NewSet(
	cache.NewRedisUserCache,
	cache.NewCodeCache,
)

var DaoProviderSet = wire.NewSet(
	dao.NewUserDao,
)

var RepositoryProviderSet = wire.NewSet(
	repository.NewUserRepository,
	repository.NewCodeRepository,
)

var SmsProviderSet = wire.NewSet(
	//wire.Bind(new(sms.SmsService), new(*sms.MemoryService)),
	sms.NewMemoryService,
	//sms.NewLimiterService,
)

var ServiceProviderSet = wire.NewSet(
	service.NewUserService,
	service.NewCodeService,
)

var WebProviderSet = wire.NewSet(
	web.NewUserHandler,
)

func InitServer() (*gin.Engine, error) {
	wire.Build(
		BootstrapProviderSet,
		CacheProviderSet,
		DaoProviderSet,
		RepositoryProviderSet,
		SmsProviderSet,
		ServiceProviderSet,
		WebProviderSet,
	)
	return &gin.Engine{}, nil
}

func InitConfig() (*bootstrap.Config, error) {
	wire.Build(
		bootstrap.NewViper,
		bootstrap.NewConfig,
	)
	return &bootstrap.Config{}, nil
}
