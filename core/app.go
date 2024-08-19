package core

import (
	"github.com/ChongYanOvO/little-blue-book/core/bootstrap"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Application struct {
	Config *bootstrap.Config
	DB     *gorm.DB
	Redis  redis.Cmdable
	Logger *zap.Logger
	Server *gin.Engine
}

// NewApplication 初始化 Application
func NewApplication(config *bootstrap.Config,
	db *gorm.DB,
	redis redis.Cmdable,
	logger *zap.Logger,
	server *gin.Engine) Application {
	return Application{
		Config: config,
		DB:     db,
		Redis:  redis,
		Logger: logger,
		Server: server,
	}
}
