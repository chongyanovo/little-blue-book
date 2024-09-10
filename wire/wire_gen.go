// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

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

// Injectors from wire.go:

func InitApp() (core.Application, error) {
	viper := bootstrap.NewViper()
	config := bootstrap.NewConfig(viper)
	logger := bootstrap.NewZap(config)
	db := bootstrap.NewMysql(config, logger)
	database := bootstrap.NewMongo(config, logger)
	cmdable := bootstrap.NewRedis(config)
	v := bootstrap.NewMiddlewares(logger)
	userDao := dao.NewUserDao(db, logger)
	userCache := cache.NewRedisUserCache(cmdable, logger)
	userRepository := repository.NewUserRepository(userDao, userCache, logger)
	userService := service.NewUserService(userRepository, logger)
	codeCache := cache.NewCodeCache(cmdable, logger)
	codeRepository := repository.NewCodeRepository(codeCache, logger)
	smsService := sms.NewMemoryService(logger)
	codeService := service.NewCodeService(codeRepository, smsService, logger)
	userHandler := handler.NewUserHandler(userService, codeService, logger)
	articleDao := article.NewArticleDao(db, logger)
	redisArticleCache := cache.NewRedisArticleCache(cmdable, logger)
	articleRepository := repository.NewArticleRepository(articleDao, redisArticleCache, logger)
	articleService := service.NewArticleService(articleRepository, logger)
	interactiveServiceImpl := service.NewInteractiveServiceImpl()
	articleHandler := handler.NewArticleHandler(articleService, interactiveServiceImpl, logger)
	engine := bootstrap.NewServer(v, userHandler, articleHandler)
	application := core.NewApplication(config, db, database, cmdable, logger, engine)
	return application, nil
}

func InitConfig() (*bootstrap.Config, error) {
	viper := bootstrap.NewViper()
	config := bootstrap.NewConfig(viper)
	return config, nil
}

func InitArticleHandler() (*handler.ArticleHandler, error) {
	viper := bootstrap.NewViper()
	config := bootstrap.NewConfig(viper)
	logger := bootstrap.NewZap(config)
	db := bootstrap.NewMysql(config, logger)
	articleDao := article.NewArticleDao(db, logger)
	cmdable := bootstrap.NewRedis(config)
	redisArticleCache := cache.NewRedisArticleCache(cmdable, logger)
	articleRepository := repository.NewArticleRepository(articleDao, redisArticleCache, logger)
	articleService := service.NewArticleService(articleRepository, logger)
	interactiveServiceImpl := service.NewInteractiveServiceImpl()
	articleHandler := handler.NewArticleHandler(articleService, interactiveServiceImpl, logger)
	return articleHandler, nil
}

func InitMysql() (*gorm.DB, error) {
	viper := bootstrap.NewViper()
	config := bootstrap.NewConfig(viper)
	logger := bootstrap.NewZap(config)
	db := bootstrap.NewMysql(config, logger)
	return db, nil
}

func InitMongo() (*mongo.Database, error) {
	viper := bootstrap.NewViper()
	config := bootstrap.NewConfig(viper)
	logger := bootstrap.NewZap(config)
	database := bootstrap.NewMongo(config, logger)
	return database, nil
}

// wire.go:

var BaseProvider = wire.NewSet(bootstrap.NewViper, bootstrap.NewConfig, bootstrap.NewMysql, bootstrap.NewMongo, bootstrap.NewRedis, bootstrap.NewZap, bootstrap.NewMiddlewares, bootstrap.NewServer, core.NewApplication)

var UserProvider = wire.NewSet(cache.NewCodeCache, cache.NewRedisUserCache, dao.NewUserDao, repository.NewCodeRepository, repository.NewUserRepository, sms.NewMemoryService, service.NewCodeService, service.NewUserService, handler.NewUserHandler)

var InteractiveProvider = wire.NewSet(cache.NewRedisInteractiveCache, wire.Bind(new(cache.InteractiveCache), new(*cache.RedisInteractiveCache)), dao.NewInteractiveDaoMysql, wire.Bind(new(dao.InteractiveDao), new(*dao.InteractiveDaoMysql)), repository.NewInteractiveRepositoryImpl, wire.Bind(new(repository.InteractiveRepository), new(*repository.InteractiveRepositoryImpl)), service.NewInteractiveServiceImpl, wire.Bind(new(service.InteractiveService), new(*service.InteractiveServiceImpl)))

var ArticleProvider = wire.NewSet(
	InteractiveProvider, article.NewArticleDao, cache.NewRedisArticleCache, wire.Bind(new(cache.ArticleCache), new(*cache.RedisArticleCache)), repository.NewArticleRepository, service.NewArticleService, handler.NewArticleHandler,
)
