package bootstrap

import (
	"context"
	"github.com/ChongYanOvO/little-blue-book/internal/handler/middleware"
	"github.com/ChongYanOvO/little-blue-book/pkg/ginx/middleware/accesslog"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewMiddleware(l *zap.Logger) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		NewLoginMiddleWare(),
		NewLoggerMiddleware(l),
	}
}

func NewLoginMiddleWare() gin.HandlerFunc {
	return middleware.NewLoginBuilder().
		IgnorePaths("/users/login").
		IgnorePaths("/users/signup").
		IgnorePaths("/users/login/code").
		Build()
}

func NewLoggerMiddleware(l *zap.Logger) gin.HandlerFunc {
	return accesslog.
		NewBuilder(func(ctx context.Context, log *accesslog.AccessLog) {
			l.Info("Http请求", zap.Any("日志", log))
		}).
		AllowRequestBody(true).
		AllowResponseBody(true).
		Build()
}
