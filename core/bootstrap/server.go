package bootstrap

import (
	"github.com/ChongYanOvO/little-blue-book/internal/handler"
	"github.com/gin-gonic/gin"
)

// ServerConfig server配置
type ServerConfig struct {
	Host string `mapstructure:"host" json:"host" yaml:"host"`
	Port int    `mapstructure:"port" json:"port" yaml:"port"`
}

type Server gin.Engine

// NewServer 创建server
func NewServer(middlewares []gin.HandlerFunc, uh *handler.UserHandler) *gin.Engine {
	server := gin.Default()

	server.Use(middlewares...)
	uh.RegisterRoutes(server)
	return server
}
