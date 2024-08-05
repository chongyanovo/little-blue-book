package main

import (
	"github.com/ChongYanOvO/little-blue-book/internal/repository"
	"github.com/ChongYanOvO/little-blue-book/internal/repository/cache"
	"github.com/ChongYanOvO/little-blue-book/internal/repository/dao"
	"github.com/ChongYanOvO/little-blue-book/internal/service"
	"github.com/ChongYanOvO/little-blue-book/internal/web"
	"github.com/ChongYanOvO/little-blue-book/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func initWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB,
		ioc.InitRedis,
		dao.NewUserDao,
		cache.NewUserCache,
		cache.NewCodeCache,
		repository.NewCodeRepository,
		repository.NewUserRepository,
		service.NewCodeService,
		service.NewUserService,
		ioc.InitSmsService,
		web.NewUserHandler,
		ioc.InitGin,
		ioc.InitMiddlewares,
	)
	return new(gin.Engine)
}
