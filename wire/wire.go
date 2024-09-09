//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/ChongYanOvO/little-blue-book/core"
	"github.com/ChongYanOvO/little-blue-book/core/bootstrap"
	"github.com/ChongYanOvO/little-blue-book/internal/handler"
	"github.com/ChongYanOvO/little-blue-book/internal/repository"
	"github.com/ChongYanOvO/little-blue-book/internal/repository/cache"
	"github.com/ChongYanOvO/little-blue-book/internal/repository/dao"
	"github.com/ChongYanOvO/little-blue-book/internal/repository/dao/article"
	"github.com/ChongYanOvO/little-blue-book/internal/service"
	"github.com/ChongYanOvO/little-blue-book/internal/service/sms"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var BaseProvider = wire.NewSet(
	bootstrap.NewViper,
	bootstrap.NewConfig,
	bootstrap.NewMysql,
	bootstrap.NewMongo,
	bootstrap.NewRedis,
	bootstrap.NewZap,
	bootstrap.NewMiddlewares,
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

var ArticleProvider = wire.NewSet(
	article.NewArticleDao,
	cache.NewRedisArticleCache,
	wire.Bind(new(cache.ArticleCache), new(*cache.RedisArticleCache)),
	repository.NewArticleRepository,
	service.NewArticleService,
	handler.NewArticleHandler,
)

func InitApp() (core.Application, error) {
	wire.Build(
		BaseProvider,
		UserProvider,
		ArticleProvider,
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

func InitArticleHandler() (*handler.ArticleHandler, error) {
	wire.Build(
		BaseProvider,
		ArticleProvider,
	)
	return &handler.ArticleHandler{}, nil
}

func InitMysql() (*gorm.DB, error) {
	wire.Build(
		bootstrap.NewViper,
		bootstrap.NewConfig,
		bootstrap.NewMysql,
		bootstrap.NewZap,
	)
	return &gorm.DB{}, nil
}

func InitMongo() (*mongo.Database, error) {
	wire.Build(
		bootstrap.NewViper,
		bootstrap.NewConfig,
		bootstrap.NewMongo,
		bootstrap.NewZap,
	)
	return &mongo.Database{}, nil
}
