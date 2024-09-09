package bootstrap

import (
	"context"
	"github.com/ChongYanOvO/little-blue-book/internal/handler/middleware"
	"github.com/ChongYanOvO/little-blue-book/pkg/ginx/middleware/accesslog"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func NewMiddlewares(l *zap.Logger) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		LoginMiddleWare(),
		LoggerMiddleware(l),
		CorsMiddleware(),
	}
}

// LoginMiddleWare 登录中间件
func LoginMiddleWare() gin.HandlerFunc {
	return middleware.NewLoginBuilder().
		IgnorePaths("/users/login").
		IgnorePaths("/users/signup").
		IgnorePaths("/users/login/code").
		Build()
}

// LoggerMiddleware 日志中间件
func LoggerMiddleware(l *zap.Logger) gin.HandlerFunc {
	return accesslog.
		NewBuilder(func(ctx context.Context, log *accesslog.AccessLog) {
			l.Info("Http请求", zap.Any("日志", log))
		}).
		AllowRequestBody(true).
		AllowResponseBody(true).
		Build()
}

// CorsMiddleware 跨域中间件
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, Token")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
